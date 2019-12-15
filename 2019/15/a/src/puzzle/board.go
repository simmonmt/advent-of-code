package puzzle

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/intmath"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type Tile int

const (
	TILE_UNKNOWN Tile = iota
	TILE_WALL
	TILE_OPEN
	TILE_GOAL
)

func (t Tile) String() string {
	switch t {
	case TILE_UNKNOWN:
		return " "
	case TILE_WALL:
		return "#"
	case TILE_OPEN:
		return "."
	case TILE_GOAL:
		return "o"
	default:
		panic("bad tile")
	}
}

type Dir int

const (
	DIR_UNKNOWN Dir = iota
	DIR_NORTH
	DIR_SOUTH
	DIR_WEST
	DIR_EAST
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

type Board struct {
	b        map[pos.P2]Tile
	min, max pos.P2
}

func NewBoard() *Board {
	return &Board{
		b: map[pos.P2]Tile{},
	}
}

func (b *Board) Set(p pos.P2, t Tile) {
	if len(b.b) == 0 {
		b.min = p
		b.max = p
	} else {
		b.min.X = intmath.IntMin(b.min.X, p.X)
		b.min.Y = intmath.IntMin(b.min.Y, p.Y)
		b.max.X = intmath.IntMax(b.max.X, p.X)
		b.max.Y = intmath.IntMax(b.max.Y, p.Y)
	}

	if _, found := b.b[p]; found {
		panic("double set")
	}
	b.b[p] = t
}

func (b *Board) Get(p pos.P2) Tile {
	return b.b[p]
}

func (b *Board) CenterAt(newCenter pos.P2) *Board {
	nb := NewBoard()
	for p, t := range b.b {
		nb.Set(pos.P2{p.X - newCenter.X, p.Y - newCenter.Y}, t)
	}
	return nb
}

func PrintBoard(b *Board, cur pos.P2) {
	for y := b.min.Y; y <= b.max.Y; y++ {
		for x := b.min.X; x <= b.max.X; x++ {
			p := pos.P2{x, y}
			if p.Equals(cur) {
				fmt.Print("D")
			} else {
				fmt.Print(b.Get(p))
			}
		}
		fmt.Println()
	}
}
