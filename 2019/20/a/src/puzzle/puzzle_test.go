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

var (
	map1 = []string{
		// 0000000011111111111
		// 2345678901234567890
		"         A           ", // 0
		"         A           ", // 1
		"  #######.#########  ", // 2
		"  #######.........#  ", // 3
		"  #######.#######.#  ", // 4
		"  #######.#######.#  ", // 5
		"  #######.#######.#  ", // 6
		"  #####  B    ###.#  ", // 7
		"BC...##  C    ###.#  ", // 8
		"  ##.##       ###.#  ", // 9
		"  ##...DE  F  ###.#  ", // 10
		"  #####    G  ###.#  ", // 11
		"  #########.#####.#  ", // 12
		"DE..#######...###.#  ", // 13
		"  #.#########.###.#  ", // 14
		"FG..#########.....#  ", // 15
		"  ###########.#####  ", // 16
		"             Z       ", // 17
		"             Z       ", // 18
	}

	map2 = []string{
		"                   A               ",
		"                   A               ",
		"  #################.#############  ",
		"  #.#...#...................#.#.#  ",
		"  #.#.#.###.###.###.#########.#.#  ",
		"  #.#.#.......#...#.....#.#.#...#  ",
		"  #.#########.###.#####.#.#.###.#  ",
		"  #.............#.#.....#.......#  ",
		"  ###.###########.###.#####.#.#.#  ",
		"  #.....#        A   C    #.#.#.#  ",
		"  #######        S   P    #####.#  ",
		"  #.#...#                 #......VT",
		"  #.#.#.#                 #.#####  ",
		"  #...#.#               YN....#.#  ",
		"  #.###.#                 #####.#  ",
		"DI....#.#                 #.....#  ",
		"  #####.#                 #.###.#  ",
		"ZZ......#               QG....#..AS",
		"  ###.###                 #######  ",
		"JO..#.#.#                 #.....#  ",
		"  #.#.#.#                 ###.#.#  ",
		"  #...#..DI             BU....#..LF",
		"  #####.#                 #.#####  ",
		"YN......#               VT..#....QG",
		"  #.###.#                 #.###.#  ",
		"  #.#...#                 #.....#  ",
		"  ###.###    J L     J    #.#.###  ",
		"  #.....#    O F     P    #.#...#  ",
		"  #.###.#####.#.#####.#####.###.#  ",
		"  #...#.#.#...#.....#.....#.#...#  ",
		"  #.#####.###.###.#.#.#########.#  ",
		"  #...#.#.....#...#.#.#.#.....#.#  ",
		"  #.###.#####.###.###.#.#.#######  ",
		"  #.#.........#...#.............#  ",
		"  #########.###.###.#############  ",
		"           B   J   C               ",
		"           U   P   P               ",
	}
)
