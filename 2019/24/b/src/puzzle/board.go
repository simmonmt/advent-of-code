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
	"github.com/simmonmt/aoc/2019/common/pos"
)

type Board struct {
	level    int
	c        [5][5]bool
	up, down *Board
}

func makeBoard(level int, up, down *Board) *Board {
	return &Board{
		level: level,
		up:    up,
		down:  down,
	}
}

func NewBoard(lines []string) *Board {
	if len(lines[0]) != 5 || len(lines) != 5 {
		panic("bad size")
	}

	b := makeBoard(0, nil, nil)
	for y := range lines {
		for x, r := range lines[y] {
			b.set(pos.P2{x, y}, r == '#')
		}
	}

	return b
}

func (b *Board) set(p pos.P2, val bool) {
	if p.X < 0 || p.Y < 0 || p.X >= 5 || p.Y >= 5 {
		panic("bad pos")
	}

	b.c[p.Y][p.X] = val
}

func (b *Board) Get(p pos.P2) bool {
	if p.X < 0 || p.Y < 0 || p.X >= 5 || p.Y >= 5 {
		return false
	}

	return b.c[p.Y][p.X]
}

func (b *Board) cloneEmptyDown() *Board {
	nb := makeBoard(b.level, nil, nil)
	if b.down != nil {
		nb.down = b.down.cloneEmptyDown()
		nb.down.up = nb
	}
	return nb
}

func (b *Board) cloneEmpty() *Board {
	curRef := b
	for curRef.up != nil {
		curRef = curRef.up
	}

	topNew := curRef.cloneEmptyDown()
	curNew := topNew
	for curRef != b {
		curRef = curRef.down
		curNew = curNew.down
	}

	return curNew
}

type neighborPos struct {
	level int
	posns []pos.P2
}

var (
	specialNeighbors = [5][5]neighborPos{
		// Row 0 A-E
		[5]neighborPos{
			neighborPos{-1, []pos.P2{pos.P2{2, 1}, pos.P2{1, 2}}}, // A: 8,12
			neighborPos{-1, []pos.P2{pos.P2{2, 1}}},               // B: 8
			neighborPos{-1, []pos.P2{pos.P2{2, 1}}},               // C: 8
			neighborPos{-1, []pos.P2{pos.P2{2, 1}}},               // D: 8
			neighborPos{-1, []pos.P2{pos.P2{2, 1}, pos.P2{3, 2}}}, // E: 8,14
		},

		// Row 1 F-J
		[5]neighborPos{
			neighborPos{-1, []pos.P2{pos.P2{1, 2}}}, // F: 12
			neighborPos{},                           // G: none
			neighborPos{1, []pos.P2{ // H/8: A,B,C,D,E
				pos.P2{0, 0}, pos.P2{1, 0}, pos.P2{2, 0}, pos.P2{3, 0},
				pos.P2{4, 0}}},
			neighborPos{},                           // I: none
			neighborPos{-1, []pos.P2{pos.P2{3, 2}}}, // J:14
		},

		// Row 2 K-O
		[5]neighborPos{
			neighborPos{-1, []pos.P2{pos.P2{1, 2}}}, // K: 12
			neighborPos{1, []pos.P2{ // L/12: A,F,K,P,U
				pos.P2{0, 0}, pos.P2{0, 1}, pos.P2{0, 2}, pos.P2{0, 3},
				pos.P2{0, 4}}},
			neighborPos{}, // ?: none
			neighborPos{1, []pos.P2{ // N/14: E,J,O,T,Y
				pos.P2{4, 0}, pos.P2{4, 1}, pos.P2{4, 2}, pos.P2{4, 3},
				pos.P2{4, 4}}},
			neighborPos{-1, []pos.P2{pos.P2{3, 2}}}, // O:14
		},

		// Row 3 P-T
		[5]neighborPos{
			neighborPos{-1, []pos.P2{pos.P2{1, 2}}}, // P: 12
			neighborPos{},                           // Q: none
			neighborPos{1, []pos.P2{ // R/18: U,V,W,X,Y
				pos.P2{0, 4}, pos.P2{1, 4}, pos.P2{2, 4}, pos.P2{3, 4},
				pos.P2{4, 4}}},
			neighborPos{},                           // S: none
			neighborPos{-1, []pos.P2{pos.P2{3, 2}}}, // T:14
		},

		// Row 4 U-Y
		[5]neighborPos{
			neighborPos{-1, []pos.P2{pos.P2{1, 2}, pos.P2{2, 3}}}, // U: 12,18
			neighborPos{-1, []pos.P2{pos.P2{2, 3}}},               // V: 18
			neighborPos{-1, []pos.P2{pos.P2{2, 3}}},               // W: 18
			neighborPos{-1, []pos.P2{pos.P2{2, 3}}},               // X: 18
			neighborPos{-1, []pos.P2{pos.P2{3, 2}, pos.P2{2, 3}}}, // Y: 14,18
		},
	}
)

func (b *Board) getSpecial(levelDelta int, p pos.P2) bool {
	var lb *Board
	if levelDelta < 0 {
		lb = b.up
	} else {
		lb = b.down
	}

	if lb == nil {
		return false
	}

	return lb.Get(p)
}

func (b *Board) Evolve() *Board {
	nb := b.cloneEmpty()

	curNew, curRef := nb, b
	for curRef.up != nil {
		curNew = curNew.up
		curRef = curRef.up
	}

	maybeRef := makeBoard(b.level-1, nil, curRef)
	maybeNew := makeBoard(b.level-1, nil, curNew)
	if numSet := evolveThis(maybeRef, maybeNew); numSet > 0 {
		curNew.up = maybeNew
	}

	for {
		//fmt.Printf("evolving %v %v\n", curRef, curNew)

		evolveThis(curRef, curNew)

		if curNew.down == nil {
			maybeNew = makeBoard(b.level+1, curNew, nil)
			maybeRef = makeBoard(b.level+1, curRef, nil)
			if numSet := evolveThis(maybeRef, maybeNew); numSet > 0 {
				curNew.down = maybeNew
			}
			break
		}

		curNew = curNew.down
		curRef = curRef.down
	}

	return nb
}

func evolveThis(ref, nb *Board) int {
	numSet := 0

	refGet := func(p pos.P2) bool {
		if ref == nil {
			return false
		}
		return ref.Get(p)
	}

	set := func(p pos.P2, v bool) {
		nb.set(p, v)
		if v {
			numSet++
		}
	}

	for y := range nb.c {
		for x := range nb.c[0] {
			if x == 2 && y == 2 {
				continue
			}

			p := pos.P2{x, y}
			neighbors := 0
			for _, dir := range dir.AllDirs {
				if np := dir.From(p); refGet(np) {
					neighbors++
				}
			}
			special := specialNeighbors[p.Y][p.X]
			for _, sp := range special.posns {
				if ref != nil && ref.getSpecial(special.level, sp) {
					neighbors++
				}
			}

			if refGet(p) {
				// existing survives iff it has 1 neighbor
				set(p, neighbors == 1)
			} else {
				// create if empty has 1 or 2 neighbors
				set(p, neighbors == 1 || neighbors == 2)
			}
		}
	}

	return numSet
}

func (b *Board) Strings() []string {
	out := []string{}
	for y := range b.c {
		line := ""
		for x := range b.c[0] {
			p := pos.P2{x, y}
			if p.Equals(pos.P2{2, 2}) {
				line += "?"
			} else if b.Get(p) {
				line += "#"
			} else {
				line += "."
			}
		}
		out = append(out, line)
	}
	return out
}

func dumpOneBoard(b *Board) {
	fmt.Printf("Level %d:\n", b.level)
	for y := range b.c {
		for x := range b.c[0] {
			p := pos.P2{x, y}
			if b.Get(p) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (b *Board) Dump() {
	top := b
	for top.up != nil {
		top = top.up
	}

	for nb := top; nb != nil; nb = nb.down {
		dumpOneBoard(nb)
	}
}

func (b *Board) NumBugs() int {
	curRef := b
	for curRef.up != nil {
		curRef = curRef.up
	}

	num := 0
	for ; curRef != nil; curRef = curRef.down {
		for _, row := range curRef.c {
			for _, v := range row {
				if v {
					num++
				}
			}
		}
	}
	return num
}
