package puzzle

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2019/common/pos"
)

func TestFindAllPathsFromPortal(t *testing.T) {
	board := NewBoard(map1)

	type TestCase struct {
		srcName string
		srcPos  pos.P2
		dests   map[pos.P2]int
	}

	testCases := []TestCase{
		TestCase{
			srcName: "AA",
			srcPos:  board.Gate("AA").Portal1(),
			dests: map[pos.P2]int{
				board.Gate("BC").Gate2(): 5,
				board.Gate("FG").Gate2(): 31,
				board.Gate("ZZ").Gate1(): 26 + 1,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			paths := FindAllPathsFromPortal(tc.srcName, tc.srcPos, board)

			got := map[pos.P2]int{}
			for _, path := range paths {
				got[path.DestPos] = path.Dist
			}

			if !reflect.DeepEqual(got, tc.dests) {
				t.Errorf("got %v want %v", got, tc.dests)
			}
		})
	}
}

func whatIs(board *Board, p pos.P2) string {
	for _, gate := range board.Gates() {
		for i, pp := range []pos.P2{gate.Portal1(), gate.Portal2()} {
			if pp.Equals(p) {
				return fmt.Sprintf("%s:p%d", gate.name, i+1)
			}
		}
		for i, gp := range []pos.P2{gate.Gate1(), gate.Gate2()} {
			if gp.Equals(p) {
				return fmt.Sprintf("%s:g%d", gate.name, i+1)
			}
		}
	}
	panic("unknown")
}

func TestFindAllPathsFromAllPortals(t *testing.T) {
	board := NewBoard(map1)

	expected := map[pos.P2]map[pos.P2]int{
		board.Gate("AA").Portal1(): map[pos.P2]int{
			board.Gate("BC").Gate2(): 5,
			board.Gate("FG").Gate2(): 31,
			board.Gate("ZZ").Gate1(): 27,
		},
		board.Gate("BC").Portal1(): map[pos.P2]int{
			board.Gate("DE").Gate2(): 7,
		},
		board.Gate("BC").Portal2(): map[pos.P2]int{
			board.Gate("AA").Gate1(): 5,
			board.Gate("FG").Gate2(): 33,
			board.Gate("ZZ").Gate1(): 29,
		},
		board.Gate("DE").Portal1(): map[pos.P2]int{
			board.Gate("FG").Gate1(): 5,
		},
		board.Gate("DE").Portal2(): map[pos.P2]int{
			board.Gate("BC").Gate1(): 7,
		},
		board.Gate("FG").Portal1(): map[pos.P2]int{
			board.Gate("DE").Gate1(): 5,
		},
		board.Gate("FG").Portal2(): map[pos.P2]int{
			board.Gate("AA").Gate1(): 31,
			board.Gate("BC").Gate2(): 33,
			board.Gate("ZZ").Gate1(): 7,
		},
		board.Gate("ZZ").Portal1(): map[pos.P2]int{
			board.Gate("AA").Gate1(): 27,
			board.Gate("BC").Gate2(): 29,
			board.Gate("FG").Gate2(): 7,
		},
	}

	gotPathMap := FindAllPathsFromAllPortals(board)

	for from := range gotPathMap {
		if _, found := expected[from]; !found {
			t.Errorf("expected has no %v (%s)", from, whatIs(board, from))
		}
	}

	for from, wantPaths := range expected {
		gotPathsForPortal, found := gotPathMap[from]
		if !found {
			t.Errorf("got has no %v (%s)", from, whatIs(board, from))
			continue
		}

		gotPaths := map[pos.P2]int{}
		for _, path := range gotPathsForPortal {
			gotPaths[path.DestPos] = path.Dist
		}

		if !reflect.DeepEqual(gotPaths, wantPaths) {
			t.Errorf("%v (%s): got %v want %v",
				from, whatIs(board, from), gotPaths, wantPaths)
		}
	}
}
