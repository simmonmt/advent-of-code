package sleigh

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRemoveElemes(t *testing.T) {
	type testCase struct {
		in       []int
		remove   []int
		expected []int
	}

	testCases := []testCase{
		testCase{[]int{1, 2, 3, 4}, []int{2, 3}, []int{1, 4}},
		testCase{[]int{1, 2, 3, 4, 1, 2, 3, 4}, []int{2, 3}, []int{1, 4, 1, 2, 3, 4}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("in=%v,remove=%v", tc.in, tc.remove),
			func(t *testing.T) {
				out := removeElems(tc.in, tc.remove)
				if !reflect.DeepEqual(tc.expected, out) {
					t.Errorf("removeElems(%v, %v) = %v, want %v",
						tc.in, tc.remove, out, tc.expected)
				}
			})
	}
}
