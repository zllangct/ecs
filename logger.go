package ecs

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})

	Debugf(fmt string, v ...interface{})
	Infof(fmt string, v ...interface{})
	Errorf(fmt string, v ...interface{})
	Fatalf(fmt string, v ...interface{})
}

var Log Logger = NewStdLog()

type StdLogLevel uint8

const (
	StdLogLevelDebug StdLogLevel = iota
	StdLogLevelInfo
	StdLogLevelError
	StdLogLevelFatal
	StdLogLevelNoPrint
)

type StdLog struct {
	logger *log.Logger
	level  StdLogLevel
}

func NewStdLog(level ...StdLogLevel) *StdLog {
	l := StdLogLevelDebug
	if len(level) > 0 {
		l = level[0]
	}
	return &StdLog{
		level:  l,
		logger: log.New(os.Stdout, "", log.Lshortfile),
	}
}

func (p StdLog) Debug(v ...interface{}) {
	p.logger.Output(2, fmt.Sprintf("[DEBUG][%d] %s", goroutineID(), fmt.Sprint(v...)))
}

func (p StdLog) Debugf(format string, v ...interface{}) {
	p.logger.Output(2, fmt.Sprintf("[DEBUG][%d] %s", goroutineID(), fmt.Sprintf(format, v...)))
}

func (p StdLog) Info(v ...interface{}) {
	if p.level > StdLogLevelInfo {
		return
	}
	p.logger.Output(2, fmt.Sprint(v...))
}

func (p StdLog) Infof(format string, v ...interface{}) {
	if p.level > StdLogLevelInfo {
		return
	}
	p.logger.Output(2, fmt.Sprintf(format, v...))
}

func (p StdLog) Error(v ...interface{}) {
	if p.level > StdLogLevelError {
		return
	}
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
	if p.level > StdLogLevelError {
		return
	}
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
	if p.level > StdLogLevelFatal {
		return
	}
	p.Error(v...)
	os.Exit(1)
}

func (p StdLog) Fatalf(format string, v ...interface{}) {
	if p.level > StdLogLevelFatal {
		return
	}
	p.Errorf(format, v...)
	os.Exit(1)
}
