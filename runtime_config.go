package ecs

import (
	"runtime"
)

type RuntimeConfig struct {
	Debug           bool   //Debug模式
	CpuNum          int    //使用的最大cpu数量
	MaxPoolThread   uint32 //线程池最大线程数量
	MaxPoolJobQueue uint32 //线程池最大任务队列长度
	Logger          Logger
}

func NewDefaultRuntimeConfig() *RuntimeConfig {
	return &RuntimeConfig{
		Debug:           true,
		CpuNum:          runtime.NumCPU(),
		Logger:          NewStdLog(),
		MaxPoolThread:   uint32(runtime.NumCPU() * 4),
		MaxPoolJobQueue: 20,
	}
}
