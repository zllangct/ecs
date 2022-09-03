package ecs

import (
	"testing"
	"time"
)

type __ShapeGetter_Test_C_1 struct {
	Component[__ShapeGetter_Test_C_1]
	Field1 int
}

type __ShapeGetter_Test_C_2 struct {
	Component[__ShapeGetter_Test_C_2]
	Field1 int
}

type __ShapeGetter_Test_Shape_1 struct {
	c1 *__ShapeGetter_Test_C_1
	c2 *__ShapeGetter_Test_C_2
}

type __ShapeGetter_Test_S_1 struct {
	System[__ShapeGetter_Test_S_1]

	getter1 *ShapeGetter[__ShapeGetter_Test_Shape_1]
}

func (t *__ShapeGetter_Test_S_1) Init(initializer *SystemInitializer) {
	t.SetRequirements(initializer, &__ShapeGetter_Test_C_1{}, &__ShapeGetter_Test_C_2{})

	var err error
	t.getter1, err = NewShapeGetter[__ShapeGetter_Test_Shape_1](initializer)
	if err != nil {
		Log.Fatal(err)
	}
}

func (t *__ShapeGetter_Test_S_1) Update(event Event) {
	Log.Infof("__ShapeGetter_Test_S_1.Update, frame:%d", event.Frame)
	iter := t.getter1.Get()
	for s := iter.Begin(); !iter.End(); s = iter.Next() {
		Log.Infof("s.c1:%+v, s.c2:%+v", s.c1, s.c2)
	}
}

func TestNewShapeGetter(t *testing.T) {
	world := NewWorld(NewDefaultWorldConfig())
	RegisterSystem[__ShapeGetter_Test_S_1](world)

	launcher := world.GetSyncLauncher()

	for i := 0; i < 3; i++ {
		e := world.NewEntity()
		e.Add(&__ShapeGetter_Test_C_1{Field1: i}, &__ShapeGetter_Test_C_2{Field1: i * 10})
	}

	launcher.Update()
	time.Sleep(time.Second)
	launcher.Update()
}
