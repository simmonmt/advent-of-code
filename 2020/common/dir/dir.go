package dir

import "github.com/simmonmt/aoc/2020/common/pos"

type Dir int

const (
	DIR_UNKNOWN Dir = iota
	DIR_NORTH
	DIR_SOUTH
	DIR_WEST
	DIR_EAST
)

var (
	AllDirs = []Dir{DIR_NORTH, DIR_SOUTH, DIR_WEST, DIR_EAST}
)

func (d Dir) String() string {
	switch d {
	case DIR_NORTH:
		return "N"
	case DIR_SOUTH:
		return "S"
	case DIR_WEST:
		return "W"
	case DIR_EAST:
		return "E"
	default:
		panic("bad dir")
	}
}

func (d Dir) Reverse() Dir {
	switch d {
	case DIR_NORTH:
		return DIR_SOUTH
	case DIR_SOUTH:
		return DIR_NORTH
	case DIR_WEST:
		return DIR_EAST
	case DIR_EAST:
		return DIR_WEST
	default:
		panic("bad dir")
	}
}

func (d Dir) From(p pos.P2) pos.P2 {
	switch d {
	case DIR_NORTH:
		return pos.P2{X: p.X, Y: p.Y - 1}
	case DIR_SOUTH:
		return pos.P2{X: p.X, Y: p.Y + 1}
	case DIR_EAST:
		return pos.P2{X: p.X + 1, Y: p.Y}
	case DIR_WEST:
		return pos.P2{X: p.X - 1, Y: p.Y}
	default:
		panic("bad dir")
	}
}
