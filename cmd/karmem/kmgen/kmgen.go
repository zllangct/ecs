package kmgen

import (
	"embed"
	"io"
	"strings"
	"text/template"

	"github.com/zllangct/ecs/cmd/karmem/kmparser"
)

//go:embed *_template.*
var templateFiles embed.FS

// Generators is a list of all generators available and registered by RegisterGenerator
var Generators []Generator

var componentSeq int32

// RegisterGenerator register the given Generator.
// You should use it on `init` function.
func RegisterGenerator(g Generator) {
	Generators = append(Generators, g)
}

type Generator interface {
	Start(file *kmparser.Content) (compiler Compiler, err error)
	Options() map[string]string
	Extensions() []string
	Language() string
	Finish(output io.Writer, buffer io.Reader) (err error)
}

type Compiler struct {
	Template []*template.Template
	Modules  []string
}

type TemplateFunctions struct {
	FromTags        func(s string) string
	FromStructTags  func(o kmparser.Structure, s string) string
	FromStructClass func(cls kmparser.StructClass) string

	HasTag func(o kmparser.Tags, s string) bool
	GetTag func(o kmparser.Tags, s string) string

	CanRef func(t kmparser.Type) bool

	ToNamePadding func(val any, root any) string

	ToStructName   func(val any) string
	ToFieldName    func(val any) string
	ToEnumName     func(val any) string
	ToFunctionName func(val any) string

	ToDefault       func(typ kmparser.Type) string
	ToPlainDefault  func(typ kmparser.Type) string
	ToType          func(typ kmparser.Type) string
	ToPlainType     func(typ kmparser.Type) string
	ToTypeView      func(typ kmparser.Type) string
	ToTypeSource    func(typ kmparser.Type) string
	ToPlainTypeView func(typ kmparser.Type) string

	IndexStruct      func(data []kmparser.Structure, index int32) kmparser.Structure
	IndexComment     func(data []kmparser.MultiComment, index int32) kmparser.MultiComment
	IndexStructField func(fields []kmparser.StructField, index int32) *kmparser.StructField
	IndexEnum        func(data []kmparser.Enumeration, index int32) kmparser.Enumeration

	StructGenerateArgs func(o kmparser.Structure, args ...string) StructGenerateArgs[kmparser.Structure]

	ComponentSeq func() int32
}

type StructGenerateArgs[T any] struct {
	Value     T
	EnableRef bool
	Map       map[string]string
}

type TagResult struct {
	Exist bool
	Value string
}

func fromGlobalTags(gen Generator, content *kmparser.Content, s string) string {
	def, ok := gen.Options()[s]
	if !ok {
		panic("invalid tag search")
	}

	if s == "package" {
		def = content.Name
	}

	tags := content.Tags
	name := gen.Language() + "." + s
	for i := range tags {
		if tags[i].Name == name {
			return tags[i].Value
		}
	}

	return def
}

func fromTags(gen Generator, content *kmparser.Content, tags kmparser.Tags, s string) string {
	def := fromGlobalTags(gen, content, s)
	name := gen.Language() + "." + s
	for i := range tags {
		if tags[i].Name == name {
			return tags[i].Value
		}
	}
	return def
}

func genStructArgs[T any](gen Generator, content *kmparser.Content, tags kmparser.Tags, o T, args ...string) StructGenerateArgs[T] {
	refTag := fromTags(gen, content, tags, "field.enable_ref")
	m := make(map[string]string)
	for i := 0; i < len(args); i += 2 {
		key := args[i]
		value := args[i+1]
		m[key] = value
	}
	return StructGenerateArgs[T]{
		Value:     o,
		EnableRef: refTag != "false",
		Map:       m,
	}
}

func NewTemplateFunctions(gen Generator, content *kmparser.Content) TemplateFunctions {
	return TemplateFunctions{
		FromTags: func(s string) string {
			return fromGlobalTags(gen, content, s)
		},
		FromStructTags: func(o kmparser.Structure, s string) string {
			def := fromGlobalTags(gen, content, s)
			tags := o.Data.Tags
			name := gen.Language() + "." + s
			for i := range tags {
				if tags[i].Name == name {
					return tags[i].Value
				}
			}
			return def
		},
		HasTag: func(o kmparser.Tags, s string) bool {
			for _, tag := range o {
				if tag.Name == s {
					return true
				}
			}
			return false
		},
		GetTag: func(o kmparser.Tags, s string) string {
			for _, tag := range o {
				if tag.Name == s {
					return tag.Value
				}
			}
			return ""
		},
		CanRef: func(t kmparser.Type) bool {
			if t.Model == kmparser.TypeModelSingle && (t.Format == kmparser.TypeFormatStruct || t.Format == kmparser.TypeFormatTable) {
				return true
			}
			return false
		},
		FromStructClass: func(cls kmparser.StructClass) string {
			switch cls {
			case kmparser.StructClassTable:
				return "table"
			case kmparser.StructClassInline:
				return "inline"
			default:
				panic("invalid struct class")
			}
		},
		ToNamePadding: func(val any, root any) string {
			var largest int
			var name string
			switch val := val.(type) {
			case kmparser.StructField:
				name = strings.TrimSpace(val.Data.Name)
				root := root.(kmparser.Structure)
				for i := range root.Data.Fields {
					if l := len(root.Data.Fields[i].Data.Name); l > largest {
						largest = l
					}
				}
			case kmparser.EnumField:
				name = strings.TrimSpace(val.Data.Name)
				root := root.(kmparser.Enumeration)
				if root.Data.IsSequential {
					return name
				}
				for i := range root.Data.Fields {
					if l := len(root.Data.Fields[i].Data.Name); l > largest {
						largest = l
					}
				}
			default:
				panic("invalid type")
			}
			return name + strings.Repeat(" ", largest-len(name))
		},
		IndexStruct: func(data []kmparser.Structure, index int32) kmparser.Structure {
			return data[index]
		},
		IndexComment: func(data []kmparser.MultiComment, index int32) kmparser.MultiComment {
			return data[index]
		},
		IndexStructField: func(fields []kmparser.StructField, index int32) *kmparser.StructField {
			fields[index].Data.Type.IsArray()
			return &fields[index]
		},
		IndexEnum: func(data []kmparser.Enumeration, index int32) kmparser.Enumeration {
			return data[index]
		},
		StructGenerateArgs: func(o kmparser.Structure, args ...string) StructGenerateArgs[kmparser.Structure] {
			return genStructArgs(gen, content, o.Data.Tags, o, args...)
		},
		ComponentSeq: func() int32 {
			componentSeq++
			return componentSeq
		},
	}
}

type TemplateData struct {
	*kmparser.Content
}

func NewTemplate(modules []string, funcs TemplateFunctions, pattern ...string) (compiler Compiler) {
	compiler.Modules = modules
	compiler.Template = make([]*template.Template, len(pattern))
	for i, v := range pattern {
		t := template.New("")
		t = t.Funcs(template.FuncMap{
			"FromTags":        funcs.FromTags,
			"FromStructTags":  funcs.FromStructTags,
			"FromStructClass": funcs.FromStructClass,

			"GetTag": funcs.GetTag,
			"HasTag": funcs.HasTag,

			"CanRef": funcs.CanRef,

			"ToDefault":      funcs.ToDefault,
			"ToPlainDefault": funcs.ToPlainDefault,

			"ToPlainType":     funcs.ToPlainType,
			"ToType":          funcs.ToType,
			"ToPlainTypeView": funcs.ToPlainTypeView,
			"ToTypeView":      funcs.ToTypeView,
			"ToTypeSource":    funcs.ToTypeSource,

			"ToNamePadding": funcs.ToNamePadding,

			"IndexStruct":        funcs.IndexStruct,
			"IndexComment":       funcs.IndexComment,
			"IndexStructField":   funcs.IndexStructField,
			"IndexEnum":          funcs.IndexEnum,
			"StructGenerateArgs": funcs.StructGenerateArgs,
			"ComponentSeq":       funcs.ComponentSeq,
		})
		t, err := t.ParseFS(templateFiles, v)
		if err != nil {
			panic(err)
		}
		compiler.Template[i] = t
	}
	return compiler
}

type generatorFinishCopy struct{}

func (*generatorFinishCopy) Finish(output io.Writer, buffer io.Reader) error {
	_, err := io.Copy(output, buffer)
	return err
}
