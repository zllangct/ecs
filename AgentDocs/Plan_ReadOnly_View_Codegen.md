# ECS 只读视图（TReadOnly）Codegen 方案设计与实施计划

> 日期：2026-09-16
> 状态：✅ 方案已批准（用户审批通过），实施中
> 关联缺陷：D08（语义已扶正）、D09（GetBuddy 只读保护空操作）、D31（IterReadOnly 拷贝开销）

---

## 一、背景

- D08 已闭环：`Dep[T]()` 默认可写、`WithDepReadOnly` 显式只读，`readonly()` 正确驱动 `isFriend` 并行分组。
- 遗留：只读依赖的访问侧保护不一致——`GetComponents` 走 `IterReadOnly`（每元素真拷贝），`GetBuddy` 返回原指针（无保护）。
- Go 无 const，编译期只读只能靠类型系统：codegen 为每个组件生成零拷贝只读视图。

## 二、用户已确认的决策点

| 决策点 | 结论 |
|--------|------|
| RO dep 调用可写 API（GetComponents/GetBuddy） | 静默返回空/nil（沿用 D19/D27 静默约定） |
| View API 适用面 | 仅限只读依赖；RW dep 调 View API 返回空 |
| View 生成深度 | 全量：嵌套 struct 生成嵌套 View；定长数组生成 Len()/At(i)；[n]byte @string() 生成字符串 getter |
| IterReadOnly | 删除，统一走零拷贝 View |
| 命名后缀 | `ReadOnly`（区别于 RockMem 的 Viewer/Modifier） |

补充分析结论：RockMem 生成的 `TViewer` 是序列化二进制（wire format）的视图，内存布局与 ECS 存储的 Go struct 不兼容（offset 0 是 4 字节 size 头，变长字段需 Reader），**不可直接复用**；TReadOnly 是 `*T` 的包装，两者独立。

## 三、生成物设计（扩展 ComponentGenerator）

对每个 `@component()` 组件、以及其字段可达的所有 struct（递归），生成：

```go
// PositionReadOnly 是 Position 的只读视图（8 字节指针包装，零拷贝，仅 getter）
type PositionReadOnly struct{ p *Position }

// 平铺基础字段
func (v PositionReadOnly) X() float32 { return v.p.X }

// 嵌套 struct 字段 → 嵌套视图
func (v TransformReadOnly) Inner() InnerReadOnly { return InnerReadOnly{p: &v.p.Inner} }

// 定长数组（基础类型）
func (v InventoryReadOnly) ItemsLen() int { return 8 }
func (v InventoryReadOnly) ItemsAt(i int) int32 { return v.p.Items[i] }

// 定长数组（struct 元素）→ 元素视图
func (v PathReadOnly) PointsAt(i int) PointReadOnly { return PointReadOnly{p: &v.p.Points[i]} }

// [n]byte @string() → 字符串 getter（复用 Get 逻辑，截断尾部 0）
func (v PlayerReadOnly) NameString() string { ... }

// 仅组件生成：泛型约束连接方法
func (x *Position) ReadOnly() PositionReadOnly { return PositionReadOnly{p: x} }
```

规则：
- 下划线字段跳过；字段名用 `common.ToPascalCase`。
- getter 返回类型：基础类型直接用 `PlainSchema`；自定义 enum 用 `PlainSchema`（跨包加包名前缀）；struct 字段返回对应 `TReadOnly`。
- 组件字段校验（inline、无 string/slice/@ptr）保证 View 树全覆盖、无泄漏点。
- 生成代码追加到现有 `components_generated.go` 同一文件（ComponentGenerator.Generate 末尾）。

## 四、ECS 运行时 API

```go
// component.go
type ReadOnlyView[TR any] interface {
    ComponentPacketIdentifier() rockmem.PacketIdentifier
    FromPtr(p unsafe.Pointer) TR
}

// ecs.go（SystemContext 方法 + 包级函数各一份，单类型参数）
func (ctx *SystemContext) GetComponentsReadOnly[TR ReadOnlyView[TR]]() iter.Seq2[EntityIndex, TR]
func (ctx *SystemContext) GetBuddyReadOnly[TR ReadOnlyView[TR]](index EntityIndex) (TR, bool)
```

调用点：

```go
for idx, pos := range ctx.GetComponentsReadOnly[components.PositionReadOnly]() {
    _ = pos.X()
}
```

### 视图自描述单类型参数（2026-09-16 终态决策）

双类型参数对调用方不友好。曾评估「单类型参数 + 结构体约束推断」：可行，但要求视图字段导出（Go 规则：不同包非导出字段名视为不同标识符，结构体约束无法匹配），导致 `view.P.X = 1` 逃逸口。**最终方案（用户提出）**：视图自描述全部类型信息——`ComponentPacketIdentifier()` 返回与源组件一致的 PacketIdentifier（依赖/组件集查找），`FromPtr(unsafe.Pointer)` 构造视图；API 不再需要 T/TP：

```go
type ReadOnlyView[TR any] interface {
    ComponentPacketIdentifier() rockmem.PacketIdentifier
    FromPtr(p unsafe.Pointer) TR
}

ctx.GetComponentsReadOnly[components.VelocityReadOnly]()
ctx.GetBuddyReadOnly[components.VelocityReadOnly](index)
```

配套：`ComponentSet` 接口新增 `iterPtr()`（unsafe.Pointer 密集遍历，替代原 `*CSet[T,TP]` 断言路径）；视图字段 `p` 保持非导出，编译期绝对防写不变。只读访问 API 统一为 `ctx.GetComponentsReadOnly[TR]()` / `ctx.GetBuddyReadOnly[TR](index)` 单一形态（曾生成的包级便捷函数 GetTReadOnly/GetTReadOnlyAt 已按用户决策移除，保持 API 统一性）。

### 运行时行为矩阵

| dep 声明 | GetComponents / GetBuddy | GetComponentsReadOnly / GetBuddyReadOnly |
|----------|--------------------------|-------------------------------------------|
| 未声明   | 空 / nil（现状）          | 空 / 零值+false                            |
| ReadWrite| 返回 *T（现状）           | 空 / false                                 |
| ReadOnly | 空 / nil（新，静默）       | 返回 TReadOnly 视图                        |

### 删除 IterReadOnly

- `sparse_array.go` 删除 `IterReadOnly`；`ecs.go` 只读分支改为返回空迭代器。
- `sparse_array_fix_test.go`、`dependency_fix_test.go` 中相关用例同步重写。

## 五、手写组件回退约定

未实现 `ReadOnly()` 方法的手写组件无法使用 View API；其 RO dep 走可写 API 会静默返回空。README 明确说明：只读保护依赖 codegen 生成，手写组件需自行实现 `ReadOnly()` 才能获得同等保护。

## 六、实施任务

1. ECS 运行时行为矩阵改造（TDD）：RO dep × 可写 API → 空；删 IterReadOnly。
2. `ReadOnlyPointer` 约束 + `GetComponentsReadOnly`/`GetBuddyReadOnly`（TDD，dummyComponent 手写视图验证）。
3. ComponentGenerator 生成 ReadOnly 视图（generator 单测：嵌套/数组/@string/连接方法）。
4. 重编 `rockmem-ecs`，重新生成 example1/example2/bench/test 组件。
5. 迁移调用点：example1、example2、bench、test/cmp；README 更新。
6. 全量测试；ECS_Defects.md 讨论记录回填（D08/D09 → 🛠）。

## 七、验证计划

- 单测：行为矩阵（未声明/RW/RO × 可写 API/View API）、生成器输出断言。
- 集成：ecs/test/cmp 用真实生成组件验证视图读取与编译期防写。
- example1/2 运行、bench 回归（RO 迭代零拷贝）。
