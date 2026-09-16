package ecs

import (
	"testing"
	"testing/synctest"
	"time"
)

// D06: delta 必须是本帧间隔，不得累积
func TestWorldUpdate_DeltaIsFrameInterval(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		w := NewWorld()
		var deltas []time.Duration
		sys := LightSystem(func(ctx *SystemContext, event Event) error {
			deltas = append(deltas, event.Delta)
			return nil
		})
		if err := w.RegisterLight(sys, WithDep[dummyComponent]()); err != nil {
			t.Fatal(err)
		}
		w.Update()
		time.Sleep(100 * time.Millisecond)
		w.Update()
		time.Sleep(100 * time.Millisecond)
		w.Update()
		synctest.Wait()

		if len(deltas) != 3 {
			t.Fatalf("want 3 deltas, got %d", len(deltas))
		}
		// 第三帧 delta 应约等于 100ms，而非从首帧累积的 200ms
		if deltas[2] > 150*time.Millisecond {
			t.Errorf("delta accumulated across frames: %v", deltas[2])
		}
	})
}

// D07: StageSyncAfterDestroy 必须被执行
// 泛型 Register 在内部创建 system 实例，外部可观测状态通过钩子函数回传
var (
	destroyStageOnDestroy          func()
	destroyStageOnSyncAfterDestroy func()
)

type destroyStageSys struct{}

func (s *destroyStageSys) Init(ctx *SystemInitContext) error {
	ctx.SetOption(WithName("destroyStageSys"), WithDep[dummyComponent]())
	return nil
}

// Update 为 SystemPointer 约束所要求的空实现
func (s *destroyStageSys) Update(ctx *SystemContext, event Event) error {
	return nil
}

func (s *destroyStageSys) Destroy(ctx *SystemContext, event Event) error {
	if destroyStageOnDestroy != nil {
		destroyStageOnDestroy()
	}
	return nil
}

func (s *destroyStageSys) SyncAfterDestroy(ctx *SystemContext, event Event) error {
	if destroyStageOnSyncAfterDestroy != nil {
		destroyStageOnSyncAfterDestroy()
	}
	return nil
}

func TestSyncAfterDestroy_Executed(t *testing.T) {
	w := NewWorld()
	destroyCalled := false
	syncAfterDestroyOk := false
	destroyStageOnDestroy = func() { destroyCalled = true }
	destroyStageOnSyncAfterDestroy = func() { syncAfterDestroyOk = true }
	defer func() {
		destroyStageOnDestroy = nil
		destroyStageOnSyncAfterDestroy = nil
	}()
	if err := w.Register[destroyStageSys](); err != nil {
		t.Fatal(err)
	}
	// 进入 Destroy 状态（当前无公开 API，测试直接置状态）
	ww := w
	for _, info := range ww.systems.systems {
		info.setState(SystemStateDestroy)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if !destroyCalled {
		t.Error("Destroy not called")
	}
	if !syncAfterDestroyOk {
		t.Error("SyncAfterDestroy not called (StageSyncAfterDestroy missing)")
	}
}

// D13: 自定义 Order 必须按升序执行
// 泛型 Register 在内部创建 system 实例，执行记录通过钩子函数回传，
// 名称与顺序经 Register 的 SystemOption 传入
var orderRecordOnUpdate func(name string)

type orderRecordSys struct{}

func (s *orderRecordSys) Init(ctx *SystemInitContext) error {
	ctx.SetOption(WithDep[dummyComponent]())
	return nil
}

func (s *orderRecordSys) Update(ctx *SystemContext, event Event) error {
	if orderRecordOnUpdate != nil {
		orderRecordOnUpdate(ctx.info.name())
	}
	return nil
}

func TestFlowRegister_CustomOrder(t *testing.T) {
	w := NewWorld()
	var seq []string
	orderRecordOnUpdate = func(name string) { seq = append(seq, name) }
	defer func() { orderRecordOnUpdate = nil }()
	// 乱序注册 10、5、7，期望执行顺序 5、7、10
	for _, spec := range []struct {
		name  string
		order Order
	}{{"s10", 10}, {"s5", 5}, {"s7", 7}} {
		if err := w.Register[orderRecordSys](WithName(spec.name), WithOrder(spec.order), WithDep[dummyComponent]()); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	want := []string{"s5", "s7", "s10"}
	if len(seq) != 3 {
		t.Fatalf("want 3 executions, got %v", seq)
	}
	for i := range want {
		if seq[i] != want[i] {
			t.Fatalf("order wrong: want %v, got %v", want, seq)
		}
	}
}

// D14: 未实现对应接口的 system 在 getSystemTask 中不得 panic
func TestGetSystemTask_NoImplNoPanic(t *testing.T) {
	w := NewWorld()
	if err := w.Register[orderRecordSys](WithName("onlyUpdate"), WithDep[dummyComponent]()); err != nil {
		t.Fatal(err)
	}
	var info SystemInfo
	for _, si := range w.systems.systems {
		info = si
	}
	// orderRecordSys 只实现 Update，state=Start 时查 StageStart 任务应返回无效而非 panic
	task, err := w.systems.getSystemTask(info, StageStart)
	if err != nil {
		t.Fatal(err)
	}
	if task.isValid {
		t.Error("task should be invalid for unimplemented stage")
	}
}
