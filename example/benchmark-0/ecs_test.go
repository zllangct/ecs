package main

import (
	"github.com/zllangct/ecs"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
	_ "unsafe"
)

const (
	PlayerCount = 1000
)

type Player struct {
	rw sync.RWMutex

	ID int64

	X                  int
	Y                  int
	Z                  int
	V                  int
	Dir                [3]int
	HP                 int
	AttackRange        int
	PhysicalBaseAttack int
	Strength           int
	CriticalChange     int
	CriticalMultiple   int
	ActionType         int
	Test1              int
	Test2              int
	Test3              int
	Test4              int
	Test5              int
}

type GameNormal struct {
	players map[int64]*Player
}

func (g *GameNormal) init() {
	DataGenerateNormal(g)
}

func (g *GameNormal) doFrame(parallel bool, frame uint64, delta time.Duration) {
	if parallel {
		wg := &sync.WaitGroup{}
		wg.Add(7)
		go func() {
			g.DoMoveParallel(delta)
			wg.Done()
		}()
		go func() {
			g.DoDamageParallel()
			wg.Done()
		}()
		go func() {
			g.OtherLoadParallel1()
			wg.Done()
		}()
		go func() {
			g.OtherLoadParallel2()
			wg.Done()
		}()
		go func() {
			g.OtherLoadParallel3()
			wg.Done()
		}()
		go func() {
			g.OtherLoadParallel4()
			wg.Done()
		}()
		go func() {
			g.OtherLoadParallel5()
			wg.Done()
		}()
		wg.Wait()
	} else {
		// 移动
		g.DoMove(delta)
		// 攻击处理
		g.DoDamage()
		// 其他负载
		g.OtherLoad1()
		g.OtherLoad2()
		g.OtherLoad3()
		g.OtherLoad4()
		g.OtherLoad5()
	}
}
func (g *GameNormal) OtherLoad1() {
	for _, p := range g.players {
		for i := 0; i < 1000; i++ {
			p.Test2 += 1
		}
	}
}
func (g *GameNormal) OtherLoad2() {
	for _, p := range g.players {
		for i := 0; i < 1000; i++ {
			p.Test2 += 1
		}
	}
}
func (g *GameNormal) OtherLoad3() {
	for _, p := range g.players {
		for i := 0; i < 1000; i++ {
			p.Test2 += 1
		}
	}
}
func (g *GameNormal) OtherLoad4() {
	for _, p := range g.players {
		for i := 0; i < 1000; i++ {
			p.Test2 += 1
		}
	}
}
func (g *GameNormal) OtherLoad5() {
	for _, p := range g.players {
		for i := 0; i < 1000; i++ {
			p.Test2 += 1
		}
	}
}

func (g *GameNormal) OtherLoadParallel1() {
	for _, p := range g.players {
		p.rw.Lock()
		for i := 0; i < 1000; i++ {
			p.Test2 += 1
		}
		p.rw.Unlock()
	}
}
func (g *GameNormal) OtherLoadParallel2() {
	for _, p := range g.players {
		p.rw.Lock()
		for i := 0; i < 1000; i++ {
			p.Test2 += 1
		}
		p.rw.Unlock()
	}
}
func (g *GameNormal) OtherLoadParallel3() {
	for _, p := range g.players {
		p.rw.Lock()
		for i := 0; i < 1000; i++ {
			p.Test2 += 1
		}
		p.rw.Unlock()
	}
}
func (g *GameNormal) OtherLoadParallel4() {
	for _, p := range g.players {
		p.rw.Lock()
		for i := 0; i < 1000; i++ {
			p.Test2 += 1
		}
		p.rw.Unlock()
	}
}
func (g *GameNormal) OtherLoadParallel5() {
	for _, p := range g.players {
		p.rw.Lock()
		for i := 0; i < 1000; i++ {
			p.Test2 += 1
		}
		p.rw.Unlock()
	}
}

func (g *GameNormal) DoMove(delta time.Duration) {
	for _, p := range g.players {
		p.X = p.X + int(float64(p.Dir[0]*p.V)*delta.Seconds())
		p.Y = p.Y + int(float64(p.Dir[1]*p.V)*delta.Seconds())
		p.Z = p.Z + int(float64(p.Dir[2]*p.V)*delta.Seconds())
	}
}

func (g *GameNormal) DoMoveParallel(delta time.Duration) {
	for _, p := range g.players {
		p.rw.Lock()
		p.X = p.X + int(float64(p.Dir[0]*p.V)*delta.Seconds())
		p.Y = p.Y + int(float64(p.Dir[1]*p.V)*delta.Seconds())
		p.Z = p.Z + int(float64(p.Dir[2]*p.V)*delta.Seconds())
		p.rw.Unlock()
	}
}

func (g *GameNormal) DoDamage() {
	for _, caster := range g.players {
		for _, target := range g.players {
			if caster.ID == target.ID {
				continue
			}

			//计算距离
			distance := int(math.Sqrt(math.Pow(2, float64(caster.X-target.X)) + math.Pow(2, float64(caster.Y-target.Y))))
			if distance > caster.AttackRange {
				continue
			}

			//伤害公式：伤害=（基础攻击+力量）+ 暴击伤害， 暴击伤害=基础攻击 * 2
			damage := caster.PhysicalBaseAttack * caster.Strength
			critical := 0
			if rand.Intn(100) < caster.CriticalChange {
				critical = caster.PhysicalBaseAttack * caster.CriticalMultiple
			}
			damage = damage + critical
			target.HP -= damage
			if target.HP < 0 {
				target.HP = 0
			}
		}
	}
}
func (g *GameNormal) DoDamageParallel() {
	for _, caster := range g.players {
		caster.rw.RLock()
		for _, target := range g.players {
			if caster.ID == target.ID {
				continue
			}

			target.rw.Lock()
			//计算距离
			distance := int(math.Sqrt(math.Pow(2, float64(caster.X-target.X)) + math.Pow(2, float64(caster.Y-target.Y))))
			if distance > caster.AttackRange {
				target.rw.Unlock()
				continue
			}

			//伤害公式：伤害=（基础攻击+力量）+ 暴击伤害， 暴击伤害=基础攻击 * 2
			damage := caster.PhysicalBaseAttack * caster.Strength
			critical := 0
			if rand.Intn(100) < caster.CriticalChange {
				critical = caster.PhysicalBaseAttack * caster.CriticalMultiple
			}
			damage = damage + critical
			target.HP -= damage
			if target.HP < 0 {
				target.HP = 0
			}
			target.rw.Unlock()
		}
		caster.rw.RUnlock()
	}
}

type GameECS struct {
	world    ecs.IWorld
	entities []ecs.Entity
}

//go:linkname doFrame github.com/zllangct/ecs.doFrameForBenchmark
func doFrame(w ecs.IWorld, frame uint64, lastDelta time.Duration)

func (g *GameECS) init() {
	ecs.RuntimeConfigure(ecs.NewDefaultRuntimeConfig())
	ecs.Run()

	g.world = ecs.CreateWorld(ecs.NewDefaultWorldConfig())

	ecs.RegisterSystem[MoveSystem](g.world)
	ecs.RegisterSystem[DamageSystem](g.world)

	DataGenerateECS(g)
}

func (g *GameECS) attack() {
	act := &Action{
		ActionType: 1,
	}
	for _, entity := range g.entities {
		info := ecs.GetEntityInfo(g.world, entity)
		info.Add(act)
	}
}

func DataGenerateNormal(normal *GameNormal) {
	players := make(map[int64]*Player)
	for i := 0; i < PlayerCount; i++ {
		p := &Player{
			ID:                 int64(i),
			X:                  0,
			Y:                  0,
			Z:                  0,
			V:                  100,
			Dir:                [3]int{1, 0, 0},
			HP:                 100,
			AttackRange:        10000,
			PhysicalBaseAttack: 10,
			Strength:           0,
			CriticalChange:     0,
			CriticalMultiple:   0,
			ActionType:         1,
		}
		players[p.ID] = p
	}
	normal.players = players
}

func DataGenerateECS(game *GameECS) {
	for i := 0; i < PlayerCount; i++ {
		p := &Position{
			X: 0,
			Y: 0,
			Z: 0,
		}
		m := &Movement{
			V:   100,
			Dir: [3]int{1, 0, 0},
		}
		h := &HealthPoint{
			HP: 100,
		}
		f := &Force{
			AttackRange:        10000,
			PhysicalBaseAttack: 10,
		}
		a := &Action{
			ActionType: 1,
		}

		t1 := &Test1{}
		t2 := &Test2{}
		t3 := &Test3{}
		t4 := &Test4{}
		t5 := &Test5{}

		e := game.world.NewEntity()
		e.Add(p, m, h, f, a, t1, t2, t3, t4, t5)
		game.entities = append(game.entities, e.Entity())
	}
}

func BenchmarkNormal(b *testing.B) {
	game := &GameNormal{
		players: make(map[int64]*Player),
	}
	game.init()
	b.ResetTimer()

	var delta time.Duration
	var ts time.Time
	var frameInterval time.Duration = time.Millisecond * 33
	for i := 0; i < b.N; i++ {
		ts = time.Now()
		game.doFrame(false, uint64(i), frameInterval)
		delta = time.Since(ts)
		if frameInterval-delta > 0 {
			delta = frameInterval
		}
	}
}

func BenchmarkNormalParallel(b *testing.B) {
	game := &GameNormal{
		players: make(map[int64]*Player),
	}
	game.init()
	b.ResetTimer()

	var delta time.Duration
	var ts time.Time
	var frameInterval time.Duration = time.Millisecond * 33
	for i := 0; i < b.N; i++ {
		ts = time.Now()
		game.doFrame(true, uint64(i), frameInterval)
		delta = time.Since(ts)
		if frameInterval-delta > 0 {
			delta = frameInterval
		}
	}
}

func BenchmarkEcs(b *testing.B) {
	game := &GameECS{}
	game.init()

	b.ResetTimer()

	var delta time.Duration
	var ts time.Time
	var frameInterval time.Duration = time.Millisecond * 33
	for i := 0; i < b.N; i++ {
		ts = time.Now()
		doFrame(game.world, uint64(i), frameInterval)
		//game.attack()
		delta = time.Since(ts)
		if frameInterval-delta > 0 {
			delta = frameInterval
		}
	}
}
