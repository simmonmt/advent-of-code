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

package main

import (
	"fmt"
)

type Fill int

const (
	FILL_ROCKY Fill = iota
	FILL_WET
	FILL_NARROW
)

type Board struct {
	cells [][]Fill
}

func NewBoard(w, h int) *Board {
	cells := make([][]Fill, h)
	for y := range cells {
		cells[y] = make([]Fill, w)
	}

	return &Board{
		cells: cells,
	}
}

func (b *Board) Get(p Pos) Fill {
	return b.cells[p.Y][p.X]
}

var (
	bfsDirs = []Pos{Pos{-1, 0}, Pos{1, 0}, Pos{0, -1}, Pos{0, 1}}
)

func (b *Board) BFS(start Pos, visitor func(p Pos)) {
	w := len(b.cells[0])
	h := len(b.cells)

	visited := map[Pos]bool{}
	todo := []Pos{start}

	for len(todo) > 0 {
		cur := todo[0]
		todo = todo[1:]
		visitor(cur)
		visited[cur] = true

		for _, dir := range bfsDirs {
			cand := Pos{cur.X + dir.X, cur.Y + dir.Y}

			if cand.X < 0 || cand.X >= w {
				continue
			} else if cand.Y < 0 || cand.Y >= h {
				continue
			} else if _, found := visited[cand]; !found {
				todo = append(todo, cand)
			}
		}
	}
}

func (b *Board) Dump(start, target Pos) {
	for y := range b.cells {
		for x := range b.cells[y] {
			p := Pos{x, y}

			var char string
			switch b.cells[y][x] {
			case p.Eq(start):
				char = "M"
			case FILL_ROCKY:
				char = "."
			case FILL_WET:
				char = "="
			case FILL_NARROW:
				char = "|"
			default:
				panic("unknown")
			}

			fmt.Print(char)
		}
		fmt.Println()
	}
}
