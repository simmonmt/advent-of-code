package lineio

import (
	"reflect"
	"strconv"
	"testing"
)

func TestBlankSeparatedGroups(t *testing.T) {
	in := []string{"a", "b", "", "d"}
	want := [][]string{[]string{"a", "b"}, []string{"d"}}
	if got := BlankSeparatedGroups(in); !reflect.DeepEqual(got, want) {
		t.Errorf(`BlankSeparatedGroupsFromLines("%v") = %v, want %v`,
			in, got, want)
	}
}

func TestNumbersFromLine(t *testing.T) {
	type TestCase struct {
		in   string
		want []int
	}

	testCases := []TestCase{
		TestCase{"1,2,3,4", []int{1, 2, 3, 4}},
		TestCase{"", nil},
		TestCase{"bad", nil},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := NumbersFromLine(tc.in, ",")
			if tc.want == nil {
				if err == nil {
					t.Errorf(`ParseNumbersFromLine("%v") = %v, %v, want _, non-nil`,
						tc.in, got, err)
				}
			} else {
				if err != nil || !reflect.DeepEqual(got, tc.want) {
					t.Errorf(`ParseNumbersFromLine("%v") = %v, %v, want %v, nil`,
						tc.in, got, err, tc.want)
				}
			}
		})
	}
}
