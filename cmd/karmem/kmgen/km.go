package kmgen

import (
	"github.com/zllangct/ecs/cmd/karmem/kmparser"
)

type KarmemSchema struct {
	content *kmparser.Content
	generatorFinishCopy
}

func KarmemSchemaGenerator() Generator {
	return &KarmemSchema{}
}

func (gen *KarmemSchema) Start(file *kmparser.Content) (compiler Compiler, err error) {
	gen.content = file
	return NewTemplate(
		[]string{`header`, `enums`, `blocks`},
		gen.functions(),
		"km_template.*",
	), nil
}

func (gen *KarmemSchema) functions() (f TemplateFunctions) {
	f = NewTemplateFunctions(gen, gen.content)
	return f
}

func (gen *KarmemSchema) Language() string { return "karmem" }

func (gen *KarmemSchema) Options() map[string]string {
	return map[string]string{}
}

func (gen *KarmemSchema) Extensions() []string { return []string{`.km`} }
