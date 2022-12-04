package area

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2022/common/pos"
)

func TestParseArea1D(t *testing.T) {
	type TestCase struct {
		s    string
		want *Area1D
	}

	testCases := []TestCase{
		TestCase{"", nil},
		TestCase{"1-2", &Area1D{1, 2}},
		TestCase{"1", nil},
		TestCase{"2-1", nil},
		TestCase{"bob-sue", nil},
		TestCase{"bob-1", nil},
		TestCase{"1-sue", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			got, err := ParseArea1D(tc.s)
			errMismatch := (err == nil) == (tc.want == nil)
			if errMismatch || (tc.want != nil && !reflect.DeepEqual(got, *tc.want)) {
				wantArea, wantErr := "_", "nil"
				if tc.want == nil {
					wantErr = "non-nil"
				} else {
					wantArea = tc.want.String()
				}

				t.Errorf(`ParseArea1D("%s") = %v, %v, want %v, %v`,
					tc.s, got, err, wantArea, wantErr)
			}
		})
	}
}

func TestArea1DContains(t *testing.T) {
	type TestCase struct {
		one, two Area1D
		want     bool
	}

	testCases := []TestCase{
		TestCase{Area1D{0, 1}, Area1D{2, 3}, false},
		TestCase{Area1D{0, 4}, Area1D{2, 3}, true},
		TestCase{Area1D{1, 2}, Area1D{2, 3}, false},
		TestCase{Area1D{1, 3}, Area1D{2, 3}, true},
		TestCase{Area1D{2, 3}, Area1D{2, 3}, true},
		TestCase{Area1D{2, 4}, Area1D{2, 3}, true},
		TestCase{Area1D{3, 4}, Area1D{2, 3}, false},
		TestCase{Area1D{4, 5}, Area1D{2, 3}, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s/%s", tc.one, tc.two), func(t *testing.T) {
			if got := tc.one.Contains(tc.two); got != tc.want {
				t.Errorf("%v.Contains(%v) = %v, want %v",
					tc.one, tc.two, got, tc.want)
			}
		})
	}
}

func TestArea1DOverlaps(t *testing.T) {
	type TestCase struct {
		one, two Area1D
		want     bool
	}

	testCases := []TestCase{
		TestCase{Area1D{0, 1}, Area1D{2, 3}, false},
		TestCase{Area1D{0, 4}, Area1D{2, 3}, true},
		TestCase{Area1D{1, 2}, Area1D{2, 3}, true},
		TestCase{Area1D{1, 3}, Area1D{2, 3}, true},
		TestCase{Area1D{2, 3}, Area1D{2, 3}, true},
		TestCase{Area1D{2, 4}, Area1D{2, 3}, true},
		TestCase{Area1D{3, 4}, Area1D{2, 3}, true},
		TestCase{Area1D{4, 5}, Area1D{2, 3}, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s/%s", tc.one, tc.two),
			func(t *testing.T) {
				if got := tc.one.Overlaps(tc.two); got != tc.want {
					t.Errorf("%v.Overlaps(%v) = %v, want %v",
						tc.one, tc.two, got, tc.want)
				}
				if got := tc.two.Overlaps(tc.one); got != tc.want {
					t.Errorf("%v.Overlaps(%v) = %v, want %v",
						tc.two, tc.one, got, tc.want)
				}
			})
	}
}

func TestArea3DContains(t *testing.T) {
	type TestCase struct {
		one, two Area3D
		want     bool
	}

	testCases := []TestCase{
		TestCase{ // completely disjoint
			Area3D{pos.P3{-5, 0, -3}, pos.P3{-3, 2, 5}},
			Area3D{pos.P3{10, 10, 10}, pos.P3{12, 12, 12}},
			false,
		},
		TestCase{ // containing
			Area3D{pos.P3{-1, -1, -1}, pos.P3{7, 8, 9}},
			Area3D{pos.P3{1, 2, 3}, pos.P3{2, 4, 6}},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s/%s", tc.one, tc.two),
			func(t *testing.T) {
				if got := tc.one.Contains(tc.two); got != tc.want {
					t.Errorf("%v.Contains(%v) = %v, want %v",
						tc.one, tc.two, got, tc.want)
				}
			})
	}
}

func TestArea3DOverlaps(t *testing.T) {
	type TestCase struct {
		one, two Area3D
		want     bool
	}

	testCases := []TestCase{
		TestCase{ // completely disjoint
			Area3D{pos.P3{-5, 0, -3}, pos.P3{-3, 2, 5}},
			Area3D{pos.P3{10, 10, 10}, pos.P3{12, 12, 12}},
			false,
		},
		TestCase{ // containing
			Area3D{pos.P3{1, 2, 3}, pos.P3{2, 4, 6}},
			Area3D{pos.P3{-1, -1, -1}, pos.P3{7, 8, 9}},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s/%s", tc.one, tc.two),
			func(t *testing.T) {
				if got := tc.one.Overlaps(tc.two); got != tc.want {
					t.Errorf("%v.Overlaps(%v) = %v, want %v",
						tc.one, tc.two, got, tc.want)
				}
				if got := tc.two.Overlaps(tc.one); got != tc.want {
					t.Errorf("%v.Overlaps(%v) = %v, want %v",
						tc.two, tc.one, got, tc.want)
				}
			})
	}
}
