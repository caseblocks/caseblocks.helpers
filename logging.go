package helpers

import (
	"io"
	"log"
	"os"
)

type Log interface {
	Error(message string)
	Info(message string)
	Debug(message string)
	Panic(message string)
	PanicErr(err error)
}

const (
	LERROR = 1
	LINFO  = 2
	LDEBUG = 3
)

type DefaultLog struct {
	Log *log.Logger
}

func (l *DefaultLog) Error(message string) {
	l.Log.Printf("ERROR-%s\n", message)
}
func (l *DefaultLog) Info(message string) {
	l.Log.Printf("INFO -%s\n", message)
}
func (l *DefaultLog) Debug(message string) {
	l.Log.Printf("DEBUG-%s\n", message)
}
func (l *DefaultLog) Panic(message string) {
	l.Log.Panicf("PANIC-%s\n", message)
}
func (l *DefaultLog) PanicErr(err error) {
	l.Log.Panicf("PANIC-%s\n", err)
}

func NewConsoleLogger() Log {
	return NewDefaultLogger(os.Stdout)
}

func NewDefaultLogger(out io.Writer) Log {
	logger := DefaultLog{Log: log.New(out, "", log.LstdFlags)}
	return &logger
}
