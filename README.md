# ECS

这是一个面向游戏服务与客户端的数据驱动 ECS（Entity-Component-System）Go 语言实现，聚焦高内聚、低耦合、连续内存布局、无锁并发调度以及整世界快照序列化。

## 目录

- [特性概览](#特性概览)
- [安装与代码生成](#安装与代码生成)
- [快速开始](#快速开始)
- [核心概念](#核心概念)
  - [World 与执行模式](#world-与执行模式)
  - [Entity 与句柄代数](#entity-与句柄代数)
  - [Component 与代码生成约束](#component-与代码生成约束)
  - [System 与执行阶段](#system-与执行阶段)
  - [只读依赖与零拷贝只读视图](#只读依赖与零拷贝只读视图)
  - [固定访问组（WithGroup / 组表存储）与 View API](#固定访问组withgroup--组表存储与-view-api)
  - [帧同步点与操作日志（OpLog）](#帧同步点与操作日志oplog)
- [世界序列化与确定性回放](#世界序列化与确定性回放)
- [运行示例与测试](#运行示例与测试)

---

## 特性概览

- **连续内存与两级存储**：组件默认存储在基于 `SparseArray` + `USet` 的连续数组中，支持 O(1) 稀疏索引与 swap-remove；通过 `WithGroup` 可将联合高频访问的组件自动合并为列式组表（`Archetype`），配合 `View2..View8` / `View2RO..View8RO` 实现同步列指针遍历。
- **自动依赖冲突调度**：`WithWorldASyncMode()` 下框架根据系统声明的 `WithDep`（`ReadWrite` / `ReadOnly`）自动划分无冲突批次并行执行，并在同 `Order` 档位内为冲突系统保持注册顺序；`Sync*` 同步阶段始终在主 goroutine 串行执行。
- **编译期只读防护**：`rockmem-ecs` 代码生成器为每个组件及其嵌套结构体生成零拷贝只读视图 `TReadOnly`，只读依赖仅可通过 `ReadOnlyView` 访问，编译期阻断误写。
- **整世界快照与确定性步进**：支持 `World.Marshal` / `MarshalTo` 只读联合导出与 `NewWorldFromData` / `Unmarshal` 原子恢复；支持 `World.Step(delta)` 固定步长确定性推进。

---

## 安装与代码生成

```shell
go get github.com/zllangct/ecs
```

### 定义 Component（`.rm` 文件）

使用 RockMem IDL 语法定义组件（`@component()` 要求结构体为 `inline`，所有字段及可达嵌套结构体必须为连续值布局，禁止 `string`、`[]T` 与 `@ptr()`；固定长度文本请使用 `[n]byte` + `@string()`）：

```text
rockmem components @package("components") @component();

struct Position inline {
    x float32;
    y float32;
    z float32;
}

struct Velocity inline {
    x float32;
    y float32;
    z float32;
}

struct Player inline {
    name [32]byte @string();
    level int32;
}

// 一次性组件：每帧结束自动清空
struct DamageEvent inline @disposable() {
    damage float32;
    target_id int64;
}
```

### 生成组件代码

```shell
go run ./cmd/rockmem-ecs generate -o ./components components.rm
```

或使用仓库预置构建工具 `bin/rockmem-ecs.exe`（可运行根目录 `generate.bat` 一键重编并重新生成）。生成代码会在 `init()` 中自动调用 `ecs.RegisterComponent[T](pkg)`，同时生成组件接口方法、`GetXString()` / `SetXString()` 辅助方法以及只读视图 `TReadOnly`。

---

## 快速开始

完整可运行代码见 [examples/example1](examples/example1/) 与 [examples/example2](examples/example2/)。

```go
package main

import (
	"fmt"
	"time"

	"github.com/zllangct/ecs"
	"github.com/zllangct/ecs/example/example1/components"
)

type MovementSystem struct{}

func (s *MovementSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("MovementSystem"),
		ecs.WithDep[components.Position](ecs.ReadWrite),
		ecs.WithDepReadOnly[components.Velocity](),
		ecs.WithGroup(ecs.Dep[components.Position](), ecs.Dep[components.Velocity]()),
	)
	return nil
}

func (s *MovementSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	dt := float32(event.Delta.Seconds())
	for pos, vel := range ecs.View2RO[components.Position, components.VelocityReadOnly](ctx) {
		pos.X += vel.X() * dt
		pos.Y += vel.Y() * dt
		pos.Z += vel.Z() * dt
	}
	return nil
}

func main() {
	world := ecs.NewWorld(ecs.WithWorldSyncMode())
	if err := world.Register[MovementSystem](); err != nil {
		panic(err)
	}

	// 注册轻量级系统（函数形式）
	if err := world.RegisterLight(func(ctx *ecs.SystemContext, event ecs.Event) error {
		q := ctx.NewQuery(ecs.WithComp[components.Position]())
		for idx := range q.Iter() {
			p, _ := ecs.QueryGet[components.Position](&q, idx)
			fmt.Printf("Entity %d Pos=(%.2f, %.2f, %.2f)\n", idx, p.X, p.Y, p.Z)
		}
		return nil
	}, ecs.WithName("PrintSystem"), ecs.WithDep[components.Position](ecs.ReadWrite)); err != nil {
		panic(err)
	}

	world.NewEntity(ecs.WithComponents(
		&components.Position{X: 1, Y: 2, Z: 3},
		&components.Velocity{X: 10, Y: 0, Z: 0},
	))

	_ = world.Step(16 * time.Millisecond)
	_ = world.Destroy()
}
```

---

## 核心概念

### World 与执行模式

通过 `ecs.NewWorld(opts...)` 创建世界：

- `ecs.WithWorldSyncMode()`（默认）：所有系统按 `Stage` 与 `Order` 在单线程线性执行。
- `ecs.WithWorldASyncMode()` 在 `Start`、`PreUpdate`、`Update`、`PostUpdate` 阶段按依赖冲突分批并发执行，同 `Order` 档位内的冲突系统依注册顺序执行。`Sync*` 和 `World.Destroy()` 的收尾回调均在调用者的 goroutine 串行执行。
- `ecs.WithWorldAutoOptimize()` 使 `Update()` 在帧末利用剩余墙钟预算整理 `CSet`。`Step(delta)` 只推进模拟，使固定输入的遍历顺序与本机帧耗时无关。
- `ecs.WithWorldDefaultUpdateRate(rate)`：设置默认目标帧率（默认 `30`）。
- `ecs.WithoutGroups()`：全局关闭 `Archetype` 组表存储，所有组件保留在 `CSet`（用于回归对照，与开启组表语义完全等价）。

外部代码可在 `Update()` / `Step()` 之外读取或装配实体：
- `world.NewEntity(opts...)` 创建实体并返回 `Entity` 值；
- `world.GetEntityInfo(e)` 获取短生命周期的 `*EntityInfo` 以调用 `Add(...)` / `Remove(...)`（注意：**切勿跨 `NewEntity` / `DestroyEntity` / 跨帧缓存 `*EntityInfo` 指针**）；
- `world.GetComponent[T, *T](e)` 在系统外按完整 `Entity` 句柄安全读取组件；
- `world.DestroyEntity(e)` 提交延迟销毁（下一帧首同步点生效）；
- `world.Destroy()` 触发 `Destroy` 阶段链并释放组件/实体/组表资源，将世界置为 `WorldStatusStop`（不可继续 `Update` 或原地 `Unmarshal`）。

### Entity 与句柄代数

`Entity` 是 `int64` 句柄，底层为 `{index EntityIndex (int32), reuse int32}`：
- 低 32 位 `index` 在实体销毁后经 `EntityIDGenerator` 批量归并（`delayCap=10`）按最小索引优先复用，控制稀疏索引数组内存上限；
- 高 32 位 `reuse` 在每次复用时递增，用于拦截已销毁旧句柄造成的 ABA 误访问。

### Component 与代码生成约束

1. **普通组件**：长期挂载在实体上，随实体或 `Remove` 销毁。
2. **Disposable 组件（`@disposable()`）**：一次性组件，每帧结束时自动从世界与实体 `compound` 中清空（适合事件通知，不可入 `WithGroup`）。
3. **Nomadic 组件（`@nomadic()`）**：游牧组件，不绑定具体实体，通过 `world.AddNomadic(comps...)` 添加并在每帧结束自动清空（不可入 `WithGroup`）。
4. **手写组件与字符串**：
   - 如不使用 `.rm` 代码生成，手写组件需实现 `ecs.Component` 并显式调用 `ecs.RegisterComponent[T](group)`；
   - 含 Go 指针（`string`、`[]T`、`*T`、`map` 等）的手写组件禁止通过 `WithGroup` 进入组表；
   - 手写无 IDL 组件可使用 `ecs.FixedString[[N]byte]` 存储定长字符串；若需支持快照非 POD 字段，需自行实现 `WriteAsRoot` / `ReadAsRoot`（建议优先使用 `rockmem-ecs` 生成）。

### System 与执行阶段

标准系统可按需实现 `InitReceiver` 及以下阶段回调，或通过 `world.RegisterLight` 注册轻量函数系统（仅响应 `Update`）：

| 阶段 | 执行次数 | 时机与并发性 |
|---|---|---|
| `Init` | 1 | `World.Register` 注册期（同步，通过 `SystemInitContext.SetOption` 配置） |
| `SyncBeforeStart` / `Start` / `SyncAfterStart` | 1 | 首帧 `Update`/`Step`（`Sync*` 同步串行，`Start` 并行） |
| `SyncBeforePreUpdate` / `PreUpdate` / `SyncAfterPreUpdate` | 每帧 | 预更新（`Sync*` 同步串行，`PreUpdate` 并行） |
| `SyncBeforeUpdate` / `Update` / `SyncAfterUpdate` | 每帧 | 主更新（`Sync*` 同步串行，`Update` 并行） |
| `SyncBeforePostUpdate` / `PostUpdate` / `SyncAfterPostUpdate` | 每帧 | 后更新（`Sync*` 同步串行，`PostUpdate` 并行） |
| `SyncBeforeDestroy` / `Destroy` / `SyncAfterDestroy` | 1 | `World.Destroy()`（线性执行销毁收尾，对应 `SyncAfterDestroyReceiver`） |

`SystemOption` 常用选项：`WithName`、`WithOrder`（支持 `OrderFront`、自定义 `Order`、默认 `OrderAppend`）、`WithDep`、`WithDepReadOnly`、`WithGroup`。

### 只读依赖与零拷贝只读视图

系统组件依赖默认为 `ecs.ReadWrite`，也可声明为 `ecs.ReadOnly`（`WithDepReadOnly[T]()`）。两系统对同一组件均只读时不冲突，可并行运行。

- `ReadWrite` 依赖：通过 `ctx.GetComponents[T]()`、`ctx.GetBuddy[T](idx)`、`ecs.QueryGet[T](&q, idx)` 或 `ViewN` 读取可写 `*T`。
- `ReadOnly` 依赖：必须通过 `ctx.GetComponentsReadOnly[TReadOnly]()`、`ctx.GetBuddyReadOnly[TReadOnly](idx)` 或 `ViewNRO` 读取只读视图。手写只读视图需实现 `ecs.ReadOnlyView[TR]` 接口（`ComponentPacketIdentifier()` 与 `FromPtr(unsafe.Pointer)`）。
- 行为约定（静默失败）：未声明依赖、权限不匹配（如对只读依赖调可写 API 或对可写依赖调只读视图 API）或组件集不存在时，均静默返回空迭代器或 `(零值, false)`。

### 固定访问组（WithGroup / 组表存储）与 View API

- **组声明**：在系统注册时通过 `ecs.WithGroup(ecs.Dep[A](), ecs.Dep[B]())` 声明联合访问组（至少 2 个不同类型、非 `disposable`/`nomadic`、无 Go 指针且非零大小）。多系统声明的重叠组在首帧自动合并为互不相交的并集组。
- **高效遍历**：
  - `ecs.View2[A, B](ctx)` / `ecs.View2RO[A, BReadOnly](ctx)`：双变量 lockstep 遍历；
  - `ecs.View3..View8` / `ecs.View3RO..View8RO`：多元组遍历，返回 `TupleN`（首组件可写、其余只读视图）；
  - `ctx.NewQuery(WithComp[A](), WithComp[B]())` + `QueryGet[T](&q, idx)`：支持任意读写混合、组表子集、跨组/非组混合查询，非组路径下各组件池按 `EntityIndex` 有序时自动启用 `merge-join` 归并连接。

### 帧同步点与操作日志（OpLog）

- 在系统执行期间请使用 `ctx.AddComponents(entity, comps...)` 写入所属系统的无锁操作队列（切勿在并行系统内直接调用非线程安全的 `world.NewEntity`）。
- 帧首先处理挂起的实体销毁，再合并组件命令。系统回调结束后清理一次性和游牧组件，再合并本帧产生的组件命令。每次合并先更新 `CSet` 与实体的组件集合，再按确定顺序收拢或拆出组行。系统内提交的新增组件在下一帧对系统可见。

---

## 世界序列化与确定性回放

在世界静止时调用 `world.Marshal()` 或 `world.MarshalTo(writer)`。此时所有系统回调均已结束，外部调用者也应暂停对世界的写入。导出只读取现有存储，遇到尚未应用的组件命令或实体销毁会返回错误。

`ecs.NewWorldFromData(data, opts...)` 创建新的世界。续跑前按原顺序注册系统，并恢复应用自行维护的系统状态。`world.Unmarshal(data)` 和 `world.UnmarshalFrom(reader)` 可在已初始化、尚未运行的世界上原位恢复，保留现有系统和它们的注册顺序，丢弃旧世界的待执行命令。首个 `Step` 或 `Update` 按所声明的组重建列存储。

当前写出的 V1 快照按每种组件的残段顺序和组行顺序编码。以相同的组件组重新装配后，首帧会沿存档顺序恢复这些行。旧 V1 快照仍可读取，旧写入器已排序的组件只能从存档保留的顺序开始。改变组声明、关闭组表或调用 `Optimize` 都可能改变后续遍历顺序。

固定步长回放使用 `world.Step(fixedDelta)`。要让顺序敏感的系统从检查点复现原结果，还需保持相同的系统配置、外部输入顺序和系统私有状态，并确保并发回调的共享副作用具备确定的顺序。

---

## 运行示例与测试

```shell
# 运行全仓单元测试（含竞态检测）
go test -race -count=1 ./...

# 运行示例 1（移动、伤害、渲染系统循环）
go -C examples/example1 run .

# 运行示例 2（世界序列化与反序列化往返）
go -C examples/example2 run .

# 运行组表与视图基准测试
go -C examples/bench test -bench "BenchmarkGroupStorage_|BenchmarkView2RO_" -benchmem -benchtime 100x -run xxx .
```
