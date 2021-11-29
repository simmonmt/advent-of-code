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
