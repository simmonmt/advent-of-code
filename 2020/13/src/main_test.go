package main

import (
	"strconv"
	"testing"
)

func TestSolveB(t *testing.T) {
	type TestCase struct {
		in   []int
		want int64
	}

	testCases := []TestCase{
		TestCase{in: []int{17, -1, 13, 19}, want: 3417},
		TestCase{in: []int{67, 7, 59, 61}, want: 754018},
		TestCase{in: []int{67, -1, 7, 59, 61}, want: 779210},
		TestCase{in: []int{67, 7, -1, 59, 61}, want: 1261476},
		TestCase{in: []int{1789, 37, 47, 1889}, want: 1202161486},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ca := findCumAlignment(tc.in)
			if ca.first != tc.want {
				t.Errorf("findCumAlignment(%v) = %v, want %v",
					tc.in, ca.first, tc.want)
			}
		})
	}
}
