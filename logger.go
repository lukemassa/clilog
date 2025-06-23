package clilog

import (
	"fmt"
	"io"
	"time"
)

// So we can override in tests
var now = time.Now

// Internal representation of the configuration of a logger.
// It's not intended that users will be creating their own loggers (see "Non Goals"),
// this is split out mostly for modularity, and to aid in testing.
// All logging functionality goes through the globalLogger
type logger struct {
	level     Level
	formatter formatter
	output    io.Writer
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

// --- Core logging function ---

func (l logger) logf(level Level, msg string) {
	if level < l.level {
		return
	}

	ts := now()

	levelCode := level.code()
	levelName := level.name()

	data := logTemplateData{
		LevelCode: levelCode,
		LevelName: levelName,
		Level:     level,
		Time:      ts,
		Message:   msg,
	}
	fmt.Fprintln(l.output, l.formatter.format(data))
}
