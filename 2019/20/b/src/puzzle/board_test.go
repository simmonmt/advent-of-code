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
			pOut: pos.P2{9, 2},
			gOut: pos.P2{9, 1},
			pIn:  pos.P2{-1, -1},
			gIn:  pos.P2{-1, -1},
		},
		Gate{
			name: "BC",
			pOut: pos.P2{2, 8},
			gOut: pos.P2{1, 8},
			pIn:  pos.P2{9, 6},
			gIn:  pos.P2{9, 7},
		},
		Gate{
			name: "DE",
			pOut: pos.P2{2, 13},
			gOut: pos.P2{1, 13},
			pIn:  pos.P2{6, 10},
			gIn:  pos.P2{7, 10},
		},
		Gate{
			name: "FG",
			pOut: pos.P2{2, 15},
			gOut: pos.P2{1, 15},
			pIn:  pos.P2{11, 12},
			gIn:  pos.P2{11, 11},
		},
		Gate{
			name: "ZZ",
			pOut: pos.P2{13, 16},
			gOut: pos.P2{13, 17},
			pIn:  pos.P2{-1, -1},
			gIn:  pos.P2{-1, -1},
		},
	}

	if got := b.Gates(); !reflect.DeepEqual(got, want) {
		t.Errorf("Gates() = %v, want %v", got, want)
	}
}
