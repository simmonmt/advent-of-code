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
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestCountBits(t *testing.T) {
	type TestCase struct {
		val uint32
		num int
	}

	testCases := []TestCase{
		TestCase{0x1, 1},
		TestCase{0x5, 2},
		TestCase{0xf, 4},
		TestCase{0x0fffffff, 28},
		TestCase{0xffffffff, 32},
		TestCase{0xaaaaaaaa, 16},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%x", tc.val), func(t *testing.T) {
			res := countBits(tc.val)
			if res != tc.num {
				t.Errorf("countBits(%x) = %v, want %v", tc.val, res, tc.num)
			}
		})
	}
}

func TestIsOpenSpace(t *testing.T) {
	type TestCase struct {
		x, y   int
		isOpen bool
	}

	testCases := []TestCase{
		TestCase{-1, -1, false},
		TestCase{0, 0, true},
		TestCase{0, 2, false},
		TestCase{1, 2, true},
		TestCase{2, 2, true},
		TestCase{3, 2, true},
		TestCase{4, 2, true},
		TestCase{5, 2, false},
		TestCase{0, 3, false},
		TestCase{1, 3, false},
		TestCase{2, 3, false},
		TestCase{3, 3, true},
		TestCase{4, 3, false},
		TestCase{5, 3, true},
		TestCase{6, 3, false},
	}

	magicNumber := 10

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%x,%x", tc.x, tc.y), func(t *testing.T) {
			res := IsOpenSpace(magicNumber, tc.x, tc.y)
			if res != tc.isOpen {
				t.Errorf("isOpenSpace(%d,%d,%d) = %v, want %v",
					magicNumber, tc.x, tc.y, res, tc.isOpen)
			}
		})
	}
}

func TestAllNeighbors(t *testing.T) {
	type TestCase struct {
		start     string
		neighbors []string
	}

	testCases := []TestCase{
		TestCase{"0,0", []string{"0,1"}},
		TestCase{"3,1", []string{"3,2", "4,1"}},
		TestCase{"5,3", []string{}},
	}

	helper := &aStarHelper{magicNumber: 10}

	for _, tc := range testCases {
		t.Run(tc.start, func(t *testing.T) {
			res := helper.AllNeighbors(tc.start)
			sort.Strings(res)

			if !reflect.DeepEqual(tc.neighbors, res) {
				t.Errorf("AllNeighbors(%v) = %v, want %v",
					tc.start, res, tc.neighbors)
			}
		})
	}
}

func TestEstimate(t *testing.T) {
	type TestCase struct {
		start string
		end   string
		dist  uint
	}

	testCases := []TestCase{
		TestCase{"1,1", "3,3", 4},
		TestCase{"-4,-2", "1,6", 13},
		TestCase{"1,-2", "-4,6", 13},
		TestCase{"-4,6", "1,-2", 13},
		TestCase{"1,6", "-4,-2", 13},
	}

	helper := &aStarHelper{magicNumber: 10}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := helper.Estimate(tc.start, tc.end)

			if tc.dist != res {
				t.Errorf("Estimate(%v,%v) = %v, want %v",
					tc.start, tc.end, res, tc.dist)
			}
		})
	}
}
