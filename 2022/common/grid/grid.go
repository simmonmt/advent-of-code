// Copyright 2022 Google LLC
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

package grid

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/simmonmt/aoc/2022/common/mtsmath"
	"github.com/simmonmt/aoc/2022/common/pos"
)

type Grid[T any] struct {
	w, h int
	a    []T
}

func New[T any](w, h int) *Grid[T] {
	return &Grid[T]{
		w: w,
		h: h,
		a: make([]T, w*h),
	}
}

func NewFromLines[T any](lines []string, cellMapper func(p pos.P2, r rune) (T, error)) (*Grid[T], error) {
	g := New[T](len(lines[0]), len(lines))
	for y, line := range lines {
		if len(line) != g.Width() {
			return nil, fmt.Errorf("uneven line")
		}

		for x, r := range line {
			p := pos.P2{X: x, Y: y}
			v, err := cellMapper(p, r)
			if err != nil {
				return nil, fmt.Errorf("%v: bad parse %v",
					p, err)
			}

			g.Set(p, v)
		}
	}
	return g, nil
}

func (g *Grid[T]) IsValid(p pos.P2) bool {
	return p.X >= 0 && p.X < g.w && p.Y >= 0 && p.Y < g.h
}

func (g *Grid[T]) Width() int {
	return g.w
}

func (g *Grid[T]) Height() int {
	return g.h
}

func (g *Grid[T]) Set(p pos.P2, v T) {
	off := p.Y*g.w + p.X
	g.a[off] = v
}

func (g *Grid[T]) Get(p pos.P2) T {
	off := p.Y*g.w + p.X
	return g.a[off]
}

func (g *Grid[T]) Walk(walker func(p pos.P2, v T)) {
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			p := pos.P2{X: x, Y: y}
			walker(p, g.Get(p))
		}
	}
}

func (g *Grid[T]) AllNeighbors(p pos.P2, includeDiag bool) []pos.P2 {
	out := []pos.P2{}
	for _, n := range p.AllNeighbors(includeDiag) {
		if n.X < 0 || n.Y < 0 {
			continue
		}
		if n.X >= g.Width() || n.Y >= g.Height() {
			continue
		}
		out = append(out, n)
	}
	return out
}

type SparseGrid struct {
	start, end pos.P2
	a          map[pos.P2]any
}

func NewSparseGrid() *SparseGrid {
	return &SparseGrid{
		a: map[pos.P2]any{},
	}
}

func (g *SparseGrid) Clone() *SparseGrid {
	o := &SparseGrid{
		start: g.start,
		end:   g.end,
		a:     map[pos.P2]any{},
	}

	for k, v := range g.a {
		o.a[k] = v
	}

	return o
}

func (g *SparseGrid) Start() pos.P2 {
	return g.start
}

func (g *SparseGrid) End() pos.P2 {
	return g.end
}

func (g *SparseGrid) Set(p pos.P2, v any) {
	if len(g.a) == 0 {
		g.start, g.end = p, p
	} else {
		if p.X < g.start.X {
			g.start.X = p.X
		}
		if p.Y < g.start.Y {
			g.start.Y = p.Y
		}
		if p.X > g.end.X {
			g.end.X = p.X
		}
		if p.Y > g.end.Y {
			g.end.Y = p.Y
		}
	}

	g.a[p] = v
}

func (g *SparseGrid) Get(p pos.P2) (v any, found bool) {
	v, found = g.a[p]
	return
}

func numDigits(num int) int {
	if num == 0 {
		return 1
	}

	n := 0
	for num > 0 {
		n++
		num /= 10
	}
	return n
}

func (g *SparseGrid) DumpTo(withCoords bool, mapper func(p pos.P2, v any, found bool) string, w io.Writer) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	yDigits := numDigits(g.end.Y)
	xDigits := numDigits(g.end.X)

	maxCellWidth := 1
	height := g.end.Y - g.start.Y + 1
	rows := make([][]string, height)
	for i := 0; i < height; i++ {
		y := g.start.Y + i
		width := g.end.X - g.start.X + 1
		row := make([]string, width)
		for j := 0; j < width; j++ {
			x := g.start.X + j
			p := pos.P2{x, y}
			v, found := g.Get(p)
			s := mapper(p, v, found)
			row[j] = s
			maxCellWidth = mtsmath.Max(maxCellWidth, len(s))
		}
		rows[i] = row
	}

	div := 1
	for i := 1; i < xDigits; i++ {
		div *= 10
	}

	for i := xDigits; i > 0; i-- {
		fmt.Fprintf(bw, "%*s ", yDigits, "")
		for x := g.start.X; x <= g.end.X; x++ {
			fmt.Fprintf(bw, "%*d", maxCellWidth, (x/div)%10)
		}
		fmt.Fprintln(bw)
		div /= 10
	}

	for i := 0; i < len(rows); i++ {
		y := g.start.Y + i
		if withCoords {
			fmt.Fprintf(bw, "%*d ", yDigits, y)
		}

		for j := 0; j < len(rows[i]); j++ {
			fmt.Fprintf(bw, "%*s", maxCellWidth, rows[i][j])
		}
		fmt.Fprintln(bw)
	}
}

func (g *SparseGrid) Dump(withCoords bool, mapper func(p pos.P2, v any, found bool) string) {
	g.DumpTo(withCoords, mapper, os.Stdout)
}
