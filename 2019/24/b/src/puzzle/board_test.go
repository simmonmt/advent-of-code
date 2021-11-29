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

package puzzle

import (
	"reflect"
	"strconv"
	"testing"
)

func TestBoard(t *testing.T) {
	lines := []string{".#...", "##...", "#..#.", ".....", "....."}
	b := NewBoard(lines)

	want := []string{".#...", "##...", "#.?#.", ".....", "....."}

	if got := b.Strings(); !reflect.DeepEqual(got, want) {
		t.Errorf("NewBoard %v, want %v", got, want)
	}
}

func makeBoardTree(bs []*Board) *Board {
	var bZero *Board

	for i, b := range bs {
		if i != 0 {
			b.up = bs[i-1]
		}
		if b.level == 0 {
			bZero = b
		}
		if i != len(bs)-2 {
			b.down = bs[i+1]
		}
	}

	return bZero
}

func TestBoardEvolution(t *testing.T) {
	type TestCase struct {
		in           []string
		runSteps     int
		expectedBugs int
	}

	testCases := []TestCase{
		TestCase{
			in: []string{
				"....#",
				"#..#.",
				"#..##",
				"..#..",
				"#....",
			},
			runSteps:     10,
			expectedBugs: 99,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b := NewBoard(tc.in)

			for i := 0; i < tc.runSteps; i++ {
				b = b.Evolve()
			}

			b.Dump()

			if got := b.NumBugs(); got != tc.expectedBugs {
				t.Errorf("NumBugs = %d, want %d", got, tc.expectedBugs)
			}
		})
	}
}
