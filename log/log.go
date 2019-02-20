package log

import (
	"github.com/astaxie/beego/logs"
)

var (
	Log *logs.BeeLogger
)

func init() {
	Log = logs.NewLogger()
	Log.SetLogger(logs.AdapterConsole)
}

func Info(format string, v ...interface{}) {
	Log.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	Log.Debug(format, v...)
}

func Warn(format string, v ...interface{}) {
	Log.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	Log.Error(format, v...)
}

// wrap my logger
type Logger struct {
	xlog *logs.BeeLogger
}

func NewLogger() *Logger {
	l := logs.NewLogger()
	l.SetLogger(logs.AdapterConsole)

	return &Logger{
		xlog: l,
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.xlog.Info(format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.xlog.Debug(format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.xlog.Warn(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.xlog.Error(format, v...)
}
