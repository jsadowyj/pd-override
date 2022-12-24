package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

type Override struct {
	// parent string, ex: M-F@0900-1700
	pStr string
	// day with time range string, ex: M@0900-1700
	str string
	// pagerduty override struct
	pdOverride pagerduty.Override
	// any error to be logged
	err error
}

// implements the Error interface
func (o Override) Error() string {
	return o.err.Error()
}

func sendOverrides(overrides []Override) {
	config := getConfig()
	client := pagerduty.NewClient(config.APIKey)
	user, err := client.GetCurrentUserWithContext(context.TODO(), pagerduty.GetCurrentUserOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	for _, override := range overrides {
		override.pdOverride.User = user.APIObject
		if ov, err := client.CreateOverrideWithContext(context.TODO(), config.ScheduleID, override.pdOverride); err != nil {
			printOverrideError(err, override.str)
		} else {
			printOverride(*ov, override.str)
		}
	}
}

func createOverrides(wdStrs string) []Override {
	var overrides []Override
	dayRanges := strings.Split(wdStrs, ",")
	for _, dayRange := range dayRanges {
		dSplit := strings.Split(dayRange, "-")
		if dayRange == "" || len(dSplit) != 2 {
			continue
		}
		start, end := parseDateRanges(dSplit[0], dSplit[1])
		end = end.AddDate(0, 0, 1)
		pdOverride := pagerduty.Override{Start: start.Format(time.RFC3339), End: end.Format(time.RFC3339)}
		// todo: remove hardcoded nil
		override := Override{pStr: wdStrs, str: dayRange, pdOverride: pdOverride, err: nil}
		overrides = append(overrides, override)
	}
	return overrides
}

func createDailyOverrides(wdStrs string, timeStrs string) []Override {
	var overrides []Override
	dayRanges := strings.Split(wdStrs, ",")
	timeRanges := strings.Split(timeStrs, ",")
	dTimes := wdToTimes(dayRanges)

	for _, dTime := range dTimes {
		for _, timeRange := range timeRanges {
			if timeRange == "" {
				continue
			}
			rSplit := strings.Split(timeRange, "-")
			hours, minutes := toDuration(rSplit[0][:2], rSplit[0][2:])
			start := dTime.Add(hours).Add(minutes).Format(time.RFC3339)
			hours, minutes = toDuration(rSplit[1][:2], rSplit[1][2:])
			end := dTime.Add(hours).Add(minutes).Format(time.RFC3339)
			// todo: remove hardcoded nil
			override := Override{pStr: wdStrs + "@" + timeStrs, str: rDow[dTime.Weekday()] + "@" + timeRange, pdOverride: pagerduty.Override{Start: start, End: end}, err: nil}
			overrides = append(overrides, override)
		}
	}
	return overrides
}

func printOverride(override pagerduty.Override, msg string) {
	layout := "2006-01-02T15:04:05Z"
	start, _ := time.Parse(layout, override.Start)
	end, _ := time.Parse(layout, override.End)
	fmt.Printf("(%s) Created Override for %s: %s-%s\n", msg, override.User.Summary, start.Local().Format("Mon, 02 Jan 2006 15:04:05"), end.Local().Format("15:04:05 MST"))
}

func printOverrideError(err error, msg string) {
	var aerr pagerduty.APIError
	if errors.As(err, &aerr) {
		fmt.Fprintf(os.Stderr, "(%s) Failed to create override: %s [%d]\n", msg, aerr.APIError.ErrorObject.Message, aerr.APIError.ErrorObject.Code)
	} else {
		log.Fatalln(err)
	}
}
