package ecs

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
)

// ============================================================================
// Phase 2 测试专用组件
// ============================================================================

type p2CompA struct{ V int32 }

func (c *p2CompA) NewComponentSet() ComponentSet              { return NewCSet[p2CompA]() }
func (c *p2CompA) PacketIdentifier() rockmem.PacketIdentifier { return 99999999999001 }
func (c *p2CompA) IsNomadic() bool                            { return false }
func (c *p2CompA) IsDisposable() bool                         { return false }

type p2CompB struct{ V int32 }

func (c *p2CompB) NewComponentSet() ComponentSet              { return NewCSet[p2CompB]() }
func (c *p2CompB) PacketIdentifier() rockmem.PacketIdentifier { return 99999999999002 }
func (c *p2CompB) IsNomadic() bool                            { return false }
func (c *p2CompB) IsDisposable() bool                         { return false }

type p2CompBView struct{ p *p2CompB }

func (v p2CompBView) ComponentPacketIdentifier() rockmem.PacketIdentifier { return 99999999999002 }
func (v p2CompBView) FromPtr(p unsafe.Pointer) p2CompBView                { return p2CompBView{p: (*p2CompB)(p)} }
func (v p2CompBView) V() int32                                            { return v.p.V }

type p2CompC struct{ V int32 }

func (c *p2CompC) NewComponentSet() ComponentSet              { return NewCSet[p2CompC]() }
func (c *p2CompC) PacketIdentifier() rockmem.PacketIdentifier { return 99999999999003 }
func (c *p2CompC) IsNomadic() bool                            { return false }
func (c *p2CompC) IsDisposable() bool                         { return false }

type p2CompTag struct{ V int32 }

func (c *p2CompTag) NewComponentSet() ComponentSet              { return NewCSet[p2CompTag]() }
func (c *p2CompTag) PacketIdentifier() rockmem.PacketIdentifier { return 99999999999004 }
func (c *p2CompTag) IsNomadic() bool                            { return false }
func (c *p2CompTag) IsDisposable() bool                         { return false }

type p2PtrComp struct {
	Label string
	V     int32
}

func (c *p2PtrComp) NewComponentSet() ComponentSet              { return NewCSet[p2PtrComp]() }
func (c *p2PtrComp) PacketIdentifier() rockmem.PacketIdentifier { return 99999999999005 }
func (c *p2PtrComp) IsNomadic() bool                            { return false }
func (c *p2PtrComp) IsDisposable() bool                         { return false }

type p2ZeroTag struct{}

func (c *p2ZeroTag) NewComponentSet() ComponentSet              { return NewCSet[p2ZeroTag]() }
func (c *p2ZeroTag) PacketIdentifier() rockmem.PacketIdentifier { return 99999999999006 }
func (c *p2ZeroTag) IsNomadic() bool                            { return false }
func (c *p2ZeroTag) IsDisposable() bool                         { return false }

func init() {
	RegisterComponent[p2CompA]("p2")
	RegisterComponent[p2CompB]("p2")
	RegisterComponent[p2CompC]("p2")
	RegisterComponent[p2CompTag]("p2")
	RegisterComponent[p2PtrComp]("p2")
	RegisterComponent[p2ZeroTag]("p2")
}

// ============================================================================
// 4.1 组表子集查询与混合查询等价性
// ============================================================================

func TestQueryAndView_MergedGroupSubsetIncludesResidual(t *testing.T) {
	run := func(disableGroups bool) (map[EntityIndex][2]int32, map[[2]int32]int, map[[2]int32]int) {
		var opts []WorldOption
		if disableGroups {
			opts = append(opts, WithoutGroups())
		}
		w := NewWorld(opts...)
		// 声明两个重叠组，Finalize 后合并为 3 列组 {A, B, C}
		w.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error { return nil }),
			WithDeps(Dep[p2CompA](), Dep[p2CompB](), Dep[p2CompC]()),
			WithGroup(Dep[p2CompA](), Dep[p2CompB]()),
			WithGroup(Dep[p2CompB](), Dep[p2CompC]()),
		)
		qRes := map[EntityIndex][2]int32{}
		v2Res := map[[2]int32]int{}
		v2roRes := map[[2]int32]int{}
		w.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error {
			q := ctx.NewQuery(WithComp[p2CompA](), WithComp[p2CompB]())
			for idx := range q.Iter() {
				a, aok := QueryGet[p2CompA](&q, idx)
				b, bok := QueryGet[p2CompB](&q, idx)
				if !aok || !bok {
					t.Fatalf("QueryGet missing for idx %d: aok=%v bok=%v", idx, aok, bok)
				}
				qRes[idx] = [2]int32{a.V, b.V}
			}
			for a, b := range View2[p2CompA, p2CompB](ctx) {
				v2Res[[2]int32{a.V, b.V}]++
			}
			for a, bv := range View2RO[p2CompA, p2CompBView](ctx) {
				v2roRes[[2]int32{a.V, bv.V()}]++
			}
			return nil
		}), WithDeps(Dep[p2CompA](), Dep[p2CompB](ReadOnly)), WithName("move"))
		// 重建一个只写 A/B 的移动系统以便测试 View2
		return qRes, v2Res, v2roRes
	}
	_, _, _ = run(false)

	checkMode := func(disableGroups bool) (map[EntityIndex][2]int32, map[[2]int32]int, map[[2]int32]int, map[[2]int32]int) {
		var opts []WorldOption
		if disableGroups {
			opts = append(opts, WithoutGroups())
		}
		w := NewWorld(opts...)
		qRes := map[EntityIndex][2]int32{}
		v2Res := map[[2]int32]int{}
		v3Res := map[[2]int32]int{}
		v2roRes := map[[2]int32]int{}
		if err := w.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error {
			q := ctx.NewQuery(WithComp[p2CompA](), WithComp[p2CompB]())
			for idx := range q.Iter() {
				a, aok := QueryGet[p2CompA](&q, idx)
				b, bok := QueryGet[p2CompB](&q, idx)
				if !aok || !bok {
					t.Fatalf("QueryGet missing idx %d", idx)
				}
				qRes[idx] = [2]int32{a.V, b.V}
			}
			for a, b := range View2[p2CompA, p2CompB](ctx) {
				v2Res[[2]int32{a.V, b.V}]++
			}
			for tup := range View3[p2CompA, p2CompB, p2CompC](ctx) {
				v3Res[[2]int32{tup.V1.V, tup.V2.V}]++
			}
			return nil
		}),
			WithDeps(Dep[p2CompA](), Dep[p2CompB](), Dep[p2CompC]()),
			WithGroup(Dep[p2CompA](), Dep[p2CompB]()),
			WithGroup(Dep[p2CompB](), Dep[p2CompC]()),
		); err != nil {
			t.Fatal(err)
		}
		if err := w.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error {
			for a, bv := range View2RO[p2CompA, p2CompBView](ctx) {
				v2roRes[[2]int32{a.V, bv.V()}]++
			}
			return nil
		}), WithDeps(Dep[p2CompA](), Dep[p2CompB](ReadOnly))); err != nil {
			t.Fatal(err)
		}

		// 玩家（A,B,C -> 进 3 列组表）
		w.NewEntity(WithComponents(&p2CompA{10}, &p2CompB{11}, &p2CompC{12}))
		// 飞箭（A,B -> 留 CSet 残段）
		w.NewEntity(WithComponents(&p2CompA{20}, &p2CompB{21}))
		// 仅 A（不命中）
		w.NewEntity(WithComponents(&p2CompA{30}))

		if err := w.Update(); err != nil {
			t.Fatal(err)
		}
		return qRes, v2Res, v3Res, v2roRes
	}

	qOn, v2On, v3On, v2roOn := checkMode(false)
	qOff, v2Off, v3Off, v2roOff := checkMode(true)
	if len(qOn) != 2 || !reflect.DeepEqual(qOn, qOff) {
		t.Fatalf("subset query mismatch: on=%v off=%v", qOn, qOff)
	}
	if len(v2On) != 2 || !reflect.DeepEqual(v2On, v2Off) {
		t.Fatalf("subset View2 mismatch: on=%v off=%v", v2On, v2Off)
	}
	if len(v3On) != 1 || !reflect.DeepEqual(v3On, v3Off) {
		t.Fatalf("full View3 mismatch: on=%v off=%v", v3On, v3Off)
	}
	if len(v2roOn) != 2 || !reflect.DeepEqual(v2roOn, v2roOff) {
		t.Fatalf("subset View2RO mismatch: on=%v off=%v", v2roOn, v2roOff)
	}
}

func TestQuery_MixedGroupedAndUngrouped(t *testing.T) {
	check := func(disableGroups bool) map[EntityIndex][2]int32 {
		var opts []WorldOption
		if disableGroups {
			opts = append(opts, WithoutGroups())
		}
		w := NewWorld(opts...)
		hits := map[EntityIndex][2]int32{}
		if err := w.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error {
			q := ctx.NewQuery(WithComp[p2CompA](), WithComp[p2CompTag]())
			for idx := range q.Iter() {
				a, aok := QueryGet[p2CompA](&q, idx)
				tg, tok := QueryGet[p2CompTag](&q, idx)
				if !aok || !tok {
					t.Fatalf("mixed QueryGet failed at %d: aok=%v tok=%v", idx, aok, tok)
				}
				hits[idx] = [2]int32{a.V, tg.V}
			}
			return nil
		}),
			WithDeps(Dep[p2CompA](), Dep[p2CompB](), Dep[p2CompTag]()),
			WithGroup(Dep[p2CompA](), Dep[p2CompB]()),
		); err != nil {
			t.Fatal(err)
		}
		// 1) 全在组表内的实体（CSet[A] 对其为空，带 Tag）
		w.NewEntity(WithComponents(&p2CompA{1}, &p2CompB{2}, &p2CompTag{100}))
		// 2) 仅先放组表实体时，校验第一帧不会误判 empty
		if err := w.Update(); err != nil {
			t.Fatal(err)
		}
		if len(hits) != 1 {
			t.Fatalf("disableGroups=%v: want 1 grouped+tag hit, got %v", disableGroups, hits)
		}
		// 3) 第二帧再加一个残段实体（只有 A + Tag）
		hits = map[EntityIndex][2]int32{}
		w.NewEntity(WithComponents(&p2CompA{3}, &p2CompTag{200}))
		if err := w.Update(); err != nil {
			t.Fatal(err)
		}
		return hits
	}
	on := check(false)
	off := check(true)
	if len(on) != 2 || !reflect.DeepEqual(on, off) {
		t.Fatalf("mixed query mismatch: on=%v off=%v", on, off)
	}
}

// ============================================================================
// 4.2 & ECS-2 并行调度器：菱形冲突、注册顺序、零依赖、Sync* 串行、PreUpdate 异步
// ============================================================================

func TestScheduler_DiamondConflictAndRegistrationOrder(t *testing.T) {
	sg := NewSystemRelatedGroups().(*SystemRelatedGroups)
	var execSeq []string
	mk := func(name string, deps ...ComponentDependency) SystemInfo {
		s := newSystem(nil, LightSystem(func(ctx *SystemContext, event Event) error {
			execSeq = append(execSeq, name)
			return nil
		}), SystemTypeLight)
		s.init(WithName(name), WithDeps(deps...))
		_ = sg.add(s)
		return s
	}

	s1 := mk("S1", Dep[p2CompA]())
	s2 := mk("S2", Dep[p2CompB]())
	s3 := mk("S3", Dep[p2CompA](), Dep[p2CompB]())
	s4 := mk("S4", Dep[p2CompB](), Dep[p2CompC]())
	s0 := mk("S0") // 零依赖

	var batches [][]string
	for b := range sg.independentGroups() {
		var names []string
		for _, sys := range b {
			names = append(names, sys.name())
		}
		batches = append(batches, names)
	}
	// 重复迭代 independentGroups 必须输出完全一致的结果
	var batches2 [][]string
	for b := range sg.independentGroups() {
		var names []string
		for _, sys := range b {
			names = append(names, sys.name())
		}
		batches2 = append(batches2, names)
	}
	if !reflect.DeepEqual(batches, batches2) {
		t.Fatalf("independentGroups not repeatable: %v vs %v", batches, batches2)
	}

	// 验证同批全对全无读写冲突，且所有系统都在（包括零依赖 S0）
	total := 0
	batchOf := map[string]int{}
	for bi, b := range batches {
		for _, name := range b {
			total++
			batchOf[name] = bi
		}
	}
	if total != 5 {
		t.Fatalf("want 5 systems scheduled (incl. zero-dep S0), got %v", batches)
	}
	// S3 与 S4 均写 B，绝不能同批；且注册序 S3 < S4 => batchOf[S3] < batchOf[S4]
	if batchOf["S3"] >= batchOf["S4"] {
		t.Fatalf("diamond conflict/order broken: S3 batch=%d, S4 batch=%d, batches=%v",
			batchOf["S3"], batchOf["S4"], batches)
	}
	if batchOf["S1"] >= batchOf["S3"] || batchOf["S2"] >= batchOf["S3"] {
		t.Fatalf("dependency order broken: %v", batches)
	}
	_ = s1
	_ = s2
	_ = s3
	_ = s4
	_ = s0
}

type p2SyncAndDestroySys struct {
	inFlight    *atomic.Int32
	overlapSeen *atomic.Bool
	syncRan     *atomic.Bool
	destroyRan  *atomic.Bool
	afterRan    *atomic.Bool
	preAsync    *atomic.Bool
}

func (s *p2SyncAndDestroySys) Init(ctx *SystemInitContext) error { return nil }
func (s *p2SyncAndDestroySys) Update(ctx *SystemContext, e Event) error {
	return nil
}
func (s *p2SyncAndDestroySys) PreUpdate(ctx *SystemContext, e Event) error {
	if s.preAsync != nil {
		s.preAsync.Store(true)
	}
	return nil
}
func (s *p2SyncAndDestroySys) SyncBeforeUpdate(ctx *SystemContext, e Event) error {
	if s.inFlight.Add(1) > 1 {
		s.overlapSeen.Store(true)
	}
	time.Sleep(5 * time.Millisecond)
	s.inFlight.Add(-1)
	s.syncRan.Store(true)
	return nil
}
func (s *p2SyncAndDestroySys) Destroy(ctx *SystemContext, e Event) error {
	s.destroyRan.Store(true)
	return nil
}
func (s *p2SyncAndDestroySys) SyncAfterDestroy(ctx *SystemContext, e Event) error {
	s.afterRan.Store(true)
	return nil
}

// 验证 SyncAfterDestroyReceiver 类型别名存在且可赋值
var _ SyncAfterDestroyReceiver = (*p2SyncAndDestroySys)(nil)

func TestParallelMode_SyncStagesSerialAndZeroDepDestroy(t *testing.T) {
	w := NewWorld(WithWorldASyncMode())
	var inFlight atomic.Int32
	var overlap atomic.Bool
	var sync1, sync2, dest1, after1, pre1 atomic.Bool

	// 注册两个无共享依赖的系统：SyncBeforeUpdate 必须串行执行，Destroy 必须触发
	if err := w.Register[p2SyncAndDestroySys](); err != nil {
		t.Fatal(err)
	}
	if err := w.Register[p2SyncAndDestroySys](); err != nil {
		t.Fatal(err)
	}
	// 注入状态指针到注册后的系统实例
	i := 0
	for _, info := range w.systems.systems {
		raw := info.getRaw().(*p2SyncAndDestroySys)
		raw.inFlight = &inFlight
		raw.overlapSeen = &overlap
		if i == 0 {
			raw.syncRan = &sync1
			raw.destroyRan = &dest1
			raw.afterRan = &after1
			raw.preAsync = &pre1
		} else {
			raw.syncRan = &sync2
			raw.destroyRan = &atomic.Bool{}
			raw.afterRan = &atomic.Bool{}
		}
		i++
	}
	// 验证 PreUpdate 为异步阶段，SyncBeforeUpdate 为同步阶段
	if isSyncStage(StagePreUpdate) || !isSyncStage(StageSyncBeforeUpdate) {
		t.Fatalf("unexpected sync stage flags: PreUpdate=%v SyncBeforeUpdate=%v",
			isSyncStage(StagePreUpdate), isSyncStage(StageSyncBeforeUpdate))
	}
	for _, info := range w.systems.systems {
		task, err := w.systems.getSystemTask(info, StagePreUpdate)
		if err != nil || !task.isValid {
			t.Fatalf("PreUpdate task should be valid: task=%+v err=%v", task, err)
		}
		break
	}

	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if overlap.Load() {
		t.Fatal("SyncBeforeUpdate ran concurrently in parallel mode")
	}
	if !sync1.Load() || !sync2.Load() || !pre1.Load() {
		t.Fatal("zero-dependency systems were dropped in parallel mode")
	}
	if err := w.Destroy(); err != nil {
		t.Fatal(err)
	}
	if !dest1.Load() || !after1.Load() {
		t.Fatal("zero-dependency Destroy / SyncAfterDestroy skipped in parallel mode")
	}
}

// ============================================================================
// 4.3 快照反序列化校验与恢复卫生
// ============================================================================

func TestSerialization_CorruptSnapshotAndRestoreHygiene(t *testing.T) {
	// 1. 预置生命值 37 实体的世界，加载含重复组件池且不同帧号的损坏存档 -> 原子失败且原世界不受影响
	w := NewWorld()
	e := w.NewEntity(WithComponents(&p2CompA{37}))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	origFrame := w.frame

	w2 := NewWorld()
	w2.NewEntity(WithComponents(&p2CompA{99}))
	if err := w2.Update(); err != nil {
		t.Fatal(err)
	}
	badSnap, err := w2.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	badSnap.Frame = 777
	badSnap.ComponentSets = append(badSnap.ComponentSets, badSnap.ComponentSets[0]) // 重复组件池

	wTarget := NewWorld()
	eTarget := wTarget.NewEntity(WithComponents(&p2CompA{37}))
	if err := wTarget.Update(); err != nil {
		t.Fatal(err)
	}
	// wTarget 处于 Running 也会拒绝；测试 UnmarshalInto Initialized 世界：
	wInit, err := NewWorldFromData(func() *SerializableWorldData {
		d, _ := w.Marshal()
		return d
	}())
	if err != nil {
		t.Fatal(err)
	}
	if err := wInit.Unmarshal(badSnap); err == nil {
		t.Fatal("duplicate component pool must fail")
	}
	if wInit.frame != origFrame {
		t.Fatalf("frame mutated on failed unmarshal: got %d want %d", wInit.frame, origFrame)
	}
	if got, ok := wInit.GetComponent[p2CompA, *p2CompA](e); !ok || got.V != 37 {
		t.Fatalf("existing entity mutated on failed unmarshal: got %+v ok=%v", got, ok)
	}
	_ = eTarget

	// 2. 空集合恢复后 EleSize 不得被置为 0 或错误宽度
	emptyC := NewCSet[p2CompA]()
	emptyData, err := emptyC.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	emptyData.USetData.EleSize = 0
	if _, err := NewCSetFromData[p2CompA](emptyData); err == nil {
		t.Fatal("empty CSet with EleSize=0 must fail")
	}

	// 3. IDL 负偏移 offsets[0] = -1 必须拒绝（不得 panic）
	idlC := NewCSet[stringyComp]()
	idlC.Add(Entity(1), &stringyComp{Name: "a", Seq: 1})
	idlData, err := idlC.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	neg := *idlData
	neg.USetData.Data = append([]byte(nil), idlData.USetData.Data...)
	binary.LittleEndian.PutUint64(neg.USetData.Data[0:], ^uint64(0)) // -1
	if _, err := NewCSetFromData[stringyComp](&neg); err == nil {
		t.Fatal("negative IDL offset[0] must fail")
	}

	// 4. UnmarshalFrom 截断缓冲必须被 ValidateAsRoot 拦截
	if err := (&World{}).UnmarshalFrom(rockmem.NewReader([]byte{1, 2, 3})); err == nil {
		t.Fatal("truncated World root buffer must fail")
	}

	// 5. 读档后再插入较小 key：isKOrder 必须因恢复了 lastKey 而翻转为 false，防止 merge-join 漏行
	sa := NewCSet[p2CompA]()
	sb := NewCSet[p2CompB]()
	sa.Add(Entity(1), &p2CompA{1})
	sa.Add(Entity(5), &p2CompA{5})
	sa.Add(Entity(10), &p2CompA{10})
	sb.Add(Entity(1), &p2CompB{1})
	sb.Add(Entity(5), &p2CompB{5})
	sb.Add(Entity(10), &p2CompB{10})
	saData, _ := sa.Marshal()
	sbData, _ := sb.Marshal()
	sa2, err := NewCSetFromData[p2CompA](saData)
	if err != nil {
		t.Fatal(err)
	}
	sb2, err := NewCSetFromData[p2CompB](sbData)
	if err != nil {
		t.Fatal(err)
	}
	sa2.Add(Entity(3), &p2CompA{3}) // key=3 < max existing 10
	sb2.Add(Entity(3), &p2CompB{3})
	if sa2.IsKeyOrdered() || sb2.IsKeyOrdered() {
		t.Fatal("adding smaller key after unmarshal must clear isKOrder")
	}
	// 验证实际查询不漏掉 Entity(3)
	wRestore := NewWorld(WithoutGroups())
	e1 := wRestore.NewEntity(WithComponents(&p2CompA{1}, &p2CompB{1}))
	e2 := wRestore.NewEntity(WithComponents(&p2CompA{2})) // 占位 index=2
	e3 := wRestore.NewEntity(WithComponents(&p2CompA{3}, &p2CompB{3}))
	_ = e1
	_ = e2
	if err := wRestore.Update(); err != nil {
		t.Fatal(err)
	}
	wRestore.DestroyEntity(e2)
	if err := wRestore.Update(); err != nil {
		t.Fatal(err)
	}
	snapOrder, err := wRestore.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	wRe, err := NewWorldFromData(snapOrder, WithoutGroups())
	if err != nil {
		t.Fatal(err)
	}
	var qHits []EntityIndex
	if err := wRe.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error {
		q := ctx.NewQuery(WithComp[p2CompA](), WithComp[p2CompB]())
		for idx := range q.Iter() {
			qHits = append(qHits, idx)
		}
		return nil
	}), WithDeps(Dep[p2CompA](), Dep[p2CompB]())); err != nil {
		t.Fatal(err)
	}
	// 让 e3 先移除 B 再在后续帧插回（key=3 < max existing 4），触发不有序插入
	eRe := wRe.NewEntity(WithComponents(&p2CompA{20}, &p2CompB{20}))
	if info2, ok := wRe.GetEntityInfo(e3); ok {
		info2.Remove(&p2CompB{})
	}
	if err := wRe.Update(); err != nil {
		t.Fatal(err)
	}
	qHits = nil
	if info2, ok := wRe.GetEntityInfo(e3); ok {
		info2.Add(&p2CompB{30})
	}
	if err := wRe.Update(); err != nil {
		t.Fatal(err)
	}
	foundE3 := false
	for _, idx := range qHits {
		if idx == e3.Index() {
			foundE3 = true
		}
	}
	if !foundE3 {
		t.Fatalf("entity %d missed by query after unmarshal+re-add: hits=%v (eRe=%d)", e3.Index(), qHits, eRe.Index())
	}
}

// ============================================================================
// 4.4 组声明原子性、同组去重、含指针/零大小组件拦截
// ============================================================================

func TestGroupDeclare_AtomicityAndValidation(t *testing.T) {
	w := NewWorld()
	// 第 1 组有效，第 2 组含 disposable（非法） -> 整次注册失败且第 1 组不得残留
	err := w.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error { return nil }),
		WithDeps(Dep[p2CompA](), Dep[p2CompB](), Dep[groupCompD]()),
		WithGroup(Dep[p2CompA](), Dep[p2CompB]()),
		WithGroup(Dep[p2CompA](), Dep[groupCompD]()),
	)
	if err == nil {
		t.Fatal("expected registration error when second group is invalid")
	}
	if len(w.archetypes.declared) != 0 {
		t.Fatalf("partial group leaked into archetypes: %v", w.archetypes.declared)
	}

	// 同组重复类型 WithGroup(A, A) 必须拒绝
	if err := w.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error { return nil }),
		WithGroup(Dep[p2CompA](), Dep[p2CompA]()),
	); err == nil {
		t.Fatal("duplicate component type in single WithGroup must fail")
	}

	// 含 Go 指针组件进入 WithGroup 必须拒绝（防 GC 逃逸）
	if err := w.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error { return nil }),
		WithGroup(Dep[p2CompA](), Dep[p2PtrComp]()),
	); err == nil || !strings.Contains(err.Error(), "pointer") {
		t.Fatalf("component with Go pointers should be rejected from WithGroup, got %v", err)
	}

	// 零大小 Tag 组件进入 WithGroup 必须拒绝并给出明确报错
	if err := w.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error { return nil }),
		WithGroup(Dep[p2CompA](), Dep[p2ZeroTag]()),
	); err == nil || !strings.Contains(err.Error(), "zero-sized") {
		t.Fatalf("zero-sized component should be rejected with clear message, got %v", err)
	}
}

// ============================================================================
// 4.5 Reuse 代数校验 & EntityIDGenerator 防护
// ============================================================================

func TestEntityReuse_ABAGuards(t *testing.T) {
	w := NewWorld()
	proj1 := w.NewEntity(WithComponents(&p2CompA{111}))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}

	// 销毁 proj1 + 补满 9 个销毁触发 delayFlush 使 index=1 被复用
	w.DestroyEntity(proj1)
	var fillers []Entity
	for i := 0; i < 9; i++ {
		fillers = append(fillers, w.NewEntity())
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	for _, f := range fillers {
		w.DestroyEntity(f)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}

	proj2 := w.NewEntity(WithComponents(&p2CompA{222}))
	if proj2.Index() != proj1.Index() || proj2 == proj1 {
		t.Fatalf("expected index reuse: proj1=%v proj2=%v", proj1, proj2)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}

	// 旧句柄全面拒绝
	if w.entities.Exist(proj1) {
		t.Fatal("EntitySet.Exist must reject stale handle")
	}
	if w.entities.Remove(proj1) != nil {
		t.Fatal("EntitySet.Remove must not remove recycled entity")
	}
	if _, ok := w.GetComponent[p2CompA, *p2CompA](proj1); ok {
		t.Fatal("World.GetComponent must reject stale handle")
	}
	// OpLog 删除旧句柄不得误删 proj2 组件
	w.componentOp(Operate{Entity: proj1, Op: ComponentOperateDelete, Comp: &p2CompA{}})
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if got, ok := w.GetComponent[p2CompA, *p2CompA](proj2); !ok || got.V != 222 {
		t.Fatalf("proj2 component damaged by stale OpLog delete: %+v ok=%v", got, ok)
	}

	// EntityIDGenerator: delayCap=0 兜底；释放未分配槽位 panic
	g := NewEntityIDGenerator(4, 0)
	id := g.NewID()
	g.FreeID(id)
	defer func() {
		if recover() == nil {
			t.Fatal("freeing unallocated slot must panic")
		}
	}()
	g.FreeID(Entity(3)) // index=3 从未分配
}

// ============================================================================
// 4.6 容器安全边界
// ============================================================================

func TestContainers_SafetyBounds(t *testing.T) {
	// SparseArray.Add(-1) 不 panic
	sa := NewSparseArray[EntityIndex, p2CompA]()
	if sa.Add(-1, &p2CompA{1}) != nil {
		t.Fatal("negative key Add should return nil")
	}

	// SparseArray.Remove 清零尾槽指针释放 GC 引用
	sap := NewSparseArray[EntityIndex, p2PtrComp]()
	sap.Add(1, &p2PtrComp{Label: "one", V: 1})
	sap.Add(2, &p2PtrComp{Label: "two", V: 2})
	sap.Remove(1)
	if sap.data[1].Label != "" {
		t.Fatalf("vacated tail slot not zeroed: %+v", sap.data[1])
	}

	// OrderedIntSet.Merge 空集合不共享切片别名
	var a OrderedIntSet[ComponentIntType]
	b := OrderedIntSet[ComponentIntType]{10, 20}
	a.Merge(b)
	a.Add(15)
	if len(b) != 2 || b[1] != 20 {
		t.Fatalf("Merge aliased caller slice: b=%v", b)
	}

	// FixedString 支持超过 1024 字节不越界
	var fs FixedString[[2048]byte]
	long := strings.Repeat("x", 1500)
	fs.Set(long)
	if got := fs.String(); got != long {
		t.Fatalf("FixedString >1024B mismatch: len=%d", len(got))
	}

	// CSet.Add typed-nil 不 panic
	c := NewCSet[p2CompA]()
	var nilPtr *p2CompA
	if c.Add(Entity(1), nilPtr) != nil {
		t.Fatal("typed-nil Add should return nil")
	}
}

// ============================================================================
// ECS-1 & ECS-4 确定性回放、Step(delta)、只读 Marshal、Destroy/Unmarshal 契约
// ============================================================================

func TestDeterministicReplay_StepAndReadOnlyMarshal(t *testing.T) {
	runReplay := func(async bool) []byte {
		opts := []WorldOption{}
		if async {
			opts = append(opts, WithWorldASyncMode())
		}
		w := NewWorld(opts...)
		w.RegisterLight(LightSystem(func(ctx *SystemContext, event Event) error {
			for a, bv := range View2RO[p2CompA, p2CompBView](ctx) {
				a.V += int32(event.Delta.Milliseconds()) + bv.V()
			}
			return nil
		}),
			WithDeps(Dep[p2CompA](), Dep[p2CompB](ReadOnly)),
			WithGroup(Dep[p2CompA](), Dep[p2CompB]()),
		)
		// 两个读者可并行，分别向私有命令队列投递同一组件。首个注册者应获胜。
		var firstFrameSubmissions [2]int
		for i, offset := range []int32{100, 200} {
			err := w.RegisterLight(func(ctx *SystemContext, event Event) error {
				for idx := range ctx.GetComponentsReadOnly[p2CompBView]() {
					info := ctx.world.entities.getByIndex(idx)
					ctx.AddComponents(info.Entity(), &p2CompTag{V: int32(idx) + offset})
					if event.Frame == 1 {
						firstFrameSubmissions[i]++
					}
				}
				return nil
			}, WithDepReadOnly[p2CompB]())
			if err != nil {
				t.Fatal(err)
			}
		}

		var ents []Entity
		for i := 0; i < 20; i++ {
			ents = append(ents, w.NewEntity(WithComponents(&p2CompA{int32(i)}, &p2CompB{int32(i * 2)})))
		}
		for frame := 0; frame < 5; frame++ {
			if err := w.Step(16 * time.Millisecond); err != nil {
				t.Fatal(err)
			}
		}
		if firstFrameSubmissions != [2]int{20, 20} {
			t.Fatalf("expected competing tag submissions, got %v", firstFrameSubmissions)
		}
		for _, e := range ents {
			tag, ok := w.GetComponent[p2CompTag](e)
			if !ok || tag.V != int32(e.Index())+100 {
				t.Fatalf("entity %v: first registered tagger should win, got %v", e, tag)
			}
		}
		// 销毁部分实体并加入新实体（组表 churn）
		w.DestroyEntity(ents[2])
		w.DestroyEntity(ents[7])
		w.NewEntity(WithComponents(&p2CompA{99}, &p2CompB{198}))
		if err := w.Step(16 * time.Millisecond); err != nil {
			t.Fatal(err)
		}

		// 验证只读 Marshal 不扰动组表行序
		a, _ := w.archetypes.owner(GetIntType[p2CompA, *p2CompA]())
		rowsBefore := append([]EntityIndex(nil), a.entities...)

		buf := rockmem.NewWriter(4096)
		if _, err := w.MarshalTo(buf); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(rowsBefore, a.entities) {
			t.Fatalf("read-only Marshal mutated group table row order: before=%v after=%v", rowsBefore, a.entities)
		}
		return append([]byte(nil), buf.Bytes()...)
	}

	b1 := runReplay(false)
	b2 := runReplay(false)
	b3 := runReplay(true)
	if !bytes.Equal(b1, b2) || !bytes.Equal(b1, b3) {
		t.Fatalf("deterministic replay bytes diverged: len(b1)=%d len(b2)=%d len(b3)=%d", len(b1), len(b2), len(b3))
	}
}

func TestWorldContract_PendingMarshalErrorAndStopUnmarshalRejected(t *testing.T) {
	// Marshal 前有未 flush 操作必须显式报错
	w := NewWorld()
	w.NewEntity(WithComponents(&p2CompA{1}))
	if _, err := w.Marshal(); err == nil {
		t.Fatal("Marshal with pending opLog must fail")
	}

	// Destroy 后释放资源，且不允许对 WorldStatusStop 再 Unmarshal
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	snap, err := w.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := w.Destroy(); err != nil {
		t.Fatal(err)
	}
	if w.components != nil {
		t.Fatal("Destroy should release component map")
	}
	if err := w.Unmarshal(snap); err == nil {
		t.Fatal("Unmarshal into WorldStatusStop must fail")
	}
}
