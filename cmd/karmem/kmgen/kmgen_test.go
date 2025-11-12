package kmgen

import (
	"bytes"
	"os"
	"testing"
	"text/template"

	"github.com/zllangct/ecs/cmd/karmem/kmparser"
)

func TestGenerator(t *testing.T) {
	//path := []string{"testdata/basic.km", "testdata/paths.km", "testdata/fieldenableref.km"}
	//path := []string{"testdata/fieldenableref.km"}
	//path := []string{"testdata/structfieldstring.km"}
	path := []string{"testdata/case_string_setter.km"}
	for _, path := range path {
		f, err := os.Open(path)
		if err != nil {
			t.Error(f)
			return
		}

		r := kmparser.NewReader(path, f)
		k, err := r.Parser()
		if err != nil {
			t.Error(err)
		}

		if len(Generators) == 0 {
			t.Error("no generator found")
		}

		for _, gen := range Generators {
			compiler, err := gen.Start(k)
			if err != nil {
				t.Error(err)
				return
			}
			for _, c := range compiler.Template {
				var buffer bytes.Buffer
				var output bytes.Buffer

				for _, n := range compiler.Modules {
					if err := c.ExecuteTemplate(&buffer, n, k); err != nil {
						t.Error(err)
						return
					}
				}

				if err := gen.Finish(&output, &buffer); err != nil {
					t.Error(err)
					return
				}

				println(output.String())
			}
		}
	}
}

func TestFormatter(t *testing.T) {
	//path := []string{"testdata/basic.km", "testdata/paths.km"}
	path := []string{"testdata/fieldenableref.km"}
	for _, path := range path {
		f, err := os.Open(path)
		if err != nil {
			t.Error(f)
			return
		}

		r := kmparser.NewReader(path, f)
		k, err := r.Parser()
		if err != nil {
			t.Fatal(err)
		}

		if len(Generators) == 0 {
			t.Error("no generator found")
		}

		gen := KarmemSchemaGenerator()
		compiler, err := gen.Start(k)
		if err != nil {
			t.Error(err)
			return
		}

		for _, c := range compiler.Template {
			var buffer bytes.Buffer
			var output bytes.Buffer

			for _, n := range compiler.Modules {
				if err := c.ExecuteTemplate(&buffer, n, k); err != nil {
					t.Error(err)
					return
				}
			}

			if err := gen.Finish(&output, &buffer); err != nil {
				t.Error(err)
				return
			}

			println(output.String())
		}
	}
}

func TestTemplate(t *testing.T) {
	tpl := `
OUTPUT:{{Index .Array 1}}
`
	funcMap := template.FuncMap{
		"Index": func(array []int, i int) int {
			return array[i]
		},
	}
	tt, err := template.New("").Funcs(funcMap).Parse(tpl)
	if err != nil {
		t.Error(err)
	}

	type Data struct {
		Array []int
		Index int
	}

	tt.Execute(os.Stdout, Data{[]int{1, 2, 3}, 1})
}

type TB struct {
	Data int32
}

type Ts struct {
	Data *TB
}

func (t *Ts) Set() {
	*t.Data = TB{1}
}

func TestOther(t *testing.T) {
}
