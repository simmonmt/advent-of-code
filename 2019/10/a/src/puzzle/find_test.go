package puzzle

import (
	"reflect"
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
