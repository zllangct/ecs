package ecs

import (
	"testing"
)

func TestFixedString_String(t *testing.T) {

	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "test",
			arg:  "hello world",
			want: "hello world",
		},
		{
			name: "test1",
			arg:  "hello 中文 ☺",
			want: "hello 中文 ☺",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FixedString[Fixed128]{}
			f.Set(tt.arg)
			if got := f.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
