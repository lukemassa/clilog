package clilog

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetFormat(t *testing.T) {
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

func TestLogging(t *testing.T) {
	originalNow := now
	defer func() { now = originalNow }() // Restore after test

	now = func() time.Time {
		return time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)
	}
	cases := []struct {
		description    string
		format         string
		currentLevel   Level
		messageLevel   Level
		message        string
		expectedOutput string
	}{
		{
			description:    "Basic log",
			format:         `{{ .LevelCode }} {{ .Time | timef "2006" }} {{ .Message }}`,
			currentLevel:   LevelInfo,
			messageLevel:   LevelInfo,
			message:        "Hello!",
			expectedOutput: "I 2025 Hello!\n",
		},
		{
			description:    "Do not log debug if set to info",
			format:         `{{ .LevelCode }} {{ .Time | timef "2006" }} {{ .Message }}`,
			currentLevel:   LevelInfo,
			messageLevel:   LevelDebug,
			message:        "Hello!",
			expectedOutput: "",
		},
		{
			description:    "Do not debug if set to debug",
			format:         `{{ .LevelCode }} {{ .Time | timef "2006" }} {{ .Message }}`,
			currentLevel:   LevelDebug,
			messageLevel:   LevelDebug,
			message:        "Hello!",
			expectedOutput: "D 2025 Hello!\n",
		},
		{
			description:    "Log in color",
			format:         `{{ .LevelCode }} {{ .Time | timef "2006" | color .Level }} {{ .Message }}`,
			currentLevel:   LevelInfo,
			messageLevel:   LevelInfo,
			message:        "Hello!",
			expectedOutput: "I \x1b[32m2025\x1b[0m Hello!\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			var buf bytes.Buffer
			var exampleLogger = logger{
				level:     tc.currentLevel,
				formatter: mustNewFormatter(tc.format),
				output:    &buf,
			}
			exampleLogger.logf(tc.messageLevel, tc.message)
			logWritten := buf.String()
			assert.Equal(t, tc.expectedOutput, logWritten)
		})
	}

}
