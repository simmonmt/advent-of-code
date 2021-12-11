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
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/grid"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	input    = flag.String("input", "", "input file")
	numSteps = flag.Int("num_steps", 100, "number of steps")
)

func gridFromInput(lines []string) *grid.IntGrid {
	g := grid.NewInt(10, 10)

	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[0]); x++ {
			p := pos.P2{X: x, Y: y}
			g.SetInt(p, int(lines[y][x]-'0'))
		}
	}

	return g
}

func readInput(path string) (*grid.IntGrid, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	return gridFromInput(lines), nil
}

func step(g *grid.IntGrid) int {
	g.WalkInt(func(p pos.P2, v int) {
		g.SetInt(p, v+1)
	})

	oldNumFlashes := 0
	flashed := map[pos.P2]bool{}

	for {
		g.WalkInt(func(p pos.P2, v int) {
			if v <= 9 {
				return // skip unflashed
			}

			if _, found := flashed[p]; found {
				return // already flashed
			}
			flashed[p] = true

			for _, n := range g.AllNeighbors(p, true) {
				g.SetInt(n, g.GetInt(n)+1)
			}
		})

		if len(flashed) == oldNumFlashes {
			break
		}
		oldNumFlashes = len(flashed)
	}

	g.WalkInt(func(p pos.P2, v int) {
		if v > 9 {
			g.SetInt(p, 0)
		}
	})

	return len(flashed)
}

func solve(g *grid.IntGrid, numSteps int) {
	numFlashes := 0
	for i := 1; i <= numSteps; i++ {
		numFlashes += step(g)
		if logger.Enabled() {
			logger.LogLn()
			logger.LogF("After step %d:", i)
			g.Dump()
		}
	}

	fmt.Println("#flashes", numFlashes)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	g, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	if logger.Enabled() {
		logger.LogLn("initial board")
		g.Dump()
	}

	solve(g, *numSteps)
}
