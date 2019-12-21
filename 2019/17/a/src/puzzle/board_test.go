package puzzle

import (
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2019/common/pos"
)

func parseBoardFromStrings(strs []string) *Board {
	b := NewBoard()
	for y := range strs {
		for x, r := range strs[y] {
			p := pos.P2{x, y}
			b.Set(p, r)
		}
	}
	return b
}

func TestFindIntersections(t *testing.T) {
	boardStrs := []string{
		"..#..........",
		"..#..........",
		"#######...###",
		"#.#...#...#.#",
		"#############",
		"..#...#...#..",
		"..#####...#..",
	}

	b := parseBoardFromStrings(boardStrs)

	want := map[pos.P2]bool{
		pos.P2{2, 2}:  true,
		pos.P2{2, 4}:  true,
		pos.P2{6, 4}:  true,
		pos.P2{10, 4}: true,
	}

	if got := FindIntersections(b); !reflect.DeepEqual(want, got) {
		t.Errorf("FindIntersections = %v, want %v", got, want)
	}
}

func TestSumAlignmentParams(t *testing.T) {
	ps := map[pos.P2]bool{
		pos.P2{2, 2}:  true,
		pos.P2{2, 4}:  true,
		pos.P2{6, 4}:  true,
		pos.P2{10, 4}: true,
	}

	if want, got := 76, SumAlignmentParams(ps); want != got {
		t.Errorf("SumAlignmentParams(%v) = %v, want %v", ps, got, want)
	}
}
