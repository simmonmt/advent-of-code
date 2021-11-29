// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lib

import (
	"fmt"
	"image"
	"image/color"
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
		xmin = intmath.IntMin(xmin, line.Xmin) - 1
		xmax = intmath.IntMax(xmax, line.Xmax) + 1
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

func (b *Board) Bounds() (xmin, xmax, ymin, ymax int) {
	return b.xmin, b.xmax, b.ymin, b.ymax
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

func (b *Board) get(pos Pos) CellType {
	return b.cells[pos.Y-b.ymin][pos.X-b.xmin]
}

func (b *Board) Get(pos Pos) CellType {
	b.checkBounds(pos)
	return b.get(pos)
}

func (b *Board) GetWithDefault(pos Pos, def CellType) CellType {
	if !b.InBounds(pos) {
		return def
	}
	return b.get(pos)
}

func (b *Board) Set(pos Pos, cell CellType) {
	b.checkBounds(pos)

	// seemsBad := func(ex, cell CellType) bool {
	// 	if ex == TYPE_OPEN {
	// 		return false
	// 	}
	// 	if ex == TYPE_WALL && cell == TYPE_WALL {
	// 		return false
	// 	}
	// 	if ex == TYPE_FLOW && cell == TYPE_FILLED {
	// 		return false
	// 	}

	// 	return true
	// }

	// if ex := b.get(pos); seemsBad(ex, cell) {
	// 	fmt.Printf("overwriting %v %s with %s\n", pos, b.get(pos), cell)
	// }

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

func (b *Board) GetACursor() (Pos, bool) {
	for c := range b.cursors {
		return c, true
	}
	return Pos{}, false
}

func (b *Board) AddCursor(pos Pos) {
	// if _, found := b.cursors[pos]; found {
	// 	panic("already added")
	// }

	b.cursors[pos] = true
}

func (b *Board) DeleteCursor(pos Pos) {
	if _, found := b.cursors[pos]; !found {
		panic(fmt.Sprintf("no cursor at %v", pos))
	}
	delete(b.cursors, pos)
}

func (b *Board) Dump() {
	b.DumpBox(b.xmin, b.xmax, b.ymin, b.ymax, Pos{-1, -1})
}

func (b *Board) DumpWithFocus(focus Pos) {
	b.DumpBox(b.xmin, b.xmax, b.ymin, b.ymax, focus)
}

func (b *Board) DumpBox(xmin, xmax, ymin, ymax int, focus Pos) {
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
			if focus.X >= 0 && focus.Y >= 0 && focus.Eq(pos) {
				fmt.Printf("\033[1m%s\033[0m", short)
			} else {
				fmt.Print(short)
			}
		}
		fmt.Println()
	}
}

func (b *Board) Score() (numFlow, numFilled int) {
	minYWall := math.MaxInt32
	numFlow, numFilled = 0, 0
	for y := range b.cells {
		for _, c := range b.cells[y] {
			if c == TYPE_WALL {
				minYWall = intmath.IntMin(minYWall, y)
			} else if c == TYPE_FLOW {
				numFlow++
			} else if c == TYPE_FILLED {
				numFilled++
			}
		}
	}

	// Remove the flow values between the spigot and the top of
	// the highest wall
	if minYWall > 1 {
		numFlow -= (minYWall - 1)
	}

	return
}

func (b *Board) Visit(want CellType, visitor func(pos Pos)) {
	for y := range b.cells {
		for x, c := range b.cells[y] {
			if c == want {
				visitor(Pos{x + b.xmin, y + b.ymin})
			}
		}
	}
}

func (b *Board) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(b.xmin, b.ymin, b.xmax+1, b.ymax+1))

	for y := range b.cells {
		for x, c := range b.cells[y] {
			var col color.RGBA
			switch c {
			case TYPE_OPEN:
				col = color.RGBA{255, 255, 255, 255}
			case TYPE_WALL:
				col = color.RGBA{0, 0, 0, 255}
			case TYPE_FLOW:
				col = color.RGBA{127, 127, 127, 255}
			case TYPE_FILLED:
				col = color.RGBA{0, 0, 255, 255}
			case TYPE_SPRING:
				col = color.RGBA{255, 0, 0, 255}
			default:
				panic("unknown")
			}

			img.Set(x+b.xmin, y+b.ymin, col)
		}
	}

	return img
}
