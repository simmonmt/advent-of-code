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

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/grid"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type CoordTranslator struct {
	hGrow, vGrow map[int]int
}

func NewCoordTranslator() *CoordTranslator {
	return &CoordTranslator{
		hGrow: map[int]int{},
		vGrow: map[int]int{},
	}
}

func (ct *CoordTranslator) InflateX(x, factor int) {
	n := ct.hGrow[x]
	ct.hGrow[x] = max(n, 1) * factor
}

func (ct *CoordTranslator) InflateY(y, factor int) {
	n := ct.vGrow[y]
	ct.vGrow[y] = max(n, 1) * factor
}

func inflate(grow map[int]int, c int) int {
	// This could be faster, I'm sure, but it's late.
	out := 0
	for i := 0; i < c; i++ {
		if sz, found := grow[i]; found {
			out += sz
		} else {
			out += 1
		}
	}
	return out
}

func (ct *CoordTranslator) RealToInflated(real pos.P2) pos.P2 {
	return pos.P2{
		X: inflate(ct.hGrow, real.X),
		Y: inflate(ct.vGrow, real.Y),
	}
}

func parseInput(lines []string) (*grid.Grid[bool], error) {
	g, err := grid.NewFromLines[bool](lines, func(p pos.P2, r rune) (bool, error) {
		if r == '.' {
			return false, nil
		} else if r == '#' {
			return true, nil
		} else {
			return true, fmt.Errorf("bad %c at %v", r, p)
		}
	})

	return g, err
}

func solve(g *grid.Grid[bool], factor int) int64 {
	xlate := NewCoordTranslator()
	stars := []pos.P2{}
	xCount := map[int]int{}
	yCount := map[int]int{}

	g.Walk(func(p pos.P2, star bool) {
		if star {
			stars = append(stars, p)
			xCount[p.X]++
			yCount[p.Y]++
		}
	})

	for x := 0; x < g.Width(); x++ {
		if _, found := xCount[x]; !found {
			xlate.InflateX(x, factor)
		}
	}
	for y := 0; y < g.Height(); y++ {
		if _, found := yCount[y]; !found {
			xlate.InflateY(y, factor)
		}
	}

	out := int64(0)
	for i := 0; i < len(stars)-1; i++ {
		for j := i + 1; j < len(stars); j++ {
			from := xlate.RealToInflated(stars[i])
			to := xlate.RealToInflated(stars[j])
			dist := from.ManhattanDistance(to)
			out += int64(dist)
		}
	}
	return out
}

func solveA(g *grid.Grid[bool]) int64 {
	return solve(g, 2)
}

func solveB(g *grid.Grid[bool]) int64 {
	return solve(g, 1000000)
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
	fmt.Println("B", solveB(input))
}
