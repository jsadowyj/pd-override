package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

func createOverrides(wdStrs string, timeStrs string) []pagerduty.Override {
	var overrides []pagerduty.Override
	dayRanges := strings.Split(wdStrs, ",")
	timeRanges := strings.Split(timeStrs, ",")
	dTimes := wdToTimes(dayRanges)

	for _, dTime := range dTimes {
		for _, timeRange := range timeRanges {
			//dTime = time.Date(dTime.Year(), dTime.Month(), dTime.Day(), 0, 0, 0, 0, time.Now().Location())
			rSplit := strings.Split(timeRange, "-")
			hours, minutes := toDuration(rSplit[0][:2], rSplit[0][2:])
			start := dTime.Add(hours).Add(minutes).Format(time.RFC3339)
			hours, minutes = toDuration(rSplit[1][:2], rSplit[1][2:])
			end := dTime.Add(hours).Add(minutes).Format(time.RFC3339)
			overrides = append(overrides, pagerduty.Override{Start: start, End: end})
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
