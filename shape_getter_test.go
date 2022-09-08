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

	getter1 *Shape[__ShapeGetter_Test_Shape_1]
}

func (t *__ShapeGetter_Test_S_1) Init(initializer SystemInitializer) {
	t.SetRequirements(initializer, &__ShapeGetter_Test_C_1{}, &__ShapeGetter_Test_C_2{})

	t.getter1 = NewShape[__ShapeGetter_Test_Shape_1](initializer)
	if t.getter1 == nil {
		initializer.SetBroken("invalid getter")
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
	world := NewSyncWorld(NewDefaultWorldConfig())
	RegisterSystem[__ShapeGetter_Test_S_1](world)

	world.Startup()

	for i := 0; i < 3; i++ {
		e := world.NewEntity()
		e.Add(&__ShapeGetter_Test_C_1{Field1: i}, &__ShapeGetter_Test_C_2{Field1: i * 10})
	}

	world.Update()
	time.Sleep(time.Second)
	world.Update()
}
