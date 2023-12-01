// Copyright 2023 Google LLC
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

package intmath

import (
	"fmt"
	"testing"

	"github.com/simmonmt/aoc/2023/common/testutils"
)

func TestGCD(t *testing.T) {
	type TestCase struct {
		vs []int
		d  int
	}

	testCases := []TestCase{
		TestCase{[]int{4, 6}, 2},
		TestCase{[]int{8, 12}, 4},
		TestCase{[]int{12, 18}, 6},
		TestCase{[]int{9, 9}, 9},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.vs), func(t *testing.T) {
			if got := GCD(tc.vs...); got != tc.d {
				t.Errorf("GCD(%v) = %d, want %d", tc.vs, got, tc.d)
			}
		})
	}

	// Make sure it lets us know when we need to add more primes
	testutils.AssertPanic(t, "too big", func() { GCD(101*103, 107*109) })
}
