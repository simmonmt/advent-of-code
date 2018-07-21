package extent

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
)

func TestExtentParse(t *testing.T) {
	type TestCase struct {
		in       string
		expected *Extent
	}

	testCases := []TestCase{
		TestCase{"1-2", &Extent{1, 2}},
		TestCase{"0-2", &Extent{0, 2}},
		TestCase{fmt.Sprintf("0-%v", uint64(math.MaxUint64)), &Extent{0, math.MaxUint64}},
		TestCase{"", nil},
		TestCase{"1", nil},
		TestCase{"2-1", nil},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			wantErr := tc.expected == nil
			if ext, err := Parse(tc.in); (err != nil) != wantErr || !reflect.DeepEqual(ext, tc.expected) {
				wantErrStr := "nil"
				if wantErr {
					wantErrStr = "err"
				}

				t.Errorf(`Parse("%v") = %v, %v, want %v, %v`,
					tc.in, ext, err, tc.expected, wantErrStr)
			}
		})
	}
}

func TestExtentRemove(t *testing.T) {
	type TestCase struct {
		start    *Extent
		remove   *Extent
		expected []*Extent
	}

	testCases := []TestCase{
		TestCase{&Extent{5, 10}, &Extent{1, 4}, []*Extent{&Extent{5, 10}}},
		TestCase{&Extent{5, 10}, &Extent{1, 5}, []*Extent{&Extent{6, 10}}},
		TestCase{&Extent{5, 10}, &Extent{1, 7}, []*Extent{&Extent{8, 10}}},
		TestCase{&Extent{5, 10}, &Extent{1, 10}, nil},
		TestCase{&Extent{5, 10}, &Extent{1, 11}, nil},
		TestCase{&Extent{5, 10}, &Extent{7, 8}, []*Extent{&Extent{5, 6}, &Extent{9, 10}}},
		TestCase{&Extent{5, 10}, &Extent{6, 8}, []*Extent{&Extent{5, 5}, &Extent{9, 10}}},
		TestCase{&Extent{5, 10}, &Extent{5, 8}, []*Extent{&Extent{9, 10}}},
		TestCase{&Extent{5, 10}, &Extent{7, 9}, []*Extent{&Extent{5, 6}, &Extent{10, 10}}},
		TestCase{&Extent{5, 10}, &Extent{7, 10}, []*Extent{&Extent{5, 6}}},
		TestCase{&Extent{5, 10}, &Extent{7, 11}, []*Extent{&Extent{5, 6}}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if res := tc.start.Remove(tc.remove); !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("%+v.Remove(%+v) = %v, want %v",
					tc.start, tc.remove, res, tc.expected)
			}
		})
	}
}

func TestExtentMerge(t *testing.T) {
	type TestCase struct {
		ext      *Extent
		other    *Extent
		expected *Extent
	}

	testCases := []TestCase{
		TestCase{&Extent{5, 10}, &Extent{1, 2}, nil},
		TestCase{&Extent{5, 10}, &Extent{1, 4}, &Extent{1, 10}},
		TestCase{&Extent{5, 10}, &Extent{1, 5}, &Extent{1, 10}},
		TestCase{&Extent{5, 10}, &Extent{1, 6}, &Extent{1, 10}},
		TestCase{&Extent{5, 10}, &Extent{1, 10}, &Extent{1, 10}},
		TestCase{&Extent{5, 10}, &Extent{1, 11}, &Extent{1, 11}},
		TestCase{&Extent{5, 10}, &Extent{1, 15}, &Extent{1, 15}},
		TestCase{&Extent{5, 10}, &Extent{5, 10}, &Extent{5, 10}},
		TestCase{&Extent{5, 10}, &Extent{6, 10}, &Extent{5, 10}},
		TestCase{&Extent{5, 10}, &Extent{6, 11}, &Extent{5, 11}},
		TestCase{&Extent{5, 10}, &Extent{7, 12}, &Extent{5, 12}},
		TestCase{&Extent{5, 10}, &Extent{9, 12}, &Extent{5, 12}},
		TestCase{&Extent{5, 10}, &Extent{10, 12}, &Extent{5, 12}},
		TestCase{&Extent{5, 10}, &Extent{11, 11}, &Extent{5, 11}},
		TestCase{&Extent{5, 10}, &Extent{12, 12}, nil},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if res := tc.ext.Merge(tc.other); !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("%+v.Merge(%+v) = %v, want %v",
					tc.ext, tc.other, res, tc.expected)
			}
		})
	}
}
