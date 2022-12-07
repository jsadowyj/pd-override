package main

import (
	"regexp"
)

func isValidInput(arg string) bool {
	re := regexp.MustCompile(`^(([\d]{8}(-[\d]{8})?|(([UMTWRFS](-[UMTWRFS])?))),?)+(@(([\d]{2,4}-[\d]{2,4}),?))?$`)
	return re.Match([]byte(arg))
}
