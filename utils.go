package ecs

import "errors"

func Try(task func()) (err error) {
	defer func() {
		r := recover()
		switch r.(type) {
		case error:
			err = r.(error)
		case string:
			err = errors.New(r.(string))
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
