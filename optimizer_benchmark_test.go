package ecs

import (
	"math/rand"
	"testing"
	"time"
)

const (
	testOptimizerDummyMaxFor = 10
	testOptimizerEntityMax   = 1000000
)

type testOptimizerComponent1 struct {
	Component[testOptimizerComponent1]
	Test1 int
}

type testOptimizerComponent2 struct {
	Component[testOptimizerComponent2]
	Test2 int
}

type testOptimizerSystem struct {
	System[testOptimizerSystem]
}

func (t *testOptimizerSystem) Init() {
	t.SetRequirements(&testOptimizerComponent1{}, &testOptimizerComponent2{})
}

func (t *testOptimizerSystem) Update(event Event) {
	iter := GetInterestedComponents[testOptimizerComponent1](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		c2 := GetRelatedComponent[testOptimizerComponent2](t, c.owner)
		if c2 == nil {
			continue
		}
		for i := 0; i < testOptimizerDummyMaxFor; i++ {
			c.Test1 += i
		}

		for i := 0; i < testOptimizerDummyMaxFor; i++ {
			c2.Test2 += i
		}
	}
}

type gameECS struct {
	world    IWorld
	entities []Entity
}

func (g *gameECS) init() {
	println("init")
	config := NewDefaultWorldConfig()
	g.world = NewWorld(config)

	RegisterSystem[testOptimizerSystem](g.world)

	for i := 0; i < testOptimizerEntityMax; i++ {
		c := &testOptimizerComponent1{}
		e := g.world.NewEntity()
		e.Add(c)
		g.entities = append(g.entities, e.Entity())
	}
	rand.Seed(0)
	rand.Shuffle(len(g.entities), func(i, j int) { g.entities[i], g.entities[j] = g.entities[j], g.entities[i] })

	for i := 0; i < testOptimizerEntityMax; i++ {
		c := &testOptimizerComponent2{}
		g.world.addComponent(g.entities[i], c)
	}
}

func BenchmarkNoOptimizer(b *testing.B) {
	//go func() {
	//	http.ListenAndServe(":6060", nil)
	//}()
	println("start")
	game := &gameECS{}
	game.init()
	game.world.Update()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		game.world.Update()
	}
}

func BenchmarkWithOptimizer(b *testing.B) {
	//go func() {
	//	http.ListenAndServe(":6060", nil)
	//}()

	game := &gameECS{}
	game.init()
	game.world.Update()

	game.world.Optimize(time.Second*10, true)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		game.world.Update()
	}
}

func BenchmarkTest(b *testing.B) {
	arr := make([]int, 0, 100)
	for i, _ := range arr {
		arr[i] = i
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {

	}
}

func BenchmarkTest2(b *testing.B) {
	type test struct {
		Name string
		Age  int
	}
	t := test{Name: "test", Age: 1}
	m := map[test]int{t: 1}
	for i := 0; i < b.N; i++ {
		_ = m[t]
	}
}
