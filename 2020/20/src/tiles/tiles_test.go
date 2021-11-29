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

package tiles

import (
	"os"
	"testing"

	"github.com/simmonmt/aoc/2020/common/dir"
	"github.com/simmonmt/aoc/2020/common/logger"
	"github.com/simmonmt/aoc/2020/common/pos"
)

var (
	refTile = mustNewTile([]string{
		"..##.#..#.",
		"##..#.....",
		"#...##..#.",
		"####.#...#",
		"##.##.###.",
		"##...#.###",
		".#.#.#..##",
		"..#....#..",
		"###...#.#.",
		"..###..###",
	})
)

func mustNewTile(body []string) *Tile {
	t, err := NewTile(1, body, len(body[0]))
	if err != nil {
		panic(err)
	}
	return t
}

func TestTile(t *testing.T) {
	testCases := map[dir.Dir]string{
		dir.DIR_NORTH: "..##.#..#.",
		dir.DIR_SOUTH: "..###..###",
		dir.DIR_WEST:  ".#####..#.",
		dir.DIR_EAST:  "...#.##..#",
	}

	for d, want := range testCases {
		if got := refTile.Side(d); got.String() != want {
			t.Errorf("refTile.Side(%v) = %v, want %v", d, got, want)
		}
	}
}

func TestOrientedTile(t *testing.T) {
	type TestCase struct {
		northSide    dir.Dir
		flipH, flipV bool
		sides        map[dir.Dir]string
		rows         []string
	}

	testCases := []TestCase{
		// !h !v
		TestCase{
			northSide: dir.DIR_NORTH,
			flipH:     false,
			flipV:     false,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: "..##.#..#.",
				dir.DIR_SOUTH: "..###..###",
				dir.DIR_WEST:  ".#####..#.",
				dir.DIR_EAST:  "...#.##..#",
			},
			rows: []string{"..##.#..#.", "##..#....."},
		},
		TestCase{
			northSide: dir.DIR_WEST,
			flipH:     false,
			flipV:     false,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: ".#..#####.",
				dir.DIR_SOUTH: "#..##.#...",
				dir.DIR_WEST:  "..###..###",
				dir.DIR_EAST:  "..##.#..#.",
			},
			rows: []string{".#..#####.", ".#.####.#."},
		},
		TestCase{
			northSide: dir.DIR_SOUTH,
			flipH:     false,
			flipV:     false,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: "###..###..",
				dir.DIR_SOUTH: ".#..#.##..",
				dir.DIR_WEST:  "#..##.#...",
				dir.DIR_EAST:  ".#..#####.",
			},
			rows: []string{"###..###..", ".#.#...###"},
		},
		TestCase{
			northSide: dir.DIR_EAST,
			flipH:     false,
			flipV:     false,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: "...#.##..#",
				dir.DIR_SOUTH: ".#####..#.",
				dir.DIR_WEST:  ".#..#.##..",
				dir.DIR_EAST:  "###..###..",
			},
			rows: []string{"...#.##..#", "#.#.###.##"},
		},

		// h !v
		TestCase{
			northSide: dir.DIR_NORTH,
			flipH:     true,
			flipV:     false,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: ".#..#.##..",
				dir.DIR_SOUTH: "###..###..",
				dir.DIR_WEST:  "...#.##..#",
				dir.DIR_EAST:  ".#####..#.",
			},
			rows: []string{
				".#..#.##..",
				".....#..##",
			},
		},
		TestCase{
			northSide: dir.DIR_WEST,
			flipH:     true,
			flipV:     false,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: ".#####..#.",
				dir.DIR_SOUTH: "...#.##..#",
				dir.DIR_WEST:  "..##.#..#.",
				dir.DIR_EAST:  "..###..###",
			},
			rows: []string{
				".#####..#.",
				".#.####.#.",
			},
		},
		TestCase{
			northSide: dir.DIR_SOUTH,
			flipH:     true,
			flipV:     false,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: "..###..###",
				dir.DIR_SOUTH: "..##.#..#.",
				dir.DIR_WEST:  ".#..#####.",
				dir.DIR_EAST:  "#..##.#...",
			},
			rows: []string{
				"..###..###",
				"###...#.#.",
			},
		},
		TestCase{
			northSide: dir.DIR_EAST,
			flipH:     true,
			flipV:     false,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: "#..##.#...",
				dir.DIR_SOUTH: ".#..#####.",
				dir.DIR_WEST:  "###..###..",
				dir.DIR_EAST:  ".#..#.##..",
			},
			rows: []string{
				"#..##.#...",
				"##.###.#.#",
			},
		},

		// !h v
		TestCase{
			northSide: dir.DIR_NORTH,
			flipH:     false,
			flipV:     true,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: "..###..###",
				dir.DIR_SOUTH: "..##.#..#.",
				dir.DIR_WEST:  ".#..#####.",
				dir.DIR_EAST:  "#..##.#...",
			},
			rows: []string{
				"..###..###",
				"###...#.#.",
			},
		},
		TestCase{
			northSide: dir.DIR_WEST,
			flipH:     false,
			flipV:     true,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: "#..##.#...",
				dir.DIR_SOUTH: ".#..#####.",
				dir.DIR_WEST:  "###..###..",
				dir.DIR_EAST:  ".#..#.##..",
			},
			rows: []string{
				"#..##.#...",
				"##.###.#.#",
			},
		},
		TestCase{
			northSide: dir.DIR_SOUTH,
			flipH:     false,
			flipV:     true,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: ".#..#.##..",
				dir.DIR_SOUTH: "###..###..",
				dir.DIR_WEST:  "...#.##..#",
				dir.DIR_EAST:  ".#####..#.",
			},
			rows: []string{
				".#..#.##..",
				".....#..##",
			},
		},
		TestCase{
			northSide: dir.DIR_EAST,
			flipH:     false,
			flipV:     true,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: ".#####..#.",
				dir.DIR_SOUTH: "...#.##..#",
				dir.DIR_WEST:  "..##.#..#.",
				dir.DIR_EAST:  "..###..###",
			},
			rows: []string{
				".#####..#.",
				".#.####.#.",
			},
		},

		// h v
		TestCase{
			northSide: dir.DIR_NORTH,
			flipH:     true,
			flipV:     true,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: "###..###..",
				dir.DIR_SOUTH: ".#..#.##..",
				dir.DIR_WEST:  "#..##.#...",
				dir.DIR_EAST:  ".#..#####.",
			},
			rows: []string{
				"###..###..",
				".#.#...###",
			},
		},
		TestCase{
			northSide: dir.DIR_WEST,
			flipH:     true,
			flipV:     true,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: "...#.##..#",
				dir.DIR_SOUTH: ".#####..#.",
				dir.DIR_WEST:  ".#..#.##..",
				dir.DIR_EAST:  "###..###..",
			},
			rows: []string{
				"...#.##..#",
				"#.#.###.##",
			},
		},
		TestCase{
			northSide: dir.DIR_SOUTH,
			flipH:     true,
			flipV:     true,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: "..##.#..#.",
				dir.DIR_SOUTH: "..###..###",
				dir.DIR_WEST:  ".#####..#.",
				dir.DIR_EAST:  "...#.##..#",
			},
			rows: []string{
				"..##.#..#.",
				"##..#.....",
			},
		},
		TestCase{
			northSide: dir.DIR_EAST,
			flipH:     true,
			flipV:     true,
			sides: map[dir.Dir]string{
				dir.DIR_NORTH: ".#..#####.",
				dir.DIR_SOUTH: "#..##.#...",
				dir.DIR_WEST:  "..###..###",
				dir.DIR_EAST:  "..##.#..#.",
			},
			rows: []string{
				".#..#####.",
				".#.####.#.",
			},
		},
	}

	for _, tc := range testCases {
		ot := NewOrientedTile(refTile, tc.northSide, tc.flipH, tc.flipV)
		t.Run(ot.String(), func(t *testing.T) {
			for d, want := range tc.sides {
				if got := ot.Side(d); got.String() != want {
					t.Errorf("ot %v Side(%v) = %v, want %v",
						ot.String(), d, got, want)
				}
			}

			for y, row := range tc.rows {
				got := ""
				for x := 0; x < refTile.Dim(); x++ {
					if ot.Get(pos.P2{X: x, Y: y}) {
						got += "#"
					} else {
						got += "."
					}
				}

				if got != row {
					t.Errorf("ot %v get y=%d = %v, want %v",
						ot.String(), y, got, row)
				}
			}

		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
