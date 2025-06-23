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

func (l logger) logf(level Level, msg string) {
	if level < l.level {
		return
	}

	ts := now()

	data := logTemplateData{
		Level:   level,
		Time:    ts,
		Message: msg,
	}
	fmt.Fprintln(l.output, l.formatter.format(data))
}
