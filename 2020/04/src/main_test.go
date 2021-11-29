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
	"fmt"
	"testing"
)

func TestValidNumber(t *testing.T) {
	type TestCase struct {
		str      string
		digits   int
		min, max uint64
		want     bool
	}

	testCases := []TestCase{
		TestCase{"1980", -1, 1920, 2002, true},
		TestCase{"1920", -1, 1920, 2002, true},
		TestCase{"1919", -1, 1920, 2002, false},
		TestCase{"2002", -1, 1920, 2002, true},
		TestCase{"2003", -1, 1920, 2002, false},
		TestCase{"1980", 3, 1920, 2002, false},
		TestCase{"1980", 4, 1920, 2002, true},
		TestCase{"1980", 5, 1920, 2002, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			if got := validNumber(tc.str, tc.digits, tc.min, tc.max); got != tc.want {
				t.Errorf("validNumber(%v,%v,%v,%v) = %v, want %v",
					tc.str, tc.digits, tc.min, tc.max, got, tc.want)
			}
		})
	}
}
