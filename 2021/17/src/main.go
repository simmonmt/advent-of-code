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
	"regexp"
	"strconv"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	inputPattern = regexp.MustCompile(
		`^target area: x=(-?\d+)\.\.(-?\d+), y=(-?\d+)\.\.(-?\d+)$`)
)

type Box struct {
	P1, P2 pos.P2
}

func NewBox(p1, p2 pos.P2) *Box {
	b := &Box{
		P1: p1,
		P2: p2,
	}
	if b.P1.X > b.P2.X {
		b.P1.X, b.P2.X = b.P2.X, b.P1.X
	}
	if b.P1.Y > b.P2.Y {
		b.P1.Y, b.P2.Y = b.P2.Y, b.P1.Y
	}
	return b
}

func (b *Box) Contains(p pos.P2) bool {
	if p.X < b.P1.X || p.X > b.P2.X {
		return false
	}
	if p.Y < b.P1.Y || p.Y > b.P2.Y {
		return false
	}
	return true
}

func parseInts(strs []string) ([]int, error) {
	out := []int{}
	for _, s := range strs {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", s, err)
		}
		out = append(out, int(v))
	}
	return out, nil
}

func readInput(path string) (*Box, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("wanted 1 line, got %v", len(lines))
	}

	line := lines[0]

	parts := inputPattern.FindStringSubmatch(line)
	if parts == nil {
		return nil, fmt.Errorf("parse error")
	}

	nums, err := parseInts(parts[1:])
	if err != nil {
		return nil, err
	}

	var min, max pos.P2
	min.X, max.X, min.Y, max.Y = nums[0], nums[1], nums[2], nums[3]

	return NewBox(min, max), nil
}

func walkCurve(xVel, yVel int, cb func(p pos.P2) bool) {
	p := pos.P2{X: 0, Y: 0}

	for {
		p.X += xVel
		p.Y += yVel

		if xVel > 0 {
			xVel--
		} else if xVel < 0 {
			xVel++
		}

		yVel--

		if !cb(p) {
			return
		}
	}
}

func evalCurve(xVel, yVel int, target *Box) (inTarget bool, yMax int) {
	inTarget = false
	yMax = 0

	walkCurve(xVel, yVel, func(p pos.P2) bool {
		// if xVel == 11 && yVel == 65 {
		// 	logger.LogF("walking at %v (target %+v)", p, target)
		// }

		if p.Y > yMax {
			yMax = p.Y
		}

		if target.Contains(p) {
			inTarget = true
			return false
		}

		if p.Y < target.P1.Y {
			return false
		}

		return true
	})

	return
}

func findInitialSolution(target *Box) (xVel, yVel, yMax int, found bool) {
	xVel = 1
	for xDist := xVel; xDist < target.P1.X; xDist += xVel {
		if xDist > target.P2.X {
			panic("x overshot")
		}
		xVel++
	}

	found = false
	for yVel = 1; yVel < 1000 && !found; yVel++ { // WAG
		found, yMax = evalCurve(xVel, yVel, target)
	}

	return
}

// This function searches for a solution. It treats the initial trajectories as
// {x,y} coordinates, and wanders through the x,y space looking for the
// maximum. It falls apart on the real input for two reasons:
//
// 1. As implemented it assumes there's only one maximum. Unfortunately the
//    input has local maxima.
// 2. A trivial modification (not resetting cands when a new ymax is found)
//    would let it work in spite of local maxima. Unfortunately it still assumes
//    that it can wander between all trajectories that'll reach the box, which
//    means it breaks if the space of trajectories contains disconnected
//    islands. It'll fully explore the island it starts on, but won't discover
//    the others.
//
// Sigh. It's unfortunate that exploration doesn't work because brute forcing is
// boring. Maybe there's a way to discover all of the islands?
func searchSolveA(target *Box) {
	xVel, yVel, yMax, found := findInitialSolution(target)
	if !found {
		panic("failed to find y")
	}

	cands := []pos.P2{pos.P2{X: xVel, Y: yVel}}
	candsYMax := yMax
	seen := map[pos.P2]bool{cands[0]: true}

	for len(cands) > 0 {
		cand := cands[0]
		cands = cands[1:]
		logger.LogF("candidate %v (%d more)", cand, len(cands))

		for _, n := range cand.AllNeighbors(true) {
			if _, found := seen[n]; found {
				continue
			}
			seen[n] = true

			found, yMax := evalCurve(n.X, n.Y, target)
			if !found {
				logger.LogF("candidate %v neighbor %v missed", cand, n)
				continue
			}

			if yMax > candsYMax {
				logger.LogF("candidate %v neighbor %v new yMax %v", cand, n, yMax)
				candsYMax = yMax
				cands = []pos.P2{n}
			} else if yMax == candsYMax {
				logger.LogF("candidate %v neighbor %v same yMax %v", cand, n, yMax)
				cands = append(cands, n)
			} else {
				logger.LogF("candidate %v neighbor %v too low yMax %v", cand, n, yMax)
			}

		}
	}

	fmt.Println("A", candsYMax)
}

func bruteSolve(target *Box) {
	yMax := 0

	seen := map[pos.P2]int{}
	for y := -1000; y < 1000; y++ {
		for x := 0; x < 1000; x++ {
			found, thisYMax := evalCurve(x, y, target)
			if found {
				seen[pos.P2{x, y}] = yMax
				if thisYMax > yMax {
					logger.LogF("new YMax %v at %v,%v", yMax, x, y)
					yMax = thisYMax
				}
			}
		}
	}

	fmt.Println("A", yMax)
	fmt.Println("B", len(seen))
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	target, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	logger.LogF("target: %+v", target)

	bruteSolve(target)
}
