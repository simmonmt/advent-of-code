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
	TYPE_SPRING
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
	case TYPE_SPRING:
		return "+"
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
	case TYPE_SPRING:
		return "spring"
	default:
		panic("unknown")
	}
}

type Pos struct {
	X, Y int
}

func (p Pos) Eq(s Pos) bool {
	return p.X == s.X && p.Y == s.Y
}

type Board struct {
	cells                  [][]CellType
	cur                    Pos
	xmin, xmax, ymin, ymax int
	cursors                map[Pos]bool
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

func NewBoard(spring Pos, lines []InputLine) *Board {
	xmin, xmax, ymin, ymax := findBounds(lines)

	ymin = intmath.IntMin(ymin, spring.Y)
	xmax = intmath.IntMax(xmax, spring.X)
	xmin = intmath.IntMin(xmin, spring.X)

	cells := make([][]CellType, ymax-ymin+1)
	for y := ymin; y <= ymax; y++ {
		cells[y-ymin] = make([]CellType, xmax-xmin+1)
	}

	b := &Board{
		cells:   cells,
		xmin:    xmin,
		xmax:    xmax,
		ymin:    ymin,
		ymax:    ymax,
		cursors: map[Pos]bool{},
	}

	b.Set(spring, TYPE_SPRING)

	for _, line := range lines {
		for y := line.Ymin; y <= line.Ymax; y++ {
			for x := line.Xmin; x <= line.Xmax; x++ {
				b.Set(Pos{x, y}, TYPE_WALL)
			}
		}
	}

	return b
}

func (b *Board) Dump() {
	b.DumpBox(b.xmin, b.xmax, b.ymin, b.ymax)
}

func (b *Board) InBounds(pos Pos) bool {
	return pos.X >= b.xmin && pos.X <= b.xmax && pos.Y >= b.ymin && pos.Y <= b.ymax
}

func (b *Board) checkBounds(pos Pos) {
	if !b.InBounds(pos) {
		panic(fmt.Sprintf("%v out of bounds %v %v", pos, Pos{b.xmin, b.ymin},
			Pos{b.xmax, b.ymax}))
	}
}

func (b *Board) Get(pos Pos) CellType {
	b.checkBounds(pos)
	return b.cells[pos.Y-b.ymin][pos.X-b.xmin]
}

func (b *Board) Set(pos Pos, cell CellType) {
	b.checkBounds(pos)
	b.cells[pos.Y-b.ymin][pos.X-b.xmin] = cell
}

func (b *Board) Cursors() []Pos {
	cursors := make([]Pos, len(b.cursors))
	i := 0
	for c := range b.cursors {
		cursors[i] = c
		i++
	}
	return cursors
}

func (b *Board) AddCursor(pos Pos) {
	b.cursors[pos] = true
}

func (b *Board) MoveCursor(old, new Pos) {
	if _, found := b.cursors[old]; found {
		panic(fmt.Sprintf("no cursor at %v", old))
	}
	delete(b.cursors, old)
	b.cursors[new] = true
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
			pos := Pos{x, y}
			short := b.Get(pos).Short()
			if _, found := b.cursors[pos]; found {
				fmt.Printf("\033[1m%s\033[0m", short)
			} else {
				fmt.Print(short)
			}
		}
		fmt.Println()
	}
}
