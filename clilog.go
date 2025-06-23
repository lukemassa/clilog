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

func Debug(args ...any)                 { globalLogger.logf(LevelDebug, fmt.Sprint(args...)) }
func Debugf(format string, args ...any) { globalLogger.logf(LevelDebug, fmt.Sprintf(format, args...)) }

func Info(args ...any)                 { globalLogger.logf(LevelInfo, fmt.Sprint(args...)) }
func Infof(format string, args ...any) { globalLogger.logf(LevelInfo, fmt.Sprintf(format, args...)) }

func Warn(args ...any)                 { globalLogger.logf(LevelWarn, fmt.Sprint(args...)) }
func Warnf(format string, args ...any) { globalLogger.logf(LevelWarn, fmt.Sprintf(format, args...)) }

func Error(args ...any)                 { globalLogger.logf(LevelError, fmt.Sprint(args...)) }
func Errorf(format string, args ...any) { globalLogger.logf(LevelError, fmt.Sprintf(format, args...)) }

func Fatal(args ...any) {
	globalLogger.logf(LevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	globalLogger.logf(LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}
