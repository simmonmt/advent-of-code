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

// B: 1710 too low
// B: 2016 too low

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2022/common/dir"
	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/grid"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func buildGrid(lines []string) (*grid.Grid, error) {
	return grid.NewFromLines(lines, func(r rune) (any, error) {
		if r >= '0' && r <= '9' {
			return int(r - '0'), nil
		}
		return nil, fmt.Errorf("bad cell %s", string(r))
	})
}

func readInput(path string) (*grid.Grid, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	return buildGrid(lines)
}

func solveA(g *grid.Grid) int {
	type Side struct {
		Start, Inc, In pos.P2
	}

	left, right := pos.P2{-1, 0}, pos.P2{1, 0}
	up, down := pos.P2{0, -1}, pos.P2{0, 1}

	sides := []Side{
		Side{ // top
			Start: pos.P2{0, 0}, // top left
			Inc:   right,
			In:    down,
		},
		Side{ // left
			Start: pos.P2{0, 0}, // top left
			Inc:   down,
			In:    right,
		},
		Side{ // right
			Start: pos.P2{g.Width() - 1, 0}, // top right
			Inc:   down,
			In:    left,
		},
		Side{ // bottom
			Start: pos.P2{0, g.Height() - 1}, // bottom left
			Inc:   right,
			In:    up,
		},
	}

	visibles := map[pos.P2]bool{}
	for _, side := range sides {
		for p := side.Start; g.IsValid(p); p.Add(side.Inc) {
			maxHeight := -1
			for in := p; g.IsValid(in) && maxHeight != 9; in.Add(side.In) {
				height := g.Get(in).(int)

				if height > maxHeight {
					maxHeight = height

					if _, found := visibles[in]; !found {
						if logger.Enabled() {
							if p.Equals(in) {
								fmt.Println("edge", in)
							} else {
								fmt.Println("interior", in)
							}
						}
						visibles[in] = true
					}
				}
			}
		}
	}

	return len(visibles)
}

func lookInDir(g *grid.Grid, center pos.P2, d dir.Dir) int {
	centerHeight := g.Get(center).(int)
	canSee := 0
	for p := d.From(center); g.IsValid(p); p = d.From(p) {
		canSee++
		if height := g.Get(p).(int); height >= centerHeight {
			break
		}
	}
	return canSee
}

func scoreTree(g *grid.Grid, center pos.P2) int {
	total := 1
	for _, d := range dir.AllDirs {
		score := lookInDir(g, center, d)
		//logger.LogF("tree %v dir %v score %v", center, d, score)
		total *= score
	}
	return total
}

func solveB(g *grid.Grid) int {
	maxScore := 0
	g.Walk(func(p pos.P2, v any) {
		if score := scoreTree(g, p); score > maxScore {
			logger.LogF("new max score at %v: %v", p, score)
			maxScore = score
		}
	})
	return maxScore
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

	fmt.Println("A", solveA(g))
	fmt.Println("B", solveB(g))
}
