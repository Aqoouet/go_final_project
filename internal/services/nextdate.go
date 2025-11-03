// Package services contains business logic such as next occurrence date calculation.
package services

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

type NextDateCalculator struct{}

func NewNextDateCalculator() *NextDateCalculator {
	return &NextDateCalculator{}
}

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

func (c *NextDateCalculator) isAfter(date, reference time.Time) bool {
	return date.Format(DateFormat) > reference.Format(DateFormat)
}

func (c *NextDateCalculator) calculateYearly(now, startDate time.Time) (string, error) {
	nextDate := startDate.AddDate(1, 0, 0)
	for !c.isAfter(nextDate, now) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}

	return nextDate.Format(DateFormat), nil
}

func (c *NextDateCalculator) calculateDaily(now, startDate time.Time, intervalStr string) (string, error) {
	interval, err := c.parseDayInterval(intervalStr)
	if err != nil {
		return "", err
	}
	nextDate := startDate.AddDate(0, 0, interval)
	for !c.isAfter(nextDate, now) {
		nextDate = nextDate.AddDate(0, 0, interval)
	}

	return nextDate.Format(DateFormat), nil
}

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

func (c *NextDateCalculator) calculateWeekly(now time.Time, weekdaysStr string) (string, error) {
	allowedWeekdays, err := c.parseWeekdays(weekdaysStr)
	if err != nil {
		return "", err
	}

	return c.findNextWeekday(now, allowedWeekdays)
}

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

func (c *NextDateCalculator) findNextWeekday(now time.Time, allowedWeekdays [8]bool) (string, error) {
	nextDate := now.AddDate(0, 0, 1)
	for i := 0; i < 7; i++ {
		weekday := c.normalizeWeekday(nextDate.Weekday())
		if allowedWeekdays[weekday] {
			return nextDate.Format(DateFormat), nil
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}

	return "", errors.New("no suitable weekday found")
}

func (c *NextDateCalculator) normalizeWeekday(weekday time.Weekday) int {
	if weekday == time.Sunday {
		return 7
	}
	return int(weekday)
}

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

func (c *NextDateCalculator) parseMonths(monthsStr string) ([13]bool, error) {
	var months [13]bool

	if monthsStr == "" {
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

func (c *NextDateCalculator) isMatchingDate(date time.Time, positiveDays [32]bool, negativeDays []int, allowedMonths [13]bool) bool {
	year, month, day := date.Date()
	monthNum := int(month)
	if !allowedMonths[monthNum] {
		return false
	}
	if positiveDays[day] {
		return true
	}
	if len(negativeDays) > 0 {
		lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
		for _, negDay := range negativeDays {
			targetDay := lastDayOfMonth + negDay + 1
			if targetDay >= 1 && day == targetDay {
				return true
			}
		}
	}

	return false
}

