package clilog

import (
	"fmt"
	"io"
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

const defaultTemplate = `{{ .Time | timef "15:04:05" }} {{ .LevelCode }} {{ .Message }}`

var globalLogger = logger{
	level:        LevelInfo,
	colorEnabled: true,
	formatter:    mustNewFormatter(defaultTemplate),
	output:       os.Stderr,
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

func SetDisableColor(disable bool) {
	globalLogger.colorEnabled = !disable
}

func SetOutput(w io.Writer) {
	if w != nil {
		globalLogger.output = w
	}
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
