package ecs

import (
	"testing"
)

type TestComponent1 struct {
	Component[TestComponent1]
	ID int
}

type TestComponent2 struct {
	Component[TestComponent2]
	ID int
}

type TestComponent3 struct {
	Component[TestComponent3]
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

func TestComponentCollection(t *testing.T) {
	c1 := TestComponent1{ID: 1}
	c2 := TestComponent1{ID: 2}
	c3 := TestComponent2{ID: 3}
	c4 := TestComponent2{ID: 4}
	c5 := TestComponent3{ID: 5}
	c6 := TestComponent3{ID: 6}

	_ = c1
	_ = c2
	_ = c3
	_ = c4
	_ = c5
	_ = c6

	cc := NewComponentCollection(16 * 4)

	_=cc
	//ptrc1 := &c1
	//println("c1 old ptr:", uintptr(unsafe.Pointer(ptrc1)))
	//ret := cc.add(ptrc1, 1)
	//println("c1 new ptr:", uintptr(unsafe.Pointer(ptrc1)), uintptr(unsafe.Pointer(&ret)))
	//
	//cc.Push(&c2, 2)
	//cc.Push(&c3, 3)
	//cc.Push(&c4, 4)
	//cc.Push(&c5, 5)
	//cc.Push(&c6, 6)
	//
	////test GetComponents
	//iter0 := cc.GetComponents(TestComponent1{})
	//for com := iter0.Next(); com != iter0.End(); com = iter0.Next() {
	//	println(((*TestComponent1)(com)).ID)
	//}
	//
	////test GetComponent
	//com2 := cc.GetComponent(TestComponent2{}, 4)
	//if com2 != nil {
	//	println(((*TestComponent2)(com2)).ID)
	//}
	//test iterator
	//iter1 := cc.GetIterator()
	//for com := iter1.Next(); com != iter1.End(); com = iter1.Next() {
	//	println(com.GetBase())
	//}
}
