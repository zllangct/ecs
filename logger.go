package ecs

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type IInternalLogger interface {
	Info(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}

type StdLogger struct {
	logger *log.Logger
}

func NewStdLogger() *StdLogger {
	return &StdLogger{
		logger: log.New(os.Stdout, "", log.Lshortfile),
	}
}

func (p StdLogger) Info(v ...interface{}) {
	p.logger.Output(2, fmt.Sprint(v...))
}

func (p StdLogger) Error(v ...interface{}) {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	s := fmt.Sprint(append(v, "\n", string(buf))...)
	p.logger.Output(2, s)
}

func (p StdLogger) Fatal(v ...interface{}) {
	p.Error(v...)
	os.Exit(1)
}
