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
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) (string, error) {
	if len(lines) != 1 {
		return "", fmt.Errorf("wrong input size %d", len(lines))
	}

	return lines[0], nil
}

type Part interface {
	ShiftLeft(g *grid.SparseGrid) bool
	ShiftRight(g *grid.SparseGrid) bool
	ShiftDown(g *grid.SparseGrid) bool
	Place(g *grid.SparseGrid)
	Pos() pos.P2
}

type partImpl struct {
	bottomLeft pos.P2
	elems      []pos.P2
}

func NewPart(bottomLeft pos.P2, elems []pos.P2) Part {
	return &partImpl{
		bottomLeft: bottomLeft,
		elems:      elems,
	}
}

func (p *partImpl) Pos() pos.P2 {
	return p.bottomLeft
}

func (p *partImpl) shift(g *grid.SparseGrid, elemOff pos.P2) bool {
	for _, elem := range p.elems {
		new := p.bottomLeft
		new.Add(elem)
		new.Add(elemOff)

		if new.X < 0 || new.X > 6 {
			return false
		}
		if new.Y < 0 {
			return false // through floor
		}

		if _, found := g.Get(new); found {
			return false // placed rock
		}
	}

	p.bottomLeft.Add(elemOff)
	return true
}

func (p *partImpl) ShiftLeft(g *grid.SparseGrid) bool {
	return p.shift(g, pos.P2{-1, 0})
}

func (p *partImpl) ShiftRight(g *grid.SparseGrid) bool {
	return p.shift(g, pos.P2{1, 0})
}

func (p *partImpl) ShiftDown(g *grid.SparseGrid) bool {
	return p.shift(g, pos.P2{0, -1})
}

func (p *partImpl) Place(g *grid.SparseGrid) {
	for _, elem := range p.elems {
		p := p.bottomLeft
		p.Add(elem)
		g.Set(p, true)
	}
}

var (
	minusElems = []pos.P2{
		pos.P2{0, 0}, pos.P2{1, 0},
		pos.P2{2, 0}, pos.P2{3, 0},
	}

	plusElems = []pos.P2{
		pos.P2{1, 0},
		pos.P2{0, 1}, pos.P2{1, 1}, pos.P2{2, 1},
		pos.P2{1, 2},
	}

	ellElems = []pos.P2{
		pos.P2{2, 2},
		pos.P2{2, 1},
		pos.P2{0, 0}, pos.P2{1, 0}, pos.P2{2, 0},
	}

	pipeElems = []pos.P2{
		pos.P2{0, 3},
		pos.P2{0, 2},
		pos.P2{0, 1},
		pos.P2{0, 0},
	}

	squareElems = []pos.P2{
		pos.P2{0, 0}, pos.P2{1, 0},
		pos.P2{0, 1}, pos.P2{1, 1},
	}

	orderedElems = [][]pos.P2{
		minusElems, plusElems, ellElems, pipeElems, squareElems,
	}
)

type PartFactory struct {
	next int
}

func NewPartFactory() *PartFactory {
	return &PartFactory{0}
}

func (f *PartFactory) Next(bottomLeft pos.P2) Part {
	part := NewPart(bottomLeft, orderedElems[f.next])
	f.next = (f.next + 1) % len(orderedElems)
	return part
}

func dumpGrid(g *grid.SparseGrid) {
	var b strings.Builder
	g.DumpTo(true, func(p pos.P2, v any, found bool) string {
		if found {
			return "#"
		} else {
			return "."
		}
	}, &b)

	lines := strings.Split(b.String(), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		fmt.Println(lines[i])
	}
}

func solveA(dirs string) int {
	partFactory := NewPartFactory()
	g := grid.NewSparseGrid()

	var part Part
	numParts := 0
	for i := 0; ; i++ {
		if part == nil {
			var bottomLeft pos.P2
			if numParts == 0 {
				bottomLeft = pos.P2{2, 3}
			} else {
				bottomLeft = pos.P2{2, g.End().Y + 4}
			}

			part = partFactory.Next(bottomLeft)
			numParts++
			//logger.LogF("%d: new at %v", numParts, part.Pos())
		}

		push := dirs[i%len(dirs)]
		if push == '<' {
			part.ShiftLeft(g)
		} else {
			part.ShiftRight(g)
		}

		//logger.LogF("%d: side push %v; new bl %v", numParts, string(push), part.Pos())

		if !part.ShiftDown(g) {
			part.Place(g)
			part = nil

			// if logger.Enabled() {
			// 	dumpGrid(g)
			// }

			if numParts == 2022 {
				break
			}
		}
	}

	return g.End().Y + 1
}

func solveB(dirs string) int {
	return -1
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

	line, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(line))
	fmt.Println("B", solveB(line))
}
