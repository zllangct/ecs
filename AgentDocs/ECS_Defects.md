# ECS 框架设计缺陷清单

> 生成日期：2026-09-14
> 状态说明：⬜ 待讨论 / 🔵 讨论中 / ✅ 方案已定 / 🛠 已修复
> 注：逐项讨论结论以文末「讨论记录」表为准

---

## 一、严重正确性 Bug（导致 panic 或错误结果）

### D01 ⬜ OpLog 分桶哈希失效
- 位置：`op_log.go:60-63`
- 问题：`hash = int64((uintptr)(unsafe.Pointer(&hash))) & c.bucket` 用局部变量自身栈地址做哈希，同一 goroutine 内哈希恒定，分桶锁形同虚设。
- 原始需求假设：OpLog 分桶是为了让并行模式下多个 system 并发提交组件操作时减少锁竞争。

### D02 ⬜ OpLog 读锁内执行写操作
- 位置：`op_log.go:87-101`
- 问题：`getOpTasks` 持 `RLock` 却调用 `list.Reset()`，与 `operate` 并发时 data race。

### D03 ⬜ SparseArray.Remove 的 idx2Key 维护错误
- 位置：`sparse_array.go:90-98`
- 问题：`s.idx2Key = s.idx2Key[:len(s.idx2Key)]` 长度不变（应为 len-1）；第 93 行写入与第 96 行 swap 互相抵消。`EntityIndexes()` 返回已删除元素，索引表与数据不一致。

### D04 ⬜ SparseArray.Swap/Sort 用位置做 indices 下标
- 位置：`sparse_array.go:146-172`
- 问题：`s.indices[Key], s.indices[i] = ...` 中 `i` 是数据位置而非 key，Sort 后索引表损坏；Sort 遍历整个 indices（O(maxKey)）且比较语义混乱（位置 vs key）。optimizer 每帧调用 Sort 等于主动破坏数据。

### D05 ⬜ SparseArray.shrink 越界风险 + 收缩未生效
- 位置：`sparse_array.go:176-198`
- 问题：收缩后 `maxKey` 被抬到 `shrinkThreshold`(1024)，后续 `Get(key)` 通过 maxKey 检查但 `indices[key]` 越界 panic；末尾 `newIndices` 分配后未写回 `s.indices`（死代码）。

### D06 ⬜ World.Update 的 delta 永不更新
- 位置：`world.go:177-204`
- 问题：从未执行 `w.lastUpdate = now`，delta 单调累积；AutoOptimize 分支的 elapsed 同样错误。

### D07 ⬜ range StageMaxIndex 丢失最后一个 Stage
- 位置：`system_flow.go:130, 388, 411, 514`
- 问题：`for stage := range StageMaxIndex` 只产生 0..n-1，`StageSyncAfterDestroy` 永远无容器、不注册、不执行。

### D08 ⬜ 只读依赖语义颠倒
- 位置：`system_dependency.go:3-52`、`system.go:64-74`
- 问题：`Dep[T]()` 默认（应为可写）算出 `readonly()==true`；`WithDepReadOnly` 反而可写；`NewItDependency` 与 `Dep` 默认语义互相矛盾。同时影响 GetComponents 迭代方式与并行分组正确性。

### D09 ⬜ GetBuddy 只读保护是空操作
- 位置：`ecs.go:94-96`
- 问题：`return &*(*T)(b)` 仍返回原指针，调用方可写"只读"组件；与 GetComponents 的 IterReadOnly（真拷贝）语义不一致。

### D10 ⬜ 迭代器不检查 yield 返回值
- 位置：`query.go:59`、`sparse_array.go:203,212`、`uset.go:135`、`system_traverser.go:48`
- 问题：Go 1.23+ range-over-func 中 yield 返回 false 后继续迭代会 panic，用户 break 查询循环即崩溃。

### D11 ⬜ EntityIDGenerator 多处缺陷
- 位置：`entity.go`
- 问题：(a) L92-94 收缩条件 `pending < len(ids)/2` 永不成立（死代码），若触发会截掉存活槽位；(b) 无 double-free 检测；(c) 反序列化 delayCap=0 时 FreeID 越界 panic；(d) delayFlush 链表归并算法无注释无测试，正确性不可审查。

### D12 ⬜ ReuseID 机制名存实亡
- 位置：`entity_set.go:18`、`component_set.go:60`
- 问题：查询只用 `entity.Index()`，不比较 reuse 代数，旧 Entity 句柄在 index 复用后仍指向新实体（ABA 悬空引用）。

### D13 ⬜ flow.register 插入位置 off-by-one
- 位置：`system_flow.go:431-432`
- 问题：应为 `sl[:i] + sg + sl[i:]`，实际用 `i-1`，自定义 Order 排序错误；`append(sl[:i-1], sg)` 写脏原底层数组。

### D14 ⬜ getSystemTask 对 nil 接口取方法值
- 位置：`system_flow.go:217-309`
- 问题：`fn = system.X` 在 `imp = ok` 之前执行，断言失败时 nil interface 取方法值直接 panic，目前仅靠 register 的隐式约定维持。

### D15 ⬜ FixedString.Set 越界
- 位置：`fixed_string.go:33-37`
- 问题：超长输入 copy 截断但 `f.len` 记录原长度，`String()` 切片越界 panic。

### D16 ⬜ NewQuery 静默忽略不存在的组件集
- 位置：`ecs.go:27-37`
- 问题：查询条件中组件集不存在时应为零匹配，实际丢弃该条件返回超集。

### D17 ⬜ getTaskExecutor nil 解引用
- 位置：`world.go:148-153`
- 问题：nil 检查后空 if 分支，随后直接调用 nil factory panic。

### D18 ⬜ 其他零散 panic 点
- `component_set.go:56` RemoveAndReturn 无 nil 检查；
- `world.go:258` clearDisposable 无 nil 检查；
- `optimizer.go:60-66` collect 中 make(len)+append 产生双倍长度的零值条目。

---

## 二、并发设计缺陷

### D19 ⬜ 并行安全完全依赖用户自觉
- 无运行时检测手段，漏声明依赖即 silent data race；配套 OpLog 锁又因 D01/D02 自身损坏。

### D20 ⬜ isFriend 命名与语义相反
- 位置：`system_related_groups.go:15`
- 实际判断的是"冲突（不可并行）"，attach/resort 整段逻辑读起来全反。

### D21 ⬜ stage→接口映射重复 4 份
- 位置：`system_event.go:74`、`system_flow.go:217,444,514`
- 新增 Stage 需同步改 4 处（shotgun surgery），应表驱动化。

---

## 三、API 与语义设计缺陷

### D22 ⬜ SystemConstraint 命名与真值相反
- 位置：`system.go:94-108`
- `isValid()` 返回 `outdated` 字段，`reset()`=有效、`setOutdated()`=无效，极易误用。

### D23 ⬜ World 状态机形同虚设
- Running/Stop 从不设置；register 的"only in world init"检查无实际约束窗口；`Destroy()` 空实现，不清理系统不释放 arena。

### D24 ⬜ 生命周期功能缺失
- 无 DestroyEntity、无 EntityInfo.Remove(comp)；`clearNomadic` 写了但从未被调用；todo.go 自认半成品。

### D25 ⬜ World 接口与能力不匹配
- 接口有 Marshal 无 Unmarshal；AddNomadic 等不在接口上；NewWorldFromData 返回后无法再反序列化。

### D26 ⬜ 全局单例 Registry + panic
- 位置：`component_registry.go:29,46`
- 进程级共享，多 World 无独立命名空间；冲突直接 panic，库应返回 error。

### D27 ⬜ 静默失败文化
- GetComponents/GetBuddy 在依赖未声明、组件集不存在、断言失败时均静默返回空，无日志无错误，排错成本高。

### D28 ⬜ 15 个 Stage 接口的 API 爆炸
- Sync/Before/After 三轴笛卡尔积，多数 runSync=true 区分度存疑，复杂度收益不成比例。

---

## 四、序列化路线矛盾

### D29 ⬜ unsafe raw memcpy 序列化架空 RockMem IDL
- 位置：`serialization.go:14-30`
- 含 string/slice/map/指针字段的组件序列化出悬空指针；无大小端/位数/eleSize/版本校验。todo.go 已自问"是否需要统一使用 karmem"。

### D30 ⬜ internal_type_mock.go 绑定 Go runtime 内存布局
- 位置：`internal_type_mock.go`、`component_set.go:36`
- 复制 iface/itab/_type 结构，从接口抠 data 指针，Go 版本升级即可能崩溃。

---

## 五、性能与工程质量

### D31 ⬜ 热路径不必要开销
- Query.Iter 每次调用都 NewFixedCompound 构造 map key（query.go:51）；getDep O(n) 线性查找（system.go:151）；IterReadOnly 按值拷贝整个组件（sparse_array.go:211）。

### D32 ⬜ optimizer 实质无优化且破坏数据
- collect() 统计结果无消费者；memTidy 只做（错误的）Sort；每帧 fmt.Printf 打 3 行日志（optimizer.go:75,84,91）。

### D33 ⬜ 并行执行无 worker pool
- 每系统一个 goroutine；flushTempTask 每帧两次全量合并克隆。

### D34 ⬜ LocalUniqueID 时间位溢出
- 位置：`utils.go:19-35`
- 纳秒时间戳 <<32 后高位截断，时间回绕周期约 4.29 秒，叠加 16 位随机，碰撞风险高。

---

## 讨论记录

（逐项讨论后在此记录：用户原始需求 → 确定方案 → 决策理由）

| 编号 | 原始需求 | 确定方案 | 状态 |
|------|----------|----------|------|
| D35（实施中新发现） | NewEntity 返回 *EntityInfo 指向 EntitySet 密集数组元素，数组扩容后指针失效（悬空），持有方读到旧数据 | 用户决策：返回 Entity 值+按需查询。NewEntity/EntityTemplate.Instance(N) 返回 Entity；World 新增 GetEntityInfo（注明指针短生命周期、勿跨 NewEntity 持有）；全量调用点与 README 同步 | 🛠 已修复 2026-09-15 |
| 全局定位 | 框架定位：生产可用，rockgo 真实业务依赖；正确性、并发安全、序列化可靠性均需达标 | 所有缺陷按生产标准修复，不接受"砍掉功能"式简化 | ✅ 已闭环 2026-09-14 |
| D03/D04/D05（SparseArray） | ① USet 是密集连续容器，用 swap-remove 删除（避免内存平移）；② SparseArray 是稀疏索引层：EntityID 作下标直达内存偏移，需在 USet swap-remove 时同步维护 indices/idx2Key；③ 元素顺序无要求，但 Sort 的目的是在 Tick 空闲时间重排密集数组，使不同 SparseArray 按相同 key 序排列、联合遍历 CPU cache 友好 | 方案A'：① 修正 swap-remove 时 indices/idx2Key 同步（去掉互相抵消的双 swap，idx2Key 截断 len-1）；② Sort 保留并修对——按 idx2Key 升序置换密集数组（置换环 O(n)，一临时元素），完成后由 idx2Key 全量重建 indices；Swap 改为 swap 后用 idx2Key 回写 indices；isKOrder 作脏标记保留；③ shrink：Get/Exist 边界改判 key>=len(indices)，maxKey 不再抬到 threshold，补 newIndices 写回，USet.Remove 加 idx>=len 上界检查；④ 补单元测试（增删交叉、swap-remove 一致性、Sort 后索引正确性、shrink 边界） | ✅ 已闭环 2026-09-14 |
| D29（序列化路线） | World 序列化用于同二进制、同机器的快照/恢复（热重载、崩溃恢复），不跨版本不跨机器；但非 POD 组件 raw 直拷仍会恢复出悬空指针 | 方案 S1：统一逐元素走 RockMem IDL 生成的 Marshal/Unmarshal，string/slice/嵌套全正确；快照低频，开销可接受 | ✅ 已闭环 2026-09-14 |
| D30（iface hack） | 用户决策：保留零断言开销，接受 Go 版本绑定风险 | 维持现状，补风险注释 | ✅ 已闭环 2026-09-14 |
| D32（optimizer） | 用户决策：重新定义优化器职责（含内存收缩、shape 分组建议等） | 不在本轮修复范围，单独出 optimizer 设计文档 | 🔵 待独立设计 |
| D33（worker pool） | 系统数量不大时 goroutine 开销可接受 | 保持逐系统 goroutine，后续压测再定 | ✅ 已闭环 2026-09-14 |
| D34/D31（ID 与热路径） | getDep 遍历小数组优于 map（dep 数量天然小），不改；Query 热路径开放更优方案 | LocalUniqueID 修为纳秒高位+原子序号安全布局；Query 的 FixedCompound key 在 NewQuery 时计算一次缓存（而非每次 Iter 构造）；其余热路径项暂缓 | ✅ 已闭环 2026-09-14 |
| D26（注册表作用域） | 全局单例够用（rockgo 单 world 场景） | 维持现状：全局单例+冲突 panic，仅补文档说明 | ✅ 已闭环 2026-09-14 |
| D28（Stage 数量） | 15 Stage 设计保留 | 不精简，只修 D07 | ✅ 已闭环 2026-09-14 |
| D23/D24/D25（生命周期与 World 接口） | 生产需要完整生命周期 | 新增 DestroyEntity（延迟到帧同步点：回收 ID+清理所有组件+compound）、EntityInfo.Remove(comp)、World.Destroy() 触发 Destroy stage 链+释放资源；状态机补 Running/Stop；World 接口补 Unmarshal/AddNomadic 等能力对齐 | ✅ 已闭环 2026-09-14 |
| D15/D16/D17/D18/D22（小 Bug 批次） | 需求自明 | D15：Set 的 len clamp 到实际拷贝长度；D16：组件集不存在→零匹配；D17：删除 getTaskExecutor 死代码及 TaskExecutor 类型；D18：RemoveAndReturn/clearDisposable/collect 补 nil 与边界检查；D22：SystemConstraint 字段与命名扶正（active/activate/deactivate） | ✅ 已闭环 2026-09-14 |
| D19/D27（越权访问可观测性） | 用户决策：未声明依赖返回空是约定行为 | 保持静默，文档补说明 | ✅ 已闭环 2026-09-14 |
| D20（isFriend 命名） | 形象命名：共同依赖的系统=好友=同组串行执行 | 不改名，补注释说明语义 | ✅ 已闭环 2026-09-14 |
| D21（stage 映射重复） | 暂不接受表驱动重构 | 只修 D14 panic，重复映射保留 | ✅ 已闭环 2026-09-14 |
| D01/D02（OpLog 并发） | OpLog 攒批+帧同步点 flush 的无锁化设计保留；分桶意图是让并行写入方近似无锁提交 | 方案 C：per-system 操作队列（注册时分配，挂 SystemContext），单写者无锁；主线程操作走 world 默认队列；flush 在 wg.Wait() 后由主线程顺序合并；顺带消除 D01 伪哈希与 D02 读锁内写 | ✅ 已闭环 2026-09-14 |
| D06/D07/D10/D13/D14（调度层无争议 Bug） | 需求自代码意图自明：delta=本帧间隔；15 个 Stage 全部可执行；迭代器支持 break；自定义 Order 按序插入；断言失败应跳过而非 panic | D06：Update 末尾更新 lastUpdate，elapsed 改用系统执行耗时；D07：4 处 range 覆盖 StageSyncAfterDestroy；D10：所有 yield 检查返回值；D13：sl[:i]+sg+sl[i:] 显式构造；D14：先判 ok 再取方法值（15 分支统一），随 D21 表驱动化重构 | ✅ 已闭环 2026-09-14 |
| D08/D09（只读依赖） | 只读语义的核心是无锁并行调度的冲突判断依据（isFriend 分组→同线程串行）；Go 无 const，拷贝是防误写的无奈之举 | 方案 R2（R1 的立项细化，2026-09-16 批准）：codegen 为每个组件及其字段可达的所有 struct 生成零拷贝只读视图 TReadOnly（8 字节指针包装，仅 getter，编译期防写）；新增 GetComponentsReadOnly/GetBuddyReadOnly（视图自描述：ComponentPacketIdentifier/FromPtr，单类型参数 API，配套 ComponentSet.iterPtr）；行为矩阵：RO dep 调可写 API 静默返回空、View API 仅限只读依赖；嵌套 struct 生成嵌套视图、定长数组 Len/At、[n]byte @string() 生成 FieldString（深层防护完整，组件字段校验保证无泄漏点）；删除 IterReadOnly 拷贝迭代；手写组件需自实现 ReadOnly() 方可使用 View API；嵌套 struct 不支持跨包（生成器显式报错）；API 统一为单形态（曾评估单参数+结构体约束推断，因要求视图字段导出产生逃逸口被否决；曾生成包级便捷函数，后按用户决策移除以保持 API 统一性）。详细设计：AgentDocs/Plan_ReadOnly_View_Codegen.md | 🛠 已修复 2026-09-16（D08 语义扶正 + TReadOnly codegen 落地，ecs/example/bench/test 全量迁移验证通过） |
| D11/D12（EntityIDGenerator/ReuseID） | ① ID 复用是为了让 index 落在小数值、控制 SparseArray 稀疏索引层内存，非防悬空；② removeDelay 攒批是无锁化设计：避免并行场景下 ID 生成/回收的互斥锁等待；③ delayFlush 的 30 行归并是为了避免全表扫描 ids——flush 代价只与 freelist 链长/批量成正比，且要求保持严格升序（小号 index 优先复用） | ① delayFlush：保留语义，重写为教科书式双有序链归并（removeDelay 排序后与 ids 内嵌 freelist 链 merge，不引入新数据结构，沿用 ids[idx].index 作 next 指针），加注释+单测；② 删除 entity.go:92-94 收缩死代码；③ FreeID 增加越界/reuse 不匹配/double-free 检查；④ EntitySet.Get/CSet.get 增加 reuse 代数校验（一次 int64 比较）防悬空句柄静默命中；⑤ 实体增删约定只在主线程/OpLog 帧同步点发生，ID 生成器保持单线程无锁，并行 system 内禁止直接 NewEntity（API 约束 + 断言） | ✅ 已闭环 2026-09-14 |
