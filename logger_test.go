package clilog

import (
	"bytes"
	"testing"

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
			_, actualErr := newFormatter(tc.format)
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
		formatter:    mustNewFormatter(`{{ .LevelCode }} {{ .Message }}`),
		output:       &buf,
	}
	exampleLogger.logf(LevelInfo, "hello")
	logWritten := buf.String()
	assert.Equal(t, "I hello\n", logWritten)
}
