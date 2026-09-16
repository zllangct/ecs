package ecs

import "testing"

// D08: Dep 默认应为可写（readonly=false）
func TestDep_DefaultWritable(t *testing.T) {
	d := Dep[dummyComponent]()
	if d.readonly() {
		t.Error("default Dep should be writable (readonly=false)")
	}
}

// D08: WithDepReadOnly 应为只读（readonly=true）
func TestDepReadOnly_IsReadonly(t *testing.T) {
	c := &SystemConfig{}
	c.initDefault()
	WithDepReadOnly[dummyComponent]()(c)
	if len(c.dependencies) != 1 {
		t.Fatalf("want 1 dep, got %d", len(c.dependencies))
	}
	if !c.dependencies[0].readonly() {
		t.Error("WithDepReadOnly should be readonly")
	}
}

// D08: WithDep(ReadWrite) 显式可写
func TestWithDep_ExplicitWritable(t *testing.T) {
	c := &SystemConfig{}
	c.initDefault()
	WithDep[dummyComponent](ReadWrite)(c)
	if c.dependencies[0].readonly() {
		t.Error("WithDep(ReadWrite) should be writable")
	}
	WithDep[dummyComponent2](ReadOnly)(c)
	if !c.dependencies[1].readonly() {
		t.Error("WithDep(ReadOnly) should be readonly")
	}
}

// D08: 并行分组冲突判断——双只读同组件不冲突，含可写则冲突
func TestIsFriend_ReadonlySemantics(t *testing.T) {
	mkNode := func(deps ...ComponentDependency) *systemTreeNode {
		info := newSystem(nil, nil, SystemTypeLight)
		info.config.dependencies = deps
		return &systemTreeNode{val: info}
	}
	roA := mkNode(Dep[dummyComponent](ReadOnly))
	roB := mkNode(Dep[dummyComponent](ReadOnly))
	rwC := mkNode(Dep[dummyComponent](ReadWrite))

	if roA.isFriend(roB) {
		t.Error("two readonly deps on same component should NOT conflict")
	}
	if !roA.isFriend(rwC) {
		t.Error("readonly vs writable on same component should conflict")
	}
}

// D08/D09（方案 R2）：只读依赖调可写 API 静默返回空，只读访问走 codegen 零拷贝视图。
// 完整行为矩阵见 readonly_access_test.go。
