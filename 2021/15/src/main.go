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
	"strconv"

	"github.com/simmonmt/aoc/2021/common/astar"
	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/grid"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	return lines, err
}

func dump(g *grid.Grid) {
	g.Walk(func(p pos.P2, v interface{}) {
		fmt.Printf("%v", v.(int))
		if p.X == g.Width()-1 {
			fmt.Println()
		}
	})
}

func decodePos(s string) pos.P2 {
	p, err := pos.P2FromString(s)
	if err != nil {
		panic("bad pos")
	}
	return p
}

func encodePos(p pos.P2) string {
	return p.String()
}

type astarClient struct {
	g *grid.Grid
}

func (c *astarClient) AllNeighbors(start string) []string {
	ns := c.g.AllNeighbors(decodePos(start), false)
	out := make([]string, len(ns))

	for i, n := range ns {
		out[i] = encodePos(n)
	}

	return out
}

func (c *astarClient) EstimateDistance(start, end string) uint {
	sp, ep := decodePos(start), decodePos(end)
	return uint(sp.ManhattanDistance(ep))
}

func (c *astarClient) NeighborDistance(n1, n2 string) uint {
	n2p := decodePos(n2)
	return uint(c.g.Get(n2p).(int))
}

func (c *astarClient) GoalReached(cand, goal string) bool {
	return cand == goal
}

func solveGrid(g *grid.Grid) int {
	start := pos.P2{X: 0, Y: 0}
	end := pos.P2{X: g.Width() - 1, Y: g.Height() - 1}
	client := &astarClient{g: g}

	result := astar.AStar(encodePos(start), encodePos(end), client)

	total := -g.Get(start).(int)
	for _, p := range result {
		total += g.Get(decodePos(p)).(int)
	}

	return total
}

func solveA(g *grid.Grid) {
	fmt.Println("A", solveGrid(g))
}

func solveB(og *grid.Grid) {
	mult := 5
	g := grid.New(og.Width()*mult, og.Height()*mult)

	for yCopy := 0; yCopy < mult; yCopy++ {
		for xCopy := 0; xCopy < mult; xCopy++ {
			og.Walk(func(p pos.P2, v interface{}) {
				np := pos.P2{
					X: p.X + og.Width()*xCopy,
					Y: p.Y + og.Height()*yCopy,
				}

				nv := v.(int)
				if xCopy > 0 || yCopy > 0 {
					nv += xCopy + yCopy
					for nv >= 10 {
						nv -= 9
					}
				}
				g.Set(np, nv)
			})
		}
	}

	fmt.Println("B", solveGrid(g))
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

	g, err := grid.NewFromLines(lines, func(r rune) (interface{}, error) {
		d, err := strconv.Atoi(string(r))
		return d, err
	})

	solveA(g)
	solveB(g)
}
