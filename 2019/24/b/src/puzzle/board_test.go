package puzzle

import (
	"reflect"
	"strconv"
	"testing"
)

func TestBoard(t *testing.T) {
	lines := []string{".#...", "##...", "#.#..", ".....", "....."}
	b := NewBoard(lines)

	if got := b.Strings(); !reflect.DeepEqual(got, lines) {
		t.Errorf("NewBoard %v, want %v", got, lines)
	}
}

func TestBoardEvolution(t *testing.T) {
	type TestCase struct {
		in   []string
		step [][]string
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
			step: [][]string{
				[]string{
					"#..#.",
					"####.",
					"###.#",
					"##.##",
					".##..",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b := NewBoard(tc.in)
			for _, step := range tc.step {
				nb := b.Evolve()

				if got := nb.Strings(); !reflect.DeepEqual(got, step) {
					t.Errorf("evolve in %v got %v, want %v", b.Strings(), got, step)
				}

				b = nb
			}
		})
	}
}
