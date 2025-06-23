package clilog

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogf(t *testing.T) {

	// Restore after tests
	originalNow := now
	originalOutput := output
	defer func() { now = originalNow }()
	defer func() { output = originalOutput }()

	now = func() time.Time {
		return time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)
	}
	cases := []struct {
		description    string
		currentLevel   Level
		messageLevel   Level
		message        string
		expectedOutput string
	}{
		{
			description:    "Basic log",
			currentLevel:   LevelInfo,
			messageLevel:   LevelInfo,
			message:        "Hello!",
			expectedOutput: "I 2025/01/01 00:00:00.000 Hello!\n",
		},
		{
			description:    "Do not log debug if set to info",
			currentLevel:   LevelInfo,
			messageLevel:   LevelDebug,
			message:        "Hello!",
			expectedOutput: "",
		},
		{
			description:    "Do debug if set to debug",
			currentLevel:   LevelDebug,
			messageLevel:   LevelDebug,
			message:        "Hello!",
			expectedOutput: "D 2025/01/01 00:00:00.000 Hello!\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			var buf bytes.Buffer
			output = &buf
			var exampleLogger = logger{
				level:     tc.currentLevel,
				formatter: mustNewFormatter(DefaultFormatNoColor),
			}
			exampleLogger.logf(tc.messageLevel, tc.message)
			logWritten := buf.String()
			assert.Equal(t, tc.expectedOutput, logWritten)
		})
	}
}
