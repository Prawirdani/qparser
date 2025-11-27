package qparser

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	timezoneFixRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{1,9})?)\s(\d{2}:\d{2})$`)

	timeFormats = []string{
		time.RFC3339, // Most common API format
		time.RFC3339Nano,
		time.DateTime, // "2006-01-02 15:04:05"
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05.000000",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05.000-07:00",
		"2006-01-02T15:04:05.000000-07:00",
		"2006-01-02T15:04:05.999999999-07:00",
		"2006-01-02 15:04:05+07:00",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05.123456789-07:00",
		time.DateOnly,
		time.TimeOnly,
	}
)

func parseTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value) // Add input sanitization

	if strings.Contains(value, " ") {
		value = fixTimezoneOffset(value)
	}

	for _, format := range timeFormats {
		if parsedTime, err := time.Parse(format, value); err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("%w: unable to parse %q with any known format", ErrInvalidValue, value)
}

func fixTimezoneOffset(value string) string {
	return timezoneFixRegex.ReplaceAllString(value, "${1}+${2}")
}
