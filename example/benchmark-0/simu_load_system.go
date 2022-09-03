package main

import "github.com/zllangct/ecs"

type Test1System struct {
	ecs.System[Test1System]
}

func (t *Test1System) Init(initializer *ecs.SystemInitializer) {
	t.SetRequirements(initializer, &Test1{})
}

func (t *Test1System) Update(event ecs.Event) {
	iter := ecs.GetInterestedComponents[Test1](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		for i := 0; i < DummyMaxFor; i++ {
			c.Test1 += i
		}
	}
}

type Test2System struct {
	ecs.System[Test2System]
}

func (t *Test2System) Init(initializer *ecs.SystemInitializer) {
	t.SetRequirements(initializer, &Test2{})
}

func (t *Test2System) Update(event ecs.Event) {
	iter := ecs.GetInterestedComponents[Test2](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		for i := 0; i < DummyMaxFor; i++ {
			c.Test2 += i
		}
	}
}

type Test3System struct {
	ecs.System[Test3System]
}

func (t *Test3System) Init(initializer *ecs.SystemInitializer) {
	t.SetRequirements(initializer, &Test1{})
}

func (t *Test3System) Update(event ecs.Event) {
	iter := ecs.GetInterestedComponents[Test3](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		for i := 0; i < DummyMaxFor; i++ {
			c.Test3 += i
		}
	}
}

type Test4System struct {
	ecs.System[Test4System]
}

func (t *Test4System) Init(initializer *ecs.SystemInitializer) {
	t.SetRequirements(initializer, &Test4{})
}

func (t *Test4System) Update(event ecs.Event) {
	iter := ecs.GetInterestedComponents[Test4](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		for i := 0; i < DummyMaxFor; i++ {
			c.Test4 += i
		}
	}
}

type Test5System struct {
	ecs.System[Test5System]
}

func (t *Test5System) Init(initializer *ecs.SystemInitializer) {
	t.SetRequirements(initializer, &Test5{})
}

func (t *Test5System) Update(event ecs.Event) {
	iter := ecs.GetInterestedComponents[Test5](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		for i := 0; i < DummyMaxFor; i++ {
			c.Test5 += i
		}
	}
}
