package main

import (
	"errors"
	_ "net/http/pprof"
	"runtime/debug"
)

func CheckError() {
	if r := recover(); r != nil {
		var str string
		switch r.(type) {
		case error:
			str = r.(error).Error()
		case string:
			str = r.(string)
		}
		err := errors.New(str + "\n" + string(debug.Stack()))
		println(err.Error())
	}
}

func main() {


	tets()

}
type IT interface {
	GetName()string
}

type T struct {
	Name string
	x1 int
	x2 int
	x3 int
	x4 int
	x5 int
}

func (p *T)GetName()string  {
	return p.Name
}

type T2 struct {
	Name int
}

func tets()  {

}






