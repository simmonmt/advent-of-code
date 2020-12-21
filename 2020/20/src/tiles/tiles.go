package tiles

import (
	"fmt"

	"github.com/simmonmt/aoc/2020/common/dir"
	"github.com/simmonmt/aoc/2020/common/pos"
)

type Side string

func (s Side) String() string {
	return string(s)
}

func (s Side) Reverse() Side {
	out := make([]rune, len(s))
	for i, r := range s {
		out[len(s)-1-i] = r
	}
	return Side(out)
}

func (s Side) Get(x int) bool {
	if x < 0 || x >= len(s) {
		panic("oob")
	}

	return s[x] == '#'
}

func ParseSide(str string) (Side, error) {
	for _, r := range str {
		if r != '#' && r != '.' {
			return "", fmt.Errorf("bad side string %v", str)
		}
	}
	return Side(str), nil
}

type Tile struct {
	num int
	dim int

	// sides are stored in increasing-coordinate order.
	sides []Side // N S W E
	body  []Side
}

func arraySide(arr []string, idx int) string {
	out := make([]byte, len(arr))
	for y := 0; y < len(arr); y++ {
		a := arr[y][idx]
		out[y] = a
	}
	return string(out)
}

func NewTile(num int, body []string, dim int) (*Tile, error) {
	sides := []Side{}

	sideStrs := []string{
		body[0],                         // N
		body[len(body)-1],               // S
		arraySide(body, 0),              // W
		arraySide(body, len(body[0])-1), // E
	}

	bodySides := []Side{}
	for _, row := range body {
		side, err := ParseSide(row)
		if err != nil {
			return nil, fmt.Errorf("bad side %v: %v", row, err)
		}
		bodySides = append(bodySides, side)
	}

	for _, str := range sideStrs {
		if len(str) != dim {
			return nil, fmt.Errorf("want dim %d found %v", dim, len(str))
		}

		side, err := ParseSide(str)
		if err != nil {
			return nil, fmt.Errorf("bad side %v: %v", str, err)
		}
		sides = append(sides, side)
	}

	return &Tile{
		num:   num,
		dim:   dim,
		sides: sides,
		body:  bodySides,
	}, nil
}

func (t *Tile) Get(pos pos.P2) bool {
	if pos.X < 0 || pos.X >= t.dim || pos.Y < 0 || pos.Y >= t.dim {
		panic("oob access")
	}

	return t.body[pos.Y].Get(pos.X)
}

func (t *Tile) Num() int {
	return t.num
}

func (t *Tile) Dim() int {
	return t.dim
}

func (t *Tile) String() string {
	return fmt.Sprintf("%d: sides N %v E %v S %v W %v", t.num,
		t.Side(dir.DIR_NORTH), t.Side(dir.DIR_EAST),
		t.Side(dir.DIR_SOUTH), t.Side(dir.DIR_WEST))
}

func dirSideOff(d dir.Dir) int {
	switch d {
	case dir.DIR_NORTH:
		return 0
	case dir.DIR_SOUTH:
		return 1
	case dir.DIR_WEST:
		return 2
	case dir.DIR_EAST:
		return 3
	default:
		panic("bad dir")
	}
}

func (t *Tile) Side(d dir.Dir) Side {
	return t.sides[dirSideOff(d)]
}

type OrientedTile struct {
	*Tile
	northSide    dir.Dir // Which side is north
	flipH, flipV bool    // Is it flipped?
}

func NewOrientedTile(tile *Tile, northSide dir.Dir, flipH, flipV bool) *OrientedTile {
	return &OrientedTile{
		Tile:      tile,
		northSide: northSide,
		flipH:     flipH,
		flipV:     flipV,
	}
}

func (ot *OrientedTile) String() string {
	hv := []rune{'_', '_'}
	if ot.flipH {
		hv[0] = 'H'
	}
	if ot.flipV {
		hv[1] = 'V'
	}

	return fmt.Sprintf("%d/%s/%s", ot.Num(), ot.northSide, string(hv))
}

var (
	// This map is used to map from a post-rotation side to where it was
	// pre-rotation. See transform() for more details. The key is the
	// direction that serves as north post-rotation. Index into the value
	// with the dirSideOff() for the side to be unrotated.
	//
	// Example: We've rotated 90 degrees CW, so west now acts as
	// north. We're interested in what's now (post-rotation) the south
	// side. Which side was it pre-rotation? The value to retrieve is
	//
	//     otRotateReverseMap[DIR_WEST][dirToOff(DIR_SOUTH)]
	//
	// In pictures
	//
	//     123          741
	//     456 => CW => 852
	//     789          963
	//
	// We want to know which side '963' was before rotation. The map tells
	// us it was the east side.
	otRotateDirMap = map[dir.Dir][]dir.Dir{
		dir.DIR_NORTH: []dir.Dir{dir.DIR_NORTH, dir.DIR_SOUTH, dir.DIR_WEST, dir.DIR_EAST},
		dir.DIR_WEST:  []dir.Dir{dir.DIR_WEST, dir.DIR_EAST, dir.DIR_SOUTH, dir.DIR_NORTH},
		dir.DIR_SOUTH: []dir.Dir{dir.DIR_SOUTH, dir.DIR_NORTH, dir.DIR_EAST, dir.DIR_WEST},
		dir.DIR_EAST:  []dir.Dir{dir.DIR_EAST, dir.DIR_WEST, dir.DIR_NORTH, dir.DIR_SOUTH},
	}

	// This map is used to indicate whether reversing needs to happen to
	// undo the effects of rotation. See transform() for more details. The
	// key is the direction that serves as north post-rotation. Index into
	// the value with the dirSideOff() for the side to be unrotated. The
	// result indicates whether reversing is required.
	//
	// Example: We've rotated 90 degrees CW, so west now acts as north. We
	// want to find out whether we need to reverse the value of the
	// post-rotation south side (which would've been east pre-rotation). The
	// value to retrieve is
	//
	//     otRotateReverseMap[DIR_WEST][dirToOff(DIR_SOUTH)]
	//
	// In pictures
	//
	//     123          741
	//     456 => CW => 852
	//     789          963
	//                  \-/--- We want this side
	//
	// The '147' side was west, now it's north. We want the '963' side,
	// which is south.
	//
	otRotateReverseMap = map[dir.Dir][]bool{
		dir.DIR_NORTH: []bool{false, false, false, false},
		dir.DIR_WEST:  []bool{true, true, false, false},
		dir.DIR_SOUTH: []bool{true, true, true, true},
		dir.DIR_EAST:  []bool{false, false, true, true},
	}
)

func (ot *OrientedTile) transform(d dir.Dir) (newDir dir.Dir, reverse bool) {
	// We need to peel back the layers of flipping and rotation, mapping the
	// side direction to the corresponding side on the embedded Tile
	// instance (which is neither flipped nor rotated).

	// Flipping has two effects -- it swaps two sides and reverses the other
	// two. Example: A tile with no rotation but which is flipped
	// horizontally (around the y axis) swaps its east and west sides. That
	// is, the east side on the unflipped tile is the west side on the
	// flipped one. The north and south sides are reversed -- a north side
	// of "###..." is "...###" after a horizontal flip.

	// Undo the swapping part of any flipping that may have occurred.
	if ot.flipH {
		if d == dir.DIR_WEST {
			d = dir.DIR_EAST
		} else if d == dir.DIR_EAST {
			d = dir.DIR_WEST
		}
	}
	if ot.flipV {
		if d == dir.DIR_NORTH {
			d = dir.DIR_SOUTH
		} else if d == dir.DIR_SOUTH {
			d = dir.DIR_NORTH
		}
	}

	// Whether the side we're going to read from the Tile needs to be
	// reversed. Negation of flipping can require reversing, as can negation
	// of rotation, and sometimes both can call for reversing. This variable
	// tracks the cumulative effects of the requests for rotation.
	reverse = false

	// Undo the reversing part of any flipping that may have occurred.
	if ot.flipH {
		if d == dir.DIR_NORTH || d == dir.DIR_SOUTH {
			reverse = !reverse
		}
	}
	if ot.flipV {
		if d == dir.DIR_EAST || d == dir.DIR_WEST {
			reverse = !reverse
		}
	}

	// Rotation also has two effects -- it always changes sides and *can*
	// reverse them. Reversal can happen because the Tile stores its sides
	// in increasing-X,increasing-Y order. That is, the west side contains
	// x=0,y=0..max *in that order*. The east side is also stored with
	// y=0..max. If the west side is rotated up into north position it'll
	// need to be reversed before it can be used.
	//
	//     Example:    1..                 321
	//                 2.. => rotate CW => ...
	//                 3..                 ...

	// Undo the reversing part of rotation
	if otRotateReverseMap[ot.northSide][dirSideOff(d)] {
		reverse = !reverse
	}

	// Undo the side-changing part of rotation
	d = otRotateDirMap[ot.northSide][dirSideOff(d)]
	return d, reverse
}

func (ot *OrientedTile) Side(d dir.Dir) Side {
	d, reverse := ot.transform(d)

	side := ot.Tile.Side(d)
	if reverse {
		side = side.Reverse()
	}

	return side
}

func (ot *OrientedTile) Get(pos pos.P2) bool {
	d, reverse := ot.transform(dir.DIR_NORTH)

	rev := func(v int) int {
		return ot.Tile.dim - 1 - v
	}

	maybeRev := func(v int) int {
		if reverse {
			return rev(v)
		}
		return v
	}

	switch d {
	case dir.DIR_SOUTH:
		pos.X, pos.Y = maybeRev(pos.X), rev(pos.Y)
	case dir.DIR_NORTH:
		pos.X, pos.Y = maybeRev(pos.X), pos.Y
	case dir.DIR_WEST:
		pos.X, pos.Y = pos.Y, maybeRev(pos.X)
	case dir.DIR_EAST:
		pos.X, pos.Y = rev(pos.Y), maybeRev(pos.X)
	}

	return ot.Tile.Get(pos)
}

func (ot *OrientedTile) DebugString() string {
	return fmt.Sprintf("%d: sides N %v E %v S %v W %v", ot.Num(),
		ot.Side(dir.DIR_NORTH), ot.Side(dir.DIR_EAST),
		ot.Side(dir.DIR_SOUTH), ot.Side(dir.DIR_WEST))
}
