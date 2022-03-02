package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var global *zap.SugaredLogger
var defaultLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)

func Logger() *zap.SugaredLogger {
	return global
}

func init() {
	log, err := zap.NewDevelopment(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.FatalLevel),
		zap.IncreaseLevel(defaultLevel),
	)
	if err != nil {
		panic(err)
	}
	global = log.Sugar()
}

// SetLevel sets level for global logger
func SetLevel(l zapcore.Level) {
	defaultLevel.SetLevel(l)
}

func Debug(args ...interface{}) {
	global.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	global.Debugf(format, args...)
}

func DebugKV(message string, kvs ...interface{}) {
	global.Debugw(message, kvs...)
}

func Info(args ...interface{}) {
	global.Info(args...)
}

func Infof(format string, args ...interface{}) {
	global.Infof(format, args...)
}

func InfoKV(message string, kvs ...interface{}) {
	global.Infow(message, kvs...)
}

func Warn(args ...interface{}) {
	global.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	global.Warnf(format, args...)
}

func WarnKV(message string, kvs ...interface{}) {
	global.Warnw(message, kvs...)
}

func Error(args ...interface{}) {
	global.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	global.Errorf(format, args...)
}

func ErrorKV(message string, kvs ...interface{}) {
	global.Errorw(message, kvs...)
}

func Fatal(args ...interface{}) {
	global.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	global.Fatalf(format, args...)
}
