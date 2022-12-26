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

	"github.com/simmonmt/aoc/2022/common/dir"
	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/mtsmath"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) (map[pos.P2]bool, error) {
	out := map[pos.P2]bool{}
	for y, line := range lines {
		for x, r := range line {
			p := pos.P2{x, y}
			if r == '#' {
				out[p] = true
			} else if r != '.' {
				return nil, fmt.Errorf("bad char at %v", p)
			}
		}
	}
	return out, nil
}

func proposeMove(cur pos.P2, elves map[pos.P2]bool, dirs []dir.Dir) (pos.P2, bool) {
	canMove := false
	for _, p := range cur.AllNeighbors(true) {
		if _, found := elves[p]; found {
			canMove = true
			break
		}
	}

	if !canMove {
		return cur, false
	}

	for _, d := range dirs {
		fwd := d.From(cur)
		tries := []pos.P2{fwd, fwd, fwd}
		switch {
		case d == dir.DIR_NORTH || d == dir.DIR_SOUTH:
			tries[0].X--
			tries[2].X++
		case d == dir.DIR_EAST || d == dir.DIR_WEST:
			tries[0].Y--
			tries[2].Y++
		}

		foundElf := false
		for _, p := range tries {
			if _, found := elves[p]; found {
				foundElf = true
				break
			}
		}

		if !foundElf {
			return fwd, true
		}
	}

	return cur, false
}

func playRound(elves map[pos.P2]bool, dirs []dir.Dir) (map[pos.P2]bool, int) {
	// first half
	dests := map[pos.P2]int{}
	proposals := map[pos.P2]pos.P2{}
	noMoves := map[pos.P2]bool{}

	for elfPos := range elves {
		p, canMove := proposeMove(elfPos, elves, dirs)
		if canMove {
			proposals[elfPos] = p
			dests[p]++
		} else {
			noMoves[p] = true
		}
	}

	// second half
	newElves := map[pos.P2]bool{}
	for p := range noMoves {
		newElves[p] = elves[p]
	}

	numMoves := 0
	for from, to := range proposals {
		if n := dests[to]; n == 1 {
			newElves[to] = true
			numMoves++
		} else {
			newElves[from] = true
		}
	}

	return newElves, numMoves
}

func playGame(elves map[pos.P2]bool, maxRounds int) (map[pos.P2]bool, int) {
	dirs := []dir.Dir{dir.DIR_NORTH, dir.DIR_SOUTH, dir.DIR_WEST, dir.DIR_EAST}
	for i := 0; maxRounds < 0 || i < maxRounds; i++ {
		var numMoves int
		elves, numMoves = playRound(elves, dirs)
		if numMoves == 0 {
			logger.LogF("stopped early after %d rounds", i+1)
			return elves, i + 1
		}

		firstDir := dirs[0]
		dirs = dirs[1:]
		dirs = append(dirs, firstDir)
	}

	logger.LogF("stopped after %d rounds", maxRounds)
	return elves, maxRounds
}

func calculateScore(elves map[pos.P2]bool) int {
	start, end := pos.P2{0, 0}, pos.P2{-1, -1}
	for elf := range elves {
		if end.X < start.X {
			start, end = elf, elf
		} else {
			start.X = mtsmath.Min(start.X, elf.X)
			start.Y = mtsmath.Min(start.Y, elf.Y)
			end.X = mtsmath.Max(end.X, elf.X)
			end.Y = mtsmath.Max(end.Y, elf.Y)
		}
	}

	logger.LogF("bounds %v %v", start, end)

	sum := 0
	for y := start.Y; y <= end.Y; y++ {
		for x := start.X; x <= end.X; x++ {
			p := pos.P2{x, y}
			if _, found := elves[p]; !found {
				sum++
			}
		}
	}
	return sum
}

func solveA(elves map[pos.P2]bool) int {
	elves, _ = playGame(elves, 10)
	return calculateScore(elves)
}

func solveB(elves map[pos.P2]bool) int {
	_, answer := playGame(elves, -1)
	return answer
}

func cloneElves(in map[pos.P2]bool) map[pos.P2]bool {
	out := map[pos.P2]bool{}
	for k, v := range in {
		out[k] = v
	}
	return out
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

	elves, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(cloneElves(elves)))
	fmt.Println("B", solveB(cloneElves(elves)))
}
