package main

import (
	"reflect"
	"testing"
)

type Com1 struct {
	ComponentBase
}
type Com2 struct {
	ComponentBase
}
type Com3 struct {
	ComponentBase
}
type Com4 struct {
	ComponentBase
}
type Com5 struct {
	ComponentBase
}
type Com6 struct {
	ComponentBase
}
type Com7 struct {
	ComponentBase
}
type Com8 struct {
	ComponentBase
}
type Com9 struct {
	ComponentBase
}
type Com10 struct {
	ComponentBase
}

type TestSystem struct {
	SystemBase
	ID int
}

func NewTestSystem(ID int,rqs ... reflect.Type) *TestSystem {
	s := &TestSystem{ID: ID}
	s.SetRequirements(rqs...)
	return s
}

func (p *TestSystem)Call(label int)interface{}  {
	switch label {
	case 1:
		println(p.ID)
	}
	return nil
}

func (t TestSystem) Filter(*Entity) {

}

func (t TestSystem) SystemUpdate() {

}

func TestNewSystemGroup(t *testing.T) {
	tests:=[]ISystem{
		NewTestSystem(1,reflect.TypeOf(&Com1{}),reflect.TypeOf(&Com2{})),
		NewTestSystem(2,reflect.TypeOf(&Com1{}),reflect.TypeOf(&Com3{})),
		NewTestSystem(3,reflect.TypeOf(&Com2{}),reflect.TypeOf(&Com5{})),
		NewTestSystem(4,reflect.TypeOf(&Com2{}),reflect.TypeOf(&Com3{}),reflect.TypeOf(&Com6{})),
		NewTestSystem(5,reflect.TypeOf(&Com7{})),
		NewTestSystem(6,reflect.TypeOf(&Com9{}),reflect.TypeOf(&Com10{})),
		NewTestSystem(7,reflect.TypeOf(&Com6{})),
	}
	sg:=NewSystemGroup()
	for _, test := range tests {
		sg.insert(test)
	}

	sg.iterInit()

	for ss:=sg.pop(); len(ss) >0 ;ss=sg.pop() {
		println("========== batch:")
		for _, s := range ss {
			s.Call(1)
		}
	}
}