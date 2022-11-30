// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filereader

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestLinesFromReader(t *testing.T) {
	in := "a\nb\n\nd\n"
	r := strings.NewReader(in)
	want := []string{"a", "b", "", "d"}
	if got, err := linesFromReader(r); err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf(`linesFromReader("%v") = %v, %v, want %v, nil`,
			strconv.Quote(in), got, err, want)
	}
}

func TestBlankSeparatedGroupsFromReader(t *testing.T) {
	in := "a\nb\n\nd\n"
	r := strings.NewReader(in)
	want := [][]string{[]string{"a", "b"}, []string{"d"}}
	if got, err := blankSeparatedGroupsFromReader(r); err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf(`blankSeparatedGroupsFromReader("%v") = %v, %v, want %v, nil`,
			strconv.Quote(in), got, err, want)
	}
}

func TestBlankSeparatedGroupsFromLines(t *testing.T) {
	in := []string{"a", "b", "", "d"}
	want := [][]string{[]string{"a", "b"}, []string{"d"}}
	if got, err := BlankSeparatedGroupsFromLines(in); err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf(`BlankSeparatedGroupsFromLines("%v") = %v, %v, want %v, nil`,
			in, got, err, want)
	}
}

func TestParseNumbersFromLine(t *testing.T) {
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
			got, err := ParseNumbersFromLine(tc.in)
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
