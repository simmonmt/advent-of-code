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

	"github.com/simmonmt/aoc/2022/common/dir"
	"github.com/simmonmt/aoc/2022/common/grid"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	sample = buildGridOrDie([]string{ //
		"30373", //
		"25512", //
		"65332", //
		"33549", //
		"35390", //
	})
)

func buildGridOrDie(lines []string) *grid.Grid[int] {
	g, err := buildGrid(lines)
	if err != nil {
		panic(fmt.Sprintf("buildGrid = %v", err))
	}
	return g
}

func TestScoreTree(t *testing.T) {
	type TestCase struct {
		p    pos.P2
		want int
	}

	testCases := []TestCase{
		TestCase{pos.P2{2, 1}, 4},
		TestCase{pos.P2{2, 3}, 8},
	}

	for _, tc := range testCases {
		t.Run(tc.p.String(), func(t *testing.T) {
			fmt.Println(tc.p)
			if score := scoreTree(sample, tc.p); score != tc.want {
				t.Errorf("scoreTree(_, %v) = %v, want %v",
					tc.p, score, tc.want)
			}
		})
	}
}

func TestSolveB(t *testing.T) {
	if got := solveB(sample); got != 8 {
		t.Errorf("solveB(_) = %v, want %v", got, 8)
	}
}

func TestLookInDir(t *testing.T) {
	type TestCase struct {
		g    *grid.Grid[int]
		p    pos.P2
		d    dir.Dir
		want int
	}

	testCases := []TestCase{
		TestCase{sample, pos.P2{2, 0}, dir.DIR_NORTH, 0},

		TestCase{sample, pos.P2{2, 1}, dir.DIR_NORTH, 1},
		TestCase{sample, pos.P2{2, 1}, dir.DIR_WEST, 1},
		TestCase{sample, pos.P2{2, 1}, dir.DIR_EAST, 2},
		TestCase{sample, pos.P2{2, 1}, dir.DIR_SOUTH, 2},

		TestCase{sample, pos.P2{2, 3}, dir.DIR_NORTH, 2},
		TestCase{sample, pos.P2{2, 3}, dir.DIR_WEST, 2},
		TestCase{sample, pos.P2{2, 3}, dir.DIR_EAST, 2},
		TestCase{sample, pos.P2{2, 3}, dir.DIR_SOUTH, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.p.String(), func(t *testing.T) {
			fmt.Println(tc.p)
			if got := lookInDir(tc.g, tc.p, tc.d); got != tc.want {
				t.Errorf("lookInDir(_, %v, %v) = %v, want %v",
					tc.p, tc.d, got, tc.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
