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

package grid

import "github.com/simmonmt/aoc/2021/common/pos"

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
