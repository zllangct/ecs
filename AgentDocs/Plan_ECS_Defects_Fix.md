# ECS 缺陷修复实施计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 按 `AgentDocs/ECS_Defects.md` 讨论记录中已确认的方案，修复 ECS 框架缺陷，达到生产可用标准。

**Architecture:** 按依赖顺序分 8 个阶段：容器层 → 实体 ID → 操作日志 → 调度器 → 只读依赖 → 生命周期 → 序列化 → 杂项。每阶段 TDD：先写失败测试，再实现。

**Tech Stack:** Go 1.22+（range-over-int、iter），rockmem IDL 序列化。

**核心约束：保证原有功能语义不变**——所有现有公开 API 签名、默认行为、stage 数量、isFriend 命名、iface hack、全局注册表、15 Stage 设计均保持；修复只针对错误行为与缺失能力。

**验证命令（每任务执行）：**
```bash
cd i:/rgwork/ecs && go build ./... && go test ./... -count=1
```

---

## Phase 1: SparseArray/USet 容器层（D03/D04/D05/D10/D18-部分）

### Task 1.1: 修复 USet.Remove 上界检查 + yield 返回值

**Files:**
- Modify: `uset.go`（Remove 增加 `idx >= u.len` 检查；Iter 检查 yield 返回值）
- Test: `uset_test.go`（新建）

**Step 1: 写失败测试** — `TestUSetRemove_OutOfRange`（idx>=len 返回 nil 不 panic）、`TestUSetIter_EarlyBreak`（break 不 panic）。
**Step 2:** `go test -run TestUSet -v` 确认失败/panic。
**Step 3:** 实现：`Remove` 加 `if idx < 0 || idx >= u.len { return nil, 0, 0 }`；`Iter` 改 `if !yield(i, &u.data[i]) { return }`。
**Step 4:** 测试通过且全仓测试绿。

### Task 1.2: 修复 SparseArray.Remove 的 swap-remove 索引同步（D03）

**Files:**
- Modify: `sparse_array.go:83-105`（Remove）
- Test: `sparse_array_test.go`（扩充）

**Step 1: 写失败测试**
- `TestSparseArrayRemove_SwapRemoveConsistency`：Add 3 个元素(key 1,5,9)，Remove(5) 后：Get(5)==nil、Get(1)/Get(9) 仍正确、Len()==2、EntityIndexes 对应 idx2Key 长度为 2 且无残留 key 5、Iter 不产生 key 5。
- `TestSparseArrayRemove_LastElement`：删除末尾元素后索引一致。

**Step 3: 实现**（标准 swap-remove 同步）：
```go
func (s *SparseArray[K, V]) Remove(key K) *V {
	if key < 0 || key >= K(len(s.indices)) {
		return nil
	}
	idx := s.indices[key] - 1
	if idx < 0 {
		return nil
	}
	lastIdx := int32(s.len) - 1
	lastKey := K(s.idx2Key[lastIdx])
	removed := s.data[idx]
	s.data[idx] = s.data[lastIdx]     // USet 层 swap-remove
	s.len--
	s.idx2Key[idx] = int32(lastKey)   // 末尾 key 移到被删位置
	s.indices[lastKey] = idx + 1      // 同步稀疏索引
	s.indices[key] = 0
	s.idx2Key = s.idx2Key[:lastIdx]   // 截断末尾
	s.shrink(key)
	return &removed
}
```
注意：USet.Remove 不再被 SparseArray 调用（改由本方法内联 swap-remove 语义），USet.Remove 保留供其他调用方。保持对外语义：返回被删元素拷贝指针。

### Task 1.3: 修复 shrink 与边界检查（D05）

**Step 1: 写失败测试** — `TestSparseArrayGet_AfterShrinkNoPanic`：添加 1500 个 key 触发扩容，删除至仅剩小 key，再 Get(中间已删的大 key) 不 panic 返回 nil。
**Step 3: 实现：**
- `Get`/`Exist`/`Remove` 边界判断统一改为 `key < 0 || key >= K(len(s.indices))`（不再信任 maxKey）；
- shrink 不再把 maxKey 抬到 shrinkThreshold（删除该段）；shrink 末尾补 `s.indices = newIndices` 写回。

### Task 1.4: 修复 Sort/Swap（D04）

**Step 1: 写失败测试** — `TestSparseArraySort_KeyOrdered`：乱序 Add 若干 key，Sort 后 Iter 按 key 升序产出，且每个 key 的 Get 仍返回原值。
**Step 3: 实现：** Sort 改为置换环算法：
```go
func (s *SparseArray[K, V]) Sort() {
	if s.isKOrder {
		return
	}
	n := int(s.len)
	// 目标：idx2Key 升序。对每个位置 i，把应属该位置的元素换过来
	for i := 0; i < n; i++ {
		for int(s.idx2Key[i]) != i /* 占位逻辑见下 */ {}
	}
	...
}
```
实际采用：先构造 `order []int`（位置按 key 升序的排列，sort.Slice），再按置换环 O(n) 应用 data 与 idx2Key，最后遍历密集区重建 `indices[idx2Key[pos]] = pos+1`。Swap 修正为 swap 后用 idx2Key 回写 indices。isKOrder 语义保持（Add/Remove 置 false，Sort 置 true）。

### Task 1.5: SparseArray.Iter/IterReadOnly yield 检查（D10）

随 1.2 一并实现：`if !yield(...) { return }`。测试 `TestSparseArrayIter_EarlyBreak`。

### Task 1.6: CSet.RemoveAndReturn nil 检查（D18-部分）

**Step 1:** `TestCSetRemoveAndReturn_NotExist` 返回 nil 不 panic。
**Step 3:** `if r := c.remove(entity); r != nil { cpy := *r; return any(&cpy).(Component) }; return nil`。

---

## Phase 2: EntityIDGenerator 与 ReuseID（D11/D12）

### Task 2.1: delayFlush 重写为有序归并

**Files:** Modify `entity.go:97-129`；Test `entity_test.go`（新建）

**Step 1: 写失败测试** — `TestEntityIDGenerator_ReuseSmallestFirst`：NewID×N，FreeID 其中若干个（凑满 delayCap 触发 flush），后续 NewID 按升序拿到最小可用 index，且 reuse 代数递增。
**Step 3: 实现**（双有序链归并，沿用 ids 内嵌 freelist，不引入新结构）：
```go
func (e *EntityIDGenerator) delayFlush() {
	sort.Slice(e.removeDelay[:e.delayFree], func(i, j int) bool {
		return e.removeDelay[i].index < e.removeDelay[j].index
	})
	// 归并：freelist 链（free 游标）与 removeDelay 升序数组 → 新链写回 ids[x].index
	dummy := EntityIndex(-1)          // 虚拟头
	tail := &dummy
	cur := e.free                     // freelist 链头（升序不变量）
	hasCur := e.free < e.pending
	for i := int32(0); i < e.delayFree; i++ {
		idx := e.removeDelay[i].index
		for hasCur && cur < idx {
			*tail = cur
			tail = &e.ids[cur].index
			next := e.ids[cur].index
			hasCur = cur < e.pending && next > 0 /* 链有效 */ 
			...
		}
		*tail = idx
		tail = &e.ids[idx].index
	}
	// 接剩余 freelist 链
	...
	e.free = dummy
	e.delayFree = 0
}
```
实现要点：保持"freelist 严格升序"不变量；freelist 链节点 x 满足 `ids[x].index` 为下一空闲 index，链尾以 >=pending 或哨兵标识。需先梳理现有 NewID/FreeID 的链表示意（`index=-1` 表空闲标记、free/pending 语义），测试驱动逐步校准。

### Task 2.2: 删除收缩死代码 + FreeID 合法性检查

- 删除 `entity.go:92-94`；
- FreeID 增加：index 越界检查、reuse 不匹配检查（`e.ids[idx].reuse != realID.reuse`）、double-free 检查（当前未占用即拒绝）。
- 测试：`TestFreeID_DoubleFree`、`TestFreeID_StaleReuse`。

### Task 2.3: reuse 代数校验接入查询（D12）

- `EntitySet.Get`：`info == nil || info.entity != entity → nil,false`；
- `CSet.Get`/`get`：增加 entity 完整值校验（需要 CSet 存 entity 代？——CSet 以 index 为 key，reuse 校验依赖 EntitySet 层即可：GetComponents/GetBuddy 路径上由调用方先经 EntitySet 校验，或 CSet 增加 reuse 数组）。**实现时以最小侵入确定**：优先在 `world.getEntityInfo`/`EntitySet.Get` 校验；组件层 `GetBuddy(ctx, index)` 走的是 index，由 query 产出的 index 天然有效。
- 测试：`TestEntitySetGet_StaleEntity`。

---

## Phase 3: OpLog per-system 队列（D01/D02）

### Task 3.1: per-system 操作队列

**Files:** Modify `op_log.go`、`system.go`（SystemInfoInstance 增加 opQueue）、`entity_info.go`（componentOp 路由）

**设计：**
- 新增 `opQueue`（每 system 一个 opTaskList）；SystemContext 持有其引用；
- `EntityInfo.Add` 当前无 ctx——**保持语义不变的最小方案**：world 保留一个全局默认队列（主线程），system 内操作经 ctx 路由。由于 EntityInfo.Add 不持有 ctx，方案调整为：OpLog 改为"按写入方注册队列"：world.newEntity/主线程 → 默认队列；system 执行期间 → flow 在调用 system 前将 world 的 currentQueue 指向该系统队列（串行执行语义下安全；并行模式下每个 goroutine 执行前切换 world.currentQueue 有竞争——因此并行模式必须在 ctx 层路由）。
- **最终方案**：`componentOp` 增加 ctx 参数重载：`EntityInfo.Add` 走 world 默认队列；ecs.go 的 GetComponents/GetBuddy 不变；system 内 Add 组件的入口保持 EntityInfo.Add（world 默认队列）→ 由于组件 flush 都在帧同步点主线程执行，**所有队列都是单写多主线程读**，锁可整体移除。
- 验证：并行模式下 EntityInfo.Add 若可从多 goroutine 调用，则默认队列仍需一把锁——保留单个轻量 Mutex（比 32 桶 RWMutex 更简单正确），主线程 flush 时取锁换出队列。
- 测试：`TestOpLog_ConcurrentOperate`、`TestOpLog_FlushOrder`。

---

## Phase 4: 调度器（D06/D07/D13/D14）

### Task 4.1: D06 delta
- `Update` 末尾 `w.lastUpdate = now`；AutoOptimize elapsed 改为 `time.Since(beforeExecute)`。
- 测试：`TestWorldUpdate_DeltaIsFrameInterval`。

### Task 4.2: D07 Stage 覆盖
- `StageMaxIndex` 保持常量值，4 处 `range StageMaxIndex` 改为 `range StageMaxIndex+1`（Go 1.22 range-over-int 对常量表达式同样适用）。
- 测试：`TestSyncAfterDestroy_Called`（注册实现 SyncAfterPostDestroyReceiver 的系统，Update 后被调用）。

### Task 4.3: D13 register 插入
- 改 `sl[:i] + sg + sl[i:]`，显式构造新切片。
- 测试：`TestFlowRegister_CustomOrder`。

### Task 4.4: D14 nil 方法值
- 15 个分支统一改 `if system, ok := sys.(XReceiver); ok { fn = system.X; imp = true }`。
- 测试：现有 ecs_test/example 回归。

---

## Phase 5: 只读依赖语义（D08/D09）

### Task 5.1: 语义扶正
- `Dep[T,TP]()`: 默认 `ReadWrite`（readonly()=false）；`Dep(writable...)` 传 ReadOnly 时 readonly()=true；`WithDepReadOnly` 传 ReadOnly；`NewItDependency` 默认对齐（默认 ReadWrite）。**注意语义翻转对现有代码的影响**：ecs.go 的 `dep.readonly()` 分支、system_related_groups.go 的 isFriend 分支随语义翻转一起核对。
- 测试：`TestDep_DefaultWritable`、`TestDepReadOnly_IsReadonly`、`TestIsFriend_TwoReadonlyNoConflict`。

### Task 5.2: GetBuddy 只读（D09）
- R1 视图落地前，GetBuddy 只读分支暂返回原指针但**语义上标注弃用**；TView 由 codegen 扩展实现（涉及 rockmem-ecs generator，属于较大工程，单独立项：`ecs-rockmem-extension` skill）。
- 本期：修正 `&*(*T)(b)` 空操作为直接返回 `(*T)(b)`，删除无意义表达式；readonly 拷贝防护随 R1 在 codegen 阶段落地。
- 测试：现有测试回归。

---

## Phase 6: 生命周期（D23/D24/D25）

### Task 6.1: DestroyEntity
- `world.DestroyEntity(entity Entity)`：产生 OpLog 删除操作（新 op 类型或复用 DeleteAll 语义扩展）；帧同步点：从所有组件集移除该实体组件、EntitySet.Remove、idGenerator.FreeID、compound 清空。
- 测试：`TestDestroyEntity_RecyclesIDAndComponents`。

### Task 6.2: EntityInfo.Remove(comp)
- 暴露组件移除入口 → ComponentOperateDelete。
- 测试：`TestEntityInfoRemove_Component`。

### Task 6.3: World.Destroy + 状态机
- Destroy：遍历系统触发 Destroy stage 链（SyncBeforeDestroy→Destroy→SyncAfterDestroy）、释放资源、置 WorldStatusStop；Update 在首帧置 Running。
- 测试：`TestWorldDestroy_TriggersDestroyStages`。

### Task 6.4: World 接口补全
- 接口增加 `DestroyEntity(entity Entity)`；序列化 Unmarshal 入口暴露（`UnmarshalFrom(reader)`）。
- clearNomadic 接入帧末流程或删除（按语义：nomadic 组件每帧清空→接入 Execute 尾部，与 clearDisposable 并列）。

---

## Phase 7: 序列化 S1（D29）

### Task 7.1: CSet 序列化改逐元素 rockmem IDL
- `CSet.Marshal`：不再 raw 字节直拷 data；改为逐元素调用组件生成的 Marshal（组件是 rockmem 生成类型，实现 rockmem 序列化接口）；SerializableSparseArrayData 的 Data 字段改为 `[][]byte` 或新增 per-element 条目结构（需同步扩展 ecs.rm IDL 与 ecs_generated.go）。
- Unmarshal 对称实现；注册表的 Unmarshaler 适配。
- 测试：`TestCSetMarshal_NonPOD`（含 string/slice 字段组件快照恢复后数据正确）、现有 serialization_test 回归。

---

## Phase 8: 杂项（D15/D16/D17/D18-余/D22/D34/Query key 缓存）

### Task 8.1: D15 FixedString clamp
### Task 8.2: D16 NewQuery 组件集缺失→空查询
### Task 8.3: D17 删除 getTaskExecutor 死代码
### Task 8.4: D18 clearDisposable/collect nil 检查
### Task 8.5: D22 SystemConstraint 命名扶正
### Task 8.6: D34 LocalUniqueID 安全布局
### Task 8.7: Query 的 FixedCompound key 构造期缓存（NewQuery 时算一次）

每项 TDD：失败测试 → 实现 → 回归。

---

## 不在本计划内
- D32 optimizer 重做（独立设计文档）
- R1 TView codegen 扩展（使用 ecs-rockmem-extension skill 单独立项，Phase 5 仅完成语义扶正）
- 明确保留项：D19/D27 静默、D20 命名、D21 不重构、D26 全局单例、D28 15 Stage、D30 iface hack、D31 getDep 遍历、D33 逐系统 goroutine
