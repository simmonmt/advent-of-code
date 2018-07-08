package sleigh

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBinPack(t *testing.T) {
	type testCase struct {
		values   []int
		cap      int
		expected []int
	}

	testCases := []testCase{
		testCase{[]int{1}, 1, []int{1}},
		testCase{[]int{2}, 1, nil},

		testCase{[]int{2, 1}, 1, []int{1}},
		testCase{[]int{13, 4, 1, 43, 5, 8}, 9, []int{4, 5}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("values=%v, cap=%v", tc.values, tc.cap),
			func(t *testing.T) {
				out := BinPack(tc.values, tc.cap)
				if !reflect.DeepEqual(tc.expected, out) {
					t.Errorf("BinPack(%v, %v) = %v, want %v", tc.values, tc.cap, out, tc.expected)
				}
			})
	}
}
