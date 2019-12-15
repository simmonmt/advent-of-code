package puzzle

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/intmath"
)

type Tile int

const (
	TILE_EMPTY Tile = iota
	TILE_WALL
	TILE_BLOCK
	TILE_HPADDLE
	TILE_BALL
)

func (t Tile) String() string {
	switch t {
	case TILE_EMPTY:
		return " "
	case TILE_WALL:
		return "#"
	case TILE_BLOCK:
		return "X"
	case TILE_HPADDLE:
		return "-"
	case TILE_BALL:
		return "o"
	default:
		return "?"
	}
}

type Board struct {
	a    []Tile
	H, W int
}

func (b *Board) Get(x, y int) Tile {
	return b.a[y*b.W+x]
}

func (b *Board) Set(x, y int, t Tile) {
	b.a[y*b.W+x] = t
	//fmt.Printf("%d %d %d\n", x, y, t)
}

func NewBoard(h, w int) *Board {
	return &Board{
		a: make([]Tile, h*w),
		H: h,
		W: w,
	}
}

func MakeBoard(vs []int64) *Board {
	maxX, maxY := 0, 0
	for i := 0; i < len(vs); i += 3 {
		x, y := int(vs[i]), int(vs[i+1])
		maxX = intmath.IntMax(maxX, x)
		maxY = intmath.IntMax(maxY, y)
	}

	b := NewBoard(maxY+1, maxX+1)

	for i := 0; i < len(vs); i += 3 {
		x, y, t := int(vs[i]), int(vs[i+1]), Tile(vs[i+2])
		b.Set(x, y, t)
	}

	return b
}

func PrintBoard(b *Board) {
	fmt.Printf("     ")
	for x := 0; x < b.W; x++ {
		fmt.Printf("%d", x%10)
	}
	fmt.Println()

	for y := 0; y < b.H; y++ {
		fmt.Printf("%3d: ", y)
		for x := 0; x < b.W; x++ {
			fmt.Print(b.Get(x, y))
		}
		fmt.Println()
	}
}
