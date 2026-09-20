package ecs

import (
	"strings"
	"testing"
)

// R19：未注册组件类型恢复必须返回 error（不得 panic，也不得静默跳过业务状态）
func TestUnmarshalUnknownComponentFails(t *testing.T) {
	// testWorldComponent 已在 serialization_test.go init 注册
	c := NewCSet[testWorldComponent]()
	c.Add(Entity(1), &testWorldComponent{Value: 1})
	csetData, err := c.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	data := &SerializableWorldData{
		Version: SnapshotVersion,
		ComponentSets: []SerializableComponentSetEntry{
			{ComponentType: 0xdeadbeef, Data: *csetData}, // 未注册的组件类型
		},
	}

	restored := &World{}
	err = restored.Unmarshal(data)
	if err == nil || !strings.Contains(err.Error(), "not registered") {
		t.Fatalf("unknown component must fail with explicit error, got %v", err)
	}
	// 原子失败：目标世界的组件表不被部分替换
	if len(restored.components) != 0 {
		t.Fatal("failed unmarshal must not partially install components")
	}
}

// R19：USet 布局不匹配（EleSize 漂移）返回 error
func TestUSetUnmarshalEleSizeMismatch(t *testing.T) {
	u := NewUSet[dummyComponent]()
	u.Add(&dummyComponent{Seq: 1})
	data := u.Marshal()
	data.EleSize += 8 // 模拟组件加字段后的旧快照/伪造输入

	restored := NewUSet[dummyComponent]()
	if err := restored.Unmarshal(data); err == nil {
		t.Fatal("element size mismatch must fail")
	}
}

// R19：IDL 偏移表损坏（截断/非单调）返回 error
func TestCSetIDLCorruptedOffsetTable(t *testing.T) {
	c := NewCSet[stringyComp]()
	c.Add(Entity(1), &stringyComp{Name: "hello", Seq: 1})
	c.Add(Entity(2), &stringyComp{Name: "world!!", Seq: 2})
	data, err := c.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// 截断 Data（偏移表不完整）
	truncated := *data
	truncated.USetData.Data = data.USetData.Data[:4]
	if _, err := NewCSetFromData[stringyComp](&truncated); err == nil {
		t.Fatal("truncated IDL header must fail")
	}

	// 非单调偏移：offsets[1]=100 而 offsets[2]=4
	bad := *data
	bad.USetData.Data = append([]byte(nil), data.USetData.Data...)
	bad.USetData.Data[8] = 100
	bad.USetData.Data[16] = 4
	if _, err := NewCSetFromData[stringyComp](&bad); err == nil {
		t.Fatal("non-monotonic offset table must fail")
	}

	// blob 截断：offsets[n] 超出数据界
	bad2 := *data
	bad2.USetData.Data = append([]byte(nil), data.USetData.Data...)
	bad2.USetData.Data[16] = 0xff
	if _, err := NewCSetFromData[stringyComp](&bad2); err == nil {
		t.Fatal("truncated IDL blob must fail")
	}
}

// R19：版本不匹配拒绝恢复
func TestUnmarshalVersionMismatch(t *testing.T) {
	w := NewWorld()
	data, err := w.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	data.Version = SnapshotVersion + 1
	restored := &World{}
	if err := restored.Unmarshal(data); err == nil {
		t.Fatal("version mismatch must fail")
	}
}

// R19：运行中的世界拒绝原地恢复（必须先静止）
func TestRunningWorldRejectsUnmarshal(t *testing.T) {
	w := NewWorld()
	w.NewEntity()
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	data, err := w.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if w.status != WorldStatusRunning {
		t.Skip("world not running")
	}
	if err := w.Unmarshal(data); err == nil {
		t.Fatal("unmarshal into running world must be rejected")
	}
}

// R19：原地恢复整体替换组件表（旧组件消失），查询缓存失效
func TestUnmarshalReplacesComponentsEntirely(t *testing.T) {
	// 世界 A：含两种组件，经 Update 吸纳入表后导出快照
	wA := NewWorld()
	eA := wA.NewEntity()
	infoA, _ := wA.GetEntityInfo(eA)
	infoA.Add(&testWorldComponent{Value: 1})
	infoA.Add(&testWorldComponent2{X: 1, Y: 2})
	if err := wA.Update(); err != nil {
		t.Fatal(err)
	}
	dataA, err := wA.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// 目标世界：Initialized 状态、组件池含 testWorldComponent2
	w, err := NewWorldFromData(dataA)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := w.components[GetIntTypeByComp(&testWorldComponent2{})]; !ok {
		t.Fatal("baseline must contain type2 component pool")
	}

	// 世界 B：只含 testWorldComponent 的快照
	wB := NewWorld()
	eB := wB.NewEntity()
	infoB, _ := wB.GetEntityInfo(eB)
	infoB.Add(&testWorldComponent{Value: 99})
	if err := wB.Update(); err != nil {
		t.Fatal(err)
	}
	dataB, err := wB.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// 原地恢复：整体替换
	if err := w.Unmarshal(dataB); err != nil {
		t.Fatal(err)
	}
	// 旧组件 testWorldComponent2 的组件池必须整体消失
	if _, ok := w.components[GetIntTypeByComp(&testWorldComponent2{})]; ok {
		t.Fatal("stale component set must be removed after full replace")
	}
	got, ok := w.GetComponent[testWorldComponent, *testWorldComponent](eB)
	if !ok || got.Value != 99 {
		t.Fatalf("new component data wrong: %+v %v", got, ok)
	}
}

// R19：实体 parallel arrays 不一致返回 error
func TestEntitySetInvalidParallelArrays(t *testing.T) {
	es := NewEntitySet()
	data := es.Marshal()
	data.EntityIds = []int64{1, 2}
	data.CompoundOffsets = []int32{0} // 长度不一致
	if err := es.Unmarshal(data); err == nil {
		t.Fatal("parallel arrays mismatch must fail")
	}
}
