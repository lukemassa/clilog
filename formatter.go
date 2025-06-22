package clilog

import (
	"fmt"
	"strings"
	"text/template"
	"time"
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
	Time      time.Time
	Message   string
}

func newFormatter(format string) (formatter, error) {

	funcMap := template.FuncMap{
		"timef": func(layout string, t time.Time) string {
			return t.Format(layout)
		},
	}

	t, err := template.New("log").Funcs(funcMap).Parse(format)
	if err != nil {
		return formatter{}, err
	}

	// Validate template by rendering with dummy data
	test := logTemplateData{
		LevelCode: "D",
		LevelName: "DEBUG ",
		Time:      now(),
		Message:   "test message",
	}

	var b strings.Builder
	if err := t.Execute(&b, test); err != nil {
		return formatter{}, fmt.Errorf("invalid template: %w", err)
	}

	return formatter{
		tmpl: t.Funcs(funcMap),
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
