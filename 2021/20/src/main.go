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
	"github.com/simmonmt/aoc/2021/common/intmath"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	input    = flag.String("input", "", "input file")
	numSteps = flag.Int("num_steps", 2, "numbef of steps")
)

type Board struct {
	m                      map[pos.P2]bool
	minX, maxX, minY, maxY int
}

func NewBoard() *Board {
	return &Board{
		m: map[pos.P2]bool{},
	}
}

func (b *Board) Set(p pos.P2, v bool) {
	if len(b.m) == 0 {
		b.minX, b.maxX, b.minY, b.maxY = p.X, p.X, p.Y, p.Y
	} else {
		b.minX = intmath.IntMin(p.X, b.minX)
		b.maxX = intmath.IntMax(p.X, b.maxX)
		b.minY = intmath.IntMin(p.Y, b.minY)
		b.maxY = intmath.IntMax(p.Y, b.maxY)
	}

	b.m[p] = v
}

func (b *Board) Get(p pos.P2) (v, found bool) {
	v, found = b.m[p]
	return
}

func (b *Board) NumSet() int {
	numSet := 0
	for _, v := range b.m {
		if v {
			numSet++
		}
	}
	return numSet
}

func (b *Board) Walk(cb func(pos.P2)) {
	for y := b.minY - 2; y <= b.maxY+2; y++ {
		for x := b.minX - 2; x <= b.maxX+2; x++ {
			cb(pos.P2{x, y})
		}
	}
}

func (b *Board) Dump() {
	for y := b.minY - 2; y <= b.maxY+2; y++ {
		for x := b.minX - 2; x <= b.maxX+2; x++ {
			p := pos.P2{X: x, Y: y}
			v, found := b.Get(p)
			if found {
				if v {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			} else {
				fmt.Print("_")
			}
		}
		fmt.Println()
	}
}

func (b *Board) allEquals(yRange, xRange []int, want bool) bool {
	if len(xRange) == 1 {
		xRange = []int{xRange[0], xRange[0]}
	}
	if len(yRange) == 1 {
		yRange = []int{yRange[0], yRange[0]}
	}

	for y := yRange[0]; y <= yRange[1]; y++ {
		for x := xRange[0]; x <= xRange[1]; x++ {
			p := pos.P2{x, y}
			v, found := b.Get(p)
			if !found {
				panic("unexpected unset")
			}
			if v != want {
				return false
			}
		}
	}
	return true
}

func (b *Board) Trim(unsetVal bool) {
	for {
		changed := false
		if b.allEquals([]int{b.minY}, []int{b.minX, b.maxX}, unsetVal) {
			logger.LogF("trim min Y")
			b.minY++
			changed = true
		}
		if b.allEquals([]int{b.maxY}, []int{b.minX, b.maxX}, unsetVal) {
			logger.LogF("trim max Y")
			b.maxY--
			changed = true
		}

		if b.allEquals([]int{b.minY, b.maxY}, []int{b.minX}, unsetVal) {
			logger.LogF("trim min X")
			b.minX++
			changed = true
		}
		if b.allEquals([]int{b.minY, b.maxY}, []int{b.maxX}, unsetVal) {
			logger.LogF("trim max X")
			b.maxX--
			changed = true
		}

		if !changed {
			break
		}
	}
}

func parseBoard(lines []string) *Board {
	board := NewBoard()
	for y := 0; y < len(lines); y = y + 1 {
		for x, r := range lines[y] {
			p := pos.P2{X: x, Y: y}
			board.Set(p, r == '#')
		}
	}
	return board
}

func readInput(path string) (enhancement []bool, board *Board, err error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, nil, err
	}

	enhancement = make([]bool, len(lines[0]))
	for i, r := range lines[0] {
		if r == '#' {
			enhancement[i] = true
		}
	}

	board = parseBoard(lines[2:])
	return
}

var (
	posNumOffsets = []pos.P2{
		pos.P2{-1, -1}, pos.P2{0, -1}, pos.P2{1, -1},
		pos.P2{-1, 0}, pos.P2{0, 0}, pos.P2{1, 0},
		pos.P2{-1, 1}, pos.P2{0, 1}, pos.P2{1, 1},
	}
)

func calcPosNum(board *Board, ctr pos.P2, unsetVal bool) int {
	num := 0
	for _, off := range posNumOffsets {
		num <<= 1

		p := ctr
		p.Add(off)

		v, found := board.Get(p)
		if (found && v) || (!found && unsetVal) {
			num |= 1
		}

	}
	return num
}

func runStep(stepNum int, enhancement []bool, board *Board) *Board {
	nb := NewBoard()

	unsetVal := false
	// step 1 read as black
	// step 2 read as enhancement[2]

	if stepNum%2 == 0 {
		unsetVal = enhancement[0]
	}

	board.Walk(func(p pos.P2) {
		posNum := calcPosNum(board, p, unsetVal)
		//logger.LogF("p %v pn %v = %v", p, posNum, enhancement[posNum])
		nb.Set(p, enhancement[posNum])
	})

	nb.Trim(false)

	return nb
}

func solve(numSteps int, enhancement []bool, board *Board) {
	if logger.Enabled() {
		fmt.Println("initial:")
		board.Dump()
	}

	for i := 1; i <= numSteps; i++ {
		board = runStep(i, enhancement, board)

		if logger.Enabled() {
			fmt.Printf("\nafter step %d:\n", i)
			board.Dump()
		}
	}

	fmt.Println(board.NumSet())
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	enhancement, board, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	solve(*numSteps, enhancement, board)
}
