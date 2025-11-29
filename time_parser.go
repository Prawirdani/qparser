package qparser

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	timezoneFixRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{1,9})?)\s([+-]?\d{2}:\d{2})$`)

	timeFormats = []string{
		time.RFC3339,
		time.RFC3339Nano,
		time.DateTime,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05.000000",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05.000-07:00",
		"2006-01-02T15:04:05.000000-07:00",
		"2006-01-02T15:04:05.999999999-07:00",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05.000-07:00",
		"2006-01-02 15:04:05.000000-07:00",
		"2006-01-02 15:04:05.999999999-07:00",

		// Formats for space-separated timezone
		"2006-01-02T15:04:05 -07:00",
		"2006-01-02T15:04:05.000 -07:00",
		"2006-01-02T15:04:05.000000 -07:00",
		"2006-01-02T15:04:05.999999999 -07:00",
		time.DateOnly,
		time.TimeOnly,
	}
)

func parseTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)

	// Handle space-separated timezone offsets by converting them to standard format
	if strings.Contains(value, " ") {
		value = fixTimezoneOffset(value)
	}

	// Fast path: try most common formats first based on heuristics
	if len(value) >= 19 {
		// RFC3339 format detection
		if value[4] == '-' && value[7] == '-' && value[10] == 'T' && value[13] == ':' && value[16] == ':' {
			if strings.HasSuffix(value, "Z") || (len(value) >= 20 && (value[19] == '+' || value[19] == '-')) {
				// Try RFC3339 formats first
				if t, err := time.Parse(time.RFC3339, value); err == nil {
					return t, nil
				}
				if t, err := time.Parse(time.RFC3339Nano, value); err == nil {
					return t, nil
				}
			}
		}
	}

	// Date only detection
	if len(value) == 10 && value[4] == '-' && value[7] == '-' {
		if t, err := time.Parse(time.DateOnly, value); err == nil {
			return t, nil
		}
	}

	// Time only detection
	if len(value) >= 8 && value[2] == ':' && value[5] == ':' {
		if t, err := time.Parse(time.TimeOnly, value); err == nil {
			return t, nil
		}
	}

	// Fallback to all formats
	for _, format := range timeFormats {
		if parsedTime, err := time.Parse(format, value); err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("%w: unable to parse with any known date format: %s", ErrInvalidValue, value)
}

func fixTimezoneOffset(value string) string {
	// Handle cases like "2025-07-04T17:12:32 07:00" -> "2025-07-04T17:12:32+07:00"
	// Also handle cases that already have +/- like "2025-07-04T17:12:32 -05:00"
	matches := timezoneFixRegex.FindStringSubmatch(value)
	if matches != nil {
		timePart := matches[1]
		offsetPart := matches[2]

		// If offset doesn't start with + or -, assume positive
		if offsetPart[0] != '+' && offsetPart[0] != '-' {
			offsetPart = "+" + offsetPart
		}

		return timePart + offsetPart
	}
	return value
}
