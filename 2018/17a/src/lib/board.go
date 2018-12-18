package lib

import (
	"fmt"
	"math"

	"intmath"
)

type CellType int

const (
	TYPE_OPEN CellType = iota
	TYPE_WALL
	TYPE_FLOW
	TYPE_FILLED
)

func (c CellType) Short() string {
	switch c {
	case TYPE_OPEN:
		return "."
	case TYPE_WALL:
		return "#"
	case TYPE_FLOW:
		return "|"
	case TYPE_FILLED:
		return "~"
	default:
		panic("unknown")
	}
}

func (c CellType) String() string {
	switch c {
	case TYPE_OPEN:
		return "open"
	case TYPE_WALL:
		return "wall"
	case TYPE_FLOW:
		return "flow"
	case TYPE_FILLED:
		return "filled"
	default:
		panic("unknown")
	}
}

type Board struct {
	cells                  [][]CellType
	xmin, xmax, ymin, ymax int
}

func findBounds(lines []InputLine) (xmin, xmax, ymin, ymax int) {
	xmin, xmax = math.MaxInt32, 0
	ymin, ymax = math.MaxInt32, 0

	for _, line := range lines {
		xmin = intmath.IntMin(xmin, line.Xmin)
		xmax = intmath.IntMax(xmax, line.Xmax)
		ymin = intmath.IntMin(ymin, line.Ymin)
		ymax = intmath.IntMax(ymax, line.Ymax)
	}

	return
}

func NewBoard(lines []InputLine) *Board {
	xmin, xmax, ymin, ymax := findBounds(lines)
	return &Board{
		cells: nil,
		xmin:  xmin,
		xmax:  xmax,
		ymin:  ymin,
		ymax:  ymax,
	}
}

func (b *Board) Dump() {
	b.DumpBox(b.xmin, b.xmax, b.ymin, b.ymax)
}

func (b *Board) DumpBox(xmin, xmax, ymin, ymax int) {
	xmin = intmath.IntMax(xmin, b.xmin)
	xmax = intmath.IntMin(xmax, b.xmax)
	ymin = intmath.IntMax(ymin, b.ymin)
	ymax = intmath.IntMin(ymax, b.ymax)

	xlines := int(math.Log10(float64(xmax)) + 1)
	ylines := int(math.Log10(float64(ymax)) + 1)

	for i := 0; i < xlines; i++ {
		fmt.Printf("%*s ", ylines, " ")
		div := int(math.Pow10(xlines - 1 - i))
		for j := xmin; j <= xmax; j++ {
			fmt.Print(j / div % 10)
		}
		fmt.Println()
	}

	for y := ymin; y <= ymax; y++ {
		fmt.Printf("%*d ", ylines, y)

		for x := xmin; x <= xmax; x++ {
			fmt.Print(b.cells[y][x].Short())
		}
		fmt.Println()
	}
}
