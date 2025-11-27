package qparser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	testCases := []string{
		"17:12:32",                            // Time only
		"2025-07-04",                          // Date only
		"2025-07-04T17:12:32",                 // DateTime
		"2025-07-04T17:12:32Z",                // DateTime with UTC timezone
		"2025-07-04 17:12:32",                 // DateTime space separator
		"2025-07-04T17:12:32.123",             // milliseconds
		"2025-07-04T17:12:32.123456",          // microseconds
		"2025-07-04T17:12:32.123456789",       // nanoseconds
		"2025-07-04T17:12:32+07:00",           // DateTime with timezone
		"2025-07-04T17:12:32 07:00",           // DateTime with timezone space separator
		"2025-07-04T17:12:32.123+07:00",       // milliseconds + timezone
		"2025-07-04T17:12:32.123456 07:00",    // microseconds + space-tz
		"2025-07-04T17:12:32.123456789 07:00", // nanoseconds + space-tz
	}

	for _, tRaw := range testCases {
		result, err := parseTime(tRaw)
		assert.NoError(t, err)
		assert.IsType(t, time.Time{}, result)
	}
}
