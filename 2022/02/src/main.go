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
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([][2]string, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	moves := [][2]string{}
	for _, line := range lines {
		opp, you, found := strings.Cut(line, " ")
		if !found {
			return nil, fmt.Errorf("bad split on line: %v", line)
		}

		if opp != "A" && opp != "B" && opp != "C" {
			return nil, fmt.Errorf("bad opponent move on line: %v",
				line)
		}

		if you != "X" && you != "Y" && you != "Z" {
			return nil, fmt.Errorf("bad you move on line: %v",
				line)
		}

		moves = append(moves, [2]string{opp, you})
	}

	return moves, nil
}

type Action int

const (
	ROCK     Action = 1
	PAPER    Action = 2
	SCISSORS Action = 3
)

var (
	oppMap = map[string]Action{
		"A": ROCK,
		"B": PAPER,
		"C": SCISSORS,
	}

	youMap = map[string]Action{
		"X": ROCK,
		"Y": PAPER,
		"Z": SCISSORS,
	}
)

func playRound(opp, you Action) int {
	if opp == you {
		return 3
	}

	if you == ROCK && opp == SCISSORS {
		return 6
	} else if you == SCISSORS && opp == PAPER {
		return 6
	} else if you == PAPER && opp == ROCK {
		return 6
	} else {
		return 0
	}
}

func partA(moves [][2]string) int {
	score := 0

	for _, move := range moves {
		opp, you := oppMap[move[0]], youMap[move[1]]

		score += playRound(opp, you)
		score += int(you)
	}

	return score
}

func partB(moves [][2]string) int {
	score := 0
	needMap := map[string]int{"X": 0, "Y": 3, "Z": 6}

	for _, move := range moves {
		opp := oppMap[move[0]]
		need := needMap[move[1]]

		for _, you := range []Action{ROCK, SCISSORS, PAPER} {
			result := playRound(opp, you)
			if result == need {
				score += result
				score += int(you)
				break
			}
		}
	}

	return score
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	moves, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", partA(moves))
	fmt.Println("B", partB(moves))
}
