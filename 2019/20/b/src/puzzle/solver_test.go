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
	"strconv"
	"testing"
)

func TestSolve(t *testing.T) {
	type TestCase struct {
		board *Board
		cost  int
	}

	testCases := []TestCase{
		TestCase{NewBoard(map1), 27},
		TestCase{NewBoard(map2), 397},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			allPaths := FindAllPathsFromAllPortals(tc.board)

			start := tc.board.Gate("AA").PortalOut()
			end := tc.board.Gate("ZZ").GateOut()

			//logger.Init(true)

			cost, found := Solve(tc.board, allPaths, start, end)
			if !found || cost != tc.cost {
				t.Errorf("cost, found = %v, %v, want %v, true",
					cost, found, tc.cost)
			}

			//logger.Init(false)
		})
	}

}
