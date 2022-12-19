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

func (f *PartFactory) Peek() int {
	return f.next
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

type DirFactory struct {
	dirs string
	next int
}

func NewDirFactory(dirs string) *DirFactory {
	return &DirFactory{dirs, 0}
}

func (f *DirFactory) Next() byte {
	r := f.dirs[f.next]
	f.next = (f.next + 1) % len(f.dirs)
	return r
}

func (f *DirFactory) Peek() int {
	return f.next
}

func runPart(g *grid.SparseGrid, partFactory *PartFactory, dirFactory *DirFactory) pos.P2 {
	var bottomLeft pos.P2
	if g.Empty() {
		bottomLeft = pos.P2{2, 3}
	} else {
		bottomLeft = pos.P2{2, g.End().Y + 4}
	}

	part := partFactory.Next(bottomLeft)

	for i := 0; ; i++ {
		if dirFactory.Next() == '<' {
			part.ShiftLeft(g)
		} else {
			part.ShiftRight(g)
		}

		if !part.ShiftDown(g) {
			part.Place(g)
			return part.Pos()
		}
	}
}

func measureHeight(dirs string, numParts int) int {
	g := grid.NewSparseGrid()
	partFactory := NewPartFactory()
	dirFactory := NewDirFactory(dirs)

	for i := 0; i < numParts; i++ {
		runPart(g, partFactory, dirFactory)
	}

	return g.End().Y + 1
}

func solveA(dirs string) int {
	return measureHeight(dirs, 2022)
}

type HistoryKey struct {
	posns   [5]pos.P2
	dirNext int
}

type HistoryEnt struct {
	lastPartIdx int
	height      int
}

func findRepeat(dirs string) (first, second HistoryEnt) {
	g := grid.NewSparseGrid()
	partFactory := NewPartFactory()
	dirFactory := NewDirFactory(dirs)

	posns := [5]pos.P2{}
	history := map[HistoryKey]HistoryEnt{}

	for i := 0; ; i++ {
		partIdx := partFactory.Peek()
		posns[partIdx] = runPart(g, partFactory, dirFactory)

		if partIdx == 4 {
			for i := 1; i < 5; i++ {
				posns[i].Y -= posns[0].Y
			}
			posns[0].Y = 0

			height := g.End().Y + 1
			key := HistoryKey{
				posns, dirFactory.Peek(),
			}

			if ent, found := history[key]; found {
				return ent, HistoryEnt{i, height}
			} else {
				history[key] = HistoryEnt{i, height}
			}
		}
	}
}

func measureTallHeight(dirs string, numParts int64) int64 {
	first, second := findRepeat(dirs)

	logger.LogF("seq part [%v-%v]>%v repeats at [%v-%v]>%v",
		first.lastPartIdx-4, first.lastPartIdx, first.height,
		second.lastPartIdx-4, second.lastPartIdx, second.height)

	// play parts 0-first => first.height

	// every second.lastPartIdx-first.lastPartIdx gets us another
	// second.height-first.height

	prologueLen := int64(first.lastPartIdx + 1)
	prologueHeight := int64(first.height)

	repLen := int64(second.lastPartIdx - first.lastPartIdx)
	repHeight := int64(second.height - first.height)

	numReps := (numParts - prologueLen) / repLen
	suffixLen := (numParts - prologueLen) % repLen

	// we can remove entire repeats entirely
	simLen := prologueLen + suffixLen

	simHeight := int64(measureHeight(dirs, int(simLen)))

	suffixHeight := simHeight - prologueHeight

	totalHeight := prologueHeight + repHeight*numReps + suffixHeight
	//fmt.Println(prologueHeight, repHeight, numReps, suffixHeight)

	return totalHeight
}

func solveB(dirs string) int64 {
	need := int64(1000000000000)
	return measureTallHeight(dirs, need)
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
