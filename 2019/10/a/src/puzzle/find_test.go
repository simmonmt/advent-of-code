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
	"sort"
	"strconv"
	"testing"
)

func TestSmall(t *testing.T) {
	type TestCase struct {
		m       map[Pos]bool
		bestPos Pos
		bestNum int
	}

	testCases := []TestCase{
		TestCase{
			m: ParseMap([]string{
				".#..#",
				".....",
				"#####",
				"....#",
				"...##",
			}),
			bestPos: Pos{3, 4},
			bestNum: 8,
		},
		TestCase{
			m: ParseMap([]string{
				"......#.#.",
				"#..#.#....",
				"..#######.",
				".#.#.###..",
				".#..#.....",
				"..#....#.#",
				"#..#....#.",
				".##.#..###",
				"##...#..#.",
				".#....####",
			}),
			bestPos: Pos{5, 8},
			bestNum: 33,
		},
		TestCase{
			m: ParseMap([]string{
				".#..##.###...#######",
				"##.############..##.",
				".#.######.########.#",
				".###.#######.####.#.",
				"#####.##.#.##.###.##",
				"..#####..#.#########",
				"####################",
				"#.####....###.#.#.##",
				"##.#################",
				"#####.##.###..####..",
				"..######..##.#######",
				"####.##.####...##..#",
				".#####..#.######.###",
				"##...#.##########...",
				"#.##########.#######",
				".####.#.###.###.#.##",
				"....##.##.###..#####",
				".#.#.###########.###",
				"#.#.#.#####.####.###",
				"###.##.####.##.#..##",
			}),
			bestPos: Pos{11, 13},
			bestNum: 210,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if p, b := FindBest(tc.m); !reflect.DeepEqual(p, tc.bestPos) || b != tc.bestNum {
				t.Errorf("FindBest = %v, %d, want %v, %d", p, b, tc.bestPos, tc.bestNum)
			}
		})
	}
}

func TestFindAll(t *testing.T) {
	type TestCase struct {
		m                map[Pos]bool
		ctr              Pos
		expected         []Pos
		expectedSpecific map[int]Pos
	}

	testCases := []TestCase{
		TestCase{
			m: ParseMap([]string{
				"..#..",
				".....",
				"#...#",
				".....",
				"..#..",
			}),
			ctr: Pos{2, 2},
			expected: []Pos{
				Pos{2, 0},
				Pos{4, 2},
				Pos{2, 4},
				Pos{0, 2},
			},
		},
		TestCase{
			m: ParseMap([]string{
				"..#..",
				".....",
				"##..#",
				".....",
				"..#..",
			}),
			ctr: Pos{2, 2},
			expected: []Pos{
				Pos{2, 0},
				Pos{4, 2},
				Pos{2, 4},
				Pos{1, 2},
				Pos{0, 2},
			},
		},
		TestCase{
			m: ParseMap([]string{
				".#....#####...#..",
				"##...##.#####..##",
				"##...#...#.#####.",
				"..#.....X...###..",
				"..#.#.....#....##",
			}),
			ctr: Pos{8, 3},
			expectedSpecific: map[int]Pos{
				0: Pos{8 + 0, 3 - 2},
				1: Pos{8 + 1, 3 - 3},
				2: Pos{8 + 1, 3 - 2},
				3: Pos{8 + 2, 3 - 3},
				4: Pos{8 + 1, 3 - 1},
				5: Pos{8 + 3, 3 - 2},
				6: Pos{8 + 4, 3 - 2},
				7: Pos{8 + 3, 3 - 1},
				8: Pos{8 + 7, 3 - 2},
			},
		},
		TestCase{
			m: ParseMap([]string{
				".#..##.###...#######",
				"##.############..##.",
				".#.######.########.#",
				".###.#######.####.#.",
				"#####.##.#.##.###.##",
				"..#####..#.#########",
				"####################",
				"#.####....###.#.#.##",
				"##.#################",
				"#####.##.###..####..",
				"..######..##.#######",
				"####.##.####...##..#",
				".#####..#.######.###",
				"##...#.##########...",
				"#.##########.#######",
				".####.#.###.###.#.##",
				"....##.##.###..#####",
				".#.#.###########.###",
				"#.#.#.#####.####.###",
				"###.##.####.##.#..##",
			}),
			ctr: Pos{11, 13},
			expectedSpecific: map[int]Pos{
				0:   Pos{11, 12},
				199: Pos{8, 2},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vapes := FindAll(tc.ctr, tc.m)
			if tc.expected != nil {
				if !reflect.DeepEqual(vapes, tc.expected) {
					t.Errorf("FindAll(%v, _) = %v, want %v",
						tc.ctr, vapes, tc.expected)
				}
			}

			if tc.expectedSpecific != nil {
				is := []int{}
				for i := range tc.expectedSpecific {
					is = append(is, i)
				}
				sort.Ints(is)

				for _, i := range is {
					if !reflect.DeepEqual(tc.expectedSpecific[i], vapes[i]) {
						t.Errorf("FindAll idx %d = %v, want %v",
							i, vapes[i], tc.expectedSpecific[i])
					}
				}
			}
		})
	}
}
