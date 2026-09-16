package ecs

import "testing"

// D11: 释放后的 ID 应按升序优先复用最小 index，且 reuse 代数递增
func TestEntityIDGenerator_ReuseSmallestFirst(t *testing.T) {
	g := NewEntityIDGenerator(4, 2)

	e1 := g.NewID()
	e2 := g.NewID()
	e3 := g.NewID()
	if e1.Index() != 1 || e2.Index() != 2 || e3.Index() != 3 {
		t.Fatalf("unexpected initial ids: %d %d %d", e1.Index(), e2.Index(), e3.Index())
	}

	// 释放 e3、e2（乱序），凑满 delayCap 触发 flush
	g.FreeID(e3)
	g.FreeID(e2)

	// 复用应优先拿最小 index，reuse 递增
	n1 := g.NewID()
	if n1.Index() != 2 || n1.toReuseID().reuse != 1 {
		t.Errorf("want index=2 reuse=1, got index=%d reuse=%d", n1.Index(), n1.toReuseID().reuse)
	}
	n2 := g.NewID()
	if n2.Index() != 3 || n2.toReuseID().reuse != 1 {
		t.Errorf("want index=3 reuse=1, got index=%d reuse=%d", n2.Index(), n2.toReuseID().reuse)
	}
	// 链耗尽后应分配新 index
	n3 := g.NewID()
	if n3.Index() != 4 || n3.toReuseID().reuse != 0 {
		t.Errorf("want index=4 reuse=0, got index=%d reuse=%d", n3.Index(), n3.toReuseID().reuse)
	}
}

// D11: 重复释放同一实体必须 panic
func TestEntityIDGenerator_DoubleFreePanics(t *testing.T) {
	g := NewEntityIDGenerator(4, 10)
	e := g.NewID()
	g.FreeID(e)
	defer func() {
		if recover() == nil {
			t.Error("double free should panic")
		}
	}()
	g.FreeID(e)
}

// D11/D12: 用旧代数的句柄释放必须 panic
func TestEntityIDGenerator_StaleFreePanics(t *testing.T) {
	g := NewEntityIDGenerator(4, 2)
	e := g.NewID()
	g.FreeID(e)
	g.FreeID(g.NewID()) // 触发 flush 并释放另一个，凑满 delayCap
	fresh := g.NewID()  // 复用 e 的 index，reuse 已递增
	if fresh.Index() != e.Index() {
		t.Fatalf("expected reuse of index %d, got %d", e.Index(), fresh.Index())
	}
	defer func() {
		if recover() == nil {
			t.Error("stale free should panic")
		}
	}()
	g.FreeID(e) // 旧句柄，reuse 不匹配
	_ = fresh
}

// D11: 反序列化出 delayCap=0 的生成器，FreeID 不得越界
func TestEntityIDGenerator_ZeroDelayCapSafe(t *testing.T) {
	g := &EntityIDGenerator{}
	g.Unmarshal(&SerializableEntityIDGeneratorData{DelayCap: 0})
	e := g.NewID()
	g.FreeID(e) // 不应 panic
}

// D12: EntitySet.Get 必须拒绝旧代数句柄
func TestEntitySetGet_StaleEntity(t *testing.T) {
	es := NewEntitySet()
	old := EntityInfo{entity: Entity(1)} // index=1, reuse=0
	es.Add(old)

	// 模拟 index 被复用：移除后以新代数重新加入
	es.Remove(Entity(1))
	newInfo := EntityInfo{entity: Entity(1 | 1<<32)} // index=1, reuse=1
	es.Add(newInfo)

	if _, ok := es.Get(Entity(1)); ok {
		t.Error("stale entity handle should not resolve")
	}
	if _, ok := es.Get(Entity(1 | 1<<32)); !ok {
		t.Error("fresh entity handle should resolve")
	}
}
