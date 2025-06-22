package clilog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) Code() string {
	switch l {
	case LevelDebug:
		return "D"
	case LevelInfo:
		return "I"
	case LevelWarn:
		return "W"
	case LevelError:
		return "E"
	case LevelFatal:
		return "F"
	default:
		return "?"
	}
}

func (l Level) Name() string {
	switch l {
	case LevelDebug:
		return "DEBUG "
	case LevelInfo:
		return "INFO  "
	case LevelWarn:
		return "WARN  "
	case LevelError:
		return "ERROR "
	case LevelFatal:
		return "FATAL "
	default:
		return "INVAL " // fallback, same width
	}
}

var (
	currentLevel           = LevelInfo
	colorEnabled           = true
	timeFormat             = "2006/01/02 15:04:05.000"
	tmpl                   = template.Must(template.New("log").Parse(defaultTemplate))
	output       io.Writer = os.Stderr
)

const defaultTemplate = `{{ .Time }} {{ .LevelCode }} {{ .Message }}`

type logTemplateData struct {
	LevelCode string
	LevelName string
	Time      string
	Message   string
}

func SetLogLevel(level Level) {
	currentLevel = level
}

func SetTimestampFormat(format string) {
	timeFormat = format
}

func SetFormat(format string) error {
	t, err := template.New("log").Parse(format)
	if err != nil {
		return err
	}

	// Validate template by rendering with dummy data
	test := logTemplateData{
		LevelCode: "D",
		LevelName: "DEBUG ",
		Time:      "2006/01/02 15:04:05",
		Message:   "test message",
	}

	var b strings.Builder
	if err := t.Execute(&b, test); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	tmpl = t
	return nil
}

func SetDisableColor(disable bool) {
	colorEnabled = !disable
}

func SetOutput(w io.Writer) {
	if w != nil {
		output = w
	}
}

// --- Core logging function ---

func logf(level Level, msg string) {
	if level < currentLevel {
		return
	}

	ts := time.Now().Format(timeFormat)

	levelCode := level.Code()
	levelName := level.Name()

	if colorEnabled {
		color := colorFor(level)
		levelCode = color + levelCode + "\033[0m"
		levelName = color + levelName + "\033[0m"
		ts = color + ts + "\033[0m"
	}

	data := logTemplateData{
		LevelCode: levelCode,
		LevelName: levelName,
		Time:      ts,
		Message:   msg,
	}

	var b strings.Builder
	_ = tmpl.Execute(&b, data)
	fmt.Fprintln(output, b.String())
}

func colorFor(level Level) string {
	switch level {
	case LevelDebug:
		return "\033[36m" // cyan
	case LevelInfo:
		return "\033[32m" // green
	case LevelWarn:
		return "\033[33m" // yellow
	case LevelError:
		return "\033[31m" // red
	case LevelFatal:
		return "\033[35m" // magenta
	default:
		return ""
	}
}

// --- Public helpers ---

func Debug(args ...any)                 { logf(LevelDebug, fmt.Sprint(args...)) }
func Debugf(format string, args ...any) { logf(LevelDebug, fmt.Sprintf(format, args...)) }

func Info(args ...any)                 { logf(LevelInfo, fmt.Sprint(args...)) }
func Infof(format string, args ...any) { logf(LevelInfo, fmt.Sprintf(format, args...)) }

func Warn(args ...any)                 { logf(LevelWarn, fmt.Sprint(args...)) }
func Warnf(format string, args ...any) { logf(LevelWarn, fmt.Sprintf(format, args...)) }

func Error(args ...any)                 { logf(LevelError, fmt.Sprint(args...)) }
func Errorf(format string, args ...any) { logf(LevelError, fmt.Sprintf(format, args...)) }

func Fatal(args ...any)                 { logf(LevelFatal, fmt.Sprint(args...)); os.Exit(1) }
func Fatalf(format string, args ...any) { logf(LevelFatal, fmt.Sprintf(format, args...)); os.Exit(1) }

