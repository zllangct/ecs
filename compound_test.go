package ecs

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func Test_getCompoundType(t *testing.T) {
	filePath := "./compound_utils.go"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	h1 := `
package ecs

import (
	"unsafe"
)

func getCompoundType(compound Compound) interface{} {
	length := len(compound)
	if length == 0 || length > 255 {
		return nil
	}
	switch length {`
	write.WriteString(h1)
	h2 := `
	case %d:
		return *(*[%d]uint16)(unsafe.Pointer(&compound[0]))`
	for i := 1; i < 256; i++ {
		write.WriteString(fmt.Sprintf(h2, i, i))
	}

	h3 := `
	}

	return nil
}`
	write.WriteString(h3)
	write.Flush()
}

func TestCompound_find(t *testing.T) {
	type args struct {
		it uint16
	}
	tests := []struct {
		name string
		c    Compound
		args args
		want int
	}{
		{
			name: "1",
			c:    Compound{1, 3, 5},
			args: args{it: 3},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Find(tt.args.it); got != tt.want {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestCompound_insertIndex(t *testing.T) {
	type args struct {
		it uint16
	}
	tests := []struct {
		name string
		c    Compound
		args args
		want int
	}{
		{
			name: "1",
			c:    Compound{1, 3, 4, 6, 7},
			args: args{it: 5},
			want: 3,
		},
		{
			name: "2",
			c:    Compound{1, 3, 4, 6, 9, 10},
			args: args{it: 5},
			want: 3,
		},
		{
			name: "3",
			c:    Compound{2, 3, 5, 5, 6},
			args: args{it: 1},
			want: 0,
		},
		{
			name: "4",
			c:    Compound{1, 2, 3, 5, 7, 8},
			args: args{it: 6},
			want: 4,
		},
		{
			name: "5",
			c:    Compound{1, 2, 3, 4, 5, 6},
			args: args{it: 7},
			want: 6,
		},
		{
			name: "6",
			c:    Compound{1, 2, 3, 4, 5, 6},
			args: args{it: 3},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.InsertIndex(tt.args.it); got != tt.want {
				t.Errorf("insertIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompound_Add(t *testing.T) {
	type args struct {
		it uint16
	}
	tests := []struct {
		name    string
		c       Compound
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			c:       Compound{1, 2, 4, 5, 6},
			args:    args{it: 3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok := tt.c.Add(tt.args.it); (ok != true) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", ok, tt.wantErr)
			}
		})
	}
}

func TestCompound_Remove(t *testing.T) {
	type args struct {
		it uint16
	}
	tests := []struct {
		name string
		c    Compound
		args args
	}{
		{
			name: "1",
			c:    Compound{1, 2, 3, 4, 5, 6},
			args: args{it: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Remove(tt.args.it)
		})
	}
}

func BenchmarkCompound_Add(b *testing.B) {
	c := Compound{}
	for i := 0; i < b.N; i++ {
		c.Add(uint16(i % 65535))
	}
}

const (
	CompoundSize = 20
)

func BenchmarkCompound_Find(b *testing.B) {
	c := Compound{}
	for i := 0; i < CompoundSize; i++ {
		c.Add(uint16(i % CompoundSize))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Find(uint16(i % CompoundSize))
	}
}

func BenchmarkCompound_MapFind(b *testing.B) {
	m := map[uint16]struct{}{}
	for i := 0; i < CompoundSize; i++ {
		m[(uint16(i % CompoundSize))] = struct{}{}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, ok := m[(uint16(i % CompoundSize))]
		_ = ok
	}
}

func BenchmarkCompound_BigMapFind(b *testing.B) {
	m := map[int]struct{}{}
	for i := 0; i < 50000; i++ {
		m[i] = struct{}{}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, ok := m[i%50000]
		_ = ok
	}
}
