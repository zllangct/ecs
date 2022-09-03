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

type __optimizer_Bench_C_1 struct {
	Component[__optimizer_Bench_C_1]
	Test1 int
}

type __optimizer_Bench_C_2 struct {
	Component[__optimizer_Bench_C_2]
	Test2 int
}

type __optimizer_Bench_S_1 struct {
	System[__optimizer_Bench_S_1]
}

func (t *__optimizer_Bench_S_1) Init(initializer *SystemInitializer) {
	t.SetRequirements(initializer, &__optimizer_Bench_C_1{}, &__optimizer_Bench_C_2{})
}

func (t *__optimizer_Bench_S_1) Update(event Event) {
	iter := GetInterestedComponents[__optimizer_Bench_C_1](t)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		c2 := GetRelatedComponent[__optimizer_Bench_C_2](t, c.owner)
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

type __optimizer_Bench_GameECS struct {
	world    IWorld
	entities []Entity
}

func (g *__optimizer_Bench_GameECS) init() {
	println("init")
	config := NewDefaultWorldConfig()
	g.world = NewWorld(config)

	RegisterSystem[__optimizer_Bench_S_1](g.world)

	for i := 0; i < testOptimizerEntityMax; i++ {
		c := &__optimizer_Bench_C_1{}
		e := g.world.NewEntity()
		e.Add(c)
		g.entities = append(g.entities, e.Entity())
	}
	rand.Seed(0)
	rand.Shuffle(len(g.entities), func(i, j int) { g.entities[i], g.entities[j] = g.entities[j], g.entities[i] })

	for i := 0; i < testOptimizerEntityMax; i++ {
		c := &__optimizer_Bench_C_2{}
		g.world.addComponent(g.entities[i], c)
	}
}

func BenchmarkNoOptimizer(b *testing.B) {
	//go func() {
	//	http.ListenAndServe(":6060", nil)
	//}()
	println("start")
	game := &__optimizer_Bench_GameECS{}
	game.init()
	game.world.update()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		game.world.update()
	}
}

func BenchmarkWithOptimizer(b *testing.B) {
	//go func() {
	//	http.ListenAndServe(":6060", nil)
	//}()

	game := &__optimizer_Bench_GameECS{}
	game.init()
	game.world.update()

	game.world.Optimize(time.Second*10, true)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		game.world.update()
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
