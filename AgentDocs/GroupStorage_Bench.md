# 组表（Archetype Group）存储与查询优化基准结论

日期：2026-09-16（二轮更新：P0-P4 查询优化落地后）
环境：Windows / amd64 / Intel i9-10900K，Go 1.27，benchtime=100x
对照：移动系统 `{Transform, Velocity}` 联合遍历

## 二轮结果（查询优化后，10 万实体）

| 访问方式 | 耗时 | vs 一轮基线（CSet 5.78ms） |
|---|---|---|
| **组表 + View2RO lockstep（P0）** | **1.77ms** | **3.3x** |
| CSet + merge-join（P2，零声明自动生效） | 4.50ms | 1.28x |
| CSet + QueryGet/缓存（P1/P3/P4） | 4.50ms | 1.28x |
| 组表 + QueryGet（P1） | 6.22ms | 0.93x |
| 扰动场景 CSet / 组表 | 5.15ms / 5.16ms | 持平 |

1 万实体：View2RO 108µs（10.8ns/entity）vs CSet 407µs —— **3.8x**。

## 一轮结果（组表存储初版，Query + GetBuddy API）

| 场景 | CSet（WithoutGroups） | 组表 | 结论 |
|---|---|---|---|
| 1 万实体，干净 | 498µs | 490µs | 持平 |
| 10 万实体，干净 | 5.78ms | 6.22ms | 组表慢约 8% |
| 10 万实体，扰动 | 5.93ms | 5.42ms | 组表快约 9% |

## 分析

- **View2RO 验证了设计预期**：组表行内零查找 + 双列指针同步递增，每实体成本从 ~60ns 降到 ~18ns（10 万）/~11ns（1 万），逼近纯内存流。
- **merge-join 让 CSet 路径"不声明也能快"**：SparseArray 保序优化（递增 Add 不破坏 isKOrder）使顺序创建实体的常规负载天然有序，查询自动走归并连接，无需任何声明。
- **P1/P3/P4 是普惠优化**：查询缓存（跨帧复用预解析）、Exist 稀疏表过滤（替代 compound 指针链）、QueryGet 绑定访问器，两条路径都受益。
- 扰动场景两路径持平：merge-join 失效回退 Exist 过滤，组表行内优势被 GetBuddy 间接层稀释——该场景应迁移到 View2/View2RO。
- 组路径 Query 跳过 `IsSubSet`（不变量：行 ⟺ 拥有组全集 ⊇ buddies，恒真）是组路径成立的关键——保留过滤时组路径慢 4 倍（闭包非内联）。
- 分配：稳态各路径一致（约 22-25 allocs/op）；组表首帧收拢有一次性分配（摊销后不进稳态）。

## 过程中修复的既有 bug

1. `OrderedIntSet.FindIndexToInsert` 二分循环结束位置遇重复元素返回 `l-1` 而非 `-1`，同帧重复 Add 同组件会腐蚀 compound 为乱序。
2. op flush 对已销毁实体仍应用 Add op（destroy 先于 op flush），CSet 残留孤儿数据，index 复用后僵尸复活。

## 复现

```
cd ecs/example/bench
go test -bench "BenchmarkGroupStorage_|BenchmarkView2RO_" -benchmem -benchtime 100x -run xxx .
```

