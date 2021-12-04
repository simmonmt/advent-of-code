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
	"fmt"
	"io"
	"os"

	"github.com/simmonmt/aoc/2021/common/pos"
)

type cell struct {
	val    int
	marked bool
}

type Board struct {
	valToPos map[int]pos.P2
	cells    map[pos.P2]*cell
}

func New(start [5][5]int) *Board {
	b := &Board{
		valToPos: map[int]pos.P2{},
		cells:    map[pos.P2]*cell{},
	}

	pos.WalkP2(5, 5, func(p pos.P2) {
		v := start[p.Y][p.X]
		if _, found := b.valToPos[v]; found {
			panic("dup value")
		}

		b.valToPos[v] = p
		b.cells[p] = &cell{val: v, marked: false}
	})

	return b
}

func (b *Board) Reset() {
	for _, cell := range b.cells {
		cell.marked = false
	}
}

func (b *Board) Mark(v int) bool {
	p, found := b.valToPos[v]
	if !found {
		return false
	}

	c := b.cells[p]
	if c.marked {
		panic("remark")
	}

	c.marked = true

	return b.checkWonX(p.X) || b.checkWonY(p.Y)
}

func (b *Board) checkWonX(x int) bool {
	for y := 0; y < 5; y++ {
		if !b.cells[pos.P2{Y: y, X: x}].marked {
			return false
		}
	}
	return true
}

func (b *Board) checkWonY(y int) bool {
	for x := 0; x < 5; x++ {
		if !b.cells[pos.P2{Y: y, X: x}].marked {
			return false
		}
	}
	return true
}

func (b *Board) Score(lastCalled int) int {
	score := 0
	for _, c := range b.cells {
		if c.marked {
			continue
		}

		score += c.val
	}
	return score * lastCalled
}

func (b *Board) DumpTo(o io.Writer) {
	pos.WalkP2(5, 5, func(p pos.P2) {
		c := b.cells[p]
		if c.marked {
			fmt.Fprintf(o, "*%2d*", c.val)
		} else {
			fmt.Fprintf(o, " %2d ", c.val)
		}

		if p.X == 4 {
			fmt.Fprintln(o)
		} else {
			fmt.Fprint(o, " ")
		}
	})
}

func (b *Board) Dump() {
	b.DumpTo(os.Stdout)
}
