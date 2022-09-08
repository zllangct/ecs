package ecs

import (
	"testing"
	"time"
)

const (
	__worldTest_Entity_Count int = 3
)

type __world_Test_C_1 struct {
	Component[__world_Test_C_1]

	Field1 int
	Field2 int
}

type __world_Test_C_2 struct {
	Component[__world_Test_C_2]

	Field1 int
	Field2 int
}

type __world_Test_C_3 struct {
	Component[__world_Test_C_3]

	Name string
}

type __world_Test_S_1 struct {
	System[__world_Test_S_1]
}

func (w *__world_Test_S_1) Init(initializer SystemInitializer) {
	w.SetRequirements(initializer, &__world_Test_C_1{}, &__world_Test_C_2{}, &__world_Test_C_3{})
	w.BindUtility(initializer, &__world_Test_U_Input{})
}

func (w *__world_Test_S_1) Update(event Event) {
	Log.Infof("Update: %d", event.Frame)
	iter := GetInterestedComponents[__world_Test_C_1](w)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		c2 := GetRelatedComponent[__world_Test_C_2](w, c.owner)
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

type __world_Test_U_Input struct {
	Utility[__world_Test_S_1]
}

func (u *__world_Test_U_Input) ChangeName(entity Entity, name string) {
	sys := u.GetSystem()
	c := GetInterestedComponent[__world_Test_C_3](sys, entity)
	if c == nil {
		return
	}
	old := c.Name
	c.Name = name
	Log.Infof("Name changed, old: %s, new:%s", old, name)
}

func Test_ecsWorld_World(t *testing.T) {
	EnableMainThreadDebug()

	config := NewDefaultWorldConfig()

	world := NewSyncWorld(config)

	RegisterSystem[__world_Test_S_1](world)

	world.Startup()

	entities := make([]Entity, __worldTest_Entity_Count)

	for i := 0; i < __worldTest_Entity_Count; i++ {
		e1 := world.NewEntity()
		e1.Add(&__world_Test_C_1{}, &__world_Test_C_2{}, &__world_Test_C_3{})
		entities[i] = e1.Entity()
	}

	world.Update()

	getter := world.GetUtilityGetter()
	u, ok := GetUtility[__world_Test_S_1, __world_Test_U_Input](getter)
	if ok {
		u.ChangeName(entities[0], "name0")
	}

	for {
		world.Update()
		time.Sleep(time.Second)
	}
}

type __world_Test_Gate struct {
	Gate[__world_Test_Gate]
}

func (g *__world_Test_Gate) Init(initializer GateInitializer) {
	initializer.EventRegister("CustomEvent1", g.CustomEvent1)
}

func (g *__world_Test_Gate) CustomEvent1(getter UtilityGetter, args []interface{}) {
	Log.Infof("CustomEvent1: %+v", args)
	entity := args[0].(Entity)
	name := args[1].(string)
	u, ok := GetUtility[__world_Test_S_1, __world_Test_U_Input](getter)
	if !ok {
		return
	}
	u.ChangeName(entity, name)
}

func (g *__world_Test_Gate) input1(entity Entity, name string) {
	g.Sync(func(getter UtilityGetter) {
		u, ok := GetUtility[__world_Test_S_1, __world_Test_U_Input](getter)
		if !ok {
			return
		}
		u.ChangeName(entity, name)
	})
}

func Test_ecsWorld_World_launcher(t *testing.T) {
	EnableMainThreadDebug()

	config := NewDefaultWorldConfig()
	config.FrameInterval = time.Second

	world := NewAsyncWorld(config)

	RegisterSystem[__world_Test_S_1](world)
	BindGate[__world_Test_Gate](world)

	world.Startup()

	entities := make([]Entity, __worldTest_Entity_Count)

	for i := 0; i < __worldTest_Entity_Count; i++ {
		e1 := world.NewEntity()
		e1.Add(&__world_Test_C_1{}, &__world_Test_C_2{}, &__world_Test_C_3{})
		entities[i] = e1.Entity()
	}

	gate := GetGate[__world_Test_Gate](world)
	if gate == nil {
		t.Fatal("gate is nil")
	}

	go world.Startup()

	time.Sleep(time.Second * 2)

	// input by gate event
	gate.Emit("CustomEvent1", entities[0], "name0")

	// input by gate method
	gate.input1(entities[0], "name1")
	time.Sleep(time.Second * 2)
	gate.input1(entities[0], "name2")

	for {
		world.Update()
		time.Sleep(time.Second)
	}
}
