package main

import (
	"reflect"
	"testing"
	"time"
)

func TestToDuration(t *testing.T) {
	type Input struct {
		hour   string
		minute string
	}
	tests := []struct {
		input      Input
		wantHour   time.Duration
		wantMinute time.Duration
	}{
		{input: Input{"09", "30"}, wantHour: 9, wantMinute: 30},
		{input: Input{"9", "30"}, wantHour: 9, wantMinute: 30},
		{input: Input{"00", "45"}, wantHour: 0, wantMinute: 45},
		{input: Input{"0", "45"}, wantHour: 0, wantMinute: 45},
		{input: Input{"23", "59"}, wantHour: 23, wantMinute: 59},
		{input: Input{"25", "67"}, wantHour: 25, wantMinute: 67},
	}

	for _, test := range tests {
		wantHour, wantMinute := time.Duration(test.wantHour)*time.Hour, time.Duration(test.wantMinute)*time.Minute
		gotHour, gotMinute := toDuration(test.input.hour, test.input.minute)
		if wantHour != gotHour || wantMinute != gotMinute {
			t.Fatalf("Expected (%s), but got (%s).", wantHour+wantMinute, gotHour+gotMinute)
		}
	}
}

func TestExpandDateTimes(t *testing.T) {
	now := time.Now()
	type Input struct {
		start, end time.Time
	}
	tests := []struct {
		input Input
		want  []time.Time
	}{
		{
			input: Input{start: time.Date(2022, 12, 5, 0, 0, 0, 0, time.UTC), end: time.Date(2022, 12, 7, 0, 0, 0, 0, time.UTC)},
			want: []time.Time{
				time.Date(2022, 12, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 12, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 12, 7, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: Input{start: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC), end: time.Date(2023, 3, 4, 0, 0, 0, 0, time.UTC)},
			want: []time.Time{
				time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 3, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 3, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 3, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: Input{start: now, end: now.AddDate(0, 0, 5)},
			want: []time.Time{
				now,
				now.AddDate(0, 0, 1),
				now.AddDate(0, 0, 2),
				now.AddDate(0, 0, 3),
				now.AddDate(0, 0, 4),
				now.AddDate(0, 0, 5),
			},
		},
	}

	for _, test := range tests {
		tRange := expandDatetimes(test.input.start, test.input.end)
		if !reflect.DeepEqual(tRange, test.want) {
			t.Fatalf("Expected %v, but got %v", test.want, tRange)
		}
	}
}

func TestParseDateTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		input string
		want  time.Time
	}{
		{input: "20220327", want: time.Date(2022, 3, 27, 0, 0, 0, 0, now.Location())},
		{input: "20220417", want: time.Date(2022, 4, 17, 0, 0, 0, 0, now.Location())},
		{input: "20220511", want: time.Date(2022, 5, 11, 0, 0, 0, 0, now.Location())},
		// {input: "20220600", want: time.Date(2022, 6, 0, 0, 0, 0, 0, now.Location())},
		{input: "20220706", want: time.Date(2022, 7, 6, 0, 0, 0, 0, now.Location())},
		// {input: "20220561", want: time.Date(2022, 5, 61, 0, 0, 0, 0, now.Location())},
	}

	for _, test := range tests {
		if parsed := parseDatetime(test.input); parsed != test.want {
			t.Fatalf("Expected %v, but got %v", test.want, parsed)
		}
	}
}

func TestParseWeekdayRange(t *testing.T) {
	now := time.Now()
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -int(now.Weekday()))
	type Input struct {
		sDay, eDay string
	}
	tests := []struct {
		input Input
		sWant time.Time
		eWant time.Time
	}{
		{input: Input{sDay: "M", eDay: "F"}, sWant: weekStart.AddDate(0, 0, 1), eWant: weekStart.AddDate(0, 0, 5)},
		{input: Input{sDay: "T", eDay: "R"}, sWant: weekStart.AddDate(0, 0, 2), eWant: weekStart.AddDate(0, 0, 4)},
		{input: Input{sDay: "F", eDay: "S"}, sWant: weekStart.AddDate(0, 0, 5), eWant: weekStart.AddDate(0, 0, 6)},
		{input: Input{sDay: "U", eDay: "U"}, sWant: weekStart.AddDate(0, 0, 0), eWant: weekStart.AddDate(0, 0, 7)},
		{input: Input{sDay: "M", eDay: "M"}, sWant: weekStart.AddDate(0, 0, 1), eWant: weekStart.AddDate(0, 0, 8)},
		{input: Input{sDay: "T", eDay: "T"}, sWant: weekStart.AddDate(0, 0, 2), eWant: weekStart.AddDate(0, 0, 9)},
		{input: Input{sDay: "W", eDay: "W"}, sWant: weekStart.AddDate(0, 0, 3), eWant: weekStart.AddDate(0, 0, 10)},
		{input: Input{sDay: "R", eDay: "R"}, sWant: weekStart.AddDate(0, 0, 4), eWant: weekStart.AddDate(0, 0, 11)},
		{input: Input{sDay: "F", eDay: "F"}, sWant: weekStart.AddDate(0, 0, 5), eWant: weekStart.AddDate(0, 0, 12)},
		{input: Input{sDay: "S", eDay: "S"}, sWant: weekStart.AddDate(0, 0, 6), eWant: weekStart.AddDate(0, 0, 13)},
		{input: Input{sDay: "S", eDay: "F"}, sWant: weekStart.AddDate(0, 0, 6), eWant: weekStart.AddDate(0, 0, 12)},
	}

	for _, test := range tests {
		sGot, eGot := parseWeekdayRange(test.input.sDay, test.input.eDay)
		if sGot != test.sWant || eGot != test.eWant {
			t.Fatalf("Expected (%s -- %s), but got (%s -- %s).", test.sWant, test.eWant, sGot, eGot)
		}
	}
}

func TestParseLongDateRange(t *testing.T) {
	now := time.Now()
	type Input struct {
		sDate, eDate string
	}
	tests := []struct {
		input Input
		sWant time.Time
		eWant time.Time
	}{
		{
			input: Input{sDate: "20220203", eDate: "20220206"},
			sWant: time.Date(2022, 2, 3, 0, 0, 0, 0, now.Location()),
			eWant: time.Date(2022, 2, 6, 0, 0, 0, 0, now.Location()),
		},
		{
			input: Input{sDate: "20221220", eDate: "20221224"},
			sWant: time.Date(2022, 12, 20, 0, 0, 0, 0, now.Location()),
			eWant: time.Date(2022, 12, 24, 0, 0, 0, 0, now.Location()),
		},
		{
			input: Input{sDate: "20221220", eDate: "20221227"},
			sWant: time.Date(2022, 12, 20, 0, 0, 0, 0, now.Location()),
			eWant: time.Date(2022, 12, 27, 0, 0, 0, 0, now.Location()),
		},
		{
			input: Input{sDate: "20220203", eDate: "20220303"},
			sWant: time.Date(2022, 2, 3, 0, 0, 0, 0, now.Location()),
			eWant: time.Date(2022, 3, 3, 0, 0, 0, 0, now.Location()),
		},
	}

	for _, test := range tests {
		sGot, eGot := parseLongDateRange(test.input.sDate, test.input.eDate)
		if sGot != test.sWant || eGot != test.eWant {
			t.Fatalf("Expected (%s -- %s), but got (%s -- %s).", test.sWant, test.eWant, sGot, eGot)
		}
	}
}

func TestParseDateRanges(t *testing.T) {
	now := time.Now()
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -int(now.Weekday()))
	type Input struct {
		startStr, endStr string
	}
	tests := []struct {
		input Input
		sWant time.Time
		eWant time.Time
	}{
		{
			input: Input{startStr: "20220203", endStr: "20220206"},
			sWant: time.Date(2022, 2, 3, 0, 0, 0, 0, now.Location()),
			eWant: time.Date(2022, 2, 6, 0, 0, 0, 0, now.Location()),
		},
		{
			input: Input{startStr: "20221220", endStr: "20221224"},
			sWant: time.Date(2022, 12, 20, 0, 0, 0, 0, now.Location()),
			eWant: time.Date(2022, 12, 24, 0, 0, 0, 0, now.Location()),
		},
		{
			input: Input{startStr: "20221220", endStr: "20221227"},
			sWant: time.Date(2022, 12, 20, 0, 0, 0, 0, now.Location()),
			eWant: time.Date(2022, 12, 27, 0, 0, 0, 0, now.Location()),
		},
		{
			input: Input{startStr: "20220203", endStr: "20220303"},
			sWant: time.Date(2022, 2, 3, 0, 0, 0, 0, now.Location()),
			eWant: time.Date(2022, 3, 3, 0, 0, 0, 0, now.Location()),
		},
		{input: Input{startStr: "M", endStr: "F"}, sWant: weekStart.AddDate(0, 0, 1), eWant: weekStart.AddDate(0, 0, 5)},
		{input: Input{startStr: "T", endStr: "R"}, sWant: weekStart.AddDate(0, 0, 2), eWant: weekStart.AddDate(0, 0, 4)},
		{input: Input{startStr: "F", endStr: "S"}, sWant: weekStart.AddDate(0, 0, 5), eWant: weekStart.AddDate(0, 0, 6)},
		{input: Input{startStr: "U", endStr: "U"}, sWant: weekStart.AddDate(0, 0, 0), eWant: weekStart.AddDate(0, 0, 7)},
		{input: Input{startStr: "M", endStr: "M"}, sWant: weekStart.AddDate(0, 0, 1), eWant: weekStart.AddDate(0, 0, 8)},
		{input: Input{startStr: "T", endStr: "T"}, sWant: weekStart.AddDate(0, 0, 2), eWant: weekStart.AddDate(0, 0, 9)},
		{input: Input{startStr: "W", endStr: "W"}, sWant: weekStart.AddDate(0, 0, 3), eWant: weekStart.AddDate(0, 0, 10)},
		{input: Input{startStr: "R", endStr: "R"}, sWant: weekStart.AddDate(0, 0, 4), eWant: weekStart.AddDate(0, 0, 11)},
		{input: Input{startStr: "F", endStr: "F"}, sWant: weekStart.AddDate(0, 0, 5), eWant: weekStart.AddDate(0, 0, 12)},
		{input: Input{startStr: "S", endStr: "S"}, sWant: weekStart.AddDate(0, 0, 6), eWant: weekStart.AddDate(0, 0, 13)},
		{input: Input{startStr: "S", endStr: "F"}, sWant: weekStart.AddDate(0, 0, 6), eWant: weekStart.AddDate(0, 0, 12)},
	}

	for _, test := range tests {
		sGot, eGot := parseDateRanges(test.input.startStr, test.input.endStr)
		if sGot != test.sWant || eGot != test.eWant {
			t.Fatalf("Expected (%s -- %s), but got (%s -- %s).", test.sWant, test.eWant, sGot, eGot)
		}
	}
}

func TestWdToTimes(t *testing.T) {
	now := time.Now()
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -int(now.Weekday()))
	tests := []struct {
		input []string
		want  []time.Time
	}{
		{input: []string{"M-F", "S", "U"}, want: []time.Time{weekStart.AddDate(0, 0, 1), weekStart.AddDate(0, 0, 2), weekStart.AddDate(0, 0, 3), weekStart.AddDate(0, 0, 4), weekStart.AddDate(0, 0, 5), weekStart.AddDate(0, 0, 6), weekStart.AddDate(0, 0, 0)}},
		{input: []string{"U-S"}, want: []time.Time{weekStart.AddDate(0, 0, 0), weekStart.AddDate(0, 0, 1), weekStart.AddDate(0, 0, 2), weekStart.AddDate(0, 0, 3), weekStart.AddDate(0, 0, 4), weekStart.AddDate(0, 0, 5), weekStart.AddDate(0, 0, 6)}},
		{input: []string{"U-W", "S"}, want: []time.Time{weekStart.AddDate(0, 0, 0), weekStart.AddDate(0, 0, 1), weekStart.AddDate(0, 0, 2), weekStart.AddDate(0, 0, 3), weekStart.AddDate(0, 0, 6)}},
		{input: []string{"W-M", "M"}, want: []time.Time{weekStart.AddDate(0, 0, 3), weekStart.AddDate(0, 0, 4), weekStart.AddDate(0, 0, 5), weekStart.AddDate(0, 0, 6), weekStart.AddDate(0, 0, 7), weekStart.AddDate(0, 0, 8), weekStart.AddDate(0, 0, 1)}},
		{input: []string{"20221106-20221108", "M"}, want: []time.Time{time.Date(2022, 11, 6, 0, 0, 0, 0, now.Location()), time.Date(2022, 11, 7, 0, 0, 0, 0, now.Location()), time.Date(2022, 11, 8, 0, 0, 0, 0, now.Location()), weekStart.AddDate(0, 0, 1)}},
		{input: []string{"20221106"}, want: []time.Time{time.Date(2022, 11, 6, 0, 0, 0, 0, now.Location())}},
	}

	for _, test := range tests {
		got := wdToTimes(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("Expected %#v, but got %#v.", test.want, got)
		}
	}
}
