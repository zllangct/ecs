package ecs

import (
	"errors"
	"sync"
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

	Name FixedString[Fixed16]
}

type __world_Test_S_1 struct {
	System[__world_Test_S_1]
}

func (w *__world_Test_S_1) Init(si SystemInitConstraint) error {
	w.SetRequirements(si, &__world_Test_C_1{}, &__world_Test_C_2{}, &__world_Test_C_3{})
	BindUtility[__world_Test_U_Input](si)
	return nil
}

func (w *__world_Test_S_1) Update(event Event) {
	Log.Infof("Update: %d", event.Frame)
	iter := GetComponentAll[__world_Test_C_1](w)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		c2 := GetRelated[__world_Test_C_2](w, c.owner)
		if c2 == nil {
			Log.Infof("Component(Owner:%d) has no related component __world_Test_C_2", c.Owner())
			continue
		}
		for i := 0; i < testOptimizerDummyMaxFor; i++ {
			c.Field1 += i
		}

		for i := 0; i < testOptimizerDummyMaxFor; i++ {
			c2.Field2 += i
		}

		c3 := GetRelated[__world_Test_C_3](w, c.owner)

		Log.Infof("Component(Owner:%d) Changed: __world_Test_C_1: %d, __world_Test_C_2: %d, __world_Test_C_3: %s", c.Owner(), c.Field1, c2.Field2, c3.Name.String())
	}
}

type __world_Test_U_Input struct {
	Utility[__world_Test_U_Input]
}

func (u *__world_Test_U_Input) ChangeName(entity Entity, name string) {
	c := GetComponent[__world_Test_C_3](u.GetSystem(), entity)
	if c == nil {
		return
	}
	old := c.Name.String()
	c.Name.Set(name)
	Log.Infof("Name changed, old: %s, new:%s", old, name)
}

func Test_ecsWorld_World(t *testing.T) {
	config := NewDefaultWorldConfig()

	world := NewSyncWorld(config)

	RegisterSystem[__world_Test_S_1](world)

	world.Startup()

	entities := make([]Entity, __worldTest_Entity_Count)

	for i := 0; i < __worldTest_Entity_Count; i++ {
		e1 := world.NewEntity()
		world.Add(e1, &__world_Test_C_1{}, &__world_Test_C_2{}, &__world_Test_C_3{})
		entities[i] = e1
	}

	world.Update()

	u, ok := GetUtility[__world_Test_U_Input](world)
	if ok {
		u.ChangeName(entities[0], "name0")
	}

	world.Update()
	world.Stop()
	for false {
		world.Update()
		time.Sleep(time.Second)
	}
}

func Test_ecsWorld_World_launcher(t *testing.T) {
	config := NewDefaultWorldConfig()
	config.FrameInterval = time.Second

	world := NewAsyncWorld(config)

	RegisterSystem[__world_Test_S_1](world)

	world.Startup()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	entities := make([]Entity, __worldTest_Entity_Count)
	world.Sync(func(gaw SyncWrapper) error {
		for i := 0; i < __worldTest_Entity_Count; i++ {
			e1 := gaw.NewEntity()
			gaw.Add(e1, &__world_Test_C_1{}, &__world_Test_C_2{}, &__world_Test_C_3{})
			entities[i] = e1
		}
		wg.Done()
		return nil
	})

	wg.Wait()

	wg.Add(1)
	world.Sync(func(gaw SyncWrapper) error {
		u, ok := GetUtility[__world_Test_U_Input](gaw)
		if !ok {
			return errors.New("Utility not found")
		}
		u.ChangeName(entities[0], "name1")
		gaw.DestroyEntity(entities[1])
		gaw.Remove(entities[2], &__world_Test_C_2{})
		wg.Done()
		return nil
	})

	wg.Wait()

	wg.Add(1)
	world.Sync(func(gaw SyncWrapper) error {
		u, ok := GetUtility[__world_Test_U_Input](gaw)
		if !ok {
			return errors.New("Utility not found")
		}
		u.ChangeName(entities[0], "name2")
		wg.Done()
		return nil
	})

	wg.Wait()
	world.Stop()
}
