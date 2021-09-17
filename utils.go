package ecs

import "errors"

func Try(task func()) {
	defer func() {
		r := recover()
		if r != nil {
		}
	}()
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
