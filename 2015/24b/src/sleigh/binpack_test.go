package sleigh

import (
	"fmt"
	"reflect"
	"testing"
)

func TestOneBinPack(t *testing.T) {
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
				out := OneBinPack(tc.values, tc.cap)
				if !reflect.DeepEqual(tc.expected, out) {
					t.Errorf("OneBinPack(%v, %v) = %v, want %v", tc.values, tc.cap, out, tc.expected)
				}
			})
	}
}

func TestAllBinPacks(t *testing.T) {
	type testCase struct {
		values   []int
		cap      int
		expected [][]int
	}

	testCases := []testCase{
		testCase{[]int{1}, 1, [][]int{[]int{1}}},
		testCase{[]int{2}, 1, [][]int{}},

		testCase{[]int{2, 1}, 1, [][]int{[]int{1}}},
		testCase{[]int{13, 4, 1, 43, 5, 8}, 9, [][]int{[]int{4, 5}, []int{1, 8}}},
		testCase{[]int{13, 4, 1, 43, 3, 2, 6}, 9,
			[][]int{[]int{4, 3, 2}, []int{1, 2, 6}, []int{3, 6}}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("values=%v, cap=%v", tc.values, tc.cap),
			func(t *testing.T) {
				out := AllBinPacks(tc.values, tc.cap)
				if !reflect.DeepEqual(tc.expected, out) {
					t.Errorf("OneBinPack(%v, %v) = %v, want %v", tc.values, tc.cap, out, tc.expected)
				}
			})
	}
}
