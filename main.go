package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	// removes timestamps from log output
	log.SetFlags(0)
	for _, arg := range os.Args[1:] {
		if !isValidInput(arg) {
			log.Fatalf("Invalid Input: %s\n", arg)
		}
		dtSplit := strings.Split(arg, "@") // "M-F@09:00-12:00,01:00-02:00" -> ["M-F", "09:00-12:00,01:00-02:00"]
		var overrides []Override
		if len(dtSplit) == 2 {
			overrides = createDailyOverrides(dtSplit[0], dtSplit[1])
		} else {
			overrides = createOverrides(dtSplit[0])
		}
		sendOverrides(overrides)
	}
}
