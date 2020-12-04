package main

import (
	"fmt"
	"testing"
)

func TestValidNumber(t *testing.T) {
	type TestCase struct {
		str      string
		digits   int
		min, max uint64
		want     bool
	}

	testCases := []TestCase{
		TestCase{"1980", -1, 1920, 2002, true},
		TestCase{"1920", -1, 1920, 2002, true},
		TestCase{"1919", -1, 1920, 2002, false},
		TestCase{"2002", -1, 1920, 2002, true},
		TestCase{"2003", -1, 1920, 2002, false},
		TestCase{"1980", 3, 1920, 2002, false},
		TestCase{"1980", 4, 1920, 2002, true},
		TestCase{"1980", 5, 1920, 2002, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			if got := validNumber(tc.str, tc.digits, tc.min, tc.max); got != tc.want {
				t.Errorf("validNumber(%v,%v,%v,%v) = %v, want %v",
					tc.str, tc.digits, tc.min, tc.max, got, tc.want)
			}
		})
	}
}
