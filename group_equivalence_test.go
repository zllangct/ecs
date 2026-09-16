package ecs

import (
	"math/rand"
	"reflect"
	"testing"
)

// equivSnapshot 一帧的观测快照：查询命中集 + 组件值
type equivSnapshot struct {
	ab    map[EntityIndex][2]int32 // Query{A,B}: idx → (A.V, B.V)
	aOnly map[EntityIndex]int32    // GetComponents[A] 全覆盖（表内+残段）
}

type equivProbe struct{ snap *equivSnapshot }

func (s *equivProbe) Update(ctx *SystemContext, event Event) error {
	snap := &equivSnapshot{ab: map[EntityIndex][2]int32{}, aOnly: map[EntityIndex]int32{}}
	for idx, a := range ctx.GetComponents[groupCompA, *groupCompA]() {
		snap.aOnly[idx] = a.V
	}
	q := ctx.NewQuery(WithComp[groupCompA, *groupCompA](), WithComp[groupCompB, *groupCompB]())
	for idx := range q.Iter() {
		a, aok := ctx.GetBuddy[groupCompA, *groupCompA](idx)
		b, bok := ctx.GetBuddy[groupCompB, *groupCompB](idx)
		if !aok || !bok {
			continue
		}
		snap.ab[idx] = [2]int32{a.V, b.V}
	}
	*s.snap = *snap
	return nil
}

// TestGroupEquivalence_RandomOps 同一随机 op 序列跑组表开/关两个 World，
// 逐帧对比查询结果集与组件值（回退不变量：散回 CSet 语义完全等价）。
func TestGroupEquivalence_RandomOps(t *testing.T) {
	for seed := int64(0); seed < 20; seed++ {
		rng := rand.New(rand.NewSource(seed))

		on := NewWorld()
		off := NewWorld(WithoutGroups())
		snapOn, snapOff := &equivSnapshot{}, &equivSnapshot{}
		probeOn, probeOff := &equivProbe{snap: snapOn}, &equivProbe{snap: snapOff}

		if err := on.RegisterLight(LightSystem(probeOn.Update),
			WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
			WithGroup(Dep[groupCompA](), Dep[groupCompB]())); err != nil {
			t.Fatal(err)
		}
		if err := off.RegisterLight(LightSystem(probeOff.Update),
			WithDeps(Dep[groupCompA](), Dep[groupCompB]())); err != nil {
			t.Fatal(err)
		}

		var entities []Entity
		for frame := 0; frame < 50; frame++ {
			// 预生成 op 序列，两侧施加相同操作
			type op struct {
				kind int // 0=新建 1=加A 2=删A 3=加B 4=删B 5=销毁
			 ei   int
				v    int32
				withA, withB bool
			}
			var ops []op
			for i := 0; i < 10; i++ {
				o := op{kind: rng.Intn(6), v: rng.Int31n(1000), withA: rng.Intn(2) == 0, withB: rng.Intn(2) == 0}
				if len(entities) > 0 {
					o.ei = rng.Intn(len(entities))
				}
				ops = append(ops, o)
			}
			for _, o := range ops {
				switch o.kind {
				case 0:
					comps := []Component{}
					if o.withA {
						comps = append(comps, newGroupCompA(o.v))
					}
					if o.withB {
						comps = append(comps, newGroupCompB(o.v))
					}
					e1 := on.NewEntity(WithComponents(comps...))
					e2 := off.NewEntity(WithComponents(comps...))
					if e1 != e2 {
						t.Fatalf("seed %d: entity id diverged %v vs %v", seed, e1, e2)
					}
					entities = append(entities, e1)
				case 5:
					if len(entities) == 0 {
						continue
					}
					ei := o.ei % len(entities)
					on.DestroyEntity(entities[ei])
					off.DestroyEntity(entities[ei])
					entities = append(entities[:ei], entities[ei+1:]...)
				default:
					if len(entities) == 0 {
						continue
					}
					e := entities[o.ei%len(entities)]
					for _, w := range []*World{on, off} {
						info, ok := w.GetEntityInfo(e)
						if !ok {
							continue
						}
						switch o.kind {
						case 1:
							info.Add(newGroupCompA(o.v))
						case 2:
							info.Remove(&groupCompA{})
						case 3:
							info.Add(newGroupCompB(o.v))
						case 4:
							info.Remove(&groupCompB{})
						}
					}
				}
			}
			if err := on.Update(); err != nil {
				t.Fatal(err)
			}
			if err := off.Update(); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(snapOn.ab, snapOff.ab) ||
				!reflect.DeepEqual(snapOn.aOnly, snapOff.aOnly) {
				t.Fatalf("seed %d frame %d diverged:\non.ab=%v\noff.ab=%v\non.aOnly=%v\noff.aOnly=%v",
					seed, frame, snapOn.ab, snapOff.ab, snapOn.aOnly, snapOff.aOnly)
			}
		}
	}
}
