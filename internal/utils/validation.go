package utils

import (
	"fmt"
	"time"
)

const (
	// DateFormat is the standard date format used throughout the application (YYYYMMDD)
	DateFormat = "20060102"
)

// IsAfter checks if date is after the reference date (comparing only dates, not time)
func IsAfter(date, reference time.Time) bool {
	return date.Format(DateFormat) > reference.Format(DateFormat)
}

// ParseDate parses a date string in YYYYMMDD format
func ParseDate(dateStr string) (time.Time, error) {
	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected YYYYMMDD: %w", err)
	}
	return date, nil
}

// FormatDate formats a time.Time to YYYYMMDD format
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

