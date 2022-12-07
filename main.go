package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
)

func main() {
	// removes timestamps from log output
	log.SetFlags(0)
	config := getConfig()
	client := pagerduty.NewClient(config.APIKey)
	user, err := client.GetCurrentUserWithContext(context.TODO(), pagerduty.GetCurrentUserOptions{})
	if err != nil {
		log.Fatalln("Unable to GET current user from PagerDuty API.")
	}
	for _, arg := range os.Args[1:] {
		if !isValidInput(arg) {
			log.Fatalf("Invalid Input: %s\n", arg)
		}
		dtSplit := strings.Split(arg, "@") // "M-F@09:00-12:00,01:00-02:00" -> ["M-F", "09:00-12:00,01:00-02:00"]
		var overrides []pagerduty.Override
		if len(dtSplit) == 2 {
			overrides = createDailyOverrides(dtSplit[0], dtSplit[1])
		} else {
			overrides = createOverrides(dtSplit[0])
		}

		for _, override := range overrides {
			override.User = user.APIObject
			if ov, err := client.CreateOverrideWithContext(context.TODO(), config.ScheduleID, override); err == nil {
				printOverride(*ov, arg)
			} else {
				printOverrideError(err, arg)
			}
		}
	}
}
