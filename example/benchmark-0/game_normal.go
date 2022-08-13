package main

import (
	"math/rand"
	"sync"
	"time"
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

func (g *GameNormal) doFrame(parallel bool, frame uint64, delta time.Duration) {
	if parallel {
		wg := &sync.WaitGroup{}
		wg.Add(7)
		//go func() {
		//	g.DoMoveParallel(delta)
		//	wg.Done()
		//}()
		//go func() {
		//	g.DoDamageParallel()
		//	wg.Done()
		//}()
		go func() {
			g.SimuLoadParallel1()
			wg.Done()
		}()
		go func() {
			g.SimuLoadParallel2()
			wg.Done()
		}()
		go func() {
			g.SimuLoadParallel3()
			wg.Done()
		}()
		go func() {
			g.SimuLoadParallel4()
			wg.Done()
		}()
		go func() {
			g.SimuLoadParallel5()
			wg.Done()
		}()
		wg.Wait()
	} else {
		// 移动
		//g.DoMove(delta)
		// 攻击处理
		//g.DoDamage()
		// 模拟其他负载
		g.SimuLoad1()
		g.SimuLoad2()
		g.SimuLoad3()
		g.SimuLoad4()
		g.SimuLoad5()
	}
}
func (g *GameNormal) SimuLoad1() {
	for _, p := range g.players {
		for i := 0; i < DummyMaxFor; i++ {
			p.Test1 += 1
		}
	}
}
func (g *GameNormal) SimuLoad2() {
	for _, p := range g.players {
		for i := 0; i < DummyMaxFor; i++ {
			p.Test2 += 1
		}
	}
}
func (g *GameNormal) SimuLoad3() {
	for _, p := range g.players {
		for i := 0; i < DummyMaxFor; i++ {
			p.Test3 += 1
		}
	}
}
func (g *GameNormal) SimuLoad4() {
	for _, p := range g.players {
		for i := 0; i < DummyMaxFor; i++ {
			p.Test4 += 1
		}
	}
}
func (g *GameNormal) SimuLoad5() {
	for _, p := range g.players {
		for i := 0; i < DummyMaxFor; i++ {
			p.Test5 += 1
		}
	}
}

func (g *GameNormal) SimuLoadParallel1() {
	for _, p := range g.players {
		p.rw.Lock()
		for i := 0; i < DummyMaxFor; i++ {
			p.Test1 += 1
		}
		p.rw.Unlock()
	}
}
func (g *GameNormal) SimuLoadParallel2() {
	for _, p := range g.players {
		p.rw.Lock()
		for i := 0; i < DummyMaxFor; i++ {
			p.Test2 += 1
		}
		p.rw.Unlock()
	}
}
func (g *GameNormal) SimuLoadParallel3() {
	for _, p := range g.players {
		p.rw.Lock()
		for i := 0; i < DummyMaxFor; i++ {
			p.Test3 += 1
		}
		p.rw.Unlock()
	}
}
func (g *GameNormal) SimuLoadParallel4() {
	for _, p := range g.players {
		p.rw.Lock()
		for i := 0; i < DummyMaxFor; i++ {
			p.Test4 += 1
		}
		p.rw.Unlock()
	}
}
func (g *GameNormal) SimuLoadParallel5() {
	for _, p := range g.players {
		p.rw.Lock()
		for i := 0; i < DummyMaxFor; i++ {
			p.Test5 += 1
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
			distance := (caster.X-target.X)*(caster.X-target.X) + (caster.Y-target.Y)*(caster.Y-target.Y)
			if distance > caster.AttackRange*caster.AttackRange {
				continue
			}

			//伤害公式：伤害=（基础攻击+力量）+ 暴击伤害， 暴击伤害=基础攻击 * 2
			damage := caster.PhysicalBaseAttack + caster.Strength
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
			distance := (caster.X-target.X)*(caster.X-target.X) + (caster.Y-target.Y)*(caster.Y-target.Y)
			if distance > caster.AttackRange*caster.AttackRange {
				continue
			}

			//伤害公式：伤害=（基础攻击+力量）+ 暴击伤害， 暴击伤害=基础攻击 * 2
			damage := caster.PhysicalBaseAttack + caster.Strength
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
