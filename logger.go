package clilog

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"
)

// Internal representation of the configuration of a logger.
// It's not intended that users will be creating their own loggers (see "Non Goals"),
// this is split out mostly for modularity, and to aid in testing.
// All logging functionality goes through the globalLogger
type logger struct {
	level        Level
	colorEnabled bool
	timeFormat   string
	tmpl         *template.Template
	output       io.Writer
}

type logTemplateData struct {
	LevelCode string
	LevelName string
	Time      string
	Message   string
}

func (l Level) code() string {
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

func (l Level) name() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO "
	case LevelWarn:
		return "WARN "
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "INVAL" // fallback, same width
	}
}
func (l *logger) setFormat(format string) error {
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

	l.tmpl = t
	return nil
}

// --- Core logging function ---

func (l logger) logf(level Level, msg string) {
	if level < l.level {
		return
	}

	ts := time.Now().Format(l.timeFormat)

	levelCode := level.code()
	levelName := level.name()

	if l.colorEnabled {
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
	_ = l.tmpl.Execute(&b, data)
	fmt.Fprintln(l.output, b.String())
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
