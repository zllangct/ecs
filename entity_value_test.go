package ecs

import "testing"

// D35: NewEntity 返回 Entity 值（不暴露内部指针），EntityInfo 按需查询
func TestNewEntity_ReturnsEntityValue(t *testing.T) {
	w := NewWorld()
	e := w.NewEntity(WithComponents(&dummyComponent{Seq: 1}))

	// 按需查询：EntitySet 扩容后仍能从最新存储获取
	for i := 0; i < 100; i++ {
		w.NewEntity()
	}
	info, ok := w.GetEntityInfo(e)
	if !ok {
		t.Fatal("GetEntityInfo should resolve after EntitySet growth")
	}
	if info.Entity() != e {
		t.Errorf("entity mismatch: %v vs %v", info.Entity(), e)
	}

	// 不存在/已销毁的实体
	if _, ok := w.GetEntityInfo(Entity(1 << 40)); ok {
		t.Error("non-existent entity should not resolve")
	}
}

// D35: EntityTemplate 返回 Entity 值
func TestEntityTemplate_ReturnsEntityValue(t *testing.T) {
	tpl := &EntityTemplate{Components: []Component{&dummyComponent{Seq: 5}}}
	w := NewWorld()
	e := tpl.Instance(w)
	info, ok := w.GetEntityInfo(e)
	if !ok {
		t.Fatal("template instance not found")
	}
	_ = info
	es := tpl.InstanceN(w, 3)
	if len(es) != 3 {
		t.Fatalf("want 3 entities, got %d", len(es))
	}
}
