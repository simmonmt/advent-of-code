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

package maze

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestAllNeighbors(t *testing.T) {
	type TestCase struct {
		start     string
		neighbors []string
	}

	testCases := []TestCase{
		TestCase{"0,0,hijkl,", []string{"0,1,hijkl,D"}},
		TestCase{"0,1,hijkl,D", []string{"0,0,hijkl,DU", "1,1,hijkl,DR"}},
		TestCase{"3,3,hijkl,", []string{"2,3,hijkl,L", "3,2,hijkl,U"}},
	}

	helper := NewHelper(4, 4)

	for _, tc := range testCases {
		t.Run(tc.start, func(t *testing.T) {
			res := helper.AllNeighbors(tc.start)
			sort.Strings(res)

			if !reflect.DeepEqual(res, tc.neighbors) {
				t.Errorf("AllNeighbors(%v) = %v, want %v",
					tc.start, res, tc.neighbors)
			}
		})
	}
}

func TestEstimateDistance(t *testing.T) {
	type TestCase struct {
		start string
		end   string
		dist  uint
	}

	testCases := []TestCase{
		TestCase{"1,1,hijkl,", "3,3,hijkl,", 4},
		TestCase{"1,1,hijkl,", "0,3,hijkl,", 3},
	}

	helper := NewHelper(4, 4)

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := helper.EstimateDistance(tc.start, tc.end)

			if tc.dist != res {
				t.Errorf("Estimate(%v,%v) = %v, want %v",
					tc.start, tc.end, res, tc.dist)
			}
		})
	}
}

func TestRunMaze(t *testing.T) {
	type TestCase struct {
		w, h     int
		passcode string
		success  bool
		result   string
	}

	testCases := []TestCase{
		TestCase{4, 4, "hijkl", false, ""},
		TestCase{4, 4, "ihgpwlah", true, "DDRRRD"},
		TestCase{4, 4, "kglvqrro", true, "DDUDRLRRUDRD"},
		TestCase{4, 4, "ulqzkmiv", true, "DRURDRUDDLLDLUURRDULRLDUUDDDRR"},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			found, last := RunMaze(tc.w, tc.h, tc.passcode)
			if found != tc.success || (found && last != tc.result) {
				wantResult := "_"
				if tc.success {
					wantResult = tc.result
				}

				t.Errorf(`RunMaze(%v,%v,"%v") = %v,%v, want %v, %v`,
					tc.w, tc.h, tc.passcode, found, last, tc.success, wantResult)
			}
		})
	}
}
