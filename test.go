package main

import (
	"errors"
	_ "net/http/pprof"
	"reflect"
	"runtime/debug"
	"time"
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

	defer CheckError()

	v:= reflect.ValueOf(tets)
	t1:=time.Now()
	for i := 0; i < 10000000; i++ {
		v.Call([]reflect.Value{reflect.ValueOf(1)})
	}
	elapse1 := time.Since(t1)

	t2:=time.Now()
	for i := 0; i < 10000000; i++ {
		tets(1)
	}
	elapse2 := time.Since(t2)
	println(elapse1.String())
	println(elapse2.String())


}

func tets(b int)  {
	//a:=1
	//_=a
}