package puzzle

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestBoard(t *testing.T) {
	lines := []string{".#.", "##.", "#.#"}
	b := NewBoard(lines)

	want := []bool{
		false, true, false, //
		true, true, false, //
		true, false, true, //
	}

	if got := b.c; !reflect.DeepEqual(got, want) {
		t.Errorf("NewBoard %v, want %v", got, want)
	}

	wantHash := strings.Join(lines, "")
	if got := b.Hash(); got != wantHash {
		t.Errorf("hash %v want %v", got, wantHash)
	}
}

func TestBoardEvolution(t *testing.T) {
	type TestCase struct {
		in    []string
		steps [][]string
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
			steps: [][]string{
				[]string{
					"#..#.",
					"####.",
					"###.#",
					"##.##",
					".##..",
				},
				[]string{
					"#####",
					"....#",
					"....#",
					"...#.",
					"#.###",
				},
				[]string{
					"#....",
					"####.",
					"...##",
					"#.##.",
					".##.#",
				},
				[]string{
					"####.",
					"....#",
					"##..#",
					".....",
					"##...",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b := NewBoard(tc.in)
			for _, step := range tc.steps {
				nb := b.Evolve()

				if got := nb.Strings(); !reflect.DeepEqual(got, step) {
					t.Errorf("evolve in %v got %v, want %v", b.Strings(), got, step)
				}

				b = nb
			}
		})
	}
}

func TestBoardBiodiversity(t *testing.T) {
	in := []string{
		".....",
		".....",
		".....",
		"#....",
		".#...",
	}

	b := NewBoard(in)
	if got, want := b.Biodiversity(), 2129920; got != want {
		t.Errorf("bio = %v, want %v", got, want)
	}
}
