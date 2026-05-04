package utils

import (
	"fmt"
	"time"
)

// ParseDate attempts to parse a date string from multiple common formats
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	
	formats := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}

	var err error
	var t time.Time
	for _, f := range formats {
		t, err = time.Parse(f, dateStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}
