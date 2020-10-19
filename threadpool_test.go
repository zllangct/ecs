package ecs

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func init() {
	println("using MAXPROC")
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
}

//TestNewPool test new goroutine pool
func TestNewPool(t *testing.T) {
	pool := NewPool(nil, 1000, 10000)
	defer pool.Release()

	iterations := 1000000
	var counter uint64

	wg := sync.WaitGroup{}
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		arg := uint64(1)
		pool.AddJob(func(ctx *JobContext, args ...interface{}) {
			defer wg.Done()
			atomic.AddUint64(&counter, arg)
		}, nil)
	}
	wg.Wait()

	counterFinal := atomic.LoadUint64(&counter)
	if uint64(iterations) != counterFinal {
		t.Errorf("iterations %v is not equal counterFinal %v", iterations, counterFinal)
	}
}
