package golang

import (
	"fmt"

	"github.com/zllangct/rockmem/common"
	"github.com/zllangct/rockmem/generator"
	"github.com/zllangct/rockmem/generator/golang"
	"github.com/zllangct/rockmem/parser"
)

// ComponentGenerator ECS组件代码生成器
// 用于为带有 @component() 标签的结构体生成 Component 接口所需的方法，
// 以及只读依赖（ReadOnly）所需的零拷贝只读视图 TReadOnly（仅 getter，编译期防写）。
type ComponentGenerator struct {
	file     *parser.File
	fileTree *parser.FileTree
	// components 存储需要生成组件代码的结构体
	components []*parser.Struct
	// viewStructs 存储需要生成只读视图的结构体：组件及其字段可达的所有 struct（去重，按发现顺序）
	viewStructs []*parser.Struct
	// isComponent 标记 viewStructs 中哪些是组件（组件额外生成 ReadOnly() 连接方法）
	isComponent map[string]bool
	// errors 存储验证错误
	errors []string
}

func (c *ComponentGenerator) Name() string {
	return "ecs-component"
}

// Init 初始化组件生成器
func (c *ComponentGenerator) Init(fileTree *parser.FileTree, file *parser.File) error {
	c.file = file
	c.fileTree = fileTree
	c.components = nil
	c.viewStructs = nil
	c.isComponent = nil
	c.errors = nil

	// 检查文件级标签是否有 @component()
	fileHelper := parser.NewTagHelper(file.Tags)
	fileHasComponent, _ := fileHelper.GetBoolValue("component")

	// 收集需要生成组件代码的结构体
	for _, decl := range file.Declarations {
		st, ok := decl.(*parser.Struct)
		if !ok {
			continue
		}

		// 检查结构体级标签是否有 @component()
		structHelper := parser.NewTagHelper(st.Tags)
		structHasComponent, _ := structHelper.GetBoolValue("component")

		// 如果文件级或结构体级有 @component() 标签，则需要生成组件代码
		if fileHasComponent || structHasComponent {
			// 验证组件字段
			if err := c.validateComponentFields(st); err != nil {
				c.errors = append(c.errors, err.Error())
				continue
			}
			c.components = append(c.components, st)

		}
	}

	// 如果有验证错误，返回错误
	if len(c.errors) > 0 {
		errMsg := "ECS component validation errors:\n"
		for _, e := range c.errors {
			errMsg += "  - " + e + "\n"
		}
		return fmt.Errorf("%s", errMsg)
	}

	// 如果有组件需要生成，添加 ecs 导入和 init 代码
	if len(c.components) > 0 {
		golang.AddCustomImport(file.Filename, "github.com/zllangct/ecs", "ecs")
		// 向 init 函数添加组件注册代码
		for _, st := range c.components {
			code := fmt.Sprintf("ecs.RegisterComponent[%s](\"%s\")", st.Name, file.Package.Name)
			golang.AddCustomInitCode(file.Filename, code)
		}
		// 收集只读视图结构体：组件及其字段可达的所有 struct（BFS）
		if err := c.collectViewStructs(fileTree); err != nil {
			return err
		}
	}

	return nil
}

// collectViewStructs 以组件为起点 BFS 收集字段可达的所有 struct。
// 跨包 struct 字段无法在本文件生成对应视图，属于不支持场景，显式报错。
func (c *ComponentGenerator) collectViewStructs(fileTree *parser.FileTree) error {
	c.isComponent = make(map[string]bool, len(c.components))
	seen := map[string]bool{}
	queue := make([]*parser.Struct, 0, len(c.components))
	for _, st := range c.components {
		c.isComponent[st.Name] = true
		queue = append(queue, st)
	}

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]
		if seen[st.Name] {
			continue
		}
		seen[st.Name] = true
		c.viewStructs = append(c.viewStructs, st)

		for _, field := range st.Fields {
			if field.Type == nil || !field.Type.IsStruct() {
				continue
			}
			if field.Type.Package != nil && !field.Type.Package.IsLocal() {
				return fmt.Errorf("component '%s' field '%s': cross-package struct type '%s' "+
					"is not supported in components (cannot generate read-only view across packages)",
					st.Name, field.Name, field.Type.QualifiedName())
			}
			nested := fileTree.FindStructByName(field.Type.PlainSchema)
			if nested == nil {
				return fmt.Errorf("component '%s' field '%s': struct type '%s' not found",
					st.Name, field.Name, field.Type.PlainSchema)
			}
			queue = append(queue, nested)
		}
	}
	return nil
}

// calculateComponentSeq 计算组件序列号
// 优先使用 @ComSeq(n) 标签指定的值，否则基于包名+结构体名生成唯一哈希
// validateComponentFields 验证组件字段
// ECS 组件要求所有字段必须是内存连续的类型
// 禁用以下类型（基于生成的 Go 类型判断，而非 IDL 类型）：
// 1. string 类型 ([]char/[n]char) - 生成 Go string，内部是指针+长度，内存不连续
// 2. slice 类型 ([]T) - 生成 Go slice，内部是指针+长度+容量，内存不连续
// 3. @ptr() 标签字段 - 生成 Go 指针类型 (*T)，内存不连续
//
// 注意：table 结构体字段生成的 Go 类型是值类型（如 TableA 而非 *TableA），内存是连续的。
// table 和 inline 的区别仅在于序列化时的二进制布局，Go 结构体层面都是连续内存。
func (c *ComponentGenerator) validateComponentFields(st *parser.Struct) error {
	// 检查结构体本身是否为 inline 类型
	if st.Kind != parser.StructKindInline {
		return fmt.Errorf("component '%s' must be inline struct. "+
			"ECS requires contiguous memory layout. Use 'struct %s inline { ... }' instead",
			st.Name, st.Name)
	}

	for _, field := range st.Fields {
		// 1. 检查 string 类型 - 生成 Go string（指针+长度）
		if field.Type.IsString() {
			return fmt.Errorf("component '%s' field '%s' cannot be string type ([]char/[n]char). "+
				"ECS requires contiguous memory layout. Use [%d]byte with @string() tag instead",
				st.Name, field.Name, field.Type.Length)
		}

		// 2. 检查 slice 类型 - 生成 Go slice（指针+长度+容量）
		if field.Type.IsSlice() {
			return fmt.Errorf("component '%s' field '%s' cannot be slice type ([]%s). "+
				"ECS requires contiguous memory layout. Use fixed-size array [n]%s instead",
				st.Name, field.Name, field.Type.PlainSchema, field.Type.PlainSchema)
		}

		// 3. 检查 @ptr() 标签 - 生成 Go 指针类型 (*T)
		if field.IsPointer {
			return fmt.Errorf("component '%s' field '%s' cannot use @ptr() tag. "+
				"ECS requires contiguous memory layout. Pointer types are not allowed",
				st.Name, field.Name)
		}
	}
	return nil
}

// Generate 生成组件代码
func (c *ComponentGenerator) Generate(fileTree *parser.FileTree, file *parser.File) ([]byte, error) {
	// 如果没有组件需要生成，返回空
	if len(c.components) == 0 {
		return []byte{}, nil
	}

	w := generator.NewCodeWriter()

	// 为每个组件生成代码
	for _, st := range c.components {
		c.generateComponentMethods(w, st)
		w.P()

		// 生成 @string() 标签字段的字符串辅助方法
		c.generateStringHelperMethods(w, st)
	}

	// 生成只读视图（组件及其字段可达的所有 struct）
	c.generateReadOnlyViews(w)

	return w.Bytes(), nil
}

// generateReadOnlyViews 为所有视图结构体生成 TReadOnly 类型；
// 组件额外生成 (*T).ReadOnly() 连接方法（供 ecs.ReadOnlyPointer 约束推断）。
func (c *ComponentGenerator) generateReadOnlyViews(w *generator.CodeWriter) {
	for _, st := range c.viewStructs {
		c.generateReadOnlyView(w, st)
		w.P()

		if c.isComponent[st.Name] {
			c.generateReadOnlySelfDescribeMethods(w, st)
			w.P()
		}
	}
}

// generateReadOnlySelfDescribeMethods 为组件视图生成自描述方法，
// 使只读 API 仅需视图类型一个类型参数（ecs.ReadOnlyView 约束）：
//   - ComponentPacketIdentifier 与源组件 PacketIdentifier 一致（依赖/组件集查找）
//   - FromPtr 从组件内存指针构造视图（零拷贝）
func (c *ComponentGenerator) generateReadOnlySelfDescribeMethods(w *generator.CodeWriter, st *parser.Struct) {
	viewName := st.Name + "ReadOnly"

	w.Printf("// ComponentPacketIdentifier 与源组件 %s.PacketIdentifier() 返回一致，用于只读 API 的依赖与组件集查找", st.Name)
	w.Printf("func (v %s) ComponentPacketIdentifier() rockmem.PacketIdentifier {", viewName)
	w.Printf("	return PacketIdentifier%s", st.Name)
	w.Printf("}")
	w.P()

	w.Printf("// FromPtr 从组件内存指针构造 %s 的只读视图（零拷贝）", st.Name)
	w.Printf("func (v %s) FromPtr(p unsafe.Pointer) %s {", viewName, viewName)
	w.Printf("	return %s{p: (*%s)(p)}", viewName, st.Name)
	w.Printf("}")
}

// generateReadOnlyView 生成单个结构体的只读视图类型与全部字段 getter
func (c *ComponentGenerator) generateReadOnlyView(w *generator.CodeWriter, st *parser.Struct) {
	w.Printf("// %sReadOnly 是 %s 的只读视图（零拷贝指针包装，仅 getter，编译期防写）", st.Name, st.Name)
	w.Printf("type %sReadOnly struct{ p *%s }", st.Name, st.Name)
	w.P()

	for _, field := range st.Fields {
		if field.Name == "_" || field.Type == nil {
			continue
		}
		c.generateViewFieldGetter(w, st, field)
		w.P()
	}
}

// generateViewFieldGetter 生成单个字段的只读 getter：
//   - 标量（基础类型/枚举）: Field() T
//   - struct 字段: Field() TReadOnly（嵌套视图，深层防护）
//   - 定长数组: FieldLen() int + FieldAt(i) T / TReadOnly
//   - [n]byte @string(): FieldString() string
func (c *ComponentGenerator) generateViewFieldGetter(w *generator.CodeWriter, st *parser.Struct, field *parser.StructField) {
	viewName := st.Name + "ReadOnly"
	fieldName := common.ToPascalCase(field.Name)
	typ := field.Type

	switch {
	case typ.IsArray():
		// [n]byte @string() → 字符串 getter
		if c.isByteArray(field) {
			fieldHelper := parser.NewTagHelper(field.Tags)
			if hasStringTag, _ := fieldHelper.GetBoolValue("string"); hasStringTag {
				w.Printf("func (v %s) %sString() string {", viewName, fieldName)
				w.Printf("	n := len(v.p.%s)", fieldName)
				w.Printf("	for n > 0 && v.p.%s[n-1] == 0 {", fieldName)
				w.Printf("		n--")
				w.Printf("	}")
				w.Printf("	return string(v.p.%s[:n])", fieldName)
				w.Printf("}")
				return
			}
		}
		w.Printf("func (v %s) %sLen() int {", viewName, fieldName)
		w.Printf("	return %d", typ.Length)
		w.Printf("}")
		w.P()
		if typ.IsStruct() {
			elemView := typ.PlainSchema + "ReadOnly"
			w.Printf("func (v %s) %sAt(i int) %s {", viewName, fieldName, elemView)
			w.Printf("	return %s{p: &v.p.%s[i]}", elemView, fieldName)
			w.Printf("}")
		} else {
			w.Printf("func (v %s) %sAt(i int) %s {", viewName, fieldName, typ.QualifiedName())
			w.Printf("	return v.p.%s[i]", fieldName)
			w.Printf("}")
		}
	case typ.IsStruct():
		w.Printf("func (v %s) %s() %sReadOnly {", viewName, fieldName, typ.PlainSchema)
		w.Printf("	return %sReadOnly{p: &v.p.%s}", typ.PlainSchema, fieldName)
		w.Printf("}")
	default:
		// 基础类型/枚举标量
		w.Printf("func (v %s) %s() %s {", viewName, fieldName, typ.QualifiedName())
		w.Printf("	return v.p.%s", fieldName)
		w.Printf("}")
	}
}

// generateStringHelperMethods 为带有 @string() 标签的 [n]byte 字段生成字符串辅助方法
func (c *ComponentGenerator) generateStringHelperMethods(w *generator.CodeWriter, st *parser.Struct) {
	for _, field := range st.Fields {
		// 只处理 [n]byte 类型的字段
		if !c.isByteArray(field) {
			continue
		}

		// 检查是否有 @string() 标签
		fieldHelper := parser.NewTagHelper(field.Tags)
		hasStringTag, _ := fieldHelper.GetBoolValue("string")
		if !hasStringTag {
			continue
		}

		c.generateStringFieldMethods(w, st, field)
		w.P()
	}
}

// isByteArray 检查字段是否为 [n]byte 类型
func (c *ComponentGenerator) isByteArray(field *parser.StructField) bool {
	// 检查是否为数组类型
	if field.Type.Category != parser.CategoryArray {
		return false
	}
	// 检查基础类型是否为 byte (uint8)
	return field.Type.PlainSchema == "byte" || field.Type.PlainSchema == "uint8"
}

// generateStringFieldMethods 为字段生成 GetXXXString() 和 SetXXXString() 方法
func (c *ComponentGenerator) generateStringFieldMethods(w *generator.CodeWriter, st *parser.Struct, field *parser.StructField) {
	structName := st.Name
	fieldName := common.ToPascalCase(field.Name)
	arrayLen := field.Type.Length

	w.Printf("// String helper methods for %s.%s (marked with @string() tag)", structName, fieldName)
	w.P()

	// 生成 GetXXXString() 方法
	w.Printf("func (x *%s) Get%sString() string {", structName, fieldName)
	w.Printf("	// Find the actual string length (excluding trailing zeros)")
	w.Printf("	n := len(x.%s)", fieldName)
	w.Printf("	for n > 0 && x.%s[n-1] == 0 {", fieldName)
	w.Printf("		n--")
	w.Printf("	}")
	w.Printf("	return string(x.%s[:n])", fieldName)
	w.Printf("}")
	w.P()

	// 生成 SetXXXString() 方法
	w.Printf("func (x *%s) Set%sString(v string) {", structName, fieldName)
	w.Printf("	// Clear the byte array first")
	w.Printf("	for i := range x.%s {", fieldName)
	w.Printf("		x.%s[i] = 0", fieldName)
	w.Printf("	}")
	w.Printf("	// Copy string content, truncate if exceeds capacity")
	w.Printf("	n := len(v)")
	w.Printf("	if n > %d {", arrayLen)
	w.Printf("		n = %d", arrayLen)
	w.Printf("	}")
	w.Printf("	copy(x.%s[:], v[:n])", fieldName)
	w.Printf("}")
}

// generateComponentMethods 为结构体生成 Component 接口所需的方法
func (c *ComponentGenerator) generateComponentMethods(w *generator.CodeWriter, st *parser.Struct) {
	structName := st.Name

	// 检查是否有 @nomadic() 标签
	helper := parser.NewTagHelper(st.Tags)
	isNomadic, _ := helper.GetBoolValue("nomadic")

	// 检查是否有 @disposable() 标签
	isDisposable, _ := helper.GetBoolValue("disposable")

	w.Printf("// Component extension for %s", structName)
	w.P()

	// 生成 NewComponentSet 方法
	w.Printf("func (x *%s) NewComponentSet() ecs.ComponentSet {", structName)
	w.Printf("	return ecs.NewCSet[%s]()", structName)
	w.Printf("}")
	w.P()

	// 生成 IsNomadic 方法
	w.Printf("func (x *%s) IsNomadic() bool {", structName)
	w.Printf("	return %v", isNomadic)
	w.Printf("}")
	w.P()

	// 生成 IsDisposable 方法
	w.Printf("func (x *%s) IsDisposable() bool {", structName)
	w.Printf("	return %v", isDisposable)
	w.Printf("}")
}

// NewComponentGenerator 创建组件生成器
func NewComponentGenerator() generator.CodeGenerator {
	return &ComponentGenerator{}
}
