package ecs

import "testing"

// D10: Query.Iter 必须响应 break
func TestQueryIter_EarlyBreak(t *testing.T) {
	w := NewWorld()
	for i := 0; i < 2; i++ {
		w.NewEntity(WithComponents(&dummyComponent{Seq: int32(i)}))
	}

	called := false
	sys := LightSystem(func(ctx *SystemContext, event Event) error {
		called = true
		q := NewQuery(ctx, WithComp[dummyComponent]())
		count := 0
		for range q.Iter() {
			count++
			break
		}
		if count != 1 {
			t.Errorf("want 1 iteration, got %d", count)
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
