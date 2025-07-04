package qparser

import (
	"reflect"
	"regexp"
	"time"
)

var timeFormats = []string{
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
	time.RFC3339,                 // "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano,             // "2006-01-02T15:04:05.999999999Z07:00"
	"2006-01-02T15:04:05+07:00",  // Alternative timezone format
	"2006-01-02T15:04:05",        // Without timezone
	"2006-01-02T15:04:05.000",    // With milliseconds, no timezone
	"2006-01-02T15:04:05.000000", // With microseconds, no timezone
	"2006-01-02 15:04:05",        // Space separator
	"2006-01-02 15:04:05+07:00",  // Space separator with timezone
}

// setTime creates a flexible time setter that supports multiple formats
func setTime() Setter {
	return func(value string) reflect.Value {
		// Fix timezone offset if + was converted to space during URL decoding
		value = fixTimezoneOffset(value)

		var parsedTime time.Time
		var err error

		// TODO: Use Regex instead?
		for _, format := range timeFormats {
			parsedTime, err = time.Parse(format, value)
			if err == nil {
				return reflect.ValueOf(parsedTime)
			}
		}

		// If all formats failed, return invalid
		return invalidReflect
	}
}

// fixTimezoneOffset fixes timezone offset that was corrupted by URL decoding
// Converts " 07:00" back to "+07:00" when it appears at the end of a datetime string
func fixTimezoneOffset(value string) string {
	// Pattern to match datetime strings with optional nanoseconds and space before timezone offset
	// Matches things like "2025-07-04T17:12:32 07:00" or "2025-07-04T17:12:32.123456789 07:00"
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{1,9})?)\s(\d{2}:\d{2})$`)
	return re.ReplaceAllString(value, "${1}+${2}")
}
