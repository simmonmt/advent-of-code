// Copyright 2021 Google LLC
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
