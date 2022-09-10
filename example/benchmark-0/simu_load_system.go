package main

import "github.com/zllangct/ecs"

type Test1System struct {
	ecs.System[Test1System]
}

func (t *Test1System) Init(si ecs.SystemInitConstraint) {
	t.SetRequirements(si, &Test1{})
}

func (t *Test1System) Update(event ecs.Event) {
	iter := ecs.GetComponentAll[Test1](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		for i := 0; i < DummyMaxFor; i++ {
			c.Test1 += i
		}
	}
}

type Test2System struct {
	ecs.System[Test2System]
}

func (t *Test2System) Init(si ecs.SystemInitConstraint) {
	t.SetRequirements(si, &Test2{})
}

func (t *Test2System) Update(event ecs.Event) {
	iter := ecs.GetComponentAll[Test2](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		for i := 0; i < DummyMaxFor; i++ {
			c.Test2 += i
		}
	}
}

type Test3System struct {
	ecs.System[Test3System]
}

func (t *Test3System) Init(si ecs.SystemInitConstraint) {
	t.SetRequirements(si, &Test3{})
}

func (t *Test3System) Update(event ecs.Event) {
	iter := ecs.GetComponentAll[Test3](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		for i := 0; i < DummyMaxFor; i++ {
			c.Test3 += i
		}
	}
}

type Test4System struct {
	ecs.System[Test4System]
}

func (t *Test4System) Init(si ecs.SystemInitConstraint) {
	t.SetRequirements(si, &Test4{})
}

func (t *Test4System) Update(event ecs.Event) {
	iter := ecs.GetComponentAll[Test4](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		for i := 0; i < DummyMaxFor; i++ {
			c.Test4 += i
		}
	}
}

type Test5System struct {
	ecs.System[Test5System]
}

func (t *Test5System) Init(si ecs.SystemInitConstraint) {
	t.SetRequirements(si, &Test5{})
}

func (t *Test5System) Update(event ecs.Event) {
	iter := ecs.GetComponentAll[Test5](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		for i := 0; i < DummyMaxFor; i++ {
			c.Test5 += i
		}
	}
}
