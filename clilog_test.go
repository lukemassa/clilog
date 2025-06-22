package clilog

import (
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
			actualErr := SetFormat(tc.format)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.ErrorContains(t, actualErr, tc.expectedErr)
			}
		})
	}
}
