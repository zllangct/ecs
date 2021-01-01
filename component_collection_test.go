package ecs

import (
	"reflect"
	"testing"
	"unsafe"
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

type ITest interface {
	Test()
	Test2()
}

type TestStruct struct {
	Name   string
	Field1 string
	Field2 string
}

func (t TestStruct) Test() {
	println(t.Name)
}

func (t TestStruct) Test2() {
	println("test2")
}
func inter(in IComponent) {
	ifaceStruct := (*iface)(unsafe.Pointer(&in))
	_ = ifaceStruct
}

func TestComponentCollection(t *testing.T) {
	c1 := TestComponent1{ID: 1}
	c2 := TestComponent1{ID: 2}
	c3 := TestComponent2{ID: 3}
	c4 := TestComponent2{ID: 4}
	c5 := TestComponent3{ID: 5}
	c6 := TestComponent3{ID: 6}

	c1.SetType(reflect.TypeOf(c1))
	c2.SetType(reflect.TypeOf(c2))
	c3.SetType(reflect.TypeOf(c3))
	c4.SetType(reflect.TypeOf(c4))
	c5.SetType(reflect.TypeOf(c5))
	c6.SetType(reflect.TypeOf(c6))

	_ = c1
	_ = c2
	_ = c3
	_ = c4
	_ = c5
	_ = c6

	inter(&c1)

	cc := NewComponentCollection(16 * 4)

	cc.Push(&c1, 1)
	cc.Push(&c2, 2)
	cc.Push(&c3, 3)
	cc.Push(&c4, 4)
	cc.Push(&c5, 5)
	cc.Push(&c6, 6)

	//test GetComponents
	iter0 := cc.GetComponents(TestComponent1{})
	for com := iter0.Next(); com != iter0.End(); com = iter0.Next() {
		println(((*TestComponent1)(com)).ID)
	}

	//test GetComponent
	com2 := cc.GetComponent(TestComponent2{}, 4)
	if com2 != nil {
		println(((*TestComponent2)(com2)).ID)
	}
	//test iterator
	iter1 := cc.GetIterator()
	for com := iter1.Next(); com != iter1.End(); com = iter1.Next() {
		println(com.GetBase())
	}
}
