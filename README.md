# ECS
这是一个ECS（Entity-Component-System）Go语言版本的实现，它聚焦于游戏领域的应用，帮助你快速构建一个高内聚、低耦合、易扩展、高性能的并行化游戏世界。
## 快速开始
### 安装
```shell
go install github.com/zllangct/ecs
```
### 简单示例
这是一个不完整的示例，但可以帮助您快速了解如何在你的程序中接入和使用ecs框架。
```go
package main

import (
    "fmt"
    "github.com/zllangct/ecs"
    "time"
)

// 定义你的系统, 需要嵌套ecs.System[T]，T为你的System类型
type TestSystem struct {
	ecs.System[TestSystem]
}

// 系统Init事件
func (w *TestSystem) Init(si SystemInitializer) {
	// 申明系统感兴趣的组件, 系统内无法获取未申明的组件
	w.SetRequirements(si, &TestComponent1{}, &TestComponent2{}, &TestComponent3{})
}

func (w *TestSystem) Update(event Event) {
	// 获取系统感兴趣的组件, 作为遍历的索引，也可以叫做key组件
	iter := ecs.GetComponentAll[TestComponent1](w)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		// 获取关联组件，他们和key组件属于同一实体，也有人称之为"兄弟组件"
		c2 := ecs.GetRelated[TestComponent2](w, c.owner)
		if c2 == nil {
			continue
		}
		
		// 一些简单的处理逻辑
		for i := 0; i < testOptimizerDummyMaxFor; i++ {
			c.Field1 += i
		}

		for i := 0; i < testOptimizerDummyMaxFor; i++ {
			c2.Field2 += i
		}
	}
}

func main() {
	// 创建一个世界需要的配置
	config := ecs.NewDefaultWorldConfig()
	// 创建一个世界
	world := ecs.NewSyncWorld(config)
	// 注册系统
	ecs.RegisterSystem[TestSystem](world)

	// 启动你的世界
	world.Startup()

	// 为你的世界添加实体
	entities := make([]Entity, count)
	for i := 0; i < count; i++ {
		e1 := world.NewEntity()
		e1.Add(&TestComponent1{}, &TestComponent2{}, &TestComponent3{})
		entities[i] = e1.Entity()
	}

	// 持续更新你的世界
	for {
		world.Update()
		time.Sleep(time.Second)
	}
}
```
示例中的 ```__world_Test_S_1``` ```__world_Test_C_1```为系统和组件，后面会详细介绍，完整的代码请移步 [ world_test.go ](./world_test.go)。
## 快速了解
### 什么是ECS？
ECS是Entity-Component-System的缩写，它是一种数据驱动的架构，它将数据和逻辑分离开来，使得数据和逻辑可以独立的扩展和复用。   
* 实体：实体代表通用对象。例如，在游戏引擎上下文中，每个粗制游戏对象都表示为实体。通常，它仅由一个唯一的ID组成。实现通常为此使用普通整数。  
* 组件：组件将实体标记为具有特定方面的标签，并保留建模该方面所需的数据。例如，每个可能造成伤害的游戏对象都可能具有与其实体相关的健康组件。实现通常使用结构，类或关联阵列。  
* 系统：系统是一个过程，可用于所有具有所需组件的实体。 例如，物理系统可能会查询具有质量，速度和位置成分的实体，并迭代对每个实体的组件集进行物理计算的结果。 

实体的行为可以在运行时由添加、删除或修改组件的系统更改。 这消除了在面向对象编程技术中难以理解、维护和扩展的深层和广泛继承层次结构的模糊问题。
常见的ECS方法与面向数据的设计技术高度兼容，并经常与之结合。 一个组件的所有实例的数据通常一起存储在物理内存中，这使得在多个实体上操作的系统能够有效地访问内存。  
(来源：[维基百科](https://en.wikipedia.org/wiki/Entity_component_system))
### 游戏领域为什么需要ECS？
总结起来大概是两个词：”瓶颈“和”契合“。  
&emsp;&emsp;此处我们忽略一些ECS的其他特性，比如解耦、易扩展等，在其他设计模式下也能够实现，这不是ECS的显著特点和优势。同时也应该提到，ECS并不是万能的灵丹妙药，
并不能解决通用编程的所有问题，作为一种设计思路，能够很好的解决一些特定的问题，就已经足以促使我们学习和使用它。谈到“瓶颈”，首先我们来看看游戏领域遇到的一些问题，当然这不是全部，我们重点关注与ECS
相关的问题。
* 游戏属于密集型运算类应用，无论是客户端还是游戏后台都需要大量的运算，比如物理、AI、视野等相关的计算。
* 游戏对时延的要求非常高，比如FPS游戏，每一帧的渲染时间都是有限的，如果超过了这个时间，就会出现卡顿的现象。我们一直在致力于创造更低时延、更流畅的游戏体验，但是随着游戏的复杂度越来越高，计算量越来越大，
这个问题也越来越严重。虽然现在硬件资源日益强大，可以满足大部分的需求，同事可以通过分布式的计算来解决一些问题，但分布式带来的网络消耗也是非常大的，调用链路复杂后，问题更明显，分布式可以解决大计算量的问题
，但延时问题并不能有效解决。为了低延时，我们甚至常常使用“直连”的方式处理客户端和服务器之间的通讯，目的仅仅是为了rpc的消耗。除了更优秀的硬件和分布式，我们能做的还有多线程，多线程可以有效的利用多核CPU的资源，
，并通过内存进行线程间的交互，非常高效，避免了分布式中rpc的网络消耗。但是多线程会使得代码变得复杂，而且线程间的交互也会带来一些问题，比如线程安全、死锁等，甚至由于数据竞争的存在，很容易写出比单线程更低效的
多线程游戏系统。游戏系统对时序的要求非常高，多线程的异步逻辑往往很难保证正确的逻辑时序，所以通常我们的游戏服务更多的是采用单线程模式，仅仅会把耗时、逻辑独立的部分逻辑使用多线程优化，主体维持单线程。于是不能
愉快的利用分布式资源和多核资源成为了游戏领域的瓶颈。
* 游戏系统的更新，大家很熟悉的就是热更新，往往采用内嵌脚本语言的方式去实现，但代价就是性能和动态语言的缺点（同时也算是优点，矛盾总是存在的）。还有一些使用共享内存方案，将游戏系统的状态数据存放在共享内存中，
重启更新游戏后台程序后从共享内存中读取并恢复游戏状态，但对于OOP的设计模式来说，这将是一个非常复杂的过程，可能你需要一个自定义的内存分配器，因为只有这样你才能处理好共享内存的分配。可能你需要定制化你所使用
到的容器，可能需要处理恢复虚表等令人头痛的问题。还有一些序列化的方案，将游戏状态在更新前序列化并保存，更新后再反序列化并恢复状态，和共享内存方案相比，这解决了一些头疼的问题，让状态的存储和恢复变得简单了一些
，但是序列化的过程也是一个非常耗时的过程，而且序列化的数据量也是非常大的，这也是一个非常大的问题。所以游戏系统的更新也是一个非常头疼的问题。
* 在大世界的游戏世界中，单一的游戏服务器往往不能满足游戏的需求，无论用怎么样的构架方式，都避免不了游戏数据在服务器之前的迁移，不同的数据分布在不同的对象中，转移中数据的处理将是一个复杂的过程，同时也往往
会在这个环节产生很多的BUG，让你的同事和你非常头疼。  

&emsp;&emsp;刻意提到上面的一些问题，是因为我们可以通过ECS的方式去改善或者解决这些问题。上面的这些“瓶颈”我们总结一下，可以归纳为以下几点：1、怎么利用多核资源 2、怎么做到优雅更新 
3、怎么做到易于数据迁移，快速的序列化状态数据 4、怎么利用分布式资源。非ECS的方案也能通过一定的方法来解决上诉问题，但我们希望一种更简单，代价更低，甚至是更好的解决办法，我们可以逐个的分析，尝试着用ECS的方式去解决这些问题。
针对多核并行这一点，我们可以先看看ECS的特点，其中我们注意到，Component中仅有数据，一个对象的抽象，将会是拆分成不同的“数据块”分布在不同的Component中，数据是离散的，数据是边界明显的，单个System中仅是处理
改系统感兴趣Component，由此可见，数据是根据System的需求组合在一起的，组合后的数据成为System逻辑意义上的“对象”，Component是有关联的最小集合，为组合提供了更多的可能性，更灵活的复用。那么Component上的明显隔离，
System的互相独立这两点为并行创造了有利条件，因为独立的系统之间不会有数据竞争，当然System之间可能会对同一Component感兴趣，这个问题很好解决，只需要计算出系统间的组件依赖，就可以知道那些系统是可以并行的，
而这一过程可以在ECS的启动阶段完成计算，不会影响到运行时的性能。 然后利用一定的算法，系统拆分成互不依赖多批次，然后再分批次并行执行。并行的过程中由于没有数据竞争，我们完全不需要使用锁的方式去保证线程安全，
因为它们天生就是安全的，无锁化的并行将给并行带来巨大的性能提升，永远不会把大量的时间消耗在锁的竞争上，不会导致多线程比单线程性能更低的结果，所以ECS的设计模式就是天生为并行而生的。 第二点，怎么优雅的更新，上面也提到了
一些更新方案会因为系统状态的复杂性导致很难做到全量状态的存储，并恢复到更新前的状态。ECS中所有的状态即Component都是单独存储的，并且结构异常简单，是通过线性数组存储，数组是内存连续的容器，很容易做到全量的存储，所以
ECS中仅需要简单的储存所有Component的容器，这个过程非常简单，你可以直接通过二进制的形式转存到文件或者共享内存中，不需要序列化，这个过程非常高效。在恢复上，你可以简单的将二进制转换为原有的Component数组，然后ECS
系统中的所有状态就恢复了，下一次Update时，系统已经恢复如初。正是这个特性，同样适合解决第三点中提到的问题，数据分离且易于转存的特性使得数据在不同的服务器间的转移逻辑异常简单。第四点，ECS的系统是独立的，数据是界限明显的，
你可以轻易的找到独立的System和Component的组合，并将他们拆分至其他节点去运行。  
&emsp;&emsp;接下来再谈谈“契合”，ECS的数据是离散的，是有明显的边界，那么在开发中涉及到的应用中都贴合这个特性么？显然不是的，比如一套复杂的ERP管理系统，每一个API都需要访问到大量的数据，
API指间需要访问的数据间有大量的耦合，我们根本无法很好的拆分出功能独立的Component，显然这跟ECS不是那么的契合。游戏世界中往往有大量的游戏对象，并且需要高频的创建或者销毁，这在ECS中是非常廉价的操作，因为
游戏对象在ECS被抽象为Entity和一些Component，然而Entity仅仅是一个整型id，Component的添加和删除也仅仅是数组的操作，我们对Component数组中Component的顺序并不敏感，利用SwapAndRemove的方式能给组件的添加和
移除带来o(1)级别的效率，创建和销毁并不需要执行额外的逻辑，这一点非常好，我们可总结为，ECS适合大量对象频繁创建销毁的使用场景。游戏世界中虽然对象数据可能非常庞大，但是他们有着非常多的共性，比如，物理是玩家、野怪、载具他们
都有移动的功能，战斗中的对象都有伤害机制，这些都可以统一抽象成同一个Component，并使用用一个System来处理，我们知道ECS中遍历同一类型的组件是非常高效的，这一点游戏领域跟ECS是贴合的。很多游戏中都有buff机制，在ECS中
一个特定的buff可以抽象为一个Component，buff的产生和实效可以抽象为ECS中的Component的创建和销毁，非常方便。ECS提倡“组合优于继承”的设计理念，“组合”这个概念在游戏世界中也非常常见，比如，移动组件、伤害组件、敌人标签组件就
可以组合成一个简单的敌对小兵，稍作改变，我们把敌人标签换成队友标签，那么我们可以得到一个为你战斗的队友，ECS中Component的删除和新增是一个基础特性，所以你可以很轻易的完成这种基于组合的创造，而不用书写多余的逻辑。此外，
游戏世界中也是具有明显的隔离特性。至此我们可以总结一下游戏世界与ECS的显著契合点, 当然列举的可能并不全面：1、大量对象 2、频繁创建销毁 3、共性行为，比如移动、伤害、buff机制 4、组合在游戏中的抽象优势 5、明确的数据隔离特性
### 基本概念
ECS的整体方向是明确的，但是根据不同的ECS实现，会有一些细节上的差异，比如，命名上的差异，Component存储结构上的差异，甚至可能会根据实际的需要或者开发语言的原因妥协一些特性或者添加一下新的特性，下面会介绍一下，我们当前框架下
的一些基本概念。
#### World
World是ECS的核心，用来存储所有的Entity和Component，并管理所有的System, 它类似于通用开发中的管理器。框架中World分为SyncWorld和AsyncWorld，根据不同的使用场景，尤其是Go语言非常轻松的写出多线程逻辑的时候，我们
在调用ECS接口时，可能是同一线程，也可能是不同的线程。如果开发的入口逻辑可以保证是同一线程，那么可以选择使用SyncWorld，如果不能，那么选用AsyncWorld，框架会自动帮你处理好ECS与外围系统的交互安全问题。举一个很常见的例子，比如
使用框架构建一个客户端系统的时候，往往会有一个明确的主线程，此时使用SyncWorld就可以了，但如果是后台服务，我们知道在Go语言的实践中，往往不同的网络消息会交由不同的线程（goroutine，轻量线程，跟协程还是有所区别，我统一叫线程）来处理，
ECS内部能够很好的处理多线程，但是ECS外部的调用，需要由开发者去控制，此时选择AsyncWorld就可以了，AsyncWorld与SyncWorld的区别在于，AsyncWorld是线程安全的。
#### Entity
Entity是ECS中的基本单元，它是一个64位的整数ID，用来标识一个对象。Entity的ID由World来分配，在任意时刻，Entity的ID都是唯一的，但是不严格保证全局唯一，所以应当注意使用途中不应该用Entity代替UID使用，原因是ID低32位被设计成可复用
，ECS系统单次运行时所申请的ID是唯一的，但重复运行后会重复。有些ECS框架会使用GUID或者使用雪花算法等来生成全局唯一的ID，但我们没有这样做，原因后续会讲到。
#### Component
Component是ECS中的数据单元，它是一个结构体，用来存储数据，每个Component都有一个唯一的类型ID，用来标识这个Component的类型，Component的类型ID由World来分配，同样的，Component的类型ID也是唯一的，但不严格保证全局唯一。
#### System
System是ECS中的逻辑单元，在其他的ECS框架实现中，可能System会简单到就是一个函数，实际上这足够满足ECS对系统的描述，但在我们的框架中，System会相对复杂一些，但是他的本质还是处理Component的逻辑，后面会介绍具体实现。
#### Utility
Utility是ECS中的工具单元，但与其他的ECS实现的不同，我们的Utility不再是公共处理函数的集合，他是System的扩展，每一个Utility都与一个特定的System绑定，我们的原则中，仅有System可以直接操作Component，Component只能被System读写，
且System总是独立的，但也总是会有例外，Utility就是这个沟通的桥梁，比如我们明确SystemA对Component1感兴趣，那么SystemA是可以光明正大获取到Component1，但其他地方无法获取到Component1，此时我们需要使用绑定在SystemA上的Utility来借助
SystemA获取Component1，并做一定的操作。Utility的作用非常简单，就是一个桥梁，系统与外界沟通的桥梁，同时也让开发者清楚的感受到各个逻辑直接的界限，突出System独立的这个概念，可以有效的让某一特定的功能内聚到同一System内，降低耦合，
这是强制的限制，但是开发者可能通过语言层面的技巧突破这个限制，但我们希望开发者能够遵守这个原则，这样可以让代码更加清晰，更加易于维护。
#### Shape
Shape是ECS中的辅助单元，用来描述一组同属同一Entity的Component，也有ECS系统称之为Archetype或Sibling, 但是在我们这里不够准确，他们描述的是同一Entity所有Component的组合，而Shape描述的是一个子集，目的是用于同时获取一组Component，
当某System的某个操作，会固定操作同一Entity的一组组件时，建议使用Shape单元, 通过Shape对象可以获取到Component组迭代器，实现了类似Filter的功能，但我们不准备提供完整的Filter体系，因为我们尽可能不做过于复杂的筛选，API一旦提供
就避免不了滥用，当确有需要的时候，通过遍历实现，也非常简单。我们的筛选机制是通过系统Requirement的组合和Shape完成。
#### FixedString
FixedString是一个固定长度的字符串，适用于Component中的字符串，语言内置string是引用类型，如果内存的方式转移组件，那么内置string会带来一些问题。
#### “下一帧生效"
这是一个非常重要的概念，我们的ECS框架中，对Entity的Component创建、删除操作都会在下一帧生效。
### 设计思路
* 组件使用连续内存结构，减少cpu cache-miss
* 快速索引，无遍历获取兄弟组件
* 无锁化开发流程，开发者无需考虑并发安全问题
* 并行化System
* 阶段化执行流程，ecs全执行周期包含多个阶段，如start、update、destroy等
* ecs有主线程概念
* 通过语法限制用户行为，而不是靠用户“自觉“
* 整个世界可序列化，能够根据序列化数据恢复世界
* 充分利用Tick空档期

## 使用教程
### 创建一个World
#### SyncWorld
```go
world := ecs.NewSyncWorld()
```
#### AsyncWorld
```go
world := ecs.NewAsyncWorld()
```
#### 创建一个Entity
```go
// NewEntity() 创建的是一个EntityInfo类型
info := world.NewEntity()
entity := info.Entity()
```
#### 创建一个Component
```go
// 创建Component时需要嵌套Component[T], T为Component的实际类型，如下面的TestComponent
type TestComponent struct {
    ecs.Component[TestComponent]
    Field1 int
    Field2 int
}
```
#### 创建一个System
```go
type TestSystem struct {
	ecs.System[TestSystem]
}

// 所有系统事件都是可选的，如果不需要某个事件，可以不实现

// 系统Init事件
func (w *TestSystem) Init(si SystemInitializer) {}

// 系统Start事件
func (w *TestSystem) Start(event Event) {}

// 系统PreUpdate事件
func (w *TestSystem) PreUpdate(event Event) {}

// 系统Start事件
func (w *TestSystem) Update(event Event) {}

// 系统SyncUpdate事件
func (w *TestSystem) SyncUpdate(event Event) {}

// 系统PostUpdate事件
func (w *TestSystem) PostUpdate(event Event) {}

```
我们的ECS是阶段化的执行流程，提供了丰富的阶段性事件，当然最常用的是Update事件，当需要额外控制System先后顺序的时候，可以选择合适的事件混合使用。  
##### 支持的事件：
* Init
* SyncBeforeStart
* Start
* SyncAfterStart
* SyncBeforePreUpdate
* PreUpdate
* SyncAfterPreUpdate
* SyncBeforeUpdate
* Update
* SyncAfterUpdate
* SyncBeforePostUpdate
* PostUpdate
* SyncAfterPostUpdate
* SyncBeforeDestroy
* Destroy
* SyncAfterDestroy_

### 使用Utility与System交互
“Component中只有数据，System只有逻辑，只有System可以操作Component”这是我们ECS设计的指导思路，当所有的输入都来源于ECS的内部，所有系统之间
都不需要通过Component以外的方式进行交互，那么一切都将会是很美好的。当通常情况下，我们的系统会越来越复杂，不可避免会打破这个规则，但是我们并不希望
破坏我们系统的独立性，逻辑的内聚性，所以我们提供了Utility来解决这个问题。Utility与System一一对应，他可以作为单例的身份存在，同时也是与System
沟通的桥梁。在[ECS与外围系统交互](#ECS与外围系统交互)中，我们列举一个常见的例子的使用场景，在开发后台系统时，网络事件收到数据包后，需要交给我们
的ECS处理，那么就需要一次沟通，将数据包安全的交给ECS，这个时候就需要Utility来配合完成这个工作。我们先看如何创建一个Utility，并与对应的系统关联。
```go
type TestUtility struct {
    Utility[TestUtility]
}

func (u *TestUtility) ChangeName(entity Entity, name string) {
	sys := u.GetSystem()
	c := ecs.GetComponent[TestComponent](sys, entity)
	if c == nil {
		return
	}
	c.Name = name
	Log.Infof("Name changed, new:%s", name)
}

func main() {
        world := ......
		
        // 获取并使用Utility
        utility, _ :=ecs.GetUtility[TestUtility](world)
        utility.ChangeName(entity, "test")
		
        ......
}
```
需要注意的是，Utility与System是一一对应的， 需要再TestSystem的Init事件中绑定Utility,当然在不需要时可省略。
```go
func (s *TestSystem) Init(si SystemInitializer) {
    // 绑定Utility
    ecs.BindUtility[TestUtility](si)
}
```
### ECS与外围系统交互
World是ECS外部系统唯一能获取到跟对象，他是外部与ECS交互的入口，World分为SyncWorld和AsyncWorld，
SyncWorld适用于同步环境，AsyncWorld适用于异步环境。通常情况下，外部系统与ECS的交互流程如下：  
&emsp;&emsp;外部-->SyncWorld/AsyncWorld-->Utility-->Component  
直接看下面的例子，可能会觉得不太好理解，根据ECS的设计，只有System能修改Component，但是我们并不能直接获取到System实例，这是有意为之，因为System处于并行状态，不能随便的被调用，
但是ECS为此设计了Utility，Utility作为System的补充和桥梁，也可以修改Component，示例中```u.ChangeName(entity, name)```正是Utility对Component的
操作，我们先忽略其内部实现。
#### SyncWorld
```go
type TestUtility struct {
	Utility[TestUtility]
}

func (u *TestUtility) ChangeName(entity Entity, name string) {
	sys := u.GetSystem()
	c := GetComponent[TestComponent1](sys, entity)
	if c == nil {
		return
	}
	old := c.Name.String()
	c.Name.Set(name)
	Log.Infof("Name changed, old: %s, new:%s", old, name)
}

func main() {
    // 获取配置并创建世界
    config := NewDefaultWorldConfig()
    world := NewSyncWorld(config)

    // 注册系统
    RegisterSystem[TestSystem1](world)

    // 启动世界
    world.Startup()

    // 添加一下测试用的实体
    entities := make([]Entity, 100)
    for i := 0; i < 100; i++ {
        e1 := world.NewEntity()
        e1.Add(&TestComponent1{}, &TestComponent2{}, &TestComponent3{})
        entities[i] = e1.Entity()
    }

    // 尝试更新世界，使实体与他们的组件生效
    world.Update()

    // 获取Utility，并调用Utility的ChangeName方法，修改实体的Name字段
    getter := world.GetUtilityGetter()
    u, ok := GetUtility[TestUtility](getter)
	if ok {
        u.ChangeName(entities[0], "name0")
    }

    // 持续更新你的世界
    for {
        world.Update()
        time.Sleep(time.Second)
    }
}
```
#### AsyncWorld
AsyncWorld与SyncWorld不同，AsyncWorld被应用于多线程环境中，具备串行化的能力。
相较于SyncWorld，AsyncWorld多了一个Sync方法，用于将异步的操作同步化。
同时AsyncWorld不在需要我们手动调用Update方法，他会自动更新，在使用SyncWorld时我们明确的清楚每次调用Update时都处于同一线程，
但是在AsyncWorld中，Update方法可能会被多个线程调用，会导致我们的ECS混乱，这也是设计AsyncWorld的主要原因。
```go
type TestUtility struct {
    Utility[TestUtility]
}

func (u *TestUtility) ChangeName(entity Entity, name string) {
    sys := u.GetSystem()
    c := GetComponent[TestComponent1](sys, entity)
    if c == nil {
        return
    }
    old := c.Name.String()
    c.Name.Set(name)
    Log.Infof("Name changed, old: %s, new:%s", old, name)
}

func main() {
    // 获取配置并创建世界
    config := NewDefaultWorldConfig()
    world := NewSyncWorld(config)

    // 注册系统
    RegisterSystem[TestSystem1](world)

    // 启动世界
    world.Startup()

    // 添加一下测试用的实体
    entities := make([]Entity, __worldTest_Entity_Count)
    world.Sync(func(gaw SyncWrapper) {
        for i := 0; i < __worldTest_Entity_Count; i++ {
            e1 := gaw.NewEntity()
            e1.Add(&__world_Test_C_1{}, &__world_Test_C_2{}, &__world_Test_C_3{})
            entities[i] = e1.Entity()
        }
    })
	
    // 获取Utility，并调用Utility的ChangeName方法，修改实体的Name字段
    time.Sleep(time.Second * 2)
    world.Sync(func(gaw SyncWrapper) {
        u, ok := GetUtility[__world_Test_U_Input](gaw)
        if !ok {
            return
        }
        u.ChangeName(entities[0], "name2")
    })
	
    // 持续更新你的世界
    for {
        time.Sleep(time.Second)
    }
}

```

### 如何处理System的执行顺序
### 如何选择不同类型的Component
#### 常规组件
#### Free组件
#### Disposable组件
* 比如我们需要获取新增的Entity，可以在创建Entity时，同时添加Disposable组件，在Update中获取该Disposable组件集合，从而获取新增的Entity。Disposable组件
会在每一帧被清空，可以活的当做一次性tag使用。
#### FreeDisposable组件
### 系统中获取组件的方式
### 系统间的数据流动
### 一个完整的例子

## Benchmark
（努力完善中）

## API文档

## 设计细节
### 架构
### 日志
### World
### Entity
* Entity复用
* 不保证全局唯一
* Entity的ID是一个64位的整数，高32位是世界ID，低32位是EntityID，这样设计的目的是为了支持多个世界的存在，比如，我们可以在一个世界中创建一个Entity，然后把这个Entity
### Component
### System
### Utility
### Shape
### Compound
### 一次Update的执行流程
### 从添加Component到生效
### 容器
* unordered_set
* sparse_array
* ordered_int_set
### 迭代器
### 协程池
### ECS序列化和反序列化
### ECS中的并行
### Entity ID的管理
### 如何行为限制
* 通过状态量，限制某些API的访问时期，比如某些API只能在初始化阶段执行，否则将收到错误警告
* 通过guard参数，某些API的调用需要接受一个guard参数，guard参数的作用域受到框架的限制，比如在System的Init事件中
, guard生命周期为Init事件的执行周期，Init执行结束，guard参数失效，传递、转存无法维持guard参数的有效性，因此限制
了某些API仅能在特定的事件中被调用。（需要补充实现细节代码介绍）
### 优化器
* 优化器的作用，尽可能使得ecs框架内部的数据结构尽可能连续，减少cpu cache-miss
* 优化器工作时间点，tick驱动时，利用每次update的富裕时间进行优化工作
* 优化器工作内容，对于每个Component的数据结构，进行连续化优化
### 统计器
* 统计器的内容，统计ecs框架内部的运行状态，比如每个System的执行时间，每个Component的内存占用等
* 作用，帮助开发期调试、性能优化，配合优化器完整优化工作
## 特别注意
* 重复添加Component，会失败
* 同一帧内，多次移除、添加、移除...操作只会保留最终结果，因为“下一帧生效”会丢失中间过程，即使不会丢失，也没有实际的意义，建议避免这样的操作。
* Component所有成员变量都应该是值类型，string是引用类型，需要字符串类型时请使用 框架内的FixedString类型。
## 存在的一些问题
* 稀疏数组的内存占用问题
* EntityInfo的修改需要再同步点进行
* 不支持不对等tick，不存在多层次tick，比如A系统tick间隔50ms，B系统tick间隔30ms
* 并行时task的拆分粒度
* 并行退化，当开发者的系统依赖混乱，会导致系统关联度过高，框架调度时，会将有数据竞争的系统放到同一个线程中执行，从而导致并行退化，
最糟糕的情况是，退化为单线程系统
* 行为限制的缺失，由于golang语言的特性，无法严格得按照ECS的设计思路限制开发者的行为，开发者必须对数据驱动有一定的了解，无法严格
的使用语法规避潜在的危险行为
## 写在最后
欢迎大家提出宝贵意见，帮助完善ecs框架，这是一个长期和持续的过程，有不足或错误的地方，欢迎指正，欢迎PR。

### TODO
* [ ] 优化器实现
* [ ] 统计器完善
* [ ] 代码覆盖率
* [ ] world序列化
* [ ] 更新Atomic


