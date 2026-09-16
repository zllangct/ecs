package ecs

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/synctest"
)

// ============================================================================
// 并行模式负载验证：
// 1. 并行/线性双世界等价性（组表 + View2RO + 多系统真并发 + 随机 op）
// 2. synctest 气泡死锁检测
// 3. -race 运行（CI/手动）：go test -race -run TestParallel ./ecs/
// ============================================================================

// parallelSnapshot 每帧观测（各系统独立字段，避免并发写竞争）
type parallelSnapshot struct {
	c        map[EntityIndex]int32
	bSumMove int32 // move 系统经 View2RO 累计（组表 lockstep 路径）
	bSumRO   int32 // rosys 系统经 Query+GetBuddyReadOnly 累计（通用路径）
	moveN    int
}

// parallelMoveSys 组表系统：View2RO lockstep 读 B 写 A（唯一写 A/B 的系统）
type parallelMoveSys struct{ snap *parallelSnapshot }

func (s *parallelMoveSys) Update(ctx *SystemContext, event Event) error {
	s.snap.moveN = 0
	s.snap.bSumMove = 0
	// 组表路径：View2RO lockstep（双变量，无元组）
	for a, bv := range View2RO[groupCompA, groupCompBView](ctx) {
		s.snap.moveN++
		s.snap.bSumMove += bv.V() // 与 rosys 通用路径同帧同值，交叉验证
		a.V++
	}
	return nil
}

// parallelCSys 独立系统：写 C（与 A/B 无共享依赖，与 move/ro 真并发）
type parallelCSys struct{ snap *parallelSnapshot }

func (s *parallelCSys) Update(ctx *SystemContext, event Event) error {
	s.snap.c = map[EntityIndex]int32{}
	for idx, c := range ctx.GetComponents[gCompC, *gCompC]() {
		c.V += 2
		s.snap.c[idx] = c.V
	}
	return nil
}

// parallelROSys 只读系统：QueryGet 读 A，GetBuddyReadOnly 读 B 求和
// （与 move 共享 B 但均为只读 → 独立批次并发）
type parallelROSys struct{ snap *parallelSnapshot }

func (s *parallelROSys) Update(ctx *SystemContext, event Event) error {
	s.snap.bSumRO = 0
	q := ctx.NewQuery(WithComp[groupCompA, *groupCompA](), WithComp[groupCompB, *groupCompB]())
	for idx := range q.Iter() {
		bv, ok := ctx.GetBuddyReadOnly[groupCompBView](idx)
		if !ok {
			continue
		}
		s.snap.bSumRO += bv.V()
	}
	return nil
}

func registerParallelSystems(w *World, snap *parallelSnapshot) error {
	if err := w.RegisterLight(LightSystem((&parallelMoveSys{snap}).Update),
		WithName("move"),
		WithDeps(Dep[groupCompA](), Dep[groupCompB](ReadOnly)),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]())); err != nil {
		return err
	}
	if err := w.RegisterLight(LightSystem((&parallelCSys{snap}).Update),
		WithName("csys"),
		WithDeps(Dep[gCompC]())); err != nil {
		return err
	}
	if err := w.RegisterLight(LightSystem((&parallelROSys{snap}).Update),
		WithName("rosys"),
		WithDeps(Dep[groupCompA](ReadOnly), Dep[groupCompB](ReadOnly))); err != nil {
		return err
	}
	return nil
}

// setupParallelEntities 按 seed 生成确定性实体布局
func setupParallelEntities(w *World, rng *rand.Rand) []Entity {
	var entities []Entity
	for i := 0; i < 200; i++ {
		comps := []Component{newGroupCompA(rng.Int31n(100))}
		if rng.Intn(4) > 0 { // 75% 带 B（组表行）
			comps = append(comps, newGroupCompB(rng.Int31n(100)))
		}
		if rng.Intn(2) == 0 {
			comps = append(comps, &gCompC{V: rng.Int31n(100)})
		}
		entities = append(entities, w.NewEntity(WithComponents(comps...)))
	}
	return entities
}

// TestParallel_GroupEquivalence 并行/线性双世界，同随机 op 序列逐帧对比。
// 注意：moveN 与 ab 聚合值、C 值、bSum 均为确定性指标（op 单写者、写值确定性）。
func TestParallel_GroupEquivalence(t *testing.T) {
	for seed := int64(0); seed < 10; seed++ {
		snapL, snapP := &parallelSnapshot{}, &parallelSnapshot{}
		wl := NewWorld()
		wp := NewWorld(WithWorldASyncMode())
		if err := registerParallelSystems(wl, snapL); err != nil {
			t.Fatal(err)
		}
		if err := registerParallelSystems(wp, snapP); err != nil {
			t.Fatal(err)
		}

		rng := rand.New(rand.NewSource(seed))
		entitiesL := setupParallelEntities(wl, rng)
		rng = rand.New(rand.NewSource(seed))
		entitiesP := setupParallelEntities(wp, rng)

		for frame := 0; frame < 30; frame++ {
			// 同序列主队列 op：销毁/增删 B（组成员资格 churn）
			rngOps := rand.New(rand.NewSource(seed*1000 + int64(frame)))
			for i := 0; i < 10; i++ {
				if len(entitiesL) == 0 {
					break
				}
				ei := rngOps.Intn(len(entitiesL))
				switch rngOps.Intn(4) {
				case 0: // 销毁
					wl.DestroyEntity(entitiesL[ei])
					wp.DestroyEntity(entitiesP[ei])
					entitiesL = append(entitiesL[:ei], entitiesL[ei+1:]...)
					entitiesP = append(entitiesP[:ei], entitiesP[ei+1:]...)
				case 1: // 加 B（可能入行）
					v := rngOps.Int31n(100)
					if info, ok := wl.GetEntityInfo(entitiesL[ei]); ok {
						info.Add(newGroupCompB(v))
					}
					if info, ok := wp.GetEntityInfo(entitiesP[ei]); ok {
						info.Add(newGroupCompB(v))
					}
				case 2: // 删 B（可能拆行）
					if info, ok := wl.GetEntityInfo(entitiesL[ei]); ok {
						info.Remove(&groupCompB{})
					}
					if info, ok := wp.GetEntityInfo(entitiesP[ei]); ok {
						info.Remove(&groupCompB{})
					}
				default: // 新建
					v1, v2 := rngOps.Int31n(100), rngOps.Int31n(100)
					comps := []Component{newGroupCompA(v1), newGroupCompB(v2)}
					entitiesL = append(entitiesL, wl.NewEntity(WithComponents(comps...)))
					entitiesP = append(entitiesP, wp.NewEntity(WithComponents(comps...)))
				}
			}
			if err := wl.Update(); err != nil {
				t.Fatalf("seed %d frame %d linear: %v", seed, frame, err)
			}
			if err := wp.Update(); err != nil {
				t.Fatalf("seed %d frame %d parallel: %v", seed, frame, err)
			}
			// 逐帧对比确定性指标
			if snapL.moveN != snapP.moveN {
				t.Fatalf("seed %d frame %d: moveN %d vs %d", seed, frame, snapL.moveN, snapP.moveN)
			}
			if snapL.bSumMove != snapP.bSumMove || snapL.bSumRO != snapP.bSumRO {
				t.Fatalf("seed %d frame %d: bSum move %d vs %d, ro %d vs %d",
					seed, frame, snapL.bSumMove, snapP.bSumMove, snapL.bSumRO, snapP.bSumRO)
			}
			// 交叉验证：同世界内两条访问路径结果一致
			if snapL.bSumMove != snapL.bSumRO || snapP.bSumMove != snapP.bSumRO {
				t.Fatalf("seed %d frame %d: path mismatch L(%d vs %d) P(%d vs %d)",
					seed, frame, snapL.bSumMove, snapL.bSumRO, snapP.bSumMove, snapP.bSumRO)
			}
			if !reflect.DeepEqual(snapL.c, snapP.c) {
				t.Fatalf("seed %d frame %d: c diverged\nL=%v\nP=%v", seed, frame, snapL.c, snapP.c)
			}
		}
	}
}

// TestParallel_SynctestNoDeadlock synctest 气泡内跑并行世界：
// 气泡内所有 goroutine 持久阻塞即判定死锁（测试失败）；Wait 返回即无泄漏 goroutine。
func TestParallel_SynctestNoDeadlock(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		snap := &parallelSnapshot{}
		w := NewWorld(WithWorldASyncMode())
		if err := registerParallelSystems(w, snap); err != nil {
			t.Fatal(err)
		}
		rng := rand.New(rand.NewSource(42))
		entities := setupParallelEntities(w, rng)

		for frame := 0; frame < 20; frame++ {
			if err := w.Update(); err != nil {
				t.Fatalf("frame %d: %v", frame, err)
			}
			// 间隔注入 churn
			if frame%3 == 0 && len(entities) > 10 {
				ei := rng.Intn(len(entities))
				w.DestroyEntity(entities[ei])
				entities = append(entities[:ei], entities[ei+1:]...)
			}
		}
		// 等待气泡内所有 goroutine 持久阻塞：并行系统/flush 的 goroutine
		// 若泄漏未退出会在此暴露（Wait 不返回则气泡死锁，测试失败）
		synctest.Wait()
		if snap.moveN < 0 {
			t.Fatal("unreachable")
		}
	})
}
