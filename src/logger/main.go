package logger

import (
	"io"
)

func Executablef(format string, args ...interface{}) string {
	if !IsReady() {
		return ""
	}
	return globalLogger.Executablef(format, args...)
}

func Tagf(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.TagSkipf(1, format, args...)
}

func Debugf(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Panicf(format, args...)
}

func Tag(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.TagSkip(1, args...)
}

func Debug(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Debug(args...)
}

func Info(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Info(args...)
}

func Warn(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Warn(args...)
}

func Error(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Error(args...)
}

func Panic(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Panic(args...)
}

func TagWrite(msg string) {
	if !IsReady() {
		return
	}
	globalLogger.TagSkip(1, msg)
}

func DebugWrite(msg string) {
	if !IsReady() {
		return
	}
	globalLogger.DebugWrite(msg)
}

func InfoWrite(msg string) {
	if !IsReady() {
		return
	}
	globalLogger.InfoWrite(msg)
}

func WarnWrite(msg string) {
	if !IsReady() {
		return
	}
	globalLogger.WarnWrite(msg)
}

func ErrorWrite(msg string) {
	if !IsReady() {
		return
	}
	globalLogger.ErrorWrite(msg)
}

func PanicWrite(msg string) {
	if !IsReady() {
		return
	}
	globalLogger.PanicWrite(msg)
}

func GetDebugWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.GetDebugWriter()
}

func GetInfoWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.GetInfoWriter()
}

func GetWarningWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.GetWarningWriter()
}

func GetTagWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.GetTagWriter()
}

func GetErrorWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.GetErrorWriter()
}

func GetPanicWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.GetPanicWriter()
}

func IsDebugTerm() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsDebugTerm()
}

func IsInfoTerm() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsDebugTerm()
}

func IsTagTerm() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsTagTerm()
}

func IsWarnTerm() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsWarnTerm()
}

func IsErrorTerm() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsErrorTerm()
}

func IsPanicTerm() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsPanicTerm()
}

func IsDebugTermNotDumb() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsDebugTerm()
}

func IsInfoTermNotDumb() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsInfoTermNotDumb()
}

func IsTagTermNotDumb() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsTagTermNotDumb()
}

func IsWarnTermNotDumb() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsWarnTermNotDumb()
}

func IsErrorTermNotDumb() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsErrorTermNotDumb()
}

func IsPanicTermNotDumb() bool {
	if !IsReady() {
		return false
	}
	return globalLogger.IsPanicTermNotDumb()
}
