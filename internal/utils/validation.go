// Package utils provides helper functions for HTTP responses and validation.
package utils

import (
	"fmt"
	"time"
)

const (
	DateFormat = "20060102"
)

func IsAfter(date, reference time.Time) bool {
	return date.Format(DateFormat) > reference.Format(DateFormat)
}

func ParseDate(dateStr string) (time.Time, error) {
	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected YYYYMMDD: %w", err)
	}
	return date, nil
}

func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

