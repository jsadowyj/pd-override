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

func main() {
	// removes timestamps from log output
	log.SetFlags(0)
	config := getConfig()
	client := pagerduty.NewClient(config.APIKey)
	user, _ := client.GetCurrentUserWithContext(context.TODO(), pagerduty.GetCurrentUserOptions{})
	for _, arg := range os.Args[1:] {
		if !isValidArgument(arg) {
			log.Fatalf("Invalid Argument: %s\n", arg)
		}
		dtSplit := strings.Split(arg, "@") // "M-F@09:00-12:00,01:00-02:00" -> ["M-F", "09:00-12:00,01:00-02:00"]
		dayRanges := strings.Split(dtSplit[0], ",")
		timeRanges := strings.Split(dtSplit[1], ",")
		days := expandWeekdayRanges(dayRanges)
		overrides := generateOverrides(days, timeRanges)

		for _, override := range overrides {
			override.User = user.APIObject
			if res, err := client.CreateOverrideWithContext(context.TODO(), config.ScheduleID, override); err == nil {
				layout := "2006-01-02T15:04:05Z"
				start, _ := time.Parse(layout, res.Start)
				end, _ := time.Parse(layout, res.End)
				fmt.Printf("Created Override for %s: %s-%s\n", res.User.Summary, start.Local().Format("Mon, 02 Jan 2006 15:04:05"), end.Local().Format("15:04:05 MST"))
			} else {
				var aerr pagerduty.APIError
				if errors.As(err, &aerr) {
					log.Fatalf("Failed Override for %s: %s [%d]", override.User.Summary, aerr.APIError.ErrorObject.Message, aerr.APIError.ErrorObject.Code)
				} else {
					log.Fatalln(err)
				}
			}
		}
	}
}
