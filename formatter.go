package clilog

import (
	"fmt"
	"strings"
	"text/template"
)

// Internal representation of the configuration of a logger.
// It's not intended that users will be creating their own loggers (see "Non Goals"),
// this is split out mostly for modularity, and to aid in testing.
// All logging functionality goes through the globalLogger
type formatter struct {
	tmpl *template.Template
}

type logTemplateData struct {
	LevelCode string
	LevelName string
	Time      string
	Message   string
}

func newFormatter(format string) (formatter, error) {
	t, err := template.New("log").Parse(format)
	if err != nil {
		return formatter{}, err
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
		return formatter{}, fmt.Errorf("invalid template: %w", err)
	}
	return formatter{
		tmpl: t,
	}, nil
}

func mustNewFormatter(format string) formatter {
	f, err := newFormatter(format)
	if err != nil {
		panic(err)
	}
	return f
}

func (f formatter) format(data logTemplateData) string {
	var b strings.Builder
	_ = f.tmpl.Execute(&b, data)
	return b.String()
}
