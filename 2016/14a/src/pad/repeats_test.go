package pad

import (
	"reflect"
	"strconv"
	"testing"
)

func TestHasRepeats(t *testing.T) {
	type TestCase struct {
		in           string
		minLen       int
		expectedReps []rune
	}

	testCases := []TestCase{
		TestCase{"cc388847a5", 2, []rune{'c', '8'}},
		TestCase{"cc388847b5", 3, []rune{'8'}},
		TestCase{"cc388847b5", 4, []rune{}},

		TestCase{"ac388817a5", 2, []rune{'8'}},
		TestCase{"ac388817a5", 3, []rune{'8'}},
		TestCase{"ac388817a5", 4, []rune{}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			reps := HasRepeats(tc.in, tc.minLen)

			if !reflect.DeepEqual(tc.expectedReps, reps) {
				t.Errorf("HasRepeated(%v, %v) = %v, want %v",
					tc.in, tc.minLen, reps, tc.expectedReps)
			}
		})
	}
}
