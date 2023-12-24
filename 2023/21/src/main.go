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

package main

import (
	"flag"
	"fmt"

	"github.com/simmonmt/aoc/2023/common/dir"
	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/grid"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) (*grid.Grid[rune], pos.P2, error) {
	var start pos.P2
	g, err := grid.NewFromLines(lines, func(p pos.P2, r rune) (rune, error) {
		if r == 'S' {
			start = p
			r = '.'
		}
		if r != '.' && r != '#' {
			return '?', fmt.Errorf("bad rune")
		}
		return r, nil
	})

	return g, start, err
}

type GridLike interface {
	AllNeighbors(p pos.P2, includeDiag bool) []pos.P2
	Get(p pos.P2) (rune, bool)
}

func iterate(g GridLike, todo []pos.P2) []pos.P2 {
	next := map[pos.P2]bool{}
	for _, p := range todo {
		for _, n := range g.AllNeighbors(p, false) {
			v, _ := g.Get(n)
			if v == '.' {
				next[n] = true
			}
		}
	}

	out := []pos.P2{}
	for p := range next {
		out = append(out, p)
	}
	return out
}

func solveA(g *grid.Grid[rune], start pos.P2, maxSteps int) int {
	todo := []pos.P2{start}
	var num int
	for steps := 0; len(todo) > 0 && steps < maxSteps; steps++ {
		logger.Infof("steps %d num %v, todo %v", steps, num, todo)
		todo = iterate(g, todo)
		num = len(todo)
	}

	return num
}

func dumpBoard(g GridLike, xMin, xMax, yMin, yMax int, filled []pos.P2) {
	m := map[pos.P2]bool{}
	for _, p := range filled {
		m[p] = true
	}

	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			p := pos.P2{X: x, Y: y}
			if _, found := m[p]; found {
				fmt.Print("O")
			} else {
				r, _ := g.Get(p)
				fmt.Print(string(r))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

type TiledGrid struct {
	g    *grid.Grid[rune]
	w, h int
}

func NewTiledGrid(g *grid.Grid[rune]) *TiledGrid {
	return &TiledGrid{g, g.Width(), g.Height()}
}

func scaledPos(p pos.P2, w, h int) pos.P2 {
	sp := pos.P2{X: p.X % w, Y: p.Y % h}
	if sp.X < 0 {
		sp.X += w
	}
	if sp.Y < 0 {
		sp.Y += h
	}

	return sp
}

func (tg *TiledGrid) Get(p pos.P2) (rune, bool) {
	return tg.g.Get(scaledPos(p, tg.w, tg.h))
}

func (tg *TiledGrid) AllNeighbors(p pos.P2, includeDiag bool) []pos.P2 {
	out := make([]pos.P2, 4)
	for i, d := range dir.AllDirs {
		out[i] = d.From(p)
	}
	return out
}

// 11967229948 too low
// 605247179063458 too high
// 605247138198755
func solveB(g *grid.Grid[rune], start pos.P2, wantSteps int64) int64 {
	repeatInterval := int64(g.Width())
	startingOffset := wantSteps % repeatInterval

	tg := NewTiledGrid(g)

	rows := [4]int64{-1, -1, -1, -1}

	todo := []pos.P2{start}
	var steps int
	for steps = 1; ; steps++ {
		todo = iterate(tg, todo)

		if int64(steps)%repeatInterval == startingOffset {
			nr := [4]int64{-1, -1, -1, -1}
			nr[0] = int64(steps)
			nr[1] = int64(len(todo))
			if rows[1] > 0 {
				nr[2] = nr[1] - rows[1]
			}
			if rows[2] > 0 {
				nr[3] = nr[2] - rows[2]
			}
			if rows[3] > 0 {
				if nr[3] == rows[3] {
					break
				}
			}
			rows = nr
		}
	}

	for rows[0] != wantSteps {
		nr := [4]int64{}
		nr[0] = rows[0] + repeatInterval

		nr[3] = rows[3]
		nr[2] = rows[2] + nr[3]
		nr[1] = rows[1] + nr[2]
		rows = nr
	}

	return rows[1]
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	g, start, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(g, start, 64))
	fmt.Println("B", solveB(g, start, 26501365))
}
