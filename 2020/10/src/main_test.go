package main

import (
	"strconv"
	"testing"
)

func TestSeekBack(t *testing.T) {
	nums := []int{3, 6, 9, 12, 15, 18}

	type TestCase struct {
		start, goal, want int
	}

	testCases := []TestCase{
		TestCase{1, 2, 0},
		TestCase{0, 2, 0},
		TestCase{3, 4, 1},
		TestCase{3, 3, 0},
		TestCase{4, 7, 2},
		TestCase{4, 6, 1},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := seekBack(nums, tc.start, tc.goal); got != tc.want {
				t.Errorf(`seekBack(%v, %v, %v) = %v, want %v`,
					nums, tc.start, tc.goal, got, tc.want)
			}
		})
	}
}
