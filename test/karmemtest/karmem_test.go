package karmemtest

import (
	"testing"
	"unsafe"

	"github.com/google/go-cmp/cmp"
	karmem "github.com/zllangct/ecs/karmem"
)

//go:generate go run ../../cmd/karmem/main.go build -golang testdata/components.km -o components_generated.go
//go:generate go run ../../cmd/karmem/main.go fmt -s testdata/components.km
//go:generate go run ../../cmd/karmem/main.go build -golang testdata/container.km -o container_generated.go
//go:generate go run ../../cmd/karmem/main.go fmt -s testdata/container.km
//go:generate go run ../../cmd/karmem/main.go build -golang testdata/basic.km
//go:generate go run ../../cmd/karmem/main.go fmt -s testdata/basic.km
//go:generate go run ../../cmd/karmem/main.go build -golang testdata/case_chararray.km
//go:generate go run ../../cmd/karmem/main.go fmt -s testdata/case_chararray.km
//go:generate go run ../../cmd/karmem/main.go build -golang testdata/case_int64array.km
//go:generate go run ../../cmd/karmem/main.go fmt -s testdata/case_int64array.km
//go:generate go run ../../cmd/karmem/main.go build -golang testdata/case_string_setter.km -o case_string_setter_generated.go
//go:generate go run ../../cmd/karmem/main.go fmt -s testdata/case_string_setter.km

type EnDer[T any] interface {
	WriteAsRoot(writer karmem.Writer) (offset uint, err error)
	ReadAsRoot(reader *karmem.Reader)
	*T
}

func EnDe[T any, TP EnDer[T]](in T) bool {
	en := TP(&in)
	writer := karmem.NewWriter(0)
	en.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())
	out := new(T)
	TP(out).ReadAsRoot(reader)

	return cmp.Equal(in, *out)
}

func TestKR0(t *testing.T) {
	b3 := BI{
		ANC:  [20]int64{1, 2, 3},
		ANCS: "hello你好",
	}
	writer := karmem.NewWriter(0)
	writer.Bytes()
	b3.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())

	ref3new := BI{}
	ref3new.ReadAsRoot(reader)

	if !cmp.Equal(b3, ref3new) {
		t.Error("not equal")
	}
}

func TestKR1(t *testing.T) {
	b3 := Base{
		Data: [10]int64{1, 2, 3},
		ANC:  "hello你",
	}
	writer := karmem.NewWriter(0)
	b3.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())

	ref3new := Base{}
	ref3new.ReadAsRoot(reader)

	if !cmp.Equal(b3, ref3new) {
		t.Error("not equal")
	}
}

//go:generate go run ../../cmd/karmem/main.go build -golang testdata/case_ref1.km
//go:generate go run ../../cmd/karmem/main.go fmt -s testdata/case_ref1.km
func TestKR2(t *testing.T) {
	b3 := BIR{7}
	b3r := BIRR{
		Data: &b3,
	}
	writer := karmem.NewWriter(0)
	b3r.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())

	ref3new := BIRR{}
	ref3new.ReadAsRoot(reader)

	if !cmp.Equal(b3r, ref3new) {
		t.Error("not equal")
	}
}

func TestKR3(t *testing.T) {
	b3 := BII{7}
	b3r := BIIR{
		Data: &b3,
	}
	writer := karmem.NewWriter(0)
	b3r.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())

	ref3new := BIIR{}
	ref3new.ReadAsRoot(reader)

	if !cmp.Equal(b3r, ref3new) {
		t.Error("not equal")
	}
}

func TestKR4(t *testing.T) {
	b3 := Base3{
		Value: 345,
		Data: &Base{
			Data: [10]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		},
	}

	writer := karmem.NewWriter(0)

	b3.WriteAsRoot(writer)

	reader := karmem.NewReader(writer.Bytes())

	ref3new := Base3{}
	ref3new.ReadAsRoot(reader)

	if !cmp.Equal(b3, ref3new) {
		t.Error("not equal")
	}
}

func TestKR(t *testing.T) {
	b := Base{
		Data: [10]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
	}

	if !EnDe(b) {
		t.Error("EnDe() failed")
	}

	b3 := Base3{
		Value: 345,
		Data:  &b,
	}

	if !EnDe(b3) {
		t.Error("EnDe() failed")
	}

	b3arr := []Base3{b3}
	ref3 := Ref3{
		ID:   312,
		Arr:  &b3,
		Arr2: b3arr,
	}

	if !EnDe(ref3) {
		t.Error("EnDe() failed")
	}

	writer := karmem.NewWriter(0)

	ref3.WriteAsRoot(writer)

	reader := karmem.NewReader(writer.Bytes())

	ref3new := Ref3{}
	ref3new.ReadAsRoot(reader)

	if !cmp.Equal(ref3, ref3new) {
		t.Error("not equal")
	}
}

func TestWriteDefault(t *testing.T) {
	r1 := R1{
		ANC: "hello你好",
	}
	writer := karmem.NewWriter(0)
	r1.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())
	r1new := R1{}
	r1new.ReadAsRoot(reader)

	if !cmp.Equal(r1new, r1) {
		t.Error("not equal")
	}

	r2 := R2{
		ANC: "hello你好",
		Data: Base{
			ANC: "world",
		},
	}
	writer2 := karmem.NewWriter(0)
	r2.WriteAsRoot(writer2)
	reader2 := karmem.NewReader(writer2.Bytes())
	r2new := R2{}
	r2new.ReadAsRoot(reader2)

	if !cmp.Equal(r2new, r2) {
		t.Error("not equal")
	}
}

func TestRB(t *testing.T) {
	r1 := B{
		ANC: "hello你好",
	}
	println(len(r1.ANC))
	writer := karmem.NewWriter(0)
	r1.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())
	r1new := B{}
	r1new.ReadAsRoot(reader)

	if !cmp.Equal(r1new, r1) {
		t.Error("not equal")
	}
}

func TestStringSetter(t *testing.T) {
	r1 := BaseS{
		ANC:      "hello你好",
		AND:      "world",
		Username: "nihao",
	}
	println(len(r1.ANC))
	writer := karmem.NewWriter(0)
	r1.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())
	s := NewBaseSSource(reader, 0)
	p := unsafe.Pointer(s.BaseSViewer)
	_ = p
	str := s.ANC(reader)
	_ = str
	s.SetANC("你好")
	s.SetUsername("n")

	r1new := BaseS{}
	r1new.ReadAsRoot(reader)

	if !cmp.Equal(r1new.Username, "n") || !cmp.Equal(r1new.ANC, "你好") {
		t.Error("not equal")
	}
}

func BenchmarkKM(b *testing.B) {
	c := &Container{}
	writer := karmem.NewWriter(0)
	c.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())

	getter := NewContainerViewer(reader, 0)

	writer = karmem.NewWriter(0)
	cNew := &Container{}
	_ = cNew
	c.WriteAsRoot(writer)
	b.ResetTimer()

	b.Run("T1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := getter.Comp()
			_ = v[2].Y()
		}
	})

	b.Run("T2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = c.Comp[2].Y
		}
	})
}

func BenchmarkKM1(b *testing.B) {
	a := []int{}
	var c map[int]int = make(map[int]int)

	for v := range 10000 {
		a = append(a, v)
		c[v] = v
	}

	b.Run("T1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = a[i%10000]
		}
	})

	b.Run("T2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = c[i%10000]
		}
	})
}
