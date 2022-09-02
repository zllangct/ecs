package ecs

import (
	"fmt"
	"testing"
)

type __systemGroup_Test_C_1 struct {
	Component[__systemGroup_Test_C_1]
}
type __systemGroup_Test_C_2 struct {
	Component[__systemGroup_Test_C_2]
}
type __systemGroup_Test_C_3 struct {
	Component[__systemGroup_Test_C_3]
}
type __systemGroup_Test_C_4 struct {
	Component[__systemGroup_Test_C_4]
}
type __systemGroup_Test_C_5 struct {
	Component[__systemGroup_Test_C_5]
}
type __systemGroup_Test_C_6 struct {
	Component[__systemGroup_Test_C_6]
}
type __systemGroup_Test_C_7 struct {
	Component[__systemGroup_Test_C_7]
}
type __systemGroup_Test_C_8 struct {
	Component[__systemGroup_Test_C_8]
}
type __systemGroup_Test_C_9 struct {
	Component[__systemGroup_Test_C_9]
}
type __systemGroup_Test_C_10 struct {
	Component[__systemGroup_Test_C_10]
}

type __systemGroup_Test_S_1 struct {
	System[__systemGroup_Test_S_1]
	Name int
}

func NewTestSystem(ID int, rqs ...IRequirement) *__systemGroup_Test_S_1 {
	s := &__systemGroup_Test_S_1{Name: ID}
	s.SetRequirements(rqs...)
	return s
}

func (p *__systemGroup_Test_S_1) Call(label int) interface{} {
	switch label {
	case 1:
		println(p.Name)
	}
	return nil
}

func (p *__systemGroup_Test_S_1) Filter(component IComponent, op CollectionOperate) {

}

func TestNewSystemGroup(t *testing.T) {
	tests := []ISystem{
		NewTestSystem(1, &__systemGroup_Test_C_1{}, &__systemGroup_Test_C_2{}),
		NewTestSystem(2, &ReadOnly[__systemGroup_Test_C_1]{}, &__systemGroup_Test_C_3{}),
		NewTestSystem(3, &__systemGroup_Test_C_2{}, &__systemGroup_Test_C_5{}),
		NewTestSystem(4, &__systemGroup_Test_C_2{}, &__systemGroup_Test_C_3{}, &__systemGroup_Test_C_6{}),
		NewTestSystem(5, &__systemGroup_Test_C_7{}),
		NewTestSystem(6, &__systemGroup_Test_C_9{}, &__systemGroup_Test_C_10{}),
		NewTestSystem(7, &__systemGroup_Test_C_6{}),
		NewTestSystem(8, &__systemGroup_Test_C_1{}, &__systemGroup_Test_C_5{}),
		NewTestSystem(9, &__systemGroup_Test_C_4{}, &__systemGroup_Test_C_6{}),
		NewTestSystem(10, &__systemGroup_Test_C_7{}, &__systemGroup_Test_C_5{}),
		NewTestSystem(11, &ReadOnly[__systemGroup_Test_C_1]{}),
	}
	sg := NewSystemGroup()
	for _, test := range tests {
		sg.insert(test)
	}

	sg.reset()

	for ss := sg.next(); len(ss) > 0; ss = sg.next() {
		println("========== batch:")
		for _, s := range ss {
			fmt.Printf("%s\n", ObjectToString(s))
		}
	}

}
