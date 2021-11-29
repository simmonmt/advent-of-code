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
	"os"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2020/common/logger"
)

func TestFindCumAlignment(t *testing.T) {
	type TestCase struct {
		in   []int
		want int64
	}

	testCases := []TestCase{
		TestCase{in: []int{17, -1, 13, 19}, want: 3417},
		TestCase{in: []int{67, 7, 59, 61}, want: 754018},
		TestCase{in: []int{67, -1, 7, 59, 61}, want: 779210},
		TestCase{in: []int{67, 7, -1, 59, 61}, want: 1261476},
		TestCase{in: []int{1789, 37, 47, 1889}, want: 1202161486},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ca := findCumAlignment(tc.in)
			if ca.first != tc.want {
				t.Errorf("findCumAlignment(%v) = %v, want %v",
					tc.in, ca.first, tc.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
