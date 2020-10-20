package ecs

import (
	"reflect"
	"testing"
)

type TestComponent1 struct {
	ComponentBase
	ID int
}

type TestComponent2 struct {
	ComponentBase
	ID int
}

type TestComponent3 struct {
	ComponentBase
	ID int
}

func TestComponentCollection(t *testing.T) {
	tests := []IComponent{
		&TestComponent1{ID: 1},
		&TestComponent1{ID: 2},
		&TestComponent2{ID: 3},
		&TestComponent2{ID: 4},
		&TestComponent3{ID: 5},
		&TestComponent3{ID: 6},
	}
	cc := NewComponentCollection(16 * 4)

	for index, value := range tests {
		cc.Push(value, uint64(index))
	}
	//test GetComponents
	com1 := cc.GetComponents(&TestComponent1{})
	for _, value := range com1 {
		println(value.(*TestComponent1).ID)
	}
	//test GetComponent
	com2 := cc.GetComponent(&TestComponent3{}, 4)
	if com2 != nil {
		println(reflect.TypeOf(com2).String())
	}
	//test iterator
	cIter := cc.GetIterator()
	for i := cIter.first(); i != nil; i = cIter.next() {
		println(reflect.TypeOf(i).String())
	}
}
