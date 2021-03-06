package ecs

import (
	"runtime"
	"time"
)

type RuntimeConfig struct {
	Debug           bool          //Debug模式
	CpuNum          int           //使用的最大cpu数量
	HashCount       int           //容器散列数量
	FrameInterval   time.Duration //帧间隔
	MaxPoolThread   int           //线程池最大线程数量
	MaxPoolJobQueue int           //线程池最大任务队列长度
}

func NewDefaultRuntimeConfig() *RuntimeConfig {
	return &RuntimeConfig{
		Debug:         true,
		CpuNum:        runtime.NumCPU(),
		HashCount:     runtime.NumCPU() * 4,
		FrameInterval: time.Millisecond * 33,
	}
}
