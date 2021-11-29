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

package main

import (
	"strconv"
	"testing"
)

func TestSeekBack(t *testing.T) {
	nums := []int{3, 6, 9, 12, 15, 18}

	type TestCase struct {
		start, goal, want int
	}

	testCases := []TestCase{
		TestCase{1, 2, 0},
		TestCase{0, 2, 0},
		TestCase{3, 4, 1},
		TestCase{3, 3, 0},
		TestCase{4, 7, 2},
		TestCase{4, 6, 1},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := seekBack(nums, tc.start, tc.goal); got != tc.want {
				t.Errorf(`seekBack(%v, %v, %v) = %v, want %v`,
					nums, tc.start, tc.goal, got, tc.want)
			}
		})
	}
}
