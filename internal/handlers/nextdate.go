package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DateFormat = "20060102"
	MaxDayInterval = 400
)

// NextDate calculates next date using rule for repeat
func NextDate(now time.Time, date string, repeat string) (string, error) {
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

	return processRepeatRule(now, startDate, parts)
}

// processRepeatRule routes to appropriate function based on repeat rule type
func processRepeatRule(now, startDate time.Time, parts []string) (string, error) {
	rule := parts[0]

	switch rule {
	case "y":
		return nextYearly(now, startDate)
	case "d":
		if len(parts) < 2 {
			return "", errors.New("missing interval for 'd' rule")
		}
		return nextDays(now, startDate, parts[1])
	case "w":
		if len(parts) < 2 {
			return "", errors.New("missing interval for 'w' rule")
		}
		return nextWeekdays(now, parts[1])
	case "m":
		if len(parts) < 2 {
			return "", errors.New("missing interval for 'm' rule")
		}
		monthsStr := ""
		if len(parts) > 2 {
			monthsStr = parts[2]
		}
		return nextMonthDays(now, parts[1], monthsStr)
	default:
		return "", fmt.Errorf("unknown rule: %s", rule)
	}
}

// afterNow checks if date > now
func afterNow(date, now time.Time) bool {
	return date.Format("20060102") > now.Format("20060102")
}

// nextYearly calculates next date for yearly cycle
func nextYearly(now, startDate time.Time) (string, error) {
	// Add at least one year to get the next occurrence
	nextDate := startDate.AddDate(1, 0, 0)
	
	// Keep adding years until we find a date after now
	for !afterNow(nextDate, now) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	
	// If the original date was Feb 29 (leap year) and next year is not a leap year,
	// the date will be adjusted to Mar 1 by Go's AddDate behavior
	
	return nextDate.Format(DateFormat), nil
}

// nextDays calculates next date for day interval (e.g., every N days)
func nextDays(now, startDate time.Time, intervalStr string) (string, error) {
	interval, err := parseDayInterval(intervalStr)
	if err != nil {
		return "", err
	}

	// Add at least one interval to get the next occurrence
	nextDate := startDate.AddDate(0, 0, interval)
	
	// Keep adding intervals until we find a date after now
	for !afterNow(nextDate, now) {
		nextDate = nextDate.AddDate(0, 0, interval)
	}
	
	return nextDate.Format(DateFormat), nil
}

// parseDayInterval validates and parses day interval string
func parseDayInterval(intervalStr string) (int, error) {
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		return 0, fmt.Errorf("invalid interval: %w", err)
	}

	if interval <= 0 || interval > MaxDayInterval {
		return 0, fmt.Errorf("interval must be between 1 and %d", MaxDayInterval)
	}

	return interval, nil
}

// nextWeekdays calculates next date for given weekdays
func nextWeekdays(now time.Time, weekdaysStr string) (string, error) {
	allowedWeekdays, err := parseWeekdays(weekdaysStr)
	if err != nil {
		return "", err
	}

	return findNextWeekday(now, allowedWeekdays)
}

// parseWeekdays parses comma-separated weekday values (1=Monday, 7=Sunday)
func parseWeekdays(weekdaysStr string) ([8]bool, error) {
	var weekdays [8]bool

	parts := strings.Split(weekdaysStr, ",")
	for _, part := range parts {
		day, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || day < 1 || day > 7 {
			return weekdays, fmt.Errorf("incorrect weekday: %s", part)
		}
		weekdays[day] = true
	}

	return weekdays, nil
}

// findNextWeekday finds the next date that matches allowed weekdays
func findNextWeekday(now time.Time, allowedWeekdays [8]bool) (string, error) {
	nextDate := now.AddDate(0, 0, 1)

	// Check next 7 days (guaranteed to find a match)
	for i := 0; i < 7; i++ {
		weekday := normalizeWeekday(nextDate.Weekday())
		if allowedWeekdays[weekday] {
			return nextDate.Format(DateFormat), nil
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}

	// Technically unreachable if at least one weekday is specified
	return "", errors.New("no suitable weekday found")
}

// normalizeWeekday converts Go's Weekday (Sunday=0) to ISO format (Monday=1, Sunday=7)
func normalizeWeekday(weekday time.Weekday) int {
	if weekday == time.Sunday {
		return 7
	}
	return int(weekday)
}

// nextMonthDays calculates next date for specified month days
func nextMonthDays(now time.Time, daysStr, monthsStr string) (string, error) {
	positiveDays, negativeDays, err := parseMonthDays(daysStr)
	if err != nil {
		return "", err
	}

	allowedMonths, err := parseMonths(monthsStr)
	if err != nil {
		return "", err
	}

	nextDate := now.AddDate(0, 0, 1)
	for i := 0; i < 366; i++ {
		if isMatchingDate(nextDate, positiveDays, negativeDays, allowedMonths) {
			return nextDate.Format(DateFormat), nil
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}

	return "", errors.New("no suitable date found within a year")
}

// parseMonthDays parses comma-separated day values (e.g., "1,15,-1")
// Returns positive days array, negative days slice, and error
func parseMonthDays(daysStr string) ([32]bool, []int, error) {
	var positiveDays [32]bool
	var negativeDays []int

	dayParts := strings.Split(daysStr, ",")
	for _, part := range dayParts {
		day, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return positiveDays, nil, fmt.Errorf("incorrect day: %s", part)
		}

		switch {
		case day >= 1 && day <= 31:
			positiveDays[day] = true
		case day >= -31 && day < 0:
			negativeDays = append(negativeDays, day)
		default:
			return positiveDays, nil, fmt.Errorf("incorrect day: %d (must be 1-31 or -1 to -31)", day)
		}
	}

	return positiveDays, negativeDays, nil
}

// parseMonths parses comma-separated month values (e.g., "1,3,12")
// Returns months array. Empty string means all months are allowed
func parseMonths(monthsStr string) ([13]bool, error) {
	var months [13]bool

	if monthsStr == "" {
		// All months allowed - fill all indices
		for i := 1; i <= 12; i++ {
			months[i] = true
		}
		return months, nil
	}

	monthParts := strings.Split(monthsStr, ",")
	for _, part := range monthParts {
		month, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || month < 1 || month > 12 {
			return months, fmt.Errorf("incorrect month: %s", part)
		}
		months[month] = true
	}

	return months, nil
}

// isMatchingDate checks if date matches the specified days and months
func isMatchingDate(date time.Time, positiveDays [32]bool, negativeDays []int, allowedMonths [13]bool) bool {
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
