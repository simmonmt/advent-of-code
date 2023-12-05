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
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Color int

const (
	RED Color = iota
	GREEN
	BLUE
)

type Game struct {
	Num      int
	MoveSets [][]Move
}

type Move struct {
	Num   int
	Color Color
}

var (
	colorNames = map[string]Color{"red": RED, "green": GREEN, "blue": BLUE}
)

func parseMove(str string) (Move, error) {
	move := Move{}
	numStr, colorName, found := strings.Cut(str, " ")
	if !found {
		return move, fmt.Errorf("bad move")
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return move, fmt.Errorf("bad move num")
	}
	move.Num = num

	color, found := colorNames[colorName]
	if !found {
		return move, fmt.Errorf("bad move color")
	}
	move.Color = color

	return move, nil
}

func parseMoveSet(str string) ([]Move, error) {
	set := []Move{}
	parts := strings.Split(str, ", ")
	for _, part := range parts {
		move, err := parseMove(part)
		if err != nil {
			return nil, err
		}
		set = append(set, move)
	}
	return set, nil
}

func parseInput(lines []string) ([]*Game, error) {
	out := []*Game{}

	for _, line := range lines {
		game := &Game{}

		before, after, found := strings.Cut(line, ": ")
		if !found {
			return nil, fmt.Errorf("no game prefix in line %v", line)
		}

		_, numStr, found := strings.Cut(before, " ")
		if !found {
			return nil, fmt.Errorf("no game number in line %v", line)
		}

		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("bad game number in line %v", line)
		}
		game.Num = num

		moveSets := strings.Split(after, "; ")
		for _, moveSetStr := range moveSets {
			moveSet, err := parseMoveSet(moveSetStr)
			if err != nil {
				return nil, fmt.Errorf("bad move in line %v: %v", line, err)
			}

			game.MoveSets = append(game.MoveSets, moveSet)
		}

		out = append(out, game)
	}

	return out, nil
}

func solveA(games []*Game) int {
	bag := map[Color]int{RED: 12, GREEN: 13, BLUE: 14}

	out := 0
	for _, game := range games {
		good := true
		for _, set := range game.MoveSets {
			for _, move := range set {
				if bag[move.Color] < move.Num {
					logger.Infof("game %d color %v short", game.Num, move.Color)
					good = false
					break
				}
			}
			if !good {
				break
			}
		}
		if good {
			out += game.Num
		}
	}

	return out
}

func calcPower(m map[Color]int) int {
	return m[RED] * m[GREEN] * m[BLUE]
}

func solveB(games []*Game) int {
	out := 0

	for _, game := range games {
		needed := map[Color]int{}
		for _, set := range game.MoveSets {
			for _, move := range set {
				needed[move.Color] = max(needed[move.Color], move.Num)
			}
		}

		power := calcPower(needed)
		logger.Infof("game %d power %d", game.Num, power)
		out += power
	}

	return out
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
