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
}

func (p *T)GetName()string  {
	return p.Name
}

func tets()  {
	i:=interface{}(&T{Name:"zhaolei"})
	tt:= i.(IT)

	println(tt.GetName())
}






