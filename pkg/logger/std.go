package logger

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	debug *log.Logger
	error *log.Logger
}

func NewLogger() Logger {
	logFlags := log.Ldate | log.Ltime | log.Lshortfile

	return Logger{
		info:  log.New(os.Stdout, "INFO: ", logFlags),
		debug: log.New(os.Stdout, "DEBUG: ", logFlags),
		error: log.New(os.Stdout, "ERROR: ", logFlags),
	}
}

var DefaultLogger = NewLogger()

func (l Logger) Info(v ...interface{}) {
	l.info.Println(v...)
}

func (l Logger) Debug(v ...interface{}) {
	l.debug.Println(v...)
}

func (l Logger) Error(v ...interface{}) {
	l.error.Println(v...)
}

func (l Logger) Infof(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

func (l Logger) Debugf(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
}

func (l Logger) Errorf(format string, v ...interface{}) {
	l.error.Printf(format, v...)
}
