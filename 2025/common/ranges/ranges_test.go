package ranges

import (
	"os"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/simmonmt/aoc/2025/common/logger"
)

func TestIncRange(t *testing.T) {
	type TestCase struct {
		A, B     IncRange
		Overlaps bool
		Merged   IncRange
	}

	testCases := []TestCase{
		TestCase{A: IncRange{1, 2}, B: IncRange{3, 4}, Overlaps: false},
		TestCase{A: IncRange{1, 3}, B: IncRange{3, 4}, Overlaps: true, Merged: IncRange{1, 4}},
		TestCase{A: IncRange{1, 4}, B: IncRange{3, 4}, Overlaps: true, Merged: IncRange{1, 4}},
		TestCase{A: IncRange{1, 5}, B: IncRange{3, 4}, Overlaps: true, Merged: IncRange{1, 5}},
		TestCase{A: IncRange{1, 6}, B: IncRange{3, 4}, Overlaps: true, Merged: IncRange{1, 6}},
		TestCase{A: IncRange{2, 6}, B: IncRange{3, 4}, Overlaps: true, Merged: IncRange{2, 6}},
		TestCase{A: IncRange{3, 6}, B: IncRange{3, 4}, Overlaps: true, Merged: IncRange{3, 6}},
		TestCase{A: IncRange{4, 6}, B: IncRange{3, 4}, Overlaps: true, Merged: IncRange{3, 6}},
		TestCase{A: IncRange{5, 6}, B: IncRange{3, 4}, Overlaps: false},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := tc.A.Overlaps(tc.B); got != tc.Overlaps {
				t.Errorf("%v overlaps %v = %v, want %v",
					tc.A, tc.B, got, tc.Overlaps)
			}
			if got := tc.B.Overlaps(tc.A); got != tc.Overlaps {
				t.Errorf("%v overlaps %v = %v, want %v",
					tc.B, tc.A, got, tc.Overlaps)
			}

			if got, ok := tc.A.Merge(tc.B); ok != tc.Overlaps {
				t.Errorf("%v.Merge(%v) = %v, %v, want _, %v",
					tc.A, tc.B, got, ok, tc.Overlaps)
			} else if diff := cmp.Diff(tc.Merged, got); diff != "" {
				t.Errorf("%v.Merge(%v) mismatch; -want,+got:\n%s\n",
					tc.Merged, got, diff)
			}

			if got, ok := tc.B.Merge(tc.A); ok != tc.Overlaps {
				t.Errorf("%v.Merge(%v) = %v, %v, want _, %v",
					tc.B, tc.A, got, ok, tc.Overlaps)
			} else if diff := cmp.Diff(tc.Merged, got); diff != "" {
				t.Errorf("%v.Merge(%v) mismatch; -want,+got:\n%s\n",
					tc.Merged, got, diff)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
