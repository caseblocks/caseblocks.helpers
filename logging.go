package helpers

import (
  "log"
  "io"
)

type Log interface {
  Error(message string)
  Info(message string)
  Debug(message string)
}

const (
  LERROR = 1
  LINFO = 2
  LDEBUG =3
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

func NewDefaultLogger(out io.Writer) Log {
  logger := DefaultLog{Log: log.New(out, "", log.LstdFlags)}
  return &logger
}