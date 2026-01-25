package main

import (
	"fmt"

	"github.com/zllangct/ecs"
	"github.com/zllangct/ecs/example/example2/components"
	rockmem "github.com/zllangct/rockmem/golang"
)

func main() {
	fmt.Println("=== ECS Example 2: World序列化与反序列化 ===")
	fmt.Println()

	// 注意：组件已经在 components 包的 init() 函数中自动注册
	// 无需手动调用 ecs.RegisterComponent

	// ========================================
	// 第一部分：创建原始World并添加数据
	// ========================================
	fmt.Println(">>> PART 1: 创建原始World")
	fmt.Println()

	world1 := ecs.NewWorld(ecs.WithWorldSyncMode())

	// 注册一个简单的System用于更新
	world1.RegisterLight(func(ctx *ecs.SystemContext, event ecs.Event) error {
		return nil
	}, ecs.WithName("DummySystem"))

	// 创建实体1：玩家
	pos1 := &components.Position{}
	pos1.X = 100.5
	pos1.Y = 200.25
	pos1.Z = 300.75

	vel1 := &components.Velocity{}
	vel1.X = 1.5
	vel1.Y = 2.5
	vel1.Z = 3.5

	data1 := &components.EntityData{}
	copy(data1.Name[:], "Player1")
	data1.TypeId = 1
	data1.Flags = 0xFF

	inv1 := &components.Inventory{}
	inv1.Count = 3
	inv1.Items[0] = 101
	inv1.Items[1] = 102
	inv1.Items[2] = 103

	entity1 := world1.NewEntity(ecs.WithComponents(pos1, vel1, data1, inv1))

	// 创建实体2：敌人
	pos2 := &components.Position{}
	pos2.X = -50.0
	pos2.Y = 75.5
	pos2.Z = 0.0

	vel2 := &components.Velocity{}
	vel2.X = -1.0
	vel2.Y = 0.0
	vel2.Z = 1.0

	data2 := &components.EntityData{}
	copy(data2.Name[:], "Enemy1")
	data2.TypeId = 2
	data2.Flags = 0x0F

	entity2 := world1.NewEntity(ecs.WithComponents(pos2, vel2, data2))

	// 创建实体3：只有位置的道具
	pos3 := &components.Position{}
	pos3.X = 0.0
	pos3.Y = 0.0
	pos3.Z = 50.0

	entity3 := world1.NewEntity(ecs.WithComponents(pos3))

	// 运行几帧让World状态稳定
	for i := 0; i < 10; i++ {
		world1.Update()
	}

	fmt.Println("原始World状态:")
	fmt.Printf("  Entity1 (ID: %d): Player1\n", entity1.Entity().ToInt64())
	fmt.Printf("  Entity2 (ID: %d): Enemy1\n", entity2.Entity().ToInt64())
	fmt.Printf("  Entity3 (ID: %d): Item\n", entity3.Entity().ToInt64())

	// ========================================
	// 第二部分：序列化World
	// ========================================
	fmt.Println()
	fmt.Println(">>> PART 2: 序列化World为二进制")
	fmt.Println()

	// 序列化World
	writer := rockmem.NewWriter(4096)
	n, err := world1.MarshalTo(writer)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal world: %v", err))
	}

	// 获取序列化后的数据
	serializedData := writer.Bytes()

	fmt.Printf("序列化完成！数据大小: %d bytes (写入: %d)\n", len(serializedData), n)

	// ========================================
	// 第三部分：销毁原有World
	// ========================================
	fmt.Println()
	fmt.Println(">>> PART 3: 销毁原有World")
	fmt.Println()

	world1.Destroy()
	fmt.Println("原有World已销毁")

	// ========================================
	// 第四部分：从二进制恢复World
	// ========================================
	fmt.Println()
	fmt.Println(">>> PART 4: 从二进制恢复World")
	fmt.Println()

	// 反序列化World
	reader := rockmem.NewReader(serializedData)
	worldData := &ecs.SerializableWorldData{}
	worldData.ReadAsRoot(reader)

	// 创建新World并恢复状态
	world2 := ecs.NewWorldFromData(worldData, ecs.WithWorldSyncMode())

	// 注册System（System不会被序列化，需要重新注册）
	world2.RegisterLight(func(ctx *ecs.SystemContext, event ecs.Event) error {
		return nil
	}, ecs.WithName("DummySystem"))

	fmt.Println("World已从二进制恢复！")

	// ========================================
	// 第五部分：验证恢复后的World状态
	// ========================================
	fmt.Println()
	fmt.Println(">>> PART 5: 验证恢复后的World状态")
	fmt.Println()

	fmt.Println("恢复后World状态:")
	fmt.Printf("  Entities restored from binary data\n")

	// 验证数据一致性
	fmt.Println()
	fmt.Println(">>> PART 6: 数据一致性验证")
	fmt.Println()

	verifyDataConsistency(world2)

	// 继续运行恢复后的World
	fmt.Println()
	fmt.Println(">>> PART 7: 继续运行恢复后的World")
	fmt.Println()

	for i := 0; i < 5; i++ {
		world2.Update()
	}

	fmt.Println("恢复后的World继续运行5帧成功！")

	// 清理
	world2.Destroy()

	fmt.Println()
	fmt.Println("=== 示例完成 ===")
}

// getBytesAsString 从字节数组获取字符串（去掉末尾空字节）
func getBytesAsString(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

// verifyDataConsistency 验证数据一致性
func verifyDataConsistency(world ecs.World) {
	// 创建一个验证System来检查数据
	verifySystem := func(ctx *ecs.SystemContext, event ecs.Event) error {
		fmt.Println("验证组件数据:")

		// 验证Position组件
		posCount := 0
		for _, pos := range ecs.GetComponents[components.Position](ctx) {
			posCount++
			fmt.Printf("  Position: (%.2f, %.2f, %.2f)\n", pos.X, pos.Y, pos.Z)
		}
		fmt.Printf("  Position组件数量: %d (预期: 3)\n", posCount)

		// 验证Velocity组件
		velCount := 0
		for _, vel := range ecs.GetComponents[components.Velocity](ctx) {
			velCount++
			fmt.Printf("  Velocity: (%.2f, %.2f, %.2f)\n", vel.X, vel.Y, vel.Z)
		}
		fmt.Printf("  Velocity组件数量: %d (预期: 2)\n", velCount)

		// 验证EntityData组件
		dataCount := 0
		for _, data := range ecs.GetComponents[components.EntityData](ctx) {
			dataCount++
			name := getBytesAsString(data.Name[:])
			fmt.Printf("  EntityData: name='%s', typeId=%d, flags=0x%X\n",
				name, data.TypeId, data.Flags)
		}
		fmt.Printf("  EntityData组件数量: %d (预期: 2)\n", dataCount)

		// 验证Inventory组件
		invCount := 0
		for _, inv := range ecs.GetComponents[components.Inventory](ctx) {
			invCount++
			fmt.Printf("  Inventory: count=%d, items=[%d, %d, %d, ...]\n",
				inv.Count, inv.Items[0], inv.Items[1], inv.Items[2])
		}
		fmt.Printf("  Inventory组件数量: %d (预期: 1)\n", invCount)

		// 验证结果
		allPassed := posCount == 3 && velCount == 2 && dataCount == 2 && invCount == 1
		if allPassed {
			fmt.Println("\n✓ 所有验证通过！数据完整恢复。")
		} else {
			fmt.Println("\n✗ 验证失败！数据不一致。")
		}

		return nil
	}

	// 注册并运行验证System
	world.RegisterLight(verifySystem,
		ecs.WithName("VerifySystem"),
		ecs.WithDep[components.Position](ecs.ReadOnly),
		ecs.WithDep[components.Velocity](ecs.ReadOnly),
		ecs.WithDep[components.EntityData](ecs.ReadOnly),
		ecs.WithDep[components.Inventory](ecs.ReadOnly),
	)

	world.Update()
}
