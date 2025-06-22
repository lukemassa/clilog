package clilog

import (
	"bytes"
	"testing"
	"text/template"

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
			format: "{{ .Time }}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.format, func(t *testing.T) {
			l := logger{}
			actualErr := l.setFormat(tc.format)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.ErrorContains(t, actualErr, tc.expectedErr)
			}
		})
	}
}

func TestBasicLogging(t *testing.T) {

	var buf bytes.Buffer
	var exampleLogger = logger{
		level:        LevelInfo,
		colorEnabled: false,
		timeFormat:   "", // Skip checking time for now
		tmpl:         template.Must(template.New("log").Parse(`{{ .LevelCode }} {{ .Message }}`)),
		output:       &buf,
	}
	exampleLogger.logf(LevelInfo, "hello")
	logWritten := buf.String()
	assert.Equal(t, "I hello\n", logWritten)
}
