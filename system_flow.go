package ecs

import (
	"sync"
)

const (
	StageStart Stage = iota
	StagePreUpdate
	StageUpdate
	StagePostUpdate
	StageDestroy
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

//system execute flow
type systemFlow struct {
	lock      sync.Mutex
	world     *ecsWorld
	stages    map[Stage]SystemGroupList
	stageList []Stage
	systems   sync.Map
	wg        *sync.WaitGroup
}

func newSystemFlow(runtime *ecsWorld) *systemFlow {
	sf := &systemFlow{
		world: runtime,
		wg:    &sync.WaitGroup{},
	}
	sf.init()
	return sf
}

//initialize the system flow
func (p *systemFlow) init() {
	p.stageList = []Stage{
		StageStart,
		StagePreUpdate,
		StageUpdate,
		StagePostUpdate,
		StageDestroy,
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

func (p *systemFlow) run(event Event) {
	p.lock.Lock()
	defer p.lock.Unlock()

	removeList := map[int64]ISystem{}

	//Log.Info("system flow # Clear Disposable #")
	p.world.components.clearDisposable()

	//Log.Info("system flow # Temp Task Execute #")
	tasks := p.world.components.getTempTasks()
	p.wg.Add(len(tasks))
	for _, task := range tasks {
		wg := p.wg
		fn := task
		Runtime.addJob(func() {
			fn()
			wg.Done()
		})
	}
	p.wg.Wait()

	p.world.components.collectorRun()

	var sq SystemGroupList

	//Log.Info("system flow # Event Dispatch #")
	for _, period := range p.stageList {
		sq = p.stages[period]
		for _, sl := range sq {
			sl.reset()
			for ss := sl.next(); len(ss) > 0; ss = sl.next() {
				if systemCount := len(ss); systemCount != 0 {
					for i := 0; i < systemCount; i++ {
						fn := ss[i].eventDispatch
						p.wg.Add(1)
						wg := p.wg
						Runtime.addJob(func() {
							defer wg.Done()
							fn()
						})
					}
				}
				p.wg.Wait()
			}
		}
	}

	//Log.Info("system flow # Logic #")
	for _, period := range p.stageList {
		sq = p.stages[period]
		for _, sl := range sq {
			sl.reset()
			for ss := sl.next(); len(ss) > 0; ss = sl.next() {
				if systemCount := len(ss); systemCount != 0 {
					for i := 0; i < systemCount; i++ {
						sys := ss[i]
						imp := false
						var fn func(event Event)
						state := ss[i].getState()
						if state == SystemStateInit {
							switch period {
							case StageStart:
								system, ok := sys.(StartReceiver)
								fn = system.Start
								imp = ok
							}
							sys.setState(SystemStateUpdate)
						} else if state == SystemStateUpdate {
							switch period {
							case StagePreUpdate:
								system, ok := sys.(PreUpdateReceiver)
								fn = system.PreUpdate
								imp = ok
							case StageUpdate:
								system, ok := sys.(UpdateReceiver)
								fn = system.Update
								imp = ok
							case StagePostUpdate:
								system, ok := sys.(PostUpdateReceiver)
								fn = system.PostUpdate
								imp = ok
							}
						} else if state == SystemStateDestroy {
							switch period {
							case StageDestroy:
								system, ok := sys.(DestroyReceiver)
								fn = system.Destroy
								imp = ok
								sys.setState(SystemStateDestroyed)
								removeList[sys.ID()] = sys
							}
						}

						if !imp {
							continue
						}

						p.wg.Add(1)
						wg := p.wg
						Runtime.addJob(func() {
							defer func() {
								sys.setExecuting(false)
								wg.Done()
							}()

							sys.setExecuting(true)
							fn(event)
						})
					}
				}
				p.wg.Wait()
			}
		}
	}

	//do something clean
	for _, system := range removeList {
		p.unregister(system)
	}
}

//register method only in world init or func init(){}
func (p *systemFlow) register(system ISystem) {
	p.lock.Lock()
	defer p.lock.Unlock()

	//init function call
	system.baseInit(p.world, system)

	order := system.Order()
	if order > OrderAppend {
		Log.Errorf("system order must less then %d, reset order to %d", OrderAppend+1, OrderAppend)
		order = OrderAppend
	}

	for _, period := range p.stageList {
		imp := false
		switch period {
		case StageStart:
			_, imp = system.(StartReceiver)
		case StagePreUpdate:
			_, imp = system.(PreUpdateReceiver)
		case StageUpdate:
			_, imp = system.(UpdateReceiver)
		case StagePostUpdate:
			_, imp = system.(PostUpdateReceiver)
		case StageDestroy:
			_, imp = system.(DestroyReceiver)
		}

		if !imp {
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

	p.systems.Store(system.Type(), system)
}

func (p *systemFlow) unregister(system ISystem) {
	p.lock.Lock()
	defer p.lock.Unlock()

	order := system.Order()
	if order > OrderAppend {
		Log.Errorf("system order must less then %d, reset order to %d", OrderAppend+1, OrderAppend)
		order = OrderAppend
	}

	for _, period := range p.stageList {
		imp := false
		switch period {
		case StageStart:
			_, imp = system.(StartReceiver)
		case StagePreUpdate:
			_, imp = system.(PreUpdateReceiver)
		case StageUpdate:
			_, imp = system.(UpdateReceiver)
		case StagePostUpdate:
			_, imp = system.(PostUpdateReceiver)
		case StageDestroy:
			_, imp = system.(DestroyReceiver)
		}

		if !imp {
			continue
		}

		sl := p.stages[period]
		for _, group := range sl {
			if group.has(system) {
				group.remove(system)
			}
		}
	}

	p.systems.Delete(system.Type())
}

func (p *systemFlow) stop() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.reset()
}
