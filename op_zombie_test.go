package ecs

import "testing"

// 销毁与组件 op 同帧入队时，op flush 不得把组件应用到已销毁实体（僵尸数据）。
// 时间线：主队列 Add B → DestroyEntity → 次帧 destroy 先清、op flush 后到。
func TestOpLog_OpOnDestroyedEntityDropped(t *testing.T) {
	w := NewWorld()
	w.RegisterLight(LightSystem(groupProbeLight),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()))
	e := w.NewEntity(WithComponents(newGroupCompA(1)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}

	info, ok := w.GetEntityInfo(e)
	if !ok {
		t.Fatal("entity missing")
	}
	info.Add(newGroupCompB(514)) // 主队列 op，次帧生效
	w.DestroyEntity(e)           // 同帧销毁，次帧先执行
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}

	// B 池不得残留已销毁实体的数据
	if s, ok := w.getComponentSet(GetIntType[groupCompB, *groupCompB]()); ok {
		if s.Exist(e.Index()) {
			t.Fatal("zombie component: B applied to destroyed entity")
		}
	}
}

// 存活的实体 op 正常生效（有效性检查不误伤）
func TestOpLog_OpOnAliveEntityApplied(t *testing.T) {
	w := NewWorld()
	w.RegisterLight(LightSystem(groupProbeLight),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()))
	e := w.NewEntity(WithComponents(newGroupCompA(1)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	info, _ := w.GetEntityInfo(e)
	info.Add(newGroupCompB(514))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	got, ok := w.GetComponent[groupCompB, *groupCompB](e)
	if !ok || got.V != 514 {
		t.Fatalf("B = %+v %v", got, ok)
	}
}
