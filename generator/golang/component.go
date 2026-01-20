package golang

import (
	"fmt"

	"github.com/zllangct/rockmem/common"
	"github.com/zllangct/rockmem/generator"
	"github.com/zllangct/rockmem/generator/golang"
	"github.com/zllangct/rockmem/parser"
)

// ComponentGenerator ECS组件代码生成器
// 用于为带有 @component() 标签的结构体生成 Component 接口所需的方法
type ComponentGenerator struct {
	file     *parser.File
	fileTree *parser.FileTree
	// components 存储需要生成组件代码的结构体
	components []*parser.Struct
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
		return fmt.Errorf(errMsg)
	}

	// 如果有组件需要生成，添加 ecs 导入和 init 代码
	if len(c.components) > 0 {
		golang.AddCustomImport(file.Filename, "github.com/zllangct/ecs", "ecs")
		// 向 init 函数添加组件注册代码
		for _, st := range c.components {
			code := fmt.Sprintf("ecs.RegisterComponent[%s](\"%s\")", st.Name, file.Package.Name)
			golang.AddCustomInitCode(file.Filename, code)
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

	return w.Bytes(), nil
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
