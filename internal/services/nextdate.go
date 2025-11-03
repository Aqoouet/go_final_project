package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// DateFormat is the standard date format used throughout the application (YYYYMMDD)
	DateFormat = "20060102"
	
	// MaxDayInterval is the maximum allowed interval for daily repetition
	MaxDayInterval = 400
)

// NextDateCalculator provides functionality to calculate next occurrence dates
// based on repetition rules
type NextDateCalculator struct{}

// NewNextDateCalculator creates a new instance of NextDateCalculator
func NewNextDateCalculator() *NextDateCalculator {
	return &NextDateCalculator{}
}

// Calculate computes the next date using the given repetition rule
// Parameters:
//   - now: current date/time reference point
//   - date: initial date in YYYYMMDD format
//   - repeat: repetition rule (e.g., "y", "d 7", "w 1,3,5", "m 1,15 1,6")
//
// Returns empty string if repeat is empty (no repetition)
func (c *NextDateCalculator) Calculate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	startDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("empty repeat rule")
	}

	return c.processRepeatRule(now, startDate, parts)
}

// processRepeatRule routes to appropriate calculation function based on repeat rule type
func (c *NextDateCalculator) processRepeatRule(now, startDate time.Time, parts []string) (string, error) {
	rule := parts[0]

	switch rule {
	case "y":
		return c.calculateYearly(now, startDate)
	case "d":
		if len(parts) < 2 {
			return "", errors.New("missing interval for 'd' rule")
		}
		return c.calculateDaily(now, startDate, parts[1])
	case "w":
		if len(parts) < 2 {
			return "", errors.New("missing interval for 'w' rule")
		}
		return c.calculateWeekly(now, parts[1])
	case "m":
		if len(parts) < 2 {
			return "", errors.New("missing interval for 'm' rule")
		}
		monthsStr := ""
		if len(parts) > 2 {
			monthsStr = parts[2]
		}
		return c.calculateMonthly(now, parts[1], monthsStr)
	default:
		return "", fmt.Errorf("unknown rule: %s", rule)
	}
}

// isAfter checks if date is after the reference date (comparing only dates, not time)
func (c *NextDateCalculator) isAfter(date, reference time.Time) bool {
	return date.Format(DateFormat) > reference.Format(DateFormat)
}

// calculateYearly computes next date for yearly repetition
func (c *NextDateCalculator) calculateYearly(now, startDate time.Time) (string, error) {
	// Add at least one year to get the next occurrence
	nextDate := startDate.AddDate(1, 0, 0)

	// Keep adding years until we find a date after now
	for !c.isAfter(nextDate, now) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}

	return nextDate.Format(DateFormat), nil
}

// calculateDaily computes next date for daily repetition (every N days)
func (c *NextDateCalculator) calculateDaily(now, startDate time.Time, intervalStr string) (string, error) {
	interval, err := c.parseDayInterval(intervalStr)
	if err != nil {
		return "", err
	}

	// Add at least one interval to get the next occurrence
	nextDate := startDate.AddDate(0, 0, interval)

	// Keep adding intervals until we find a date after now
	for !c.isAfter(nextDate, now) {
		nextDate = nextDate.AddDate(0, 0, interval)
	}

	return nextDate.Format(DateFormat), nil
}

// parseDayInterval validates and parses day interval string
func (c *NextDateCalculator) parseDayInterval(intervalStr string) (int, error) {
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		return 0, fmt.Errorf("invalid interval: %w", err)
	}

	if interval <= 0 || interval > MaxDayInterval {
		return 0, fmt.Errorf("interval must be between 1 and %d", MaxDayInterval)
	}

	return interval, nil
}

// calculateWeekly computes next date for weekly repetition on specific weekdays
func (c *NextDateCalculator) calculateWeekly(now time.Time, weekdaysStr string) (string, error) {
	allowedWeekdays, err := c.parseWeekdays(weekdaysStr)
	if err != nil {
		return "", err
	}

	return c.findNextWeekday(now, allowedWeekdays)
}

// parseWeekdays parses comma-separated weekday values (1=Monday, 7=Sunday)
func (c *NextDateCalculator) parseWeekdays(weekdaysStr string) ([8]bool, error) {
	var weekdays [8]bool

	parts := strings.Split(weekdaysStr, ",")
	for _, part := range parts {
		day, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || day < 1 || day > 7 {
			return weekdays, fmt.Errorf("invalid weekday: %s", part)
		}
		weekdays[day] = true
	}

	return weekdays, nil
}

// findNextWeekday finds the next date that matches allowed weekdays
func (c *NextDateCalculator) findNextWeekday(now time.Time, allowedWeekdays [8]bool) (string, error) {
	nextDate := now.AddDate(0, 0, 1)

	// Check next 7 days (guaranteed to find a match)
	for i := 0; i < 7; i++ {
		weekday := c.normalizeWeekday(nextDate.Weekday())
		if allowedWeekdays[weekday] {
			return nextDate.Format(DateFormat), nil
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}

	return "", errors.New("no suitable weekday found")
}

// normalizeWeekday converts Go's Weekday (Sunday=0) to ISO format (Monday=1, Sunday=7)
func (c *NextDateCalculator) normalizeWeekday(weekday time.Weekday) int {
	if weekday == time.Sunday {
		return 7
	}
	return int(weekday)
}

// calculateMonthly computes next date for monthly repetition on specific days
func (c *NextDateCalculator) calculateMonthly(now time.Time, daysStr, monthsStr string) (string, error) {
	positiveDays, negativeDays, err := c.parseMonthDays(daysStr)
	if err != nil {
		return "", err
	}

	allowedMonths, err := c.parseMonths(monthsStr)
	if err != nil {
		return "", err
	}

	nextDate := now.AddDate(0, 0, 1)
	for i := 0; i < 366; i++ {
		if c.isMatchingDate(nextDate, positiveDays, negativeDays, allowedMonths) {
			return nextDate.Format(DateFormat), nil
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}

	return "", errors.New("no suitable date found within a year")
}

// parseMonthDays parses comma-separated day values (e.g., "1,15,-1")
// Returns positive days array, negative days slice, and error
func (c *NextDateCalculator) parseMonthDays(daysStr string) ([32]bool, []int, error) {
	var positiveDays [32]bool
	var negativeDays []int

	dayParts := strings.Split(daysStr, ",")
	for _, part := range dayParts {
		day, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return positiveDays, nil, fmt.Errorf("invalid day: %s", part)
		}

		switch {
		case day >= 1 && day <= 31:
			positiveDays[day] = true
		case day >= -31 && day < 0:
			negativeDays = append(negativeDays, day)
		default:
			return positiveDays, nil, fmt.Errorf("day must be 1-31 or -1 to -31, got: %d", day)
		}
	}

	return positiveDays, negativeDays, nil
}

// parseMonths parses comma-separated month values (e.g., "1,3,12")
// Returns months array. Empty string means all months are allowed
func (c *NextDateCalculator) parseMonths(monthsStr string) ([13]bool, error) {
	var months [13]bool

	if monthsStr == "" {
		// All months allowed
		for i := 1; i <= 12; i++ {
			months[i] = true
		}
		return months, nil
	}

	monthParts := strings.Split(monthsStr, ",")
	for _, part := range monthParts {
		month, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || month < 1 || month > 12 {
			return months, fmt.Errorf("invalid month: %s", part)
		}
		months[month] = true
	}

	return months, nil
}

// isMatchingDate checks if date matches the specified days and months
func (c *NextDateCalculator) isMatchingDate(date time.Time, positiveDays [32]bool, negativeDays []int, allowedMonths [13]bool) bool {
	year, month, day := date.Date()
	monthNum := int(month)

	// Check if month matches
	if !allowedMonths[monthNum] {
		return false
	}

	// Check positive days (e.g., day 1, 15, 31)
	if positiveDays[day] {
		return true
	}

	// Check negative days (e.g., last day, second-to-last day)
	if len(negativeDays) > 0 {
		lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
		for _, negDay := range negativeDays {
			// negDay is negative (e.g., -1 means last day)
			targetDay := lastDayOfMonth + negDay + 1
			if targetDay >= 1 && day == targetDay {
				return true
			}
		}
	}

	return false
}

