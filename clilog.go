package clilog

import (
	"fmt"
	"os"
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
	DefaultTemplate        = `{{ .LevelCode }} {{ .Time | timef "2006/01/02 15:04:05.000" | color .Level }} {{ .Message }}`
	DefaultTemplateNoColor = `{{ .LevelCode }} {{ .Time | timef "2006/01/02 15:04:05.000" }} {{ .Message }}`
)

var globalLogger = logger{
	level:     LevelInfo,
	formatter: mustNewFormatter(DefaultTemplate),
	output:    os.Stderr,
}

// --- Configuration ---

func SetLogLevel(level Level) {
	globalLogger.level = level
}

func SetFormat(format string) error {
	f, err := newFormatter(format)
	if err != nil {
		return err
	}
	globalLogger.formatter = f
	return nil
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
