package ecs

import (
	"errors"
	"fmt"
	"sync"
)

// Stage system execute period:start->pre_update->update->pre_destroy->destroy
type Stage uint8

const (
	StageSyncBeforeStart Stage = iota
	StageStart
	StageSyncAfterStart

	StageSyncBeforePreUpdate
	StagePreUpdate
	StageSyncAfterPreUpdate

	StageSyncBeforeUpdate
	StageUpdate
	StageSyncAfterUpdate

	StageSyncBeforePostUpdate
	StagePostUpdate
	StageSyncAfterPostUpdate

	StageSyncBeforeDestroy
	StageDestroy
	StageSyncAfterDestroy
)

const StageMaxIndex = StageSyncAfterDestroy

// Order default suborder of system
type Order int32

const (
	OrderFront   Order = -1
	OrderInvalid Order = 0
	OrderAppend  Order = 99999999
	OrderDefault Order = OrderAppend
)

type ExecuteErrors struct {
	lock    sync.Mutex
	Err     error
	SubErrs []error
}

func (e *ExecuteErrors) Error() string {
	return e.Err.Error()
}

func (e *ExecuteErrors) Append(err error) {
	e.SubErrs = append(e.SubErrs, err)
}

func (e *ExecuteErrors) AppendWithLock(err error) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.Append(err)
}

type SystemTraverserList []SystemTraverser

type FlowConfig struct {
	ExecuteMode WorldExecuteMode
}

func (f *FlowConfig) initDefault() {
	f.ExecuteMode = ExecuteModeLinear
}

type FlowOption func(*FlowConfig)

func WithFlowSyncMode() FlowOption {
	return func(config *FlowConfig) {
		config.ExecuteMode = ExecuteModeLinear
	}
}

func WithFlowASyncMode() FlowOption {
	return func(config *FlowConfig) {
		config.ExecuteMode = ExecuteModeParallel
	}
}

type SystemTaskContext struct {
	isValid bool
	fn      func(ctx *SystemContext, event Event) error
}

// system execute flow
type flow struct {
	config           *FlowConfig
	world            *World
	traverserFactory func() SystemTraverser
	stages           [StageMaxIndex + 1]SystemTraverserList
	plan             []stageGroup
	workers          sync.WaitGroup
	systems          map[uint64]SystemInfo
}

func newSystemFlow(world *World, opt ...FlowOption) *flow {
	config := &FlowConfig{}
	config.initDefault()

	for _, option := range opt {
		option(config)
	}

	f := &flow{
		config:  config,
		world:   world,
		systems: map[uint64]SystemInfo{},
	}
	f.reset()
	return f
}

func (p *flow) reset() {
	p.plan = nil
	if p.config.ExecuteMode == ExecuteModeLinear {
		p.traverserFactory = NewSystemDefaultList
	} else {
		p.traverserFactory = NewSystemRelatedGroups
	}
	for stage := range StageMaxIndex + 1 {
		p.stages[stage] = SystemTraverserList{}
		stlFront := p.traverserFactory()
		stlFront.setOrder(OrderFront)
		stlAppend := p.traverserFactory()
		stlAppend.setOrder(OrderAppend)
		p.stages[stage] = append(p.stages[stage], stlFront, stlAppend)
	}
}

func (p *flow) flushTempTask() error {
	tasks, syncTasks, clean := p.world.flushPendingOperate()
	defer clean()

	if err := p.runOpTasks(tasks); err != nil {
		return err
	}
	// 第二相：组表成员资格 sync，依赖第一相完成后的 CSet/compound 终态
	return p.runOpTasks(syncTasks)
}

func (p *flow) runOpTasks(tasks []func()) error {
	if len(tasks) == 0 {
		return nil
	}
	wg := &sync.WaitGroup{}
	switch p.config.ExecuteMode {
	case ExecuteModeLinear:
		for _, task := range tasks {
			task()
		}
	case ExecuteModeParallel:
		wg.Add(len(tasks))
		for _, task := range tasks {
			fn := func() {
				task()
				wg.Done()
			}
			go fn()
		}
		wg.Wait()
	default:
		return fmt.Errorf("invalid flow mode: %d", p.config.ExecuteMode)
	}
	return nil
}

func (p *flow) Execute(event Event) error {
	var err error
	// 帧同步点：先处理实体销毁，再合并组件操作
	p.world.flushPendingDestroy()
	err = p.flushTempTask()
	if err != nil {
		return err
	}
	switch p.config.ExecuteMode {
	case ExecuteModeLinear:
		err = p.executeLinear(event)
	case ExecuteModeParallel:
		err = p.executeParallel(event)
	default:
		err = fmt.Errorf("invalid flow mode: %d", p.config.ExecuteMode)
	}
	if err != nil {
		return err
	}
	err = p.world.clearDisposable()
	if err != nil {
		return err
	}
	// 游牧组件每帧末清空
	err = p.world.clearNomadic()
	if err != nil {
		return err
	}
	err = p.flushTempTask()
	if err != nil {
		return err
	}
	return nil
}

func (p *flow) getSystemTask(info SystemInfo, stage Stage) (SystemTaskContext, error) {
	if !info.isValid() {
		return SystemTaskContext{}, errors.New("invalid system")
	}
	state := info.getState()
	if state == SystemStateStart && stage > StageSyncAfterStart {
		state = SystemStateUpdate
		info.setState(state)
	}
	switch state {
	case SystemStateStart:
		if stage > StageSyncAfterStart {
			return SystemTaskContext{}, nil
		}
	case SystemStateUpdate:
		if stage < StageSyncBeforePreUpdate || stage > StageSyncAfterPostUpdate {
			return SystemTaskContext{}, nil
		}
	case SystemStateDestroy:
		if stage < StageSyncBeforeDestroy {
			return SystemTaskContext{}, nil
		}
	default:
		return SystemTaskContext{}, nil
	}
	fn := info.callback(stage)
	if fn != nil && stage == StageSyncAfterDestroy {
		info.setState(SystemStateDestroyed)
	}
	return SystemTaskContext{fn: fn, isValid: fn != nil}, nil
}

func (p *flow) runSystem(info SystemInfo, stage Stage, event Event) error {
	task, err := p.getSystemTask(info, stage)
	if err != nil || !task.isValid {
		return err
	}
	ctx := info.getContext()
	ctx.constraint.activate()
	defer ctx.constraint.deactivate()
	return task.fn(ctx, event)
}

func (p *flow) executeLinear(event Event) error {
	var errs *ExecuteErrors
	for _, group := range p.schedule() {
		for _, info := range group.systems {
			if err := p.runSystem(info, group.stage, event); err != nil {
				if errs == nil {
					errs = &ExecuteErrors{Err: errors.New("errors found in executeLinear")}
				}
				errs.Append(err)
			}
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}

func isSyncStage(stage Stage) bool {
	switch stage {
	case StageSyncBeforeStart, StageSyncAfterStart,
		StageSyncBeforePreUpdate, StageSyncAfterPreUpdate,
		StageSyncBeforeUpdate, StageSyncAfterUpdate,
		StageSyncBeforePostUpdate, StageSyncAfterPostUpdate,
		StageSyncBeforeDestroy, StageSyncAfterDestroy:
		return true
	default:
		return false
	}
}

func (p *flow) executeParallel(event Event) error {
	var errs *ExecuteErrors
	record := func(err error) {
		if err != nil {
			if errs == nil {
				errs = &ExecuteErrors{Err: errors.New("errors found in executeParallel")}
			}
			errs.Append(err)
		}
	}
	for _, group := range p.schedule() {
		if isSyncStage(group.stage) {
			for _, info := range group.systems {
				record(p.runSystem(info, group.stage, event))
			}
			continue
		}
		for _, batch := range group.batches {
			if len(batch.systems) == 1 {
				record(p.runSystem(batch.systems[0], group.stage, event))
				continue
			}
			// The caller executes one independent system while workers handle its peers.
			for i := 1; i < len(batch.systems); i++ {
				p.workers.Go(func() {
					batch.results[i] = p.runSystem(batch.systems[i], group.stage, event)
				})
			}
			record(p.runSystem(batch.systems[0], group.stage, event))
			p.workers.Wait()
			for i := 1; i < len(batch.results); i++ {
				record(batch.results[i])
				batch.results[i] = nil
			}
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}

type systemBatch struct {
	systems []SystemInfo
	results []error
}

type stageGroup struct {
	stage   Stage
	systems []SystemInfo
	batches []systemBatch
}

// Registration is closed once the world starts. Reuse the stage/Order schedule
// and parallel conflict groups until another system is registered.
func (p *flow) schedule() []stageGroup {
	if p.plan == nil && len(p.systems) != 0 {
		for stage := range StageMaxIndex + 1 {
			for _, list := range p.stages[stage] {
				if list.count() == 0 {
					continue
				}
				group := stageGroup{stage: stage, systems: list.all()}
				if p.config.ExecuteMode == ExecuteModeParallel && !isSyncStage(stage) {
					for systems := range list.independentGroups() {
						batch := systemBatch{systems: systems}
						if len(systems) > 1 {
							batch.results = make([]error, len(systems))
						}
						group.batches = append(group.batches, batch)
					}
				}
				p.plan = append(p.plan, group)
			}
		}
	}
	return p.plan
}

// register method only in world init or func init(){}
func (p *flow) register(system SystemInfo) {
	if p.world.getStatus() != WorldStatusInitialized {
		panic("system register only in world init")
	}

	order := system.getOrder()
	if order > OrderAppend {
		order = OrderAppend
	}

	for stage := range StageMaxIndex + 1 {

		if system.callback(stage) == nil {
			continue
		}

		sl := p.stages[stage]
		if order == OrderFront {
			p.stages[stage][0].add(system)
		} else if order == OrderAppend {
			p.stages[stage][len(sl)-1].add(system)
		} else {
			for i, v := range sl {
				if order == v.getOrder() {
					v.add(system)
					break
				} else if order < v.getOrder() {
					sg := p.traverserFactory()
					sg.setOrder(order)
					sg.add(system)
					// 在 i 前插入新分组，显式构造新切片避免污染原底层数组
					ns := make(SystemTraverserList, 0, len(sl)+1)
					ns = append(ns, sl[:i]...)
					ns = append(ns, sg)
					ns = append(ns, sl[i:]...)
					p.stages[stage] = ns
					break
				}
			}
		}
	}

	p.plan = nil
	p.systems[system.id()] = system
	p.world.opLog.registerSystem(system.id())

	system.setState(SystemStateStart)
}

func (p *flow) stop() {
	p.reset()
}

func (p *flow) DebugInfo() {
	m := map[Stage]string{
		StageSyncBeforeStart: "StageSyncBeforeStart",
		StageStart:           "StageStart",
		StageSyncAfterStart:  "StageSyncAfterStart",

		StageSyncBeforePreUpdate: "StageSyncBeforePreUpdate",
		StagePreUpdate:           "StagePreUpdate",
		StageSyncAfterPreUpdate:  "StageSyncAfterPreUpdate",

		StageSyncBeforeUpdate: "StageSyncBeforeUpdate",
		StageUpdate:           "StageUpdate",
		StageSyncAfterUpdate:  "StageSyncAfterUpdate",

		StageSyncBeforePostUpdate: "StageSyncBeforePostUpdate",
		StagePostUpdate:           "StagePostUpdate",
		StageSyncAfterPostUpdate:  "StageSyncAfterPostUpdate",

		StageSyncBeforeDestroy: "StageSyncBeforeDestroy",
		StageDestroy:           "StageDestroy",
		StageSyncAfterDestroy:  "StageSyncAfterDestroy",
	}
	var debugStr string
	debugStr += "┌──────────────── # System Info # ─────────────────\n"
	debugStr += fmt.Sprintf("├─ Total: %d\n", len(p.systems))

	var output []string
	var sq SystemTraverserList
	for stage := range StageMaxIndex + 1 {
		var slContent []string
		sq = p.stages[stage]
		for i, sl := range sq {
			batchTotal := sl.getBatchCount()
			batch := 0
			var batchContent []string
			for ss := range sl.independentGroups() {
				if systemCount := len(ss); systemCount != 0 {
					str := "│     │  └─ "
					if batch == batchTotal-1 {
						str = "│        └─ "
					}
					for i := 0; i < systemCount; i++ {
						str += fmt.Sprintf("%s ", ss[i].name())
					}
					if batch == batchTotal-1 {
						batchContent = append(batchContent, fmt.Sprintf("│     └─ Batch %d", batch))
					} else {
						batchContent = append(batchContent, fmt.Sprintf("│     ├─ Batch %d", batch))
					}
					batchContent = append(batchContent, str)
				}
				batch++
			}
			if len(batchContent) > 0 {
				s := make([]string, 0, len(batchContent)+1)
				if i == len(sq)-1 {
					s = append(s, fmt.Sprintf("│  └─ Order %d", i))
				} else {
					s = append(s, fmt.Sprintf("│  ├─ Order %d", i))
				}
				s = append(s, batchContent...)
				slContent = append(slContent, s...)
			}
		}
		if len(slContent) > 0 {
			s := make([]string, 0, len(slContent)+1)
			if stage == StageMaxIndex {
				s = append(s, fmt.Sprintf("└─ Stage %s", m[stage]))
			} else {
				s = append(s, fmt.Sprintf("├─ Stage %s", m[stage]))
			}
			s = append(s, slContent...)
			output = append(output, s...)
		}
	}

	for _, v := range output {
		debugStr += v + "\n"
	}
	debugStr += "└────────────── # System Info End # ───────────────\n"
}
