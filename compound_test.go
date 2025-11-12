package ecs

import (
	"testing"
)

func TestCompound_find(t *testing.T) {
	type args struct {
		it ComponentIntType
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
			if got, _ := tt.c.Find(tt.args.it); got != tt.want {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestCompound_insertIndex(t *testing.T) {
	type args struct {
		it ComponentIntType
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
			if got := tt.c.FindIndexToInsert(tt.args.it, 0); got != tt.want {
				t.Errorf("insertIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompound_Add(t *testing.T) {
	type args struct {
		it ComponentIntType
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
		it ComponentIntType
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
		c.Add(ComponentIntType(i % 65535))
	}
}

func BenchmarkCompound_Find(b *testing.B) {
	var compoundSize = 100
	c := Compound{}
	for i := 0; i < compoundSize; i++ {
		c.Add(ComponentIntType(i % compoundSize))
	}
	m := map[ComponentIntType]struct{}{}
	for i := 0; i < compoundSize; i++ {
		m[(ComponentIntType(i % compoundSize))] = struct{}{}
	}

	b.Run("c", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c.Find(ComponentIntType(i % compoundSize))
		}
	})

	b.Run("m", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, ok := m[(ComponentIntType(i % compoundSize))]
			_ = ok
		}
	})
}
