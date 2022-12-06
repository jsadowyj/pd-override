package main

import (
	"regexp"
)

func makeRange(min, max int) []int {
	// todo: make ranges go into the next week. This is just lazy.
	if min > max {
		min, max = max, min
	}
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func isValidArgument(arg string) bool {
	re := regexp.MustCompile(`^(([\d]{8}(-[\d]{8})?|(([UMTWRFS](-[UMTWRFS])?))),?)+@(([\d]{2,4}-[\d]{2,4}),?)+$`)
	return re.Match([]byte(arg))
}
