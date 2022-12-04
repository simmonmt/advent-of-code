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

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/simmonmt/aoc/2022/common/logger"
)

func TestOverlap(t *testing.T) {
	type TestCase struct {
		r, o Range
		want bool
	}

	testCases := []TestCase{
		TestCase{Range{0, 1}, Range{2, 3}, false},
		TestCase{Range{0, 4}, Range{2, 3}, true},
		TestCase{Range{1, 2}, Range{2, 3}, true},
		TestCase{Range{1, 3}, Range{2, 3}, true},
		TestCase{Range{2, 3}, Range{2, 3}, true},
		TestCase{Range{2, 4}, Range{2, 3}, true},
		TestCase{Range{3, 4}, Range{2, 3}, true},
		TestCase{Range{4, 5}, Range{2, 3}, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%d/%d_%d", tc.r.From, tc.r.To, tc.o.From, tc.o.To),
			func(t *testing.T) {
				if got := tc.r.Overlaps(tc.o); got != tc.want {
					t.Errorf("%v.Overlaps(%v) = %v, want %v",
						tc.r, tc.o, got, tc.want)
				}
				if got := tc.o.Overlaps(tc.r); got != tc.want {
					t.Errorf("%v.Overlaps(%v) = %v, want %v",
						tc.o, tc.r, got, tc.want)
				}
			})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
