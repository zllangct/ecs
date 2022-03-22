package ecs

import (
	"errors"
	"math/rand"
	"reflect"
	"runtime/debug"
	"sync/atomic"
	"time"
	"unsafe"
)

var Empty struct{} = struct{}{}

var seq uint32
var timestamp int64

func init() {
	rand.Seed(time.Now().UnixNano())
}

func UniqueID() int64 {
	tNow := int64(time.Now().UnixNano()) << 32
	tTemp := atomic.LoadInt64(&timestamp)
	if tTemp != tNow {
		atomic.StoreUint32(&seq, 0)
		for {
			if atomic.CompareAndSwapInt64(&timestamp, tTemp, tNow) {
				break
			} else {
				tTemp = atomic.LoadInt64(&timestamp)
				tNow = int64(time.Now().UnixNano()) << 32
			}
		}
	}
	s := atomic.AddUint32(&seq, 1)
	return tNow + int64((s<<16)&0xFFFF0000+rand.Uint32()&0x0000FFFF)
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

func TypeOf[T any]() reflect.Type {
	ins := (*T)(nil)
	return reflect.TypeOf(ins).Elem()
}

func memcmp(a unsafe.Pointer, b unsafe.Pointer, len uintptr) (ret bool) {
	for i := uintptr(0); i < len; i++ {
		if *(*byte)(unsafe.Pointer(uintptr(a) + i)) != *(*byte)(unsafe.Pointer(uintptr(b) + i)) {
			ret = false
			return
		}
	}
	ret = true
	Log.Infof("memory compare: %v, %v, %v, equal:%v", a, b, len, ret)
	return
}
