package ecs

import (
	"fmt"
	"reflect"
	"sync"
)

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

// Stage system execute period:start->pre_update->update->pre_destroy->destroy
type Stage uint32

// Order default suborder of system
type Order int32

const (
	OrderFront   Order = -1
	OrderInvalid Order = 0
	OrderAppend  Order = 99999999
	OrderDefault Order = OrderAppend
)

// SystemGroupList extension of system group slice
type SystemGroupList []*SystemGroup

// system execute flow
type systemFlow struct {
	world     *ecsWorld
	stages    map[Stage]SystemGroupList
	stageList []Stage
	systems   map[reflect.Type]ISystem
	wg        *sync.WaitGroup
}

func newSystemFlow(runtime *ecsWorld) *systemFlow {
	sf := &systemFlow{
		world:   runtime,
		systems: map[reflect.Type]ISystem{},
		wg:      &sync.WaitGroup{},
	}
	sf.init()
	return sf
}

// initialize the system flow
func (p *systemFlow) init() {
	p.stageList = []Stage{
		StageSyncBeforeStart,
		StageStart,
		StageSyncAfterStart,

		StageSyncBeforePreUpdate,
		StagePreUpdate,
		StageSyncAfterPreUpdate,

		StageSyncBeforeUpdate,
		StageUpdate,
		StageSyncAfterUpdate,

		StageSyncBeforePostUpdate,
		StagePostUpdate,
		StageSyncAfterPostUpdate,

		StageSyncBeforeDestroy,
		StageDestroy,
		StageSyncAfterDestroy,
	}
	p.reset()
}

func (p *systemFlow) reset() {
	p.stages = make(map[Stage]SystemGroupList)
	for _, value := range p.stageList {
		p.stages[value] = SystemGroupList{}
		sgFront := NewSystemGroup()
		sgFront.order = OrderFront
		sgAppend := NewSystemGroup()
		sgAppend.order = OrderAppend
		p.stages[value] = append(p.stages[value], sgFront, sgAppend)
	}
}

func (p *systemFlow) flushTempTask() {
	tasks := p.world.components.getTempTasks()
	p.wg.Add(len(tasks))
	for _, task := range tasks {
		wg := p.wg
		fn := task
		p.world.addJob(func() {
			fn()
			wg.Done()
		})
	}
	p.wg.Wait()
}

func (p *systemFlow) systemUpdate(event Event) {
	var sq SystemGroupList
	var sys ISystem
	var imp bool = false
	var runSync bool = false
	var fn func(event Event)
	for _, period := range p.stageList {
		sq = p.stages[period]
		for _, sl := range sq {
			if sl.systemCount() == 0 {
				continue
			}
			for ss := sl.Begin(); !sl.End(); ss = sl.Next() {
				if systemCount := len(ss); systemCount != 0 {
					for i := 0; i < systemCount; i++ {
						sys = ss[i]

						if !sys.isValid() {
							continue
						}

						imp = false
						runSync = false
						state := ss[i].getState()

						if period > StageSyncAfterStart {
							if state == SystemStateStart {
								state = SystemStateUpdate
								sys.setState(SystemStateUpdate)
							}
						}

						if state == SystemStateStart {
							if period > StageSyncAfterStart {
								continue
							}
							switch period {
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
							if period < StageSyncBeforePreUpdate || period > StageSyncAfterPostUpdate {
								continue
							}
							switch period {
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
							if period < StageSyncBeforeDestroy {
								continue
							}
							switch period {
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

								sys.setState(SystemStateDestroyed)
							}
						}

						if !imp {
							continue
						}
						if runSync {
							sys.setExecuting(true)
							sys.setSecurity(true)
							fn(event)
							sys.setSecurity(false)
							sys.setExecuting(false)
						} else {
							wrapper := func(fn func(event2 Event), e Event) func() {
								sys.setExecuting(true)
								return func() {
									defer func() {
										sys.setExecuting(false)
										p.wg.Done()
									}()
									fn(e)
								}
							}
							p.wg.Add(1)
							p.world.addJob(wrapper(fn, event))
						}
					}
				}
				p.wg.Wait()
			}
		}
	}
}

func (p *systemFlow) run(event Event) {
	reporter := p.world.metrics.NewReporter("system_flow_run")
	reporter.Start()

	//Log.Info("system flow # Temp Task Execute #")
	p.flushTempTask()
	reporter.Sample("Temp Task Execute")

	//Log.Info("system flow # Logic #")
	p.systemUpdate(event)
	reporter.Sample("system execute")

	//Log.Info("system flow # Clear Disposable #")
	p.world.components.clearDisposable()
	reporter.Sample("Clear Disposable")

	p.flushTempTask()
	reporter.Sample("Temp Task Execute")

	reporter.Stop()
	reporter.Print()
}

// register method only in world init or func init(){}
func (p *systemFlow) register(system ISystem) {
	if p.world.getStatus() != WorldStatusInitialized {
		panic("system register only in world init")
	}

	//init function call
	system.baseInit(p.world, system)

	order := system.Order()
	if order > OrderAppend {
		Log.Errorf("system order must less then %d, resort order to %d", OrderAppend+1, OrderAppend)
		order = OrderAppend
	}

	for _, period := range p.stageList {

		if !p.isImpEvent(system, period) {
			continue
		}

		sl := p.stages[period]
		if order == OrderFront {
			p.stages[period][0].insert(system)
		} else if order == OrderAppend {
			p.stages[period][len(sl)-1].insert(system)
		} else {
			for i, v := range sl {
				if order == v.order {
					v.insert(system)
					break
				} else if order < v.order {
					sg := NewSystemGroup()
					sg.order = order
					sg.insert(system)
					temp := append(SystemGroupList{}, sl[i-1:]...)
					p.stages[period] = append(append(sl[:i-1], sg), temp...)
					break
				}
			}
		}
	}

	p.systems[system.Type()] = system
}

func (p *systemFlow) isImpEvent(system ISystem, period Stage) bool {
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

func (p *systemFlow) stop() {
	p.reset()
}

func (p *systemFlow) SystemInfoPrint() {
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
	Log.Infof("┌──────────────── # System Info # ─────────────────")
	Log.Infof("├─ Total: %d", len(p.systems))

	var output []string
	var sq SystemGroupList
	for pi, period := range p.stageList {
		var slContent []string
		sq = p.stages[period]
		for i, sl := range sq {
			sl.resort()
			batchTotal := sl.batchCount()
			batch := 0
			var batchContent []string
			for ss := sl.Begin(); !sl.End(); ss = sl.Next() {
				if systemCount := len(ss); systemCount != 0 {

					str := "│     │  └─ "
					if batch == batchTotal-1 {
						str = "│        └─ "
					}
					for i := 0; i < systemCount; i++ {
						str += fmt.Sprintf("%s ", ss[i].Type().Name())
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
			if pi == len(p.stageList)-1 {
				s = append(s, fmt.Sprintf("└─ Stage %s", m[period]))
			} else {
				s = append(s, fmt.Sprintf("├─ Stage %s", m[period]))
			}
			s = append(s, slContent...)
			output = append(output, s...)
		}
	}

	for _, v := range output {
		Log.Info(v)
	}
	Log.Infof("└────────────── # System Info End # ───────────────")
}
