package puzzle

import (
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2019/common/pos"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard(map1)

	want := []Gate{
		Gate{
			name: "AA",
			p1:   pos.P2{9, 2},
			g1:   pos.P2{9, 1},
			p2:   pos.P2{-1, -1},
			g2:   pos.P2{-1, -1},
		},
		Gate{
			name: "BC",
			p1:   pos.P2{2, 8},
			g1:   pos.P2{1, 8},
			p2:   pos.P2{9, 6},
			g2:   pos.P2{9, 7},
		},
		Gate{
			name: "DE",
			p1:   pos.P2{2, 13},
			g1:   pos.P2{1, 13},
			p2:   pos.P2{6, 10},
			g2:   pos.P2{7, 10},
		},
		Gate{
			name: "FG",
			p1:   pos.P2{2, 15},
			g1:   pos.P2{1, 15},
			p2:   pos.P2{11, 12},
			g2:   pos.P2{11, 11},
		},
		Gate{
			name: "ZZ",
			p1:   pos.P2{13, 16},
			g1:   pos.P2{13, 17},
			p2:   pos.P2{-1, -1},
			g2:   pos.P2{-1, -1},
		},
	}

	if got := b.Gates(); !reflect.DeepEqual(got, want) {
		t.Errorf("Gates() = %v, want %v", got, want)
	}
}
