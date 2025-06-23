package clilog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFormatter(t *testing.T) {
	cases := []struct {
		format      string
		expectedErr string
	}{
		{
			format:      "foo{{",
			expectedErr: "unclosed action",
		},
		{
			format:      "{{ .Foobar }}",
			expectedErr: "can't evaluate field Foobar",
		},
		{
			format:      "{{ .Time | foobar }}",
			expectedErr: "function \"foobar\" not defined",
		},
		{
			format: "{{ .Time }}",
		},
		{
			format: `{{ .Time | timef "2006" }}`,
		},
		{
			format:      `{{ .Time | timef 123 }}`,
			expectedErr: "expected string; found 123",
		},
	}
	for _, tc := range cases {
		t.Run(tc.format, func(t *testing.T) {
			_, actualErr := newFormatter(tc.format)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.ErrorContains(t, actualErr, tc.expectedErr)
			}
		})
	}
}

func TestFormatter(t *testing.T) {

	jan1 := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)
	cases := []struct {
		description    string
		format         string
		level          Level
		time           time.Time
		message        string
		expectedOutput string
	}{
		{
			description:    "Basic log",
			format:         `{{ .Level | abbrev }} {{ .Time | timef "2006" }} {{ .Message }}`,
			level:          LevelInfo,
			time:           jan1,
			message:        "Hello!",
			expectedOutput: "I 2025 Hello!",
		},
		{
			description:    "Log in color",
			format:         `{{ .Level | abbrev }} {{ .Time | timef "2006" | color .Level }} {{ .Message }}`,
			level:          LevelInfo,
			time:           jan1,
			message:        "Hello!",
			expectedOutput: "I \x1b[32m2025\x1b[0m Hello!",
		},
		{
			description:    "Log with level name",
			format:         `{{ .Level }} {{ .Message }}`,
			level:          LevelInfo,
			time:           jan1,
			message:        "Hello!",
			expectedOutput: "INFO  Hello!",
		},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {

			formatter := mustNewFormatter(tc.format)

			actualOutput := formatter.format(logTemplateData{
				Level:   tc.level,
				Message: tc.message,
				Time:    tc.time,
			})

			assert.Equal(t, tc.expectedOutput, actualOutput)
		})
	}

}
