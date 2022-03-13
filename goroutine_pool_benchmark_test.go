package ecs

import (
	"sync"
	"testing"
	"time"
)

const (
	runTimes  = 1000
	poolSize  = 50
	queueSize = 50
)

func demoTask() {
	time.Sleep(time.Nanosecond * 10)
}

//BenchmarkGoroutine benchmark the goroutine doing tasks.
func BenchmarkGoroutine(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(runTimes)

		for j := 0; j < runTimes; j++ {
			go func() {
				defer wg.Done()
				demoTask()
			}()
		}

		wg.Wait()
	}
}

//BenchmarkGpool benchmarks the goroutine pool.
func BenchmarkGpool(b *testing.B) {
	pool := NewPool(poolSize, queueSize)
	pool.Start()

	defer pool.Release()
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(runTimes)

		for j := 0; j < runTimes; j++ {
			pool.Add(func() {
				defer wg.Done()
				demoTask()
			})
		}

		wg.Wait()
	}
}
