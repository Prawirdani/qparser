package qparser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "Time only",
			input:    "17:12:32",
			expected: time.Date(0, 1, 1, 17, 12, 32, 0, time.UTC),
		},
		{
			name:     "Date only",
			input:    "2025-07-04",
			expected: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "DateTime without timezone",
			input:    "2025-07-04T17:12:32",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 0, time.UTC),
		},
		{
			name:     "RFC3339 with Z",
			input:    "2025-07-04T17:12:32Z",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 0, time.UTC),
		},
		{
			name:     "RFC3339 with offset",
			input:    "2025-07-04T17:12:32+07:00",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 0, time.FixedZone("+07:00", 7*3600)),
		},
		{
			name:     "RFC3339Nano with Z",
			input:    "2025-07-04T17:12:32.123456789Z",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 123456789, time.UTC),
		},
		{
			name:     "RFC3339Nano with offset",
			input:    "2025-07-04T17:12:32.123456789+07:00",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 123456789, time.FixedZone("+07:00", 7*3600)),
		},
		{
			name:     "DateTime with space separator",
			input:    "2025-07-04 17:12:32",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 0, time.UTC),
		},
		{
			name:     "DateTime with milliseconds",
			input:    "2025-07-04T17:12:32.123",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 123000000, time.UTC),
		},
		{
			name:     "DateTime with microseconds",
			input:    "2025-07-04T17:12:32.123456",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 123456000, time.UTC),
		},
		{
			name:     "DateTime with nanoseconds",
			input:    "2025-07-04T17:12:32.123456789",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 123456789, time.UTC),
		},
		{
			name:     "DateTime with space-separated offset",
			input:    "2025-07-04T17:12:32 07:00",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 0, time.FixedZone("+07:00", 7*3600)),
		},
		{
			name:     "DateTime with milliseconds and offset",
			input:    "2025-07-04T17:12:32.123+07:00",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 123000000, time.FixedZone("+07:00", 7*3600)),
		},
		{
			name:     "DateTime with microseconds and space-separated offset",
			input:    "2025-07-04T17:12:32.123456 07:00",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 123456000, time.FixedZone("+07:00", 7*3600)),
		},
		{
			name:     "DateTime with nanoseconds and space-separated offset",
			input:    "2025-07-04T17:12:32.123456789 07:00",
			expected: time.Date(2025, 7, 4, 17, 12, 32, 123456789, time.FixedZone("+07:00", 7*3600)),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseTime(tc.input)
			assert.NoError(t, err)
			// Compare the actual time values instead of the exact objects
			assert.True(t, tc.expected.Equal(result))
		})
	}
}
