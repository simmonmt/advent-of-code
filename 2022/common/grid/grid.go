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
	"fmt"
	"io"
	"os"

	"github.com/simmonmt/aoc/2022/common/pos"
)

type Grid struct {
	w, h int
	a    []interface{}
}

func New(w, h int) *Grid {
	return &Grid{
		w: w,
		h: h,
		a: make([]interface{}, w*h),
	}
}

func NewFromLines(lines []string, cellMapper func(r rune) (interface{}, error)) (*Grid, error) {
	g := New(len(lines[0]), len(lines))
	for y, line := range lines {
		if len(line) != g.Width() {
			return nil, fmt.Errorf("uneven line")
		}

		for x, r := range line {
			p := pos.P2{X: x, Y: y}
			v, err := cellMapper(r)
			if err != nil {
				return nil, fmt.Errorf("%v: bad parse %v",
					p, err)
			}

			g.Set(p, v)
		}
	}
	return g, nil
}

func (g *Grid) Width() int {
	return g.w
}

func (g *Grid) Height() int {
	return g.h
}

func (g *Grid) Set(p pos.P2, v interface{}) {
	off := p.Y*g.w + p.X
	g.a[off] = v
}

func (g *Grid) Get(p pos.P2) interface{} {
	off := p.Y*g.w + p.X
	return g.a[off]
}

func (g *Grid) Walk(walker func(p pos.P2, v interface{})) {
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			p := pos.P2{X: x, Y: y}
			walker(p, g.Get(p))
		}
	}
}

func (g *Grid) AllNeighbors(p pos.P2, includeDiag bool) []pos.P2 {
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

type IntGrid struct {
	Grid
}

func NewInt(w, h int) *IntGrid {
	return &IntGrid{
		Grid: *New(w, h),
	}
}

func (g *IntGrid) SetInt(p pos.P2, v int) {
	g.Set(p, v)
}

func (g *IntGrid) GetInt(p pos.P2) int {
	v, _ := g.Get(p).(int)
	return v
}

func (g *IntGrid) WalkInt(walker func(p pos.P2, v int)) {
	g.Walk(func(p pos.P2, v interface{}) {
		n, _ := v.(int)
		walker(p, n)
	})
}

func (g *IntGrid) DumpTo(w io.Writer) {
	g.WalkInt(func(p pos.P2, v int) {
		fmt.Fprintf(w, "%d", v)
		if p.X == g.w-1 {
			fmt.Fprintln(w)
		}
	})
}

func (g *IntGrid) Dump() {
	g.DumpTo(os.Stdout)
}
