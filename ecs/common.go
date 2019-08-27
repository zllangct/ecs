package ecs

import (
	"errors"
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

func ConcurrentTest(fns ... func())  {
	CheckError()

	for  _,fn := range fns {
		for i := 0; i < 10; i++ {
			go func() {
				for {
					fn()
				}
			}()
		}
	}
}