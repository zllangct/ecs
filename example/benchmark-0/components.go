package main

import (
	"github.com/zllangct/ecs"
)

type Position struct {
	ecs.Component[Position]
	X int
	Y int
	Z int
}

type Movement struct {
	ecs.Component[Movement]
	V   int
	Dir [3]int
}

type HealthPoint struct {
	ecs.Component[HealthPoint]
	HP int
}

type Force struct {
	ecs.Component[Force]
	AttackRange        int
	PhysicalBaseAttack int
	Strength           int
	CriticalChange     int
	CriticalMultiple   int
}

type Action struct {
	ecs.DisposableComponent[Action]
	ActionType int
}

type Test1 struct {
	ecs.Component[Test1]
	Test1 int
}

type Test2 struct {
	ecs.Component[Test2]
	Test2 int
}

type Test3 struct {
	ecs.Component[Test3]
	Test3 int
}

type Test4 struct {
	ecs.Component[Test4]
	Test4 int
}

type Test5 struct {
	ecs.Component[Test5]
	Test5 int
}
