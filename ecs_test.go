package ecs

import (
	"fmt"
	"testing"
)

//go:generate go run ./cmd/karmem/main.go build -golang ecs.km
//go:generate go run ./cmd/karmem/main.go fmt -s ecs.km
func TestComponent(t *testing.T) {
}

type namee struct {
	n string
}

func (n *namee) Name() {
	*n = namee{n: "hello"}
}

func TestName(t *testing.T) {
	n := &namee{}
	fmt.Printf("%v\n", n)
	n.Name()
	fmt.Printf("%v\n", n)
}

func TestOther(t *testing.T) {
	a := make([]int, 1024)
	fmt.Printf("%v\n", cap(a))
	a = a[:1]
	fmt.Printf("%v\n", cap(a))
}
