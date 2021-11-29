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

package board

import (
	"io"

	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/pos"
)

type Board struct {
	m     map[pos.P4]bool
	fourD bool
	min   pos.P4
	max   pos.P4
}

func new(fourD bool) *Board {
	return &Board{
		m:     map[pos.P4]bool{},
		fourD: fourD,
	}
}

func New(lines []string, fourD bool) *Board {
	b := new(fourD)

	for y := range lines {
		for x, r := range lines[y] {
			if r == '#' {
				b.Set(pos.P4{X: x, Y: y, Z: 0, W: 0}, true)
			}
		}
	}

	return b
}

func (b *Board) Set(p pos.P4, val bool) {
	if len(b.m) == 0 {
		b.min, b.max = p, p
	} else {
		b.min.X = intmath.IntMin(b.min.X, p.X)
		b.min.Y = intmath.IntMin(b.min.Y, p.Y)
		b.min.Z = intmath.IntMin(b.min.Z, p.Z)
		b.min.W = intmath.IntMin(b.min.W, p.W)
		b.max.X = intmath.IntMax(b.max.X, p.X)
		b.max.Y = intmath.IntMax(b.max.Y, p.Y)
		b.max.Z = intmath.IntMax(b.max.Z, p.Z)
		b.max.W = intmath.IntMax(b.max.W, p.W)
	}

	b.m[p] = val
}

func (b *Board) Get(p pos.P4) bool {
	return b.m[p]
}

func (b *Board) ZBounds() (min, max int) {
	return b.min.Z, b.max.Z
}

func (b *Board) WBounds() (min, max int) {
	return b.min.W, b.max.W
}

func (b *Board) DumpZW(z, w int, wr io.Writer) {
	for y := b.min.Y; y <= b.max.Y; y++ {
		for x := b.min.X; x <= b.max.X; x++ {
			if b.m[pos.P4{X: x, Y: y, Z: z, W: w}] {
				wr.Write([]byte{'#'})
			} else {
				wr.Write([]byte{'.'})
			}
		}
		wr.Write([]byte{'\n'})
	}
}

func (b *Board) Visit(visitor func(pos.P4, bool)) {
	for p, v := range b.m {
		visitor(p, v)
	}
}

func calcNext(b *Board, p pos.P4) bool {
	cur := b.Get(p)

	numActiveNeighbors := 0
	for _, n := range p.AllNeighbors() {
		if !b.fourD && n.W != 0 {
			continue
		}

		if b.Get(n) {
			numActiveNeighbors++
			if numActiveNeighbors > 3 {
				break
			}
		}
	}

	if cur {
		return numActiveNeighbors == 2 || numActiveNeighbors == 3
	} else {
		return numActiveNeighbors == 3
	}
}

func (b *Board) Evolve() *Board {
	nb := new(b.fourD)

	b.Visit(func(p pos.P4, v bool) {
		cands := p.AllNeighbors()
		cands = append(cands, p)

		for _, cand := range cands {
			if !b.fourD && cand.W != 0 {
				continue
			}

			next := calcNext(b, cand)
			if next {
				nb.Set(cand, next)
			}
		}
	})

	return nb
}
