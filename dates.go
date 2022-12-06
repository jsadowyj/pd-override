package main

import (
	"strings"
	"time"
)

var daysOfWeek = map[string]time.Weekday{
	"U": time.Sunday,
	"M": time.Monday,
	"T": time.Tuesday,
	"W": time.Wednesday,
	"R": time.Thursday,
	"F": time.Friday,
	"S": time.Saturday,
}

func makeWeekdayRange(min, max time.Weekday) []time.Weekday {
	a := makeRange(int(min), int(max))
	b := make([]time.Weekday, len(a))
	for i, day := range a {
		b[i] = time.Weekday(day)
	}
	return b
}

func expandWeekdayRanges(dayRanges []string) []time.Weekday {
	var days []time.Weekday

	for _, dayRange := range dayRanges {
		split := strings.Split(dayRange, "-")
		if len(split) == 2 {
			start := daysOfWeek[split[0]]
			end := daysOfWeek[split[1]]
			days = append(days, makeWeekdayRange(start, end)...)
		} else {
			days = append(days, daysOfWeek[dayRange])
		}
	}

	return days
}
