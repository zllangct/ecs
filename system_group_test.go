package ecs

import (
	"testing"
)

type Com1 struct {
	Component[Com1]
}
type Com2 struct {
	Component[Com2]
}
type Com3 struct {
	Component[Com3]
}
type Com4 struct {
	Component[Com4]
}
type Com5 struct {
	Component[Com5]
}
type Com6 struct {
	Component[Com6]
}
type Com7 struct {
	Component[Com7]
}
type Com8 struct {
	Component[Com8]
}
type Com9 struct {
	Component[Com9]
}
type Com10 struct {
	Component[Com10]
}

type TestSystem struct {
	System[TestSystem]
	ID int
}

func NewTestSystem(ID int, rqs ...IComponentTemplate) *TestSystem {
	s := &TestSystem{ID: ID}
	s.SetRequirements(rqs...)
	return s
}

func (p *TestSystem) Call(label int) interface{} {
	switch label {
	case 1:
		println(p.ID)
	}
	return nil
}

func (p *TestSystem) Filter(component IComponent, op CollectionOperate) {

}

func TestNewSystemGroup(t *testing.T) {
	tests := []ISystem{
		NewTestSystem(1, &Com1{}, &Com2{}),
		NewTestSystem(2, &Com1{}, &Com3{}),
		NewTestSystem(3, &Com2{}, &Com5{}),
		NewTestSystem(4, &Com2{}, &Com3{}, &Com6{}),
		NewTestSystem(5, &Com7{}),
		NewTestSystem(6, &Com9{}, &Com10{}),
		NewTestSystem(7, &Com6{}),
		NewTestSystem(8, &Com1{}, &Com5{}),
		NewTestSystem(9, &Com4{}, &Com6{}),
		NewTestSystem(10, &Com7{}, &Com5{}),
	}
	sg := NewSystemGroup()
	for _, test := range tests {
		sg.insert(test)
	}

	sg.reset()

	for ss := sg.next(); len(ss) > 0; ss = sg.next() {
		println("========== batch:")
		for _, s := range ss {
			s.Call(1)
		}
	}
}
