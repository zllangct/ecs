package ecs

import (
	"errors"
	"reflect"
	"runtime/debug"
	"sync/atomic"
	"time"
)

var idSeq uint64

// LocalUniqueID 进程内唯一 ID：高 42 位毫秒时间戳（约 139 年回绕），
// 低 22 位原子序号（每毫秒 400 万容量），单调不减。
func LocalUniqueID() uint64 {
	ms := uint64(time.Now().UnixMilli())
	seq := atomic.AddUint64(&idSeq, 1) - 1
	return (ms << 22) | (seq & 0x3FFFFF)
}

func TypeOf[T any]() reflect.Type {
	ins := (*T)(nil)
	return reflect.TypeOf(ins).Elem()
}

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

func TryAndReport(task func() error) (err error) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		switch typ := r.(type) {
	//		case error:
	//			err = r.(error)
	//		case string:
	//			err = errors.New(r.(string))
	//		default:
	//			_ = typ
	//		}
	//	}
	//}()
	err = task()
	return err
}
