package main

import (
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/simmonmt/aoc/2024/common/dir"
	"github.com/simmonmt/aoc/2024/common/grid"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/pos"
	"github.com/simmonmt/aoc/2024/common/testutils"
)

var (
	//go:embed combined_samples.txt
	rawSample       string
	sampleTestCases = []testutils.SampleTestCase{
		testutils.SampleTestCase{ // sample_1.txt
			WantA: 2028, WantB: -1,
		},
		testutils.SampleTestCase{ // sample_2.txt
			WantA: 10092, WantB: 9021,
		},
	}
)

func TestParseInput(t *testing.T) {
	for _, tc := range sampleTestCases {
		if tc.WantInput == nil {
			continue
		}

		t.Run(tc.File, func(t *testing.T) {
			input, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.WantInput, input); diff != "" {
				t.Errorf("parseInput mismatch; -want,+got:\n%s\n", diff)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	for _, tc := range sampleTestCases {
		if tc.WantA == -1 {
			continue
		}

		t.Run(tc.File, func(t *testing.T) {
			input, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			if got := solveA(input); got != tc.WantA {
				t.Errorf("solveA(sample) = %v, want %v", got, tc.WantA)
			}
		})
	}
}

func TestSolveB(t *testing.T) {
	for _, tc := range sampleTestCases {
		if tc.WantB == -1 {
			continue
		}

		t.Run(tc.File, func(t *testing.T) {
			input, err := parseInput(tc.Body)
			if err != nil {
				t.Fatal(err)
			}

			if got := solveB(input); got != tc.WantB {
				t.Errorf("solveB(sample) = %v, want %v", got, tc.WantB)
			}
		})
	}
}

func makeTestGrid(lines []string) (*grid.Grid[rune], map[pos.P2]*Box) {
	g, err := grid.NewFromLines[rune](lines, func(p pos.P2, r rune) (rune, error) {
		return r, nil
	})
	if err != nil {
		panic("bad grid")
	}

	return g, findBoxes(g)
}

func TestMoveEW(t *testing.T) {
	g, boxes := makeTestGrid([]string{
		// 234567890123
		"##############",
		"##[]......[]##",
		"##[][]..[][]##",
		"##############",
	})

	type TestCase struct {
		box                       pos.P2
		blockedLeft, blockedRight bool
		toMoveLeft, toMoveRight   []pos.P2
	}

	testCases := []TestCase{
		TestCase{
			box:         pos.P2{X: 2, Y: 1},
			blockedLeft: true, blockedRight: false,
			toMoveLeft: nil, toMoveRight: nil,
		},
		TestCase{
			box:         pos.P2{X: 2, Y: 2},
			blockedLeft: true, blockedRight: false,
			toMoveLeft: nil, toMoveRight: []pos.P2{pos.P2{X: 4, Y: 2}},
		},
		TestCase{
			box:         pos.P2{X: 4, Y: 2},
			blockedLeft: true, blockedRight: false,
			toMoveLeft: nil, toMoveRight: nil,
		},
		TestCase{
			box:         pos.P2{X: 10, Y: 1},
			blockedLeft: false, blockedRight: true,
			toMoveLeft: nil, toMoveRight: nil,
		},
		TestCase{
			box:         pos.P2{X: 10, Y: 2},
			blockedLeft: false, blockedRight: true,
			toMoveLeft: []pos.P2{pos.P2{X: 8, Y: 2}}, toMoveRight: nil,
		},
		TestCase{
			box:         pos.P2{X: 8, Y: 2},
			blockedLeft: false, blockedRight: true,
			toMoveLeft: nil, toMoveRight: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%d", tc.box.X, tc.box.Y), func(t *testing.T) {
			b, found := boxes[tc.box]
			if !found {
				t.Fatalf("can't find box")
			}

			var toMoveLeft []*Box
			for _, p := range tc.toMoveLeft {
				toMoveLeft = append(toMoveLeft, boxes[p])
			}

			if blocked, toMove := b.Blocked(g, boxes, dir.DIR_WEST); blocked != tc.blockedLeft {
				t.Errorf("Blocked(DIR_WEST) = %v,%v; want %v,_", blocked, toMove,
					tc.blockedLeft)
			} else if diff := cmp.Diff(toMoveLeft, toMove); diff != "" {
				t.Errorf("toMoveLeft mismatch; -want,+got:\n%s\n", diff)
			}

			var toMoveRight []*Box
			for _, p := range tc.toMoveRight {
				toMoveRight = append(toMoveRight, boxes[p])
			}

			if blocked, toMove := b.Blocked(g, boxes, dir.DIR_EAST); blocked != tc.blockedRight {
				t.Errorf("Blocked(DIR_EAST) = %v,%v; want %v,_", blocked, toMove,
					tc.blockedRight)
			} else if diff := cmp.Diff(toMoveRight, toMove); diff != "" {
				t.Errorf("toMoveRight mismatch; -want,+got:\n%s\n", diff)
			}
		})
	}
}

func TestMoveNS(t *testing.T) {
	g, boxes := makeTestGrid([]string{
		// 0000000011111111112222222222
		// 2345678901234567890123456789
		"################################", // 0
		"##[]..[].......[].[][].....[].##", // 1
		"##....[][][][][][].[]..[].[][]##", // 2
		"##..[]..[].[].........[][].[].##", // 3
		"##########################[]..##", // 4
		"##............................##", // 5
		"################################",
	})

	type TestCase struct {
		box                    pos.P2
		blockedUp, blockedDown bool
		toMoveUp, toMoveDown   []pos.P2
	}

	testCases := []TestCase{
		TestCase{
			box:       pos.P2{X: 2, Y: 1},
			blockedUp: true, blockedDown: false,
			toMoveUp: nil, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 4, Y: 3},
			blockedUp: false, blockedDown: true,
			toMoveUp: nil, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 6, Y: 1},
			blockedUp: true, blockedDown: false,
			toMoveUp: nil, toMoveDown: []pos.P2{pos.P2{X: 6, Y: 2}},
		},
		TestCase{
			box:       pos.P2{X: 6, Y: 2},
			blockedUp: true, blockedDown: false,
			toMoveUp: nil, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 8, Y: 3},
			blockedUp: false, blockedDown: true,
			toMoveUp: []pos.P2{pos.P2{X: 8, Y: 2}}, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 8, Y: 2},
			blockedUp: false, blockedDown: true,
			toMoveUp: nil, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 10, Y: 2},
			blockedUp: false, blockedDown: true,
			toMoveUp: nil, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 12, Y: 2},
			blockedUp: false, blockedDown: true,
			toMoveUp: nil, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 11, Y: 3},
			blockedUp: false, blockedDown: true,
			toMoveUp: []pos.P2{pos.P2{X: 10, Y: 2}, pos.P2{X: 12, Y: 2}}, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 14, Y: 2},
			blockedUp: true, blockedDown: false,
			toMoveUp: nil, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 16, Y: 2},
			blockedUp: true, blockedDown: false,
			toMoveUp: nil, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 15, Y: 1},
			blockedUp: true, blockedDown: false,
			toMoveUp: nil, toMoveDown: []pos.P2{pos.P2{X: 14, Y: 2}, pos.P2{X: 16, Y: 2}},
		},
		TestCase{
			box:       pos.P2{X: 18, Y: 1},
			blockedUp: true, blockedDown: false,
			toMoveUp: nil, toMoveDown: []pos.P2{pos.P2{X: 19, Y: 2}},
		},
		TestCase{
			box:       pos.P2{X: 20, Y: 1},
			blockedUp: true, blockedDown: false,
			toMoveUp: nil, toMoveDown: []pos.P2{pos.P2{X: 19, Y: 2}},
		},
		TestCase{
			box:       pos.P2{X: 19, Y: 2},
			blockedUp: true, blockedDown: false,
			toMoveUp: nil, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 22, Y: 3},
			blockedUp: false, blockedDown: true,
			toMoveUp: []pos.P2{pos.P2{X: 23, Y: 2}}, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 24, Y: 3},
			blockedUp: false, blockedDown: true,
			toMoveUp: []pos.P2{pos.P2{X: 23, Y: 2}}, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 23, Y: 2},
			blockedUp: false, blockedDown: true,
			toMoveUp: nil, toMoveDown: nil,
		},
		TestCase{
			box:       pos.P2{X: 27, Y: 1},
			blockedUp: true, blockedDown: false,
			toMoveUp: nil,
			toMoveDown: []pos.P2{
				pos.P2{X: 26, Y: 2}, pos.P2{X: 28, Y: 2},
				pos.P2{X: 27, Y: 3}, pos.P2{X: 26, Y: 4},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%d", tc.box.X, tc.box.Y), func(t *testing.T) {
			b, found := boxes[tc.box]
			if !found {
				t.Fatalf("can't find box")
			}

			var toMoveUp []*Box
			for _, p := range tc.toMoveUp {
				toMoveUp = append(toMoveUp, boxes[p])
			}

			if blocked, toMove := b.Blocked(g, boxes, dir.DIR_NORTH); blocked != tc.blockedUp {
				t.Errorf("Blocked(DIR_NORTH) = %v,%v; want %v,_", blocked, toMove,
					tc.blockedUp)
			} else if diff := cmp.Diff(toMoveUp, toMove); diff != "" {
				t.Errorf("toMoveUp mismatch; -want,+got:\n%s\n", diff)
			}

			var toMoveDown []*Box
			for _, p := range tc.toMoveDown {
				toMoveDown = append(toMoveDown, boxes[p])
			}

			if blocked, toMove := b.Blocked(g, boxes, dir.DIR_SOUTH); blocked != tc.blockedDown {
				t.Errorf("Blocked(DIR_SOUTH) = %v,%v; want %v,_", blocked, toMove,
					tc.blockedDown)
			} else if diff := cmp.Diff(toMoveDown, toMove); diff != "" {
				t.Errorf("toMoveDown mismatch; -want,+got:\n%s\n", diff)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	testutils.PopulateTestCases(rawSample, sampleTestCases)
	os.Exit(m.Run())
}
