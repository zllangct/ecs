### 重要设计原则
* Component内只有数据
* Component存储内存连续，减少cpu cache-miss
* 仅System可修改Component数据
* Utility与System一一对应, Utility可缺省
* System仅通过对应的Utility与系统外交互
* World管理Entity、System
* World不可被ecs外部并发访问，可通过Gate间接访问
* Gate提供串行化能力，用于ecs外部安全访问ecs
* ”下一帧生效“
* 尽可能开发流程中无锁化，并不是ecs框架内绝对无锁
* 尽可能并行化System
* 阶段化执行流程，ecs全执行周期包含多个阶段，start、update、destroy等
* ecs全周期内，同步点位于主线程
* ecs api分为同步点api和异步点api
* 尽可能通过语法限制用户行为，而不是靠用户“自觉“