package ecs

import (
	"testing"
)

func TestUniqueID(t *testing.T) {
	m:=make(map[uint64]struct{})
	count := 0
	for i := 0; i<5000000;i++  {
		id:=UniqueID()
		if _,ok:=m[id];ok {
			count+=1
			println("repeat:",count)
		}else{
			m[id] = struct{}{}
		}
	}
}