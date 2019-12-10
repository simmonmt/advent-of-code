package puzzle

import (
	"fmt"
	"testing"

	"github.com/simmonmt/aoc/2019/common/testutils"
)

func TestFactor(t *testing.T) {
	type TestCase struct {
		n, d   int
		en, ed int
	}

	testCases := []TestCase{
		TestCase{4, 6, 2, 3},
		TestCase{12, 18, 2, 3},
		TestCase{-4, -2, -2, -1},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d/%d", tc.n, tc.d), func(t *testing.T) {
			if n, d := Factor(tc.n, tc.d); n != tc.en || d != tc.ed {
				t.Errorf("Factor(%d,%d) = %d,%d, want %d,%d",
					tc.n, tc.d, n, d, tc.en, tc.ed)
			}
		})
	}

	testutils.AssertPanic(t, "too large", func() { Factor(999, 999) })
}
