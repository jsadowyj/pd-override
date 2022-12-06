package main

import (
	"log"
	"strconv"
	"strings"
	"time"
)

var dow = map[string]time.Weekday{
	"U": time.Sunday,
	"M": time.Monday,
	"T": time.Tuesday,
	"W": time.Wednesday,
	"R": time.Thursday,
	"F": time.Friday,
	"S": time.Saturday,
}

func toDuration(hours, minutes string) (hrs, mins time.Duration) {
	h, err := strconv.Atoi(hours)
	if err != nil {
		log.Fatalln(err)
	}
	m, err := strconv.Atoi(minutes)
	if err != nil {
		log.Fatalln(err)
	}
	hrs = time.Duration(h) * time.Hour
	mins = time.Duration(m) * time.Minute
	return hrs, mins
}

func createTimeRange(start time.Time, end time.Time) []time.Time {
	var times []time.Time
	for start.Before(end) || start == end {
		times = append(times, start)
		start = start.AddDate(0, 0, 1)
	}
	return times
}

func wdToTimes(dayRanges []string) []time.Time {
	var times []time.Time
	now := time.Now()
	// starts week at sunday @ 00:00
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -int(now.Weekday()))
	for _, dayRange := range dayRanges {
		split := strings.Split(dayRange, "-")
		layout := "20060102"
		if len(split) == 2 {
			if len(dayRange) < (17) {
				start := weekStart.AddDate(0, 0, int(dow[split[0]]))
				end := weekStart.AddDate(0, 0, int(dow[split[1]]))
				times = append(times, createTimeRange(start, end)...)
			} else {
				sParsed, err := time.Parse(layout, split[0])
				if err != nil {
					log.Fatalln(err)
				}
				eParsed, err := time.Parse(layout, split[1])
				if err != nil {
					log.Fatalln(err)
				}
				start := time.Date(sParsed.Year(), sParsed.Month(), sParsed.Day(), 0, 0, 0, 0, now.Location())
				end := time.Date(eParsed.Year(), eParsed.Month(), eParsed.Day(), 0, 0, 0, 0, now.Location())
				times = append(times, createTimeRange(start, end)...)
			}
		} else {
			day := dayRange
			if len(dayRange) < 8 {
				t := weekStart.AddDate(0, 0, int(dow[day]))
				times = append(times, t)
			} else {
				parsed, err := time.Parse(layout, day)
				if err != nil {
					log.Fatalln(err)
				}
				t := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, now.Location())
				times = append(times, t)
			}
		}
	}
	return times
}
