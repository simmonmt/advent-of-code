package intmath

import (
	"fmt"
	"testing"

	"github.com/simmonmt/aoc/2019/common/testutils"
)

func TestGCD(t *testing.T) {
	type TestCase struct {
		vs []int
		d  int
	}

	testCases := []TestCase{
		TestCase{[]int{4, 6}, 2},
		TestCase{[]int{8, 12}, 4},
		TestCase{[]int{12, 18}, 6},
		TestCase{[]int{9, 9}, 9},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.vs), func(t *testing.T) {
			if got := GCD(tc.vs...); got != tc.d {
				t.Errorf("GCD(%v) = %d, want %d", tc.vs, got, tc.d)
			}
		})
	}

	// Make sure it lets us know when we need to add more primes
	testutils.AssertPanic(t, "too big", func() { GCD(101*103, 107*109) })
}
