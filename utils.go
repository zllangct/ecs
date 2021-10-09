package ecs

import (
	"errors"
	"runtime/debug"
)

func Try(task func(), catch ...func(error)) {
	defer (func() {
		if r := recover(); r != nil {
			var str string
			switch r.(type) {
			case error:
				str = r.(error).Error()
			case string:
				str = r.(string)
			}
			err := errors.New(str + "\n" + string(debug.Stack()))
			if len(catch) > 0 {
				catch[0](err)
			}
		}
	})()
	task()
}

func TryAndReport(task func()) (err error) {
	defer func() {
		r := recover()
		switch typ := r.(type) {
		case error:
			err = r.(error)
		case string:
			err = errors.New(r.(string))
		default:
			_ = typ
		}
	}()
	task()
	return nil
}

func StrHash(str string, groupCount int) int {
	total := 0
	for i := 0; i < len(str); i++ {
		total += int(str[i])
	}
	return total % groupCount
}
