package logger

import (
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/config"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	"github.com/mattn/go-isatty"
	"io"
	"os"
)

type LoggerLevel string

const (
	LevelDebug LoggerLevel = "debug"
	LevelInfo  LoggerLevel = "info"
	LevelWarn  LoggerLevel = "warn"
	LevelError LoggerLevel = "error"
	LevelPanic LoggerLevel = "panic"
	LevelNone  LoggerLevel = "none"
)

type loggerLevel int64

const (
	levelDebug loggerLevel = 1
	levelInfo  loggerLevel = 2
	levelWarn  loggerLevel = 3
	levelError loggerLevel = 4
	levelPanic loggerLevel = 5
	levelNone  loggerLevel = 6
)

var levelMap = map[LoggerLevel]loggerLevel{
	LevelDebug: levelDebug,
	LevelInfo:  levelInfo,
	LevelWarn:  levelWarn,
	LevelError: levelError,
	LevelPanic: levelPanic,
	LevelNone:  levelNone,
}

type Logger struct {
	level      LoggerLevel
	logLevel   loggerLevel
	logTag     bool
	warnWriter io.Writer
	errWriter  io.Writer
	args0      string
	args0Name  string
}

var globalLogger *Logger = nil
var DefaultWarnWriter = os.Stdout
var DefaultErrorWriter = os.Stderr

func InitLogger(warnWriter, errWriter io.Writer) error {
	if !config.IsReady() {
		panic("config is not ready")
	}

	level := LoggerLevel(config.GetConfig().GlobalConfig.LogLevel)
	logLevel, ok := levelMap[level]
	if !ok {
		return fmt.Errorf("invalid log level: %s", level)
	}

	if warnWriter == nil {
		warnWriter = DefaultWarnWriter
	}

	if errWriter == nil {
		errWriter = DefaultErrorWriter
	}

	logger := &Logger{
		level:      level,
		logLevel:   logLevel,
		logTag:     config.GetConfig().LogTag.ToBool(true),
		warnWriter: os.Stdout,
		errWriter:  os.Stderr,
		args0:      utils.GetArgs0(),
		args0Name:  utils.GetArgs0Name(),
	}

	globalLogger = logger
	return nil
}

func IsReady() bool {
	return globalLogger != nil
}

func (l *Logger) Executablef(format string, args ...interface{}) string {
	str := fmt.Sprintf(format, args...)
	if str == "" {
		_, _ = fmt.Fprintf(l.warnWriter, "[Executable]: %s\n", l.args0)
	} else {
		_, _ = fmt.Fprintf(l.warnWriter, "{Executable %s]: %s\n", l.args0, str)
	}
	return l.args0
}

func (l *Logger) Tagf(format string, args ...interface{}) {
	l.TagSkipf(1, format, args...)
}

func (l *Logger) TagSkipf(skip int, format string, args ...interface{}) {
	if !l.logTag {
		return
	}

	funcName, file, _, line := utils.GetCallingFunctionInfo(skip + 1)

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.warnWriter, "[Tag %s]: %s %s %s:%d\n", l.args0Name, str, funcName, file, line)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.logLevel > levelDebug {
		return
	}

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.warnWriter, "[Debug %s]: %s\n", l.args0Name, str)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.logLevel > levelInfo {
		return
	}

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.warnWriter, "[Info %s]: %s\n", l.args0Name, str)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.logLevel > levelWarn {
		return
	}

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.warnWriter, "[Warning %s]: %s\n", l.args0Name, str)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.logLevel > levelError {
		return
	}

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.errWriter, "[Error %s]: %s\n", l.args0Name, str)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	if l.logLevel > levelPanic {
		return
	}

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.errWriter, "[Panic %s]: %s\n", l.args0Name, str)
}

func (l *Logger) Tag(args ...interface{}) {
	l.TagSkip(1, args...)
}

func (l *Logger) TagSkip(skip int, args ...interface{}) {
	if !l.logTag {
		return
	}

	funcName, file, _, line := utils.GetCallingFunctionInfo(skip + 1)

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.warnWriter, "[Tag %s]: %s %s %s:%d\n", l.args0Name, str, funcName, file, line)
}

func (l *Logger) Debug(args ...interface{}) {
	if l.logLevel > levelDebug {
		return
	}

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.warnWriter, "[Debug %s]: %s\n", l.args0Name, str)
}

func (l *Logger) Info(args ...interface{}) {
	if l.logLevel > levelInfo {
		return
	}

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.warnWriter, "[Info %s]: %s\n", l.args0Name, str)
}

func (l *Logger) Warn(args ...interface{}) {
	if l.logLevel > levelWarn {
		return
	}

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.warnWriter, "[Warning %s]: %s\n", l.args0Name, str)
}

func (l *Logger) Error(args ...interface{}) {
	if l.logLevel > levelError {
		return
	}

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.errWriter, "[Error %s]: %s\n", l.args0Name, str)
}

func (l *Logger) Panic(args ...interface{}) {
	if l.logLevel > levelPanic {
		return
	}

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.errWriter, "[Panic %s]: %s\n", l.args0Name, str)
}

func (l *Logger) TagWrite(msg string) {
	l.TagSkipWrite(1, msg)
}

func (l *Logger) TagSkipWrite(skip int, msg string) {
	if !l.logTag {
		return
	}

	funcName, file, _, line := utils.GetCallingFunctionInfo(skip + 1)

	_, _ = fmt.Fprintf(l.warnWriter, "[Debug %s]: %s %s %s:%d\n", l.args0Name, msg, funcName, file, line)
}

func (l *Logger) DebugWrite(msg string) {
	if l.logLevel > levelDebug {
		return
	}

	_, _ = fmt.Fprintf(l.warnWriter, "[Debug %s]: %s\n", l.args0Name, msg)
}

func (l *Logger) InfoWrite(msg string) {
	if l.logLevel > levelInfo {
		return
	}

	_, _ = fmt.Fprintf(l.warnWriter, "[Info %s]: %s\n", l.args0Name, msg)
}

func (l *Logger) WarnWrite(msg string) {
	if l.logLevel > levelWarn {
		return
	}

	_, _ = fmt.Fprintf(l.warnWriter, "[Warning %s]: %s\n", l.args0Name, msg)
}

func (l *Logger) ErrorWrite(msg string) {
	if l.logLevel > levelError {
		return
	}

	_, _ = fmt.Fprintf(l.errWriter, "[Error %s]: %s\n", l.args0Name, msg)
}

func (l *Logger) PanicWrite(msg string) {
	if l.logLevel > levelPanic {
		return
	}

	_, _ = fmt.Fprintf(l.errWriter, "[Panic %s]: %s\n", l.args0Name, msg)
}

func (l *Logger) GetDebugWriter() io.Writer {
	return l.warnWriter
}

func (l *Logger) GetInfoWriter() io.Writer {
	return l.warnWriter
}

func (l *Logger) GetWarningWriter() io.Writer {
	return l.warnWriter
}

func (l *Logger) GetTagWriter() io.Writer {
	return l.warnWriter
}

func (l *Logger) GetErrorWriter() io.Writer {
	return l.errWriter
}

func (l *Logger) GetPanicWriter() io.Writer {
	return l.errWriter
}

func (l *Logger) isWarnWriterTerm() bool {
	w, ok := l.warnWriter.(*os.File)
	if !ok {
		return false
	} else if !isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()) { // 非终端
		return false
	}
	return true
}

func (l *Logger) isErrWriterTerm() bool {
	w, ok := l.errWriter.(*os.File)
	if !ok {
		return false
	} else if !isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()) { // 非终端
		return false
	}
	return true
}

func (l *Logger) isTermDump() bool {
	// TERM为dump表示终端为基础模式，不支持高级显示
	return os.Getenv("TERM") == "dumb"
}

func (l *Logger) IsDebugTerm() bool {
	return l.isWarnWriterTerm()
}

func (l *Logger) IsInfoTerm() bool {
	return l.isWarnWriterTerm()
}

func (l *Logger) IsWarnTerm() bool {
	return l.isWarnWriterTerm()
}

func (l *Logger) IsTagTerm() bool {
	return l.isWarnWriterTerm()
}

func (l *Logger) IsErrorTerm() bool {
	return l.isErrWriterTerm()
}

func (l *Logger) IsPanicTerm() bool {
	return l.isErrWriterTerm()
}

func (l *Logger) IsDebugTermNotDumb() bool {
	return l.isWarnWriterTerm() && !l.isTermDump()
}

func (l *Logger) IsInfoTermNotDumb() bool {
	return l.isWarnWriterTerm() && !l.isTermDump()
}

func (l *Logger) IsWarnTermNotDumb() bool {
	return l.isWarnWriterTerm() && !l.isTermDump()
}

func (l *Logger) IsTagTermNotDumb() bool {
	return l.isWarnWriterTerm() && !l.isTermDump()
}

func (l *Logger) IsErrorTermNotDumb() bool {
	return l.isErrWriterTerm() && !l.isTermDump()
}

func (l *Logger) IsPanicTermNotDumb() bool {
	return l.isErrWriterTerm() && !l.isTermDump()
}
