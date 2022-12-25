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

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/grid"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/mtsmath"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Contents int

const (
	START Contents = 0
	WALL  Contents = 1
	SAND  Contents = 2
)

func dumpMap(p pos.P2, c Contents, found bool) string {
	if !found {
		return "."
	}

	switch c {
	case SAND:
		return "o"
	case START:
		return "+"
	case WALL:
		return "#"
	default:
		return "?"
	}
}

func drawLine(g *grid.SparseGrid[Contents], from, to pos.P2) {
	inc := pos.P2{}
	if from.X != to.X {
		delta := to.X - from.X
		inc.X = delta / mtsmath.Abs(delta)
	}
	if from.Y != to.Y {
		delta := to.Y - from.Y
		inc.Y = delta / mtsmath.Abs(delta)
	}

	for p := from; !p.Equals(to); p.Add(inc) {
		g.Set(p, WALL)
	}
	g.Set(to, WALL)
}

func parseInput(lines []string) (*grid.SparseGrid[Contents], error) {
	g := grid.NewSparseGrid[Contents]()

	for i, line := range lines {
		coords := []pos.P2{}
		for _, s := range strings.Split(line, " -> ") {
			p, err := pos.P2FromString(s)
			if err != nil {
				return nil, fmt.Errorf("%d: bad coord: %v", i+1, err)
			}

			coords = append(coords, p)
		}

		for i := 1; i < len(coords); i++ {
			drawLine(g, coords[i-1], coords[i])
		}
	}

	return g, nil
}

func addSand(g *grid.SparseGrid[Contents], start pos.P2, floorY int) (cameToRest bool) {
	isOpen := func(p pos.P2) bool {
		if floorY != -1 && p.Y == floorY {
			return false
		}

		v, found := g.Get(p)
		return !found || v == START
	}

	p := start
	for i := 0; i < 10000; i++ {
		wants := []pos.P2{
			pos.P2{p.X, p.Y + 1},
			pos.P2{p.X - 1, p.Y + 1},
			pos.P2{p.X + 1, p.Y + 1},
		}

		var want pos.P2
		for _, w := range wants {
			if isOpen(w) {
				want = w
				break
			}
		}

		if want.Equals(pos.P2{}) {
			g.Set(p, SAND)
			return true // It came to rest
		}

		if floorY == -1 && want.Y > g.End().Y {
			return false // Falling forever
		}

		p = want
		logger.LogF("new want %v", p)
	}

	panic("too many")
}

func solveA(g *grid.SparseGrid[Contents], start pos.P2) int {
	var num int
	for num = 0; addSand(g, start, -1); num++ {
		if logger.Enabled() {
			fmt.Println("after grain", num+1)
			g.Dump(true, dumpMap)
		}

		if num > 10000 {
			panic("too many sands")
		}
	}

	return num
}

func solveB(g *grid.SparseGrid[Contents], start pos.P2) int {
	floorY := g.End().Y + 2

	num := 0
	for {
		if !addSand(g, start, floorY) {
			return num
		}
		num++ // we just added a grain

		if logger.Enabled() {
			fmt.Println("after grain", num)
			g.Dump(true, dumpMap)
		}

		if v, found := g.Get(start); found && v == SAND {
			return num
		}
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	g, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	start := pos.P2{500, 0}
	g.Set(start, START)

	fmt.Println("A", solveA(g.Clone(), start))
	fmt.Println("B", solveB(g.Clone(), start))
}
