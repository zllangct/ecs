package ecs

import "testing"

// Go 1.27 泛型方法：SystemContext 直接提供类型安全的组件访问 API，
// 行为须与同名包级泛型函数一致。
func TestSystemContext_GenericMethods(t *testing.T) {
	w := NewWorld()
	for i := 0; i < 3; i++ {
		w.NewEntity(WithComponents(&dummyComponent{Seq: int32(i)}))
	}

	called := false
	sys := LightSystem(func(ctx *SystemContext, event Event) error {
		called = true

		// 泛型方法：GetComponents
		count := 0
		var sum int32
		for index, comp := range ctx.GetComponents[dummyComponent]() {
			count++
			sum += comp.Seq

			// 泛型方法：GetBuddy，与包级函数结果一致
			buddy, ok := ctx.GetBuddy[dummyComponent](index)
			if !ok || buddy != comp {
				t.Errorf("GetBuddy mismatch at index %d", index)
			}
			buddyPkg, okPkg := GetBuddy[dummyComponent](ctx, index)
			if okPkg != ok || buddyPkg != buddy {
				t.Errorf("GetBuddy parity with package func broken at index %d", index)
			}
		}
		if count != 3 || sum != 3 {
			t.Errorf("want count=3 sum=3, got count=%d sum=%d", count, sum)
		}

		// 方法形态：NewQuery，与包级函数迭代结果一致
		q := ctx.NewQuery(WithComp[dummyComponent]())
		qCount := 0
		for range q.Iter() {
			qCount++
		}
		if qCount != 3 {
			t.Errorf("want 3 query results, got %d", qCount)
		}

		return nil
	})
	if err := w.RegisterLight(sys, WithDep[dummyComponent]()); err != nil {
		t.Fatal(err)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("system not called")
	}
}
