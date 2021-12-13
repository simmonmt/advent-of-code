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
	"io"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/grid"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	instructionPattern = regexp.MustCompile(`^fold along (.)=([0-9]+)$`)
)

type Axis int

func (a Axis) String() string {
	if a == X_AXIS {
		return "x"
	} else {
		return "y"
	}
}

const (
	X_AXIS Axis = iota
	Y_AXIS
)

type Instruction struct {
	Axis  Axis
	Coord int
}

func dumpTo(w io.Writer, g *grid.Grid) {
	g.Walk(func(p pos.P2, value interface{}) {
		if value == true {
			fmt.Fprint(w, "#")
		} else {
			fmt.Fprint(w, ".")
		}
		if p.X == g.Width()-1 {
			fmt.Fprint(w, "\n")
		}
	})
}

func dump(g *grid.Grid) {
	dumpTo(os.Stdout, g)
}

func readInput(path string) (*grid.Grid, []Instruction, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, nil, err
	}

	maxX, maxY := 0, 0
	ps := []pos.P2{}

	lineNum := 0
	for _, line := range lines {
		lineNum++

		if line == "" {
			break
		}

		nums, err := filereader.ParseNumbersFromLine(line)
		if err != nil {
			return nil, nil, fmt.Errorf("%d: bad parse: %v",
				lineNum, err)
		}
		x, y := nums[0], nums[1]

		ps = append(ps, pos.P2{X: x, Y: y})
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}

	g := grid.New(maxX+1, maxY+1)
	for _, p := range ps {
		g.Set(p, true)
	}

	insts := []Instruction{}
	for _, line := range lines[lineNum:] {
		lineNum++

		parts := instructionPattern.FindStringSubmatch(line)
		if parts == nil {
			return nil, nil,
				fmt.Errorf("%d: bad instruction parse: %v",
					lineNum, line)
		}

		var axis Axis
		if parts[1] == "x" {
			axis = X_AXIS
		} else if parts[1] == "y" {
			axis = Y_AXIS
		} else {
			return nil, nil, fmt.Errorf("%d: bad axis %v",
				lineNum, parts[1])
		}

		coord, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, nil, fmt.Errorf("%d: bad fold coord: %v",
				lineNum, parts[2])
		}

		insts = append(insts, Instruction{
			Axis:  axis,
			Coord: coord,
		})
	}

	return g, insts, nil
}

func walkYFold(g *grid.Grid, coord int, other bool, cb func(p pos.P2)) {
	var minY, maxY int
	if other {
		minY, maxY = coord+1, g.Height()-1
	} else {
		minY, maxY = 0, coord-1
	}

	logger.LogF("other %v minY %d maxY %d\n", other, minY, maxY)

	for y := minY; y <= maxY; y++ {
		for x := 0; x < g.Width(); x++ {
			cb(pos.P2{X: x, Y: y})
		}
	}
}

func walkXFold(g *grid.Grid, coord int, other bool, cb func(p pos.P2)) {
	var minX, maxX int
	if other {
		minX, maxX = coord+1, g.Width()-1
	} else {
		minX, maxX = 0, coord-1
	}

	for y := 0; y < g.Height(); y++ {
		for x := minX; x <= maxX; x++ {
			cb(pos.P2{X: x, Y: y})
		}
	}
}

func mapYFold(coord int, p pos.P2) pos.P2 {
	y := (coord - 1) - (p.Y - coord - 1)
	return pos.P2{X: p.X, Y: y}
}

func mapXFold(coord int, p pos.P2) pos.P2 {
	x := (coord - 1) - (p.X - coord - 1)
	return pos.P2{X: x, Y: p.Y}
}

type FoldWalker func(*grid.Grid, int, bool, func(p pos.P2))
type FoldMapper func(int, pos.P2) pos.P2

func performFold(g *grid.Grid, inst Instruction) *grid.Grid {
	var walker FoldWalker
	var mapper FoldMapper

	if inst.Axis == X_AXIS {
		walker = walkXFold
		mapper = mapXFold
	} else {
		walker = walkYFold
		mapper = mapYFold
	}

	newWidth, newHeight := g.Width(), g.Height()
	if inst.Axis == X_AXIS {
		newWidth = inst.Coord
	} else {
		newHeight = inst.Coord
	}
	logger.LogF("g was w %v h %v, ng is w %v h %v",
		g.Width(), g.Height(), newWidth, newHeight)

	ng := grid.New(newWidth, newHeight)
	walker(g, inst.Coord, false, func(p pos.P2) {
		ng.Set(p, g.Get(p))
	})

	walker(g, inst.Coord, true, func(p pos.P2) {
		v := g.Get(p)
		mp := mapper(inst.Coord, p)
		mv := g.Get(mp)
		r := v == true || mv == true

		// logger.LogF("copying from %v (%v) to %v (%v) = %v",
		// 	p, v, mp, mv, r)

		ng.Set(mp, r)
	})

	return ng
}

func solveA(g *grid.Grid, insts []Instruction) {
	logger.LogF("inst %v", insts[0])
	g = performFold(g, insts[0])

	num := 0
	g.Walk(func(p pos.P2, v interface{}) {
		if v == true {
			num++
		}
	})

	fmt.Println("A", num)
}

func solveB(g *grid.Grid, insts []Instruction) {
	for _, inst := range insts {
		logger.LogF("%v", inst)
		g = performFold(g, inst)
	}

	fmt.Println("B:")
	dump(g)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	g, insts, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	solveA(g, insts)
	solveB(g, insts)
}
