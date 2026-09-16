package ecs

import (
	"testing"
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
)

// dummyComponentReadOnly 手写只读视图（模拟 codegen 产物），验证运行时行为矩阵。
// 视图自描述：携带组件类型标识（ComponentPacketIdentifier）与构造方法（FromPtr），
// 使 API 无需再显式书写组件类型参数。
type dummyComponentReadOnly struct{ p *dummyComponent }

func (v dummyComponentReadOnly) Seq() int32 { return v.p.Seq }

func (v dummyComponentReadOnly) ComponentPacketIdentifier() rockmem.PacketIdentifier {
	return 65530 // 与 dummyComponent.PacketIdentifier() 一致
}

func (v dummyComponentReadOnly) FromPtr(p unsafe.Pointer) dummyComponentReadOnly {
	return dummyComponentReadOnly{p: (*dummyComponent)(p)}
}

func newWorldWithDummies(t *testing.T, seqs ...int32) *World {
	t.Helper()
	w := NewWorld()
	for _, s := range seqs {
		w.NewEntity(WithComponents(&dummyComponent{Seq: s}))
	}
	return w
}

// 只读依赖调用可写 API：GetComponents 返回空迭代器，GetBuddy 返回 nil,false（静默）
func TestReadOnlyDep_WritableAPIReturnsEmpty(t *testing.T) {
	w := newWorldWithDummies(t, 1, 2, 3)

	called := false
	sys := LightSystem(func(ctx *SystemContext, event Event) error {
		called = true
		for range ctx.GetComponents[dummyComponent]() {
			t.Error("GetComponents on readonly dep should be empty")
		}
		for i := EntityIndex(0); i < 3; i++ {
			if b, ok := ctx.GetBuddy[dummyComponent](i); ok || b != nil {
				t.Errorf("GetBuddy on readonly dep should return nil,false, got %v,%v", b, ok)
			}
		}
		return nil
	})
	if err := w.RegisterLight(sys, WithDepReadOnly[dummyComponent]()); err != nil {
		t.Fatal(err)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("system not called")
	}
}

// 只读依赖调用 View API：返回零拷贝视图，读取值正确
func TestReadOnlyDep_ViewAPIReturnsViews(t *testing.T) {
	w := newWorldWithDummies(t, 10, 20, 30)

	sys := LightSystem(func(ctx *SystemContext, event Event) error {
		count := 0
		var sum int32
		for index, view := range ctx.GetComponentsReadOnly[dummyComponentReadOnly]() {
			count++
			sum += view.Seq()
			buddy, ok := ctx.GetBuddyReadOnly[dummyComponentReadOnly](index)
			if !ok || buddy.Seq() != view.Seq() {
				t.Errorf("GetBuddyReadOnly mismatch at index %d", index)
			}
			// 包级函数形态一致
			buddyPkg, okPkg := GetBuddyReadOnly[dummyComponentReadOnly](ctx, index)
			if okPkg != ok || buddyPkg.Seq() != buddy.Seq() {
				t.Errorf("GetBuddyReadOnly parity with package func broken at index %d", index)
			}
		}
		if count != 3 || sum != 60 {
			t.Errorf("want count=3 sum=60, got count=%d sum=%d", count, sum)
		}
		// 不存在的实体索引
		if _, ok := ctx.GetBuddyReadOnly[dummyComponentReadOnly](999); ok {
			t.Error("GetBuddyReadOnly on missing index should return false")
		}
		return nil
	})
	if err := w.RegisterLight(sys, WithDepReadOnly[dummyComponent]()); err != nil {
		t.Fatal(err)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
}

// 可写依赖调用 View API：返回空（View 仅限只读依赖）
func TestWritableDep_ViewAPIReturnsEmpty(t *testing.T) {
	w := newWorldWithDummies(t, 1, 2, 3)

	sys := LightSystem(func(ctx *SystemContext, event Event) error {
		for range ctx.GetComponentsReadOnly[dummyComponentReadOnly]() {
			t.Error("GetComponentsReadOnly on writable dep should be empty")
		}
		if _, ok := ctx.GetBuddyReadOnly[dummyComponentReadOnly](0); ok {
			t.Error("GetBuddyReadOnly on writable dep should return false")
		}
		// 可写 API 正常工作
		count := 0
		for range ctx.GetComponents[dummyComponent]() {
			count++
		}
		if count != 3 {
			t.Errorf("GetComponents on writable dep: want 3, got %d", count)
		}
		return nil
	})
	if err := w.RegisterLight(sys, WithDep[dummyComponent]()); err != nil {
		t.Fatal(err)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
}

// 未声明依赖：View API 返回空（与可写 API 一致的静默约定）
func TestUndeclaredDep_ViewAPIReturnsEmpty(t *testing.T) {
	w := newWorldWithDummies(t, 1)

	sys := LightSystem(func(ctx *SystemContext, event Event) error {
		for range ctx.GetComponentsReadOnly[dummyComponentReadOnly]() {
			t.Error("GetComponentsReadOnly on undeclared dep should be empty")
		}
		if _, ok := ctx.GetBuddyReadOnly[dummyComponentReadOnly](0); ok {
			t.Error("GetBuddyReadOnly on undeclared dep should return false")
		}
		return nil
	})
	// 声明另一个组件的依赖，保证 dummyComponent 未声明
	if err := w.RegisterLight(sys, WithDep[dummyComponent2]()); err != nil {
		t.Fatal(err)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
}
