package ecs

import (
	"errors"
	"math/rand"
	"reflect"
	"runtime/debug"
	"sync/atomic"
	"time"
)

var seq uint32
var timestamp uint64

func init() {
	rand.Seed(time.Now().UnixNano())
}

func LocalUniqueID() uint64 {
	tNow := uint64(time.Now().UnixNano()) << 32
	tTemp := atomic.LoadUint64(&timestamp)
	if tTemp != tNow {
		atomic.StoreUint32(&seq, 0)
		for {
			if atomic.CompareAndSwapUint64(&timestamp, tTemp, tNow) {
				break
			} else {
				tTemp = atomic.LoadUint64(&timestamp)
				tNow = uint64(time.Now().UnixNano()) << 32
			}
		}
	}
	s := atomic.AddUint32(&seq, 1)
	return tNow + uint64((s<<16)&0xFFFF0000+rand.Uint32()&0x0000FFFF)
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
