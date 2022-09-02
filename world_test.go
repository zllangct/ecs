package ecs

import (
	"testing"
	"time"
)

const (
	__worldTest_Entity_Count int = 3
)

type __worldTest_Com1 struct {
	Component[__worldTest_Com1]

	Field1 int
	Field2 int
}

type __worldTest_Com2 struct {
	Component[__worldTest_Com2]

	Field1 int
	Field2 int
}

type __worldTest_Com3 struct {
	Component[__worldTest_Com3]

	Name string
}

type __worldTest_System1 struct {
	System[__worldTest_System1]
}

func (w *__worldTest_System1) Init() {
	w.SetRequirements(&__worldTest_Com1{}, &__worldTest_Com2{}, &__worldTest_Com3{})
	w.SetUtility(&__worldTest_U_Input{})
}

func (w *__worldTest_System1) Update(event Event) {
	Log.Infof("Update: %d", event.Frame)
	iter := GetInterestedComponents[__worldTest_Com1](w)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		c2 := GetRelatedComponent[__worldTest_Com2](w, c.owner)
		if c2 == nil {
			continue
		}
		for i := 0; i < testOptimizerDummyMaxFor; i++ {
			c.Field1 += i
		}

		for i := 0; i < testOptimizerDummyMaxFor; i++ {
			c2.Field2 += i
		}

		Log.Infof("Component Changed: %d, %d", c.Field1, c2.Field2)
	}
}

type __worldTest_U_Input struct {
	Utility[__worldTest_System1]
}

func (u *__worldTest_U_Input) ChangeName(entity Entity, name string) {
	sys := u.GetSystem()
	c := GetInterestedComponent[__worldTest_Com3](sys, entity)
	if c == nil {
		return
	}
	old := c.Name
	c.Name = name
	Log.Infof("Name changed, old: %s, new:%s", old, name)
}

func Test_ecsWorld_World(t *testing.T) {
	config := NewDefaultWorldConfig()
	world := NewWorld(config)

	RegisterSystem[__worldTest_System1](world)

	entities := make([]Entity, __worldTest_Entity_Count)

	for i := 0; i < __worldTest_Entity_Count; i++ {
		e1 := world.NewEntity()
		e1.Add(&__worldTest_Com1{}, &__worldTest_Com2{}, &__worldTest_Com3{})
		entities[i] = e1.Entity()
	}

	world.Update()

	for {
		world.Update()
		time.Sleep(time.Second)
	}
}

type __worldTest_Gate struct {
	Gate[__worldTest_Gate]
}

func (g *__worldTest_Gate) Init(initializer GateInitializer) {
	initializer.EventRegister("CustomEvent1", g.CustomEvent1)
}

func (g *__worldTest_Gate) CustomEvent1(api *GateApi, args []interface{}) {
	Log.Infof("CustomEvent1: %+v", args)
	entity := args[0].(Entity)
	name := args[1].(string)
	u, ok := GetUtilityByGate[__worldTest_System1, __worldTest_U_Input](api)
	if !ok {
		return
	}
	u.ChangeName(entity, name)
}

func (g *__worldTest_Gate) input1(entity Entity, name string) {
	g.Sync(func(api *GateApi) {
		u, ok := GetUtilityByGate[__worldTest_System1, __worldTest_U_Input](api)
		if !ok {
			return
		}
		u.ChangeName(entity, name)
	})
}

func Test_ecsWorld_World_launcher(t *testing.T) {
	config := NewDefaultWorldConfig()
	config.FrameInterval = time.Second

	world := NewWorld(config)

	RegisterSystem[__worldTest_System1](world)

	entities := make([]Entity, __worldTest_Entity_Count)

	for i := 0; i < __worldTest_Entity_Count; i++ {
		e1 := world.NewEntity()
		e1.Add(&__worldTest_Com1{}, &__worldTest_Com2{}, &__worldTest_Com3{})
		entities[i] = e1.Entity()
	}

	launcher := NewAsyncWorldLauncher(world)
	iGate := launcher.SetGate(&__worldTest_Gate{})
	gate := iGate.(*__worldTest_Gate)

	go launcher.Run()

	time.Sleep(time.Second * 2)

	// input by gate event
	gate.Emit("CustomEvent1", entities[0], "name0")

	// input by gate method
	gate.input1(entities[0], "name1")

	time.Sleep(time.Second * 2)

	gate.input1(entities[0], "name2")

	for {
		time.Sleep(time.Second)
	}
}
