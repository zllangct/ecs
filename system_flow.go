package ecs

import (
	"errors"
	"fmt"
	"iter"
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
	isAsync bool
}

// system execute flow
type flow struct {
	config           *FlowConfig
	world            *world
	traverserFactory func() SystemTraverser
	stages           map[Stage]SystemTraverserList
	systems          map[uint64]SystemInfo
}

func newSystemFlow(world *world, opt ...FlowOption) *flow {
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
	p.stages = make(map[Stage]SystemTraverserList)
	if p.config.ExecuteMode == ExecuteModeLinear {
		p.traverserFactory = NewSystemDefaultList
	} else {
		p.traverserFactory = NewSystemRelatedGroups
	}
	for stage := range StageMaxIndex {
		p.stages[stage] = SystemTraverserList{}
		stlFront := p.traverserFactory()
		stlFront.setOrder(OrderFront)
		stlAppend := p.traverserFactory()
		stlAppend.setOrder(OrderAppend)
		p.stages[stage] = append(p.stages[stage], stlFront, stlAppend)
	}
}

func (p *flow) flushTempTask() error {
	wg := &sync.WaitGroup{}
	tasks, clean := p.world.flushPendingOperate()
	defer clean()

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
	err = p.flushTempTask()
	if err != nil {
		return err
	}
	return nil
}

func (p *flow) getSystemTask(info SystemInfo, stage Stage) (ctx SystemTaskContext, err error) {
	if !info.isValid() {
		err = errors.New("invalid system")
		return
	}
	var imp bool = false
	var runSync bool = false
	var fn func(ctx *SystemContext, event Event) error

	state := info.getState()
	sys := info.getRaw()

	if stage > StageSyncAfterStart {
		if state == SystemStateStart {
			state = SystemStateUpdate
			info.setState(SystemStateUpdate)
		}
	}

	if state == SystemStateStart {
		if stage > StageSyncAfterStart {
			return
		}
		switch stage {
		case StageSyncBeforeStart:
			system, ok := sys.(SyncBeforeStartReceiver)
			fn = system.SyncBeforeStart
			imp = ok
			runSync = true
		case StageStart:
			system, ok := sys.(StartReceiver)
			fn = system.Start
			imp = ok
			runSync = false
		case StageSyncAfterStart:
			system, ok := sys.(SyncAfterStartReceiver)
			fn = system.SyncAfterStart
			imp = ok
			runSync = true
		}
	} else if state == SystemStateUpdate {
		if stage < StageSyncBeforePreUpdate || stage > StageSyncAfterPostUpdate {
			return
		}
		switch stage {
		case StageSyncBeforePreUpdate:
			system, ok := sys.(SyncBeforePreUpdateReceiver)
			fn = system.SyncBeforePreUpdate
			imp = ok
			runSync = true
		case StagePreUpdate:
			system, ok := sys.(PreUpdateReceiver)
			fn = system.PreUpdate
			imp = ok
			runSync = true
		case StageSyncAfterPreUpdate:
			system, ok := sys.(SyncAfterPreUpdateReceiver)
			fn = system.SyncAfterPreUpdate
			imp = ok
			runSync = true

		case StageSyncBeforeUpdate:
			system, ok := sys.(SyncBeforeUpdateReceiver)
			fn = system.SyncBeforeUpdate
			imp = ok
			runSync = true
		case StageUpdate:
			system, ok := sys.(UpdateReceiver)
			fn = system.Update
			imp = ok
			runSync = false
		case StageSyncAfterUpdate:
			system, ok := sys.(SyncAfterUpdateReceiver)
			fn = system.SyncAfterUpdate
			imp = ok
			runSync = true

		case StageSyncBeforePostUpdate:
			system, ok := sys.(SyncBeforePostUpdateReceiver)
			fn = system.SyncBeforePostUpdate
			imp = ok
			runSync = true
		case StagePostUpdate:
			system, ok := sys.(PostUpdateReceiver)
			fn = system.PostUpdate
			imp = ok
			runSync = false
		case StageSyncAfterPostUpdate:
			system, ok := sys.(SyncAfterPostUpdateReceiver)
			fn = system.SyncAfterPostUpdate
			imp = ok
			runSync = true
		}
	} else if state == SystemStateDestroy {
		if stage < StageSyncBeforeDestroy {
			return
		}
		switch stage {
		case StageSyncBeforeDestroy:
			system, ok := sys.(SyncBeforeDestroyReceiver)
			fn = system.SyncBeforeDestroy
			imp = ok
			runSync = true
		case StageDestroy:
			system, ok := sys.(DestroyReceiver)
			fn = system.Destroy
			imp = ok
			runSync = false
		case StageSyncAfterDestroy:
			system, ok := sys.(SyncAfterPostDestroyReceiver)
			fn = system.SyncAfterDestroy
			imp = ok
			runSync = true

			info.setState(SystemStateDestroyed)
		}
	}

	if !imp {
		return
	}

	return SystemTaskContext{
		fn:      fn,
		isAsync: !runSync,
		isValid: true,
	}, nil
}

func (p *flow) executeLinear(event Event) error {
	errs := &ExecuteErrors{}
	for stage, sl := range p.traverserIter() {
		for _, info := range sl.all() {
			task, err := p.getSystemTask(info, stage)
			if err != nil {
				errs.Append(err)
				continue
			}
			if task.isValid {
				ctx := info.getContext()
				ctx.constraint.reset()
				err := task.fn(ctx, event)
				ctx.constraint.setOutdated()
				if err != nil {
					errs.Append(err)
					continue
				}
			}
		}
	}
	if len(errs.SubErrs) > 0 {
		errs.Err = errors.New("errors found in executeLinear")
		return errs
	}
	return nil
}

func (p *flow) executeParallel(event Event) error {
	errs := &ExecuteErrors{}
	wg := &sync.WaitGroup{}
	for stage, sl := range p.traverserIter() {
		for batch := range sl.independentGroups() {
			for _, info := range batch {
				task, err := p.getSystemTask(info, stage)
				if err != nil {
					errs.AppendWithLock(err)
					continue
				}
				if task.isValid {
					ctx := info.getContext()
					wg.Add(1)
					go func() {
						defer wg.Done()
						ctx.constraint.reset()
						err := task.fn(ctx, event)
						ctx.constraint.setOutdated()
						if err != nil {
							errs.AppendWithLock(err)
						}
					}()
				}
			}
			wg.Wait()
		}
	}
	if len(errs.SubErrs) > 0 {
		errs.Err = errors.New("errors found in executeParallel")
		return errs
	}
	return nil
}

func (p *flow) traverserIter() iter.Seq2[Stage, SystemTraverser] {
	return func(yield func(Stage, SystemTraverser) bool) {
		for stage := range StageMaxIndex {
			sq := p.stages[stage]
			for _, sl := range sq {
				if sl.count() == 0 {
					continue
				}
				yield(stage, sl)
			}
		}
	}
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

	for stage := range StageMaxIndex {

		if !system.impl(stage) {
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
					temp := append(SystemTraverserList{}, sl[i-1:]...)
					p.stages[stage] = append(append(sl[:i-1], sg), temp...)
					break
				}
			}
		}
	}

	p.systems[system.id()] = system

	system.setState(SystemStateStart)
}

func (p *flow) isImpEvent(info SystemInfo, period Stage) bool {
	system := info.getRaw()
	imp := false
	switch period {
	case StageSyncBeforeStart:
		_, imp = system.(SyncBeforeStartReceiver)
	case StageStart:
		_, imp = system.(StartReceiver)
	case StageSyncAfterStart:
		_, imp = system.(SyncAfterStartReceiver)
	case StageSyncBeforePreUpdate:
		_, imp = system.(SyncBeforePreUpdateReceiver)
	case StagePreUpdate:
		_, imp = system.(PreUpdateReceiver)
	case StageSyncAfterPreUpdate:
		_, imp = system.(SyncAfterPreUpdateReceiver)
	case StageSyncBeforeUpdate:
		_, imp = system.(SyncBeforeUpdateReceiver)
	case StageUpdate:
		_, imp = system.(UpdateReceiver)
	case StageSyncAfterUpdate:
		_, imp = system.(SyncAfterUpdateReceiver)
	case StageSyncBeforePostUpdate:
		_, imp = system.(SyncBeforePostUpdateReceiver)
	case StagePostUpdate:
		_, imp = system.(PostUpdateReceiver)
	case StageSyncAfterPostUpdate:
		_, imp = system.(SyncAfterPostUpdateReceiver)
	case StageSyncBeforeDestroy:
		_, imp = system.(SyncBeforeDestroyReceiver)
	case StageDestroy:
		_, imp = system.(DestroyReceiver)
	case StageSyncAfterDestroy:
		_, imp = system.(SyncAfterPostDestroyReceiver)
	}
	return imp
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
	for stage := range StageMaxIndex {
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
