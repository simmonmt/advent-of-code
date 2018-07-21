package extent

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestExtentsSort(t *testing.T) {
	extents := []*Extent{&Extent{5, 8}, &Extent{0, 2}, &Extent{4, 7}}
	sort.Sort(Extents(extents))

	expected := []*Extent{&Extent{0, 2}, &Extent{4, 7}, &Extent{5, 8}}
	if !reflect.DeepEqual(extents, expected) {
		t.Errorf("sort = %+v, expected %+v", extents, expected)
	}
}

func TestExtentsMerge(t *testing.T) {
	type TestCase struct {
		in       Extents
		expected Extents
	}

	testCases := []TestCase{
		TestCase{
			in:       Extents{&Extent{5, 10}},
			expected: Extents{&Extent{5, 10}},
		},
		TestCase{
			in:       Extents{&Extent{5, 10}, &Extent{15, 20}},
			expected: Extents{&Extent{5, 10}, &Extent{15, 20}},
		},
		TestCase{
			in:       Extents{&Extent{1, 4}, &Extent{5, 10}},
			expected: Extents{&Extent{1, 10}},
		},
		TestCase{
			in:       Extents{&Extent{1, 5}, &Extent{5, 10}},
			expected: Extents{&Extent{1, 10}},
		},
		TestCase{
			in:       Extents{&Extent{1, 6}, &Extent{5, 10}},
			expected: Extents{&Extent{1, 10}},
		},
		TestCase{
			in:       Extents{&Extent{1, 4}, &Extent{5, 10}, &Extent{11, 15}, &Extent{20, 25}},
			expected: Extents{&Extent{1, 15}, &Extent{20, 25}},
		},
		TestCase{
			in:       Extents{&Extent{0, 2}, &Extent{4, 5}, &Extent{5, 10}, &Extent{11, 15}, &Extent{20, 25}},
			expected: Extents{&Extent{0, 2}, &Extent{4, 15}, &Extent{20, 25}},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if res := tc.in.Merge(); !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("%v Merge = %v, want %v", tc.in, res, tc.expected)
			}
		})
	}
}
