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

func parseInput(lines []string) (*grid.Grid[rune], error) {
	return grid.NewFromLines[rune](lines, func(p pos.P2, r rune) (rune, error) {
		if (r != 'O') && (r != '.') && (r != '#') {
			return '?', fmt.Errorf("bad rune %c at %v", r, p)
		}
		return r, nil
	})
}

func dumpGrid(g *grid.Grid[rune]) {
	g.Dump(true, func(p pos.P2, v rune, _ bool) string {
		return string(v)
	})
}

func moveCell(g *grid.Grid[rune], src pos.P2, d dir.Dir) {
	r, _ := g.Get(src)

	if r != 'O' {
		return
	}

	dest := src
	for {
		next := d.From(dest)
		if !g.IsValid(next) {
			break
		}
		if v, _ := g.Get(next); v != '.' {
			break
		}
		dest = next
	}

	if !src.Equals(dest) {
		//logger.Infof("moving %v to %v", src, dest)
		g.Set(src, '.')
		g.Set(dest, 'O')
	}
}

func pushNS(g *grid.Grid[rune], d dir.Dir, topFirst bool) {
	rows := []int{}
	if topFirst {
		for i := 0; i < g.Height(); i++ {
			rows = append(rows, i)
		}
	} else {
		for i := g.Height() - 1; i >= 0; i-- {
			rows = append(rows, i)
		}
	}

	for _, y := range rows {
		for x := 0; x < g.Width(); x++ {
			moveCell(g, pos.P2{X: x, Y: y}, d)
		}
	}
}

func pushEW(g *grid.Grid[rune], d dir.Dir, leftFirst bool) {
	cols := []int{}
	if leftFirst {
		for i := 0; i < g.Width(); i++ {
			cols = append(cols, i)
		}
	} else {
		for i := g.Width() - 1; i >= 0; i-- {
			cols = append(cols, i)
		}
	}

	for _, x := range cols {
		for y := 0; y < g.Height(); y++ {
			moveCell(g, pos.P2{X: x, Y: y}, d)
		}
	}
}

func countWeight(g *grid.Grid[rune]) int {
	sum := 0
	g.Walk(func(src pos.P2, r rune) {
		if r != 'O' {
			return
		}

		height := g.Height() - src.Y
		sum += height
	})
	return sum
}

func solveA(g *grid.Grid[rune]) int {
	pushNS(g, dir.DIR_NORTH, true)
	return countWeight(g)
}

func spin(g *grid.Grid[rune]) {
	pushNS(g, dir.DIR_NORTH, true)
	pushEW(g, dir.DIR_WEST, true)
	pushNS(g, dir.DIR_SOUTH, false)
	pushEW(g, dir.DIR_EAST, false)
}

func serialize(g *grid.Grid[rune]) string {
	out := make([]byte, g.Height()*g.Width())
	g.Walk(func(p pos.P2, r rune) {
		off := p.Y*g.Width() + p.X
		out[off] = byte(r)
	})
	return string(out)
}

func solveB(g *grid.Grid[rune]) int {
	saved := g.Clone()

	m := map[string]int{}
	var loopStart, loopCycle int
	for i := 0; i < 100000; i++ {
		s := serialize(g)
		if num, found := m[s]; found {
			logger.Infof("i=%d also i=%d off=%d\n", i, num, i-num)
			loopStart = num
			loopCycle = i - num
			break
		}
		m[s] = i

		spin(g)
	}

	need := 1000000000
	need -= loopStart
	need = need % loopCycle

	for i := 0; i < loopStart+need; i++ {
		spin(saved)
	}
	return countWeight(saved)
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

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))

	input, _ = parseInput(lines)
	fmt.Println("B", solveB(input))
}
