package ecs

import "testing"

// World 泛型方法：system 外的装配/检查代码可类型安全地直读组件，
// 无需经过 ComponentSet 接口断言。注意 system 执行期内应优先使用
// (*SystemContext).GetBuddy（带 WithDep 依赖校验）。
func TestWorld_GetComponent(t *testing.T) {
	w := NewWorld()
	e := w.NewEntity(WithComponents(&dummyComponent{Seq: 7}))
	empty := w.NewEntity()

	w.Update() // 帧同步点：组件操作生效

	// 命中：类型安全直读，无断言
	comp, ok := w.GetComponent[dummyComponent](e)
	if !ok || comp == nil {
		t.Fatal("expected component hit")
	}
	if comp.Seq != 7 {
		t.Fatalf("want Seq=7, got %d", comp.Seq)
	}

	// 返回指针指向组件集内数据，写穿透
	comp.Seq = 8
	again, _ := w.GetComponent[dummyComponent](e)
	if again.Seq != 8 {
		t.Fatalf("expected write-through pointer, got Seq=%d", again.Seq)
	}

	// 未命中：实体无此组件
	if _, ok := w.GetComponent[dummyComponent](empty); ok {
		t.Fatal("expected miss for entity without component")
	}

	// 未命中：组件集不存在（无任何实体持有该类型）
	if _, ok := w.GetComponent[nomadicComp](empty); ok {
		t.Fatal("expected miss for nonexistent component set")
	}
}
