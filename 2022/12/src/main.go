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

	"github.com/simmonmt/aoc/2022/common/astar"
	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/grid"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func parseInput(lines []string) (g *grid.Grid[int], start, end pos.P2, err error) {
	g, err = grid.NewFromLines(lines, func(p pos.P2, r rune) (int, error) {
		if r == 'S' {
			start = p
			return 1, nil
		} else if r == 'E' {
			end = p
			return 26, nil
		} else if r >= 'a' && r <= 'z' {
			return int(r - 'a' + 1), nil
		}
		return -1, fmt.Errorf("bad value %d", r)
	})

	return
}

type astarClient struct {
	g *grid.Grid[int]
}

func (ac *astarClient) AllNeighbors(start string) []string {
	p, err := pos.P2FromString(start)
	if err != nil {
		panic("bad start")
	}
	pHeight := ac.g.Get(p)

	out := []string{}
	for _, n := range ac.g.AllNeighbors(p, false) {
		if nHeight := ac.g.Get(n); nHeight <= pHeight+1 {
			out = append(out, n.String())
		}
	}
	return out
}

func (ac *astarClient) EstimateDistance(start, end string) uint {
	sp, err := pos.P2FromString(start)
	if err != nil {
		panic("bad start")
	}

	ep, err := pos.P2FromString(end)
	if err != nil {
		panic("bad end")
	}

	return uint(sp.ManhattanDistance(ep))
}

func (ac *astarClient) NeighborDistance(n1, n2 string) uint {
	return 1
}

func (ac *astarClient) GoalReached(cand, goal string) bool {
	return cand == goal
}

func solveA(g *grid.Grid[int], start, end pos.P2) int {
	ac := &astarClient{g}
	path := astar.AStar(start.String(), end.String(), ac)
	return len(path) - 1
}

func solveB(g *grid.Grid[int], end pos.P2) int {
	ac := &astarClient{g}

	starts := []string{}
	ac.g.Walk(func(p pos.P2, height int) {
		if height == 1 {
			starts = append(starts, p.String())
		}
	})

	shortest := -1
	for _, start := range starts {
		path := astar.AStar(start, end.String(), ac)
		if len(path) == 0 {
			continue
		}

		if shortest == -1 || len(path)-1 < shortest {
			shortest = len(path) - 1
		}
	}

	return shortest
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	g, start, end, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(g, start, end))
	fmt.Println("B", solveB(g, end))
}
