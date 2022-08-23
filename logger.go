package ecs

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type Logger interface {
	Info(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})

	Infof(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Fatalf(fmt string, args ...interface{})
}

var Log Logger = NewStdLog()

type StdLog struct {
	logger *log.Logger
}

func NewStdLog() *StdLog {
	return &StdLog{
		logger: log.New(os.Stdout, "", log.Lshortfile),
	}
}

func (p StdLog) Info(v ...interface{}) {
	p.logger.Output(2, fmt.Sprint(v...))
}

func (p StdLog) Infof(format string, v ...interface{}) {
	p.logger.Output(2, fmt.Sprintf(format, v...))
}

func (p StdLog) Error(v ...interface{}) {
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

func (p StdLog) Errorf(format string, v ...interface{}) {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	s := fmt.Sprint(fmt.Sprintf(format, v...), "\n", string(buf))
	p.logger.Output(2, s)
}

func (p StdLog) Fatal(v ...interface{}) {
	p.Error(v...)
	os.Exit(1)
}

func (p StdLog) Fatalf(format string, v ...interface{}) {
	p.Errorf(format, v...)
	os.Exit(1)
}
