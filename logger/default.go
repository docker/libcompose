package logger

import (
	"log"
)

type DefaultLogger struct {
}

func (l *DefaultLogger) Errorf(format string, values ...interface{}) {
	log.Printf("err:"+format, values...)
}

func (l *DefaultLogger) Error(values ...interface{}) {
	log.Print("err:", values)
}

func (l *DefaultLogger) Infof(format string, values ...interface{}) {
	log.Printf("info:"+format, values...)
}

func (l *DefaultLogger) Info(values ...interface{}) {
	log.Print("info:", values)
}

func (l *DefaultLogger) Debugf(format string, values ...interface{}) {
	log.Printf("debug:"+format, values...)
}

func (l *DefaultLogger) Debug(values ...interface{}) {
	log.Print("debug:", values)
}

func (l *DefaultLogger) Warnf(format string, values ...interface{}) {
	log.Printf("warn:"+format, values...)
}

func (l *DefaultLogger) Warn(values ...interface{}) {
	log.Print("warn:", values)
}
