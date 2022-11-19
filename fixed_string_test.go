package ecs

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestFixedString_Generate(t *testing.T) {
	return
	filePath := "./fixed_string_utils.go"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	h1 := `package ecs

`
	write.WriteString(h1)
	h2 := "type Fixed%d [%d]byte\n"
	for i := 1; i < 129; i++ {
		write.WriteString(fmt.Sprintf(h2, i, i))
	}
	write.Flush()
}

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
