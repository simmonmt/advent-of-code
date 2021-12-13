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
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2021/common/grid"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

func parseGrid(in []string) *grid.Grid {
	g := grid.New(len(in[0]), len(in))

	for y, line := range in {
		for x, r := range line {
			p := pos.P2{X: x, Y: y}
			if r == '#' {
				g.Set(p, true)
			}
		}
	}

	return g
}

func TestPerformFold(t *testing.T) {
	type TestCase struct {
		in    []string
		insts []Instruction
		want  []string
	}

	testCases := []TestCase{
		TestCase{
			in: []string{
				"...#..#..#.", // 0
				"....#......", // 1
				"...........", // 2
				"#..........", // 3
				"...#....#.#", // 4
				"...........", // 5
				"...........", // 6
				"...........", // 7
				"...........", // 8
				"...........", // 9
				".#....#.##.", // 10
				"....#......", // 11
				"......#...#", // 12
				"#..........", // 13
				"#.#........", // 14
			},
			insts: []Instruction{
				Instruction{Axis: Y_AXIS, Coord: 7},
			},
			want: []string{
				"#.##..#..#.",
				"#...#......",
				"......#...#",
				"#...#......",
				".#.#..#.###",
				"...........",
				"...........",
			},
		},
		// TestCase{
		// 	in: []string{
		// 		"#.##..#..#.",
		// 		"#...#......",
		// 		"......#...#",
		// 		"#...#......",
		// 		".#.#..#.###",
		// 		"...........",
		// 		"...........",
		// 	},
		// 	insts: []Instruction{
		// 		Instruction{Axis: X_AXIS, Coord: 5},
		// 	},
		// 	want: []string{
		// 		"#####",
		// 		"#...#",
		// 		"#...#",
		// 		"#...#",
		// 		"#####",
		// 		".....",
		// 		".....",
		// 	},
		// },
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g := parseGrid(tc.in)

			for _, inst := range tc.insts {
				g = performFold(g, inst)
			}

			buf := bytes.Buffer{}
			dumpTo(&buf, g)
			got := buf.String()

			wantStr := strings.Join(tc.want, "\n") + "\n"

			if got != wantStr {
				t.Errorf("unexpected board. got:\n%v\nwant:\n%v",
					got, wantStr)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
