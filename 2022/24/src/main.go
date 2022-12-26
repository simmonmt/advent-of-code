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
	"io"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2022/common/astar"
	"github.com/simmonmt/aoc/2022/common/dir"
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

func parseDir(r rune) dir.Dir {
	switch r {
	case '.':
		return dir.DIR_UNKNOWN
	case '<':
		return dir.DIR_WEST
	case '>':
		return dir.DIR_EAST
	case 'v':
		return dir.DIR_SOUTH
	case '^':
		return dir.DIR_NORTH
	default:
		panic("bad rune")
	}
}

func parseInput(lines []string) (g *grid.Grid[dir.Dir], startX, endX int) {
	for x, r := range lines[0] {
		if r == '.' {
			startX = x - 1
			break
		}
	}
	for x, r := range lines[len(lines)-1] {
		if r == '.' {
			endX = x - 1
		}
	}

	g = grid.New[dir.Dir](len(lines[0])-2, len(lines)-2)
	for y, line := range lines[1 : len(lines)-1] {
		for x, r := range line[1 : len(line)-1] {
			p := pos.P2{x, y}
			g.Set(p, parseDir(r))
		}
	}

	return
}

func encodeGrid(g *grid.Grid[dir.Dir]) map[dir.Dir][]*big.Int {
	out := map[dir.Dir][]*big.Int{
		dir.DIR_NORTH: make([]*big.Int, g.Width()),
		dir.DIR_SOUTH: make([]*big.Int, g.Width()),
		dir.DIR_WEST:  make([]*big.Int, g.Height()),
		dir.DIR_EAST:  make([]*big.Int, g.Height()),
	}

	for d, a := range out {
		for i := range a {
			out[d][i] = &big.Int{}
		}
	}

	g.Walk(func(p pos.P2, d dir.Dir) {
		if d == dir.DIR_UNKNOWN {
			return
		}

		a := out[d]
		if d == dir.DIR_NORTH || d == dir.DIR_SOUTH {
			a[p.X].SetBit(a[p.X], p.Y, 1)
		} else {
			a[p.Y].SetBit(a[p.Y], p.X, 1)
		}
	})

	return out
}

func encodingToGrid(enc map[dir.Dir][]*big.Int) *grid.Grid[any] {
	g := grid.New[any](len(enc[dir.DIR_NORTH]), len(enc[dir.DIR_WEST]))

	for d, a := range enc {
		var p pos.P2
		var aDest, nDest *int
		if d == dir.DIR_NORTH || d == dir.DIR_SOUTH {
			aDest, nDest = &p.X, &p.Y
		} else {
			aDest, nDest = &p.Y, &p.X
		}

		for i, n := range a {
			*aDest = i
			for j := 0; j < n.BitLen(); j++ {
				*nDest = j
				if n.Bit(j) == 0 {
					continue
				}

				v, _ := g.Get(p)
				if v == nil {
					g.Set(p, d)
				} else if _, ok := v.(dir.Dir); ok {
					g.Set(p, 2)
				} else {
					g.Set(p, v.(int)+1)
				}
			}
		}
	}

	return g
}

func dumpEncodedTo(enc map[dir.Dir][]*big.Int, w io.Writer) {
	g := encodingToGrid(enc)
	g.DumpTo(false, func(p pos.P2, v any, found bool) string {
		if v == nil {
			return "."
		} else if d, ok := v.(dir.Dir); ok {
			switch d {
			case dir.DIR_NORTH:
				return "^"
			case dir.DIR_SOUTH:
				return "V"
			case dir.DIR_WEST:
				return "<"
			case dir.DIR_EAST:
				return ">"
			default:
				panic("bad dir")
			}
		} else if n, ok := v.(int); ok {
			if n < 0 || n > 9 {
				panic("bad num")
			}
			return strconv.Itoa(n)
		} else {
			panic("bad cell")
		}
	}, w)
}

func dumpEncoded(enc map[dir.Dir][]*big.Int) {
	dumpEncodedTo(enc, os.Stdout)
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func lcm(a, b int) int {
	if a == 0 && b == 0 {
		return 0
	}
	return mtsmath.Abs(a) * (mtsmath.Abs(b) / gcd(a, b))
}

func advanceEncoding(inEnc map[dir.Dir][]*big.Int) map[dir.Dir][]*big.Int {
	h, w := len(inEnc[dir.DIR_WEST]), len(inEnc[dir.DIR_NORTH])
	outEnc := map[dir.Dir][]*big.Int{}

	for d, in := range inEnc {
		outEnc[d] = make([]*big.Int, len(in))

		for i, inNum := range in {
			outNum := &big.Int{}

			if d == dir.DIR_NORTH || d == dir.DIR_WEST {
				highBitNum := h - 1
				if d == dir.DIR_WEST {
					highBitNum = w - 1
				}

				carry := inNum.Bit(0)
				outNum.Rsh(inNum, 1)
				if carry != 0 {
					outNum.SetBit(outNum, highBitNum, 1)
				}
			} else {
				highBitNum := h - 1
				if d == dir.DIR_EAST {
					highBitNum = w - 1
				}

				carry := inNum.Bit(highBitNum)
				outNum.SetBit(inNum, highBitNum, 0)
				outNum.Lsh(outNum, 1)
				if carry != 0 {
					outNum.SetBit(outNum, 0, 1)
				}
			}

			outEnc[d][i] = outNum
		}
	}

	return outEnc
}

func makeAllEncodings(initial map[dir.Dir][]*big.Int) []map[dir.Dir][]*big.Int {
	h, w := len(initial[dir.DIR_WEST]), len(initial[dir.DIR_NORTH])
	numEnc := lcm(h, w)

	logger.LogF("making %d encodings for h %v w %v", numEnc, h, w)
	out := make([]map[dir.Dir][]*big.Int, numEnc)
	out[0] = initial

	for i := 1; i < numEnc; i++ {
		out[i] = advanceEncoding(out[i-1])
	}

	return out
}

type BoardStep struct {
	id int
	g  *grid.Grid[any]
}

func NewBoardStep(id int, enc map[dir.Dir][]*big.Int) *BoardStep {
	g := encodingToGrid(enc)

	return &BoardStep{
		id: id,
		g:  g,
	}
}

func (s *BoardStep) OpenNeighbors(p pos.P2) []pos.P2 {
	out := []pos.P2{}
	for _, n := range s.g.AllNeighbors(p, false) {
		if v, _ := s.g.Get(n); v == nil {
			out = append(out, n)
		}
	}

	return out
}

func (s *BoardStep) IsOpen(p pos.P2) bool {
	v, _ := s.g.Get(p)
	return v == nil
}

func encodeNodeName(stepID int, cur pos.P2) string {
	return fmt.Sprintf("%d/%v", stepID, cur)
}

func decodeNodeName(name string) (stepID int, cur pos.P2) {
	left, right, ok := strings.Cut(name, "/")
	if !ok {
		panic("bad cut")
	}

	var err error
	stepID, err = strconv.Atoi(left)
	if err != nil {
		panic("bad step ID")
	}

	cur, err = pos.P2FromString(right)
	if err != nil {
		panic("bad pos")
	}

	return
}

type astarClient struct {
	startPos, endPos pos.P2
	steps            []*BoardStep
}

func (c *astarClient) AllNeighbors(node string) []string {
	curStepID, curPos := decodeNodeName(node)

	nextStepID := (curStepID + 1) % len(c.steps)
	nextStep := c.steps[nextStepID]

	var mustMove bool
	if curPos.Equals(c.startPos) {
		mustMove = false
	} else {
		mustMove = !nextStep.IsOpen(curPos)
	}

	neighbors := nextStep.OpenNeighbors(curPos)
	if curPos.ManhattanDistance(c.endPos) == 1 {
		neighbors = append(neighbors, c.endPos)
	}

	logger.LogF("cur %v must %v neighbors %v",
		node, mustMove, neighbors)

	out := []string{}
	if !mustMove {
		out = append(out, encodeNodeName(nextStepID, curPos))
	}
	for _, n := range neighbors {
		out = append(out, encodeNodeName(nextStepID, n))
	}

	return out
}

func (c *astarClient) EstimateDistance(start, end string) uint {
	_, curPos := decodeNodeName(start)
	return uint(curPos.ManhattanDistance(c.endPos))
}

func (c *astarClient) NeighborDistance(n1, n2 string) uint { return 1 }

func (c *astarClient) GoalReached(cand, goal string) bool {
	_, curPos := decodeNodeName(cand)
	return curPos.Equals(c.endPos)
}

func solve(startPos, endPos pos.P2, startStep int, steps []*BoardStep) int {
	startStep = startStep % len(steps)

	startNode := encodeNodeName(startStep, startPos)
	client := &astarClient{
		startPos: startPos,
		endPos:   endPos,
		steps:    steps,
	}

	path := astar.AStar(startNode, "", client)
	logger.LogF("path is %v %v", len(path), path)

	return len(path) - 1
}

func solveA(g *grid.Grid[dir.Dir], startX, endX int) int {
	allEncodings := makeAllEncodings(encodeGrid(g))
	steps := make([]*BoardStep, len(allEncodings))
	for i, enc := range allEncodings {
		steps[i] = NewBoardStep(i, enc)
	}

	return solve(pos.P2{startX, -1}, pos.P2{endX, g.Height()}, 0, steps)
}

func solveB(g *grid.Grid[dir.Dir], startX, endX int) int {
	allEncodings := makeAllEncodings(encodeGrid(g))
	steps := make([]*BoardStep, len(allEncodings))
	for i, enc := range allEncodings {
		steps[i] = NewBoardStep(i, enc)
	}

	sum := 0
	lens := []int{}

	top, bottom := pos.P2{startX, -1}, pos.P2{endX, g.Height()}
	numSteps := solve(top, bottom, 0, steps)
	sum += numSteps
	lens = append(lens, numSteps)

	startStep := sum
	numSteps = solve(bottom, top, startStep, steps)
	sum += numSteps
	lens = append(lens, numSteps)

	startStep = sum
	numSteps = solve(top, bottom, startStep, steps)
	sum += numSteps
	lens = append(lens, numSteps)

	return sum
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

	g, startX, endX := parseInput(lines)

	fmt.Println("A", solveA(g, startX, endX))
	fmt.Println("B", solveB(g, startX, endX))
}
