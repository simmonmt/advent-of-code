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

package pad

import (
	"reflect"
	"strconv"
	"testing"
)

func TestHasRepeats(t *testing.T) {
	type TestCase struct {
		in           string
		minLen       int
		expectedReps []rune
	}

	testCases := []TestCase{
		TestCase{"cc388847a5", 2, []rune{'c', '8'}},
		TestCase{"cc388847b5", 3, []rune{'8'}},
		TestCase{"cc388847b5", 4, []rune{}},

		TestCase{"ac388817a5", 2, []rune{'8'}},
		TestCase{"ac388817a5", 3, []rune{'8'}},
		TestCase{"ac388817a5", 4, []rune{}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			reps := HasRepeats(tc.in, tc.minLen)

			if !reflect.DeepEqual(tc.expectedReps, reps) {
				t.Errorf("HasRepeated(%v, %v) = %v, want %v",
					tc.in, tc.minLen, reps, tc.expectedReps)
			}
		})
	}
}
