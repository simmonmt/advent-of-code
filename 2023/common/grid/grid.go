// Copyright 2023 Google LLC
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

	"github.com/simmonmt/aoc/2023/common/mtsmath"
	"github.com/simmonmt/aoc/2023/common/pos"
)

type dumpableGrid[T any] interface {
	Start() pos.P2
	End() pos.P2
	Get(pos.P2) (T, bool)
}

func dumpTo[V any, T dumpableGrid[V]](g T, withCoords bool, mapper func(p pos.P2, v V, found bool) string, w io.Writer) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	start, end := g.Start(), g.End()
	yDigits := numDigits(end.Y)
	xDigits := numDigits(end.X)

	maxCellWidth := 1
	height := end.Y - start.Y + 1
	rows := make([][]string, height)
	for i := 0; i < height; i++ {
		y := start.Y + i
		width := end.X - start.X + 1
		row := make([]string, width)
		for j := 0; j < width; j++ {
			x := start.X + j
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

	if withCoords {
		for i := xDigits; i > 0; i-- {
			fmt.Fprintf(bw, "%*s ", yDigits, "")
			for x := start.X; x <= end.X; x++ {
				fmt.Fprintf(bw, "%*d", maxCellWidth, (x/div)%10)
			}
			fmt.Fprintln(bw)
			div /= 10
		}
	}

	for i := 0; i < len(rows); i++ {
		y := start.Y + i
		if withCoords {
			fmt.Fprintf(bw, "%*d ", yDigits, y)
		}

		for j := 0; j < len(rows[i]); j++ {
			fmt.Fprintf(bw, "%*s", maxCellWidth, rows[i][j])
		}
		fmt.Fprintln(bw)
	}
}

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

func (g *Grid[T]) Start() pos.P2 {
	return pos.P2{0, 0}
}

func (g *Grid[T]) End() pos.P2 {
	return pos.P2{g.w - 1, g.h - 1}
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

func (g *Grid[T]) Get(p pos.P2) (T, bool) {
	off := p.Y*g.w + p.X
	if off < 0 || off >= len(g.a) {
		return g.a[0], false
	}
	return g.a[off], true
}

func (g *Grid[T]) Walk(walker func(p pos.P2, v T)) {
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			p := pos.P2{X: x, Y: y}
			v, _ := g.Get(p)
			walker(p, v)
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

func (g *Grid[T]) DumpTo(withCoords bool, mapper func(p pos.P2, v T, found bool) string, w io.Writer) {
	dumpTo[T](g, withCoords, mapper, w)
}

func (g *Grid[T]) Dump(withCoords bool, mapper func(p pos.P2, v T, found bool) string) {
	g.DumpTo(withCoords, mapper, os.Stdout)
}

type SparseGrid[T any] struct {
	start, end pos.P2
	a          map[pos.P2]T
}

func NewSparseGrid[T any]() *SparseGrid[T] {
	return &SparseGrid[T]{
		a: map[pos.P2]T{},
	}
}

func (g *SparseGrid[T]) Clone() *SparseGrid[T] {
	o := &SparseGrid[T]{
		start: g.start,
		end:   g.end,
		a:     map[pos.P2]T{},
	}

	for k, v := range g.a {
		o.a[k] = v
	}

	return o
}

func (g *SparseGrid[T]) Start() pos.P2 {
	return g.start
}

func (g *SparseGrid[T]) End() pos.P2 {
	return g.end
}

func (g *SparseGrid[T]) Empty() bool {
	return g.start.Equals(g.end)
}

func (g *SparseGrid[T]) Set(p pos.P2, v T) {
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

func (g *SparseGrid[T]) Get(p pos.P2) (v T, found bool) {
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

func (g *SparseGrid[T]) DumpTo(withCoords bool, mapper func(p pos.P2, v T, found bool) string, w io.Writer) {
	dumpTo[T](g, withCoords, mapper, w)
}

func (g *SparseGrid[T]) Dump(withCoords bool, mapper func(p pos.P2, v T, found bool) string) {
	g.DumpTo(withCoords, mapper, os.Stdout)
}
