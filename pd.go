package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

func generateOverrides(days []time.Weekday, timeRanges []string) []pagerduty.Override {
	var overrides []pagerduty.Override
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

	for _, day := range days {
		for _, timeRange := range timeRanges {
			rSplit := strings.Split(timeRange, "-")
			hour, _ := strconv.Atoi((rSplit[0])[:2])
			minute, _ := strconv.Atoi((rSplit[0])[2:])
			start := weekStart.AddDate(0, 0, int(day)).Add(time.Duration(hour) * time.Hour).Add(time.Duration(minute) * time.Minute).Format(time.RFC3339)
			hour, _ = strconv.Atoi((rSplit[1])[:2])
			minute, _ = strconv.Atoi((rSplit[1])[2:])
			end := weekStart.AddDate(0, 0, int(day)).Add(time.Duration(hour) * time.Hour).Add(time.Duration(minute) * time.Minute).Format(time.RFC3339)

			overrides = append(overrides, pagerduty.Override{Start: start, End: end})

		}
	}

	return overrides
}
