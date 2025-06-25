package clilog

import (
	"fmt"
	"io"
	"os"
	"time"
)

// So we can override in tests
var (
	now              = time.Now
	output io.Writer = os.Stderr
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	DefaultFormat        = `{{ .Level | abbrev }} {{ .Time | timef "2006/01/02 15:04:05.000" | color .Level }} {{ .Message }}`
	DefaultFormatNoColor = `{{ .Level | abbrev }} {{ .Time | timef "2006/01/02 15:04:05.000" }} {{ .Message }}`
)

// Internal representation of the configuration of a logger.
// It's not intended that users will be creating their own loggers (see "Non Goals"), it is split out mostly for modularity, and to aid in testing.
// All logging functionality goes through the globalLogger
type logger struct {
	level     Level
	formatter formatter
}

// Specific "logger" we use in all the public logging functions
var globalLogger = logger{
	level:     LevelInfo,
	formatter: mustNewFormatter(DefaultFormat),
}

func (l logger) logf(level Level, msg string) {
	if level < l.level {
		return
	}

	data := logTemplateData{
		Level:   level,
		Time:    now(),
		Message: msg,
	}
	fmt.Fprintln(output, l.formatter.format(data))
}

// --- Configuration ---

// SetLogLevel sets the log level, so only logs at or above this level will be displayed
func SetLogLevel(level Level) {
	globalLogger.level = level
}

// SetFormat sets the format of the log
func SetFormat(format string) error {
	f, err := newFormatter(format)
	if err != nil {
		return err
	}
	globalLogger.formatter = f
	return nil
}

func MustSetFormat(format string) {
	err := SetFormat(format)
	if err != nil {
		panic(fmt.Sprintf("Failed to set log format: %v", err))
	}
}

// --- Public logging helpers ---

// Debug logs a message at DEBUG level using fmt.Sprint.
func Debug(args ...any) { globalLogger.logf(LevelDebug, fmt.Sprint(args...)) }

// Debugf logs a formatted message at DEBUG level using fmt.Sprintf.
func Debugf(format string, args ...any) { globalLogger.logf(LevelDebug, fmt.Sprintf(format, args...)) }

// Info logs a message at INFO level using fmt.Sprint.
func Info(args ...any) { globalLogger.logf(LevelInfo, fmt.Sprint(args...)) }

// Infof logs a formatted message at INFO level using fmt.Sprintf.
func Infof(format string, args ...any) { globalLogger.logf(LevelInfo, fmt.Sprintf(format, args...)) }

// Warn logs a message at WARN level using fmt.Sprint.
func Warn(args ...any) { globalLogger.logf(LevelWarn, fmt.Sprint(args...)) }

// Warnf logs a formatted message at WARN level using fmt.Sprintf.
func Warnf(format string, args ...any) { globalLogger.logf(LevelWarn, fmt.Sprintf(format, args...)) }

// Error logs a message at ERROR level using fmt.Sprint.
func Error(args ...any) { globalLogger.logf(LevelError, fmt.Sprint(args...)) }

// Errorf logs a formatted message at ERROR level using fmt.Sprintf.
func Errorf(format string, args ...any) { globalLogger.logf(LevelError, fmt.Sprintf(format, args...)) }

// Fatal logs a message at FATAL level using fmt.Sprint and then exits with status code 1.
func Fatal(args ...any) {
	globalLogger.logf(LevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}

// Fatalf logs a formatted message at FATAL level using fmt.Sprintf and then exits with status code 1.
func Fatalf(format string, args ...any) {
	globalLogger.logf(LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}
