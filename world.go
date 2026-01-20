package ecs

import (
	"runtime"
	"time"

	rockmem "github.com/zllangct/rockmem/golang"
)

type WorldStatus int

const (
	WorldStatusInitializing WorldStatus = iota
	WorldStatusInitialized
	WorldStatusRunning
	WorldStatusStop
)

type WorldExecuteMode uint8

const (
	ExecuteModeInvalid WorldExecuteMode = iota
	ExecuteModeLinear
	ExecuteModeParallel
)

type TaskExecutor func(task func()) error

type TaskExecutorFactory func() TaskExecutor

type World interface {
	Update() error
	Destroy() error
	NewEntity(opts ...EntityOption) *EntityInfo
	RegisterLight(system LightSystem, Option ...SystemOption) error
	RegisterStandard(system SystemStandard) error
	Optimize(t time.Duration, force bool) error
	MarshalTo(writer rockmem.Writer) (uint, error)
	Marshal() *SerializableWorldData
}

type WorldOption func(config *WorldConfig)

type WorldConfig struct {
	ExecuteMode  WorldExecuteMode
	AutoOptimize bool
	Rate         int
}

func (s *WorldConfig) initDefault() {
	s.ExecuteMode = ExecuteModeLinear
	s.AutoOptimize = false
	s.Rate = 30
}

func WithWorldSyncMode() WorldOption {
	return func(config *WorldConfig) {
		config.ExecuteMode = ExecuteModeLinear
	}
}

func WithWorldASyncMode() WorldOption {
	return func(config *WorldConfig) {
		config.ExecuteMode = ExecuteModeParallel
	}
}

func WithWorldAutoOptimize() WorldOption {
	return func(config *WorldConfig) {
		config.AutoOptimize = true
	}
}

func WithWorldDefaultUpdateRate(rate int) WorldOption {
	return func(config *WorldConfig) {
		config.Rate = rate
	}
}

type serializableWorld struct {
	idGenerator     *EntityIDGenerator
	entities        *EntitySet
	components      map[ComponentIntType]ComponentSet
	disposableTypes []ComponentIntType
	nomadicTypes    []ComponentIntType
	frame           uint64
}

type world struct {
	serializableWorld
	status              WorldStatus
	config              *WorldConfig
	taskExecutorFactory TaskExecutorFactory
	optimizer           *optimizer
	opLog               *OpLog

	systems *flow

	lastUpdate time.Time
	delta      time.Duration
}

func NewWorld(opts ...WorldOption) World {
	c := &WorldConfig{}
	c.initDefault()
	for _, opt := range opts {
		opt(c)
	}

	w := &world{}
	w.status = WorldStatusInitializing
	w.config = c
	w.idGenerator = NewEntityIDGenerator(1024, 10)
	w.entities = NewEntitySet()
	w.components = map[ComponentIntType]ComponentSet{}
	flowExecOpt := WithFlowSyncMode()
	if c.ExecuteMode == ExecuteModeParallel {
		flowExecOpt = WithFlowASyncMode()
	}
	w.systems = newSystemFlow(w, flowExecOpt)
	w.opLog = NewOpLog(w, runtime.NumCPU())
	w.optimizer = newOptimizer(w)
	w.status = WorldStatusInitialized
	return w
}

func (w *world) NewEntity(opts ...EntityOption) *EntityInfo {
	return w.newEntity(opts...)
}

func (w *world) newEntity(opts ...EntityOption) *EntityInfo {
	config := EntityConfig{}
	config.initDefault()
	for _, opt := range opts {
		opt(&config)
	}
	ins := w.entities.Add(EntityInfo{
		world:    w,
		entity:   w.idGenerator.NewID(),
		compound: NewCompound(4),
	})
	if len(config.comps) > 0 {
		ins.Add(config.comps...)
	}
	return ins
}

func (w *world) getTaskExecutor() TaskExecutor {
	if w.taskExecutorFactory != nil {

	}
	return w.taskExecutorFactory()
}

func (w *world) getComponentSet(it ComponentIntType) (ComponentSet, bool) {
	s, ok := w.components[it]
	return s, ok
}

func (w *world) componentOp(op Operate) {
	w.opLog.operate(op)
}

func (w *world) AddNomadic(comps ...Component) {
	for _, comp := range comps {
		if !comp.IsNomadic() {
			continue
		}
		op := Operate{
			Op:   ComponentOperateAdd,
			Comp: comp,
		}
		w.componentOp(op)
	}
}

func (w *world) Update() error {
	if w.lastUpdate.IsZero() {
		w.lastUpdate = time.Now()
	}

	w.frame++
	now := time.Now()
	w.delta = now.Sub(w.lastUpdate)
	e := Event{
		Frame: w.frame,
		Delta: w.delta,
	}
	err := w.systems.Execute(e)
	if err != nil {
		return err
	}

	if w.config.AutoOptimize {
		elapsed := now.Sub(w.lastUpdate)
		t := (time.Second / time.Duration(w.config.Rate)) - elapsed
		err = w.Optimize(t, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *world) Destroy() error {
	return nil
}

func (w *world) Optimize(t time.Duration, force bool) error {
	if w.optimizer == nil || t < 0 {
		return nil
	}
	w.optimizer.optimize(t, force)
	return nil
}

func (w *world) RegisterStandard(system SystemStandard) error {
	s := newSystem(w, system, SystemTypeStandard)
	s.setState(SystemStateInit)

	ctx := &SystemInitContext{
		b: *s.getContext(),
	}
	ctx.b.constraint.reset()
	err := TryAndReport(func() error {
		return system.Init(ctx)
	})
	ctx.b.constraint.setOutdated()
	if err != nil {
		return err
	}

	s.init(ctx.opts...)
	w.systems.register(s)
	return nil
}

func (w *world) RegisterLight(system LightSystem, opts ...SystemOption) error {
	s := newSystem(w, system, SystemTypeLight)
	s.setState(SystemStateInit)
	s.init(opts...)
	w.systems.register(s)
	return nil
}

func (w *world) flushPendingOperate() ([]func(), func()) {
	return w.opLog.getOpTasks()
}

func (w *world) clearDisposable() error {
	for _, it := range w.disposableTypes {
		set, ok := w.getComponentSet(it)
		if !ok {
			continue
		}
		for _, index := range set.EntityIndexes() {
			e := w.entities.getByIndex(index)
			e.compound.Remove(it)
		}

		set.Reset()
	}
	return nil
}

func (w *world) clearNomadic() error {
	for _, it := range w.nomadicTypes {
		set, ok := w.getComponentSet(it)
		if !ok {
			continue
		}
		set.Reset()
	}
	return nil
}

func (w *world) setStatus(status WorldStatus) {
	w.status = status
}

func (w *world) getStatus() WorldStatus {
	return w.status
}

func (w *world) getEntityInfo(entity Entity) (*EntityInfo, bool) {
	return w.entities.Get(entity)
}
