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

func expandDatetimes(start time.Time, end time.Time) []time.Time {
	var times []time.Time
	for start.Before(end) || start == end {
		times = append(times, start)
		start = start.AddDate(0, 0, 1)
	}
	return times
}

func parseDatetime(str string) (dt time.Time) {
	now := time.Now()
	// starts week at sunday @ 00:00
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -int(now.Weekday()))
	layout := "20060102"
	if len(str) < 8 {
		dt = weekStart.AddDate(0, 0, int(dow[str]))
	} else {
		parsed, err := time.Parse(layout, str)
		if err != nil {
			log.Fatalln(err)
		}
		dt = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, now.Location())
	}
	return dt
}

func parseWeekdayRange(sDay, eDay string) (start, end time.Time) {
	now := time.Now()
	// starts week at sunday @ 00:00
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -int(now.Weekday()))
	start = weekStart.AddDate(0, 0, int(dow[sDay]))
	end = weekStart.AddDate(0, 0, int(dow[eDay]))
	// handles range between weeks
	overflowDays := 0
	if start.After(end) || start == end {
		overflowDays += 7
	}
	end = end.AddDate(0, 0, overflowDays)

	return start, end
}

func parseLongDateRange(sDate, eDate string) (start, end time.Time) {
	now := time.Now()
	layout := "20060102"
	sParsed, err := time.Parse(layout, sDate)
	if err != nil {
		log.Fatalln(err)
	}
	eParsed, err := time.Parse(layout, eDate)
	if err != nil {
		log.Fatalln(err)
	}
	start = time.Date(sParsed.Year(), sParsed.Month(), sParsed.Day(), 0, 0, 0, 0, now.Location())
	end = time.Date(eParsed.Year(), eParsed.Month(), eParsed.Day(), 0, 0, 0, 0, now.Location())
	return start, end
}

func parseDateRanges(startStr, endStr string) (start, end time.Time) {
	if len(startStr) == 1 && len(endStr) == 1 {
		start, end = parseWeekdayRange(startStr, endStr)
	} else if len(startStr) == 8 && len(endStr) == 8 {
		start, end = parseLongDateRange(startStr, endStr)
	} else {
		// this should theoretically never happen; but it's here just in case.
		log.Fatalf("parseDays(): invalid input -- %s-%s", startStr, endStr)
	}
	return start, end
}

func wdToTimes(dtRanges []string) []time.Time {
	var times []time.Time
	for _, dtRange := range dtRanges {
		if dtRange == "" {
			continue
		}
		dSplit := strings.Split(dtRange, "-")
		if len(dSplit) == 2 {
			start, end := parseDateRanges(dSplit[0], dSplit[1])
			times = append(times, expandDatetimes(start, end)...)
		} else {
			datetime := parseDatetime(dSplit[0])
			times = append(times, datetime)
		}
	}
	return times
}
