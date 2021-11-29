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

package puzzle

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/dir"
	"github.com/simmonmt/aoc/2019/common/intmath"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type Board struct {
	m      map[pos.P2]rune
	tl, br pos.P2
}

func NewBoard() *Board {
	return &Board{
		m:  map[pos.P2]rune{},
		tl: pos.P2{-1, -1},
		br: pos.P2{-1, -1},
	}
}

func (b *Board) Set(p pos.P2, ch rune) {
	if b.tl.X == -1 {
		b.tl = p
		b.br = p
	} else {
		b.tl.X = intmath.IntMin(b.tl.X, p.X)
		b.tl.Y = intmath.IntMin(b.tl.Y, p.Y)

		b.br.X = intmath.IntMax(b.br.X, p.X)
		b.br.Y = intmath.IntMax(b.br.Y, p.Y)
	}

	b.m[p] = ch
}

func (b *Board) Get(p pos.P2) rune {
	if ch, found := b.m[p]; found {
		return ch
	}
	return '.'
}

func DumpBoard(b *Board, vac *Vac, intersections map[pos.P2]bool) {
	for y := 0; y <= b.br.Y; y++ {
		for x := 0; x <= b.br.X; x++ {
			p := pos.P2{x, y}
			if p.Equals(vac.Pos) {
				fmt.Print(vac.Dir)
			} else if intersections != nil && intersections[p] {
				fmt.Print("O")
			} else {
				fmt.Printf("%c", b.Get(p))
			}
		}
		fmt.Println()
	}
}

func FindIntersections(b *Board) map[pos.P2]bool {
	intersections := map[pos.P2]bool{}

	for y := 0; y <= b.br.Y; y++ {
		for x := 0; x <= b.br.X; x++ {
			p := pos.P2{x, y}
			if b.Get(p) != '#' {
				continue
			}

			numNeighborScaffolds := 0
			for _, d := range dir.AllDirs {
				if b.Get(d.From(p)) == '#' {
					numNeighborScaffolds++
				}
			}

			if numNeighborScaffolds > 2 {
				intersections[p] = true
			}
		}
	}

	return intersections
}

func SumAlignmentParams(ps map[pos.P2]bool) int {
	align := 0
	for p := range ps {
		align += p.X * p.Y
	}
	return align
}
