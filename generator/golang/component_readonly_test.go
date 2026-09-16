package golang

import (
	"strings"
	"testing"

	"github.com/zllangct/rockmem/parser"
)

// 运行 ComponentGenerator 生成 readonly_components.rm 的组件扩展代码
func generateReadonlyComponents(t *testing.T) string {
	t.Helper()
	ft, err := parser.Parse("testdata/readonly_components.rm")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	file := ft.GetInputFile()
	if file == nil {
		t.Fatal("input file not found in file tree")
	}
	g := NewComponentGenerator()
	if err := g.Init(ft, file); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	out, err := g.Generate(ft, file)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	return string(out)
}

func assertContains(t *testing.T, code string, fragments ...string) {
	t.Helper()
	for _, f := range fragments {
		if !strings.Contains(code, f) {
			t.Errorf("generated code missing fragment:\n%s", f)
		}
	}
}

// 只读视图类型与自描述方法：每个组件视图携带 ComponentPacketIdentifier 与 FromPtr
func TestComponentGenerator_ReadOnlyViewTypes(t *testing.T) {
	code := generateReadonlyComponents(t)
	assertContains(t, code,
		"type Vec3ReadOnly struct",
		"type TransformReadOnly struct",
		"type InventoryReadOnly struct",
		"type PathReadOnly struct",
		"type AIReadOnly struct",
		"type PlayerReadOnly struct",
		"func (v TransformReadOnly) ComponentPacketIdentifier() rockmem.PacketIdentifier {",
		"return PacketIdentifierTransform",
		"func (v TransformReadOnly) FromPtr(p unsafe.Pointer) TransformReadOnly {",
		"return TransformReadOnly{p: (*Transform)(p)}",
		"func (v PlayerReadOnly) ComponentPacketIdentifier() rockmem.PacketIdentifier {",
		"func (v PlayerReadOnly) FromPtr(p unsafe.Pointer) PlayerReadOnly {",
	)
}

// 平铺基础字段 getter
func TestComponentGenerator_ReadOnlyScalarGetter(t *testing.T) {
	code := generateReadonlyComponents(t)
	assertContains(t, code,
		"func (v Vec3ReadOnly) X() float32",
		"return v.p.X",
		"func (v TransformReadOnly) Scale() float32",
		"func (v AIReadOnly) TargetId() int64",
	)
}

// 嵌套 struct 字段返回嵌套只读视图（深层防护）
func TestComponentGenerator_ReadOnlyNestedStructGetter(t *testing.T) {
	code := generateReadonlyComponents(t)
	assertContains(t, code,
		"func (v TransformReadOnly) Pos() Vec3ReadOnly",
		"return Vec3ReadOnly{p: &v.p.Pos}",
		"func (v TransformReadOnly) Rot() Vec3ReadOnly",
	)
}

// 定长数组：基础类型生成 Len/At，struct 元素生成元素视图 At
func TestComponentGenerator_ReadOnlyArrayGetter(t *testing.T) {
	code := generateReadonlyComponents(t)
	assertContains(t, code,
		"func (v InventoryReadOnly) ItemsLen() int",
		"return 8",
		"func (v InventoryReadOnly) ItemsAt(i int) int32",
		"return v.p.Items[i]",
		"func (v PathReadOnly) PointsLen() int",
		"return 4",
		"func (v PathReadOnly) PointsAt(i int) Vec3ReadOnly",
		"return Vec3ReadOnly{p: &v.p.Points[i]}",
	)
}

// 枚举字段 getter 返回枚举类型
func TestComponentGenerator_ReadOnlyEnumGetter(t *testing.T) {
	code := generateReadonlyComponents(t)
	assertContains(t, code,
		"func (v AIReadOnly) State() State",
		"return v.p.State",
	)
}

// @string() 字段生成字符串 getter（截断尾部 0）
func TestComponentGenerator_ReadOnlyStringGetter(t *testing.T) {
	code := generateReadonlyComponents(t)
	assertContains(t, code,
		"func (v PlayerReadOnly) NameString() string",
		"v.p.Name[n-1] == 0",
		"return string(v.p.Name[:n])",
	)
}

// 不生成包级便捷函数：只读访问统一走 ctx.GetComponentsReadOnly[TR]() 泛型 API
func TestComponentGenerator_NoConvenienceFuncs(t *testing.T) {
	code := generateReadonlyComponents(t)
	if strings.Contains(code, "func GetTransformReadOnly(") || strings.Contains(code, "func GetTransformReadOnlyAt(") {
		t.Error("convenience funcs should not be generated")
	}
}
