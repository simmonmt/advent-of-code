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
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2021/04/src/board"
	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseLine(line string) ([]int, error) {
	vals := []int{}
	for _, str := range strings.Split(line, " ") {
		if str == "" {
			continue
		}
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return vals, nil
}

func readInput(path string) ([]int, []*board.Board, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, nil, err
	}

	moves := []int{}
	for _, str := range strings.Split(lines[0], ",") {
		if str == "" {
			continue
		}
		move, err := strconv.Atoi(str)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to parse move %v: %v", str, err)
		}

		moves = append(moves, move)
	}

	if (len(lines)-2)%6 != 5 {
		return nil, nil, fmt.Errorf("input doesn't divide into boards")
	}

	boards := []*board.Board{}
	for i := 2; i < len(lines); i += 6 {
		start := [5][5]int{}
		for j := 0; j < 5; j++ {
			lineNum := i + j
			line := lines[lineNum]
			vals, err := parseLine(line)
			if err != nil || len(vals) != 5 {
				return nil, nil, fmt.Errorf(
					"bad line %v: %v", lineNum, err)
			}

			for k, val := range vals {
				start[j][k] = val
			}
		}

		boards = append(boards, board.New(start))
	}

	return moves, boards, nil
}

func playAllBoardsGame(moves []int, boards []*board.Board) (lastMove int, winningBoard *board.Board) {
	for _, move := range moves {
		for _, b := range boards {
			if b.Mark(move) {
				return move, b
			}
		}
	}

	return -1, nil
}

func solveA(moves []int, boards []*board.Board) {
	lastMove, winningBoard := playAllBoardsGame(moves, boards)
	if winningBoard == nil {
		panic("no win")
	}

	winningBoard.Dump()
	fmt.Println("winning move", lastMove)

	fmt.Println("A", winningBoard.Score(lastMove))
}

func playOneBoardGame(moves []int, b *board.Board) (numMoves int) {
	for i, move := range moves {
		if b.Mark(move) {
			return i + 1
		}
	}
	return -1
}

func solveB(moves []int, boards []*board.Board) {
	longestNumMoves := -1
	longestScore := -1

	for _, board := range boards {
		numMoves := playOneBoardGame(moves, board)
		if numMoves < 0 {
			panic("no win")
		}

		if numMoves > longestNumMoves {
			longestNumMoves = numMoves
			longestScore = board.Score(moves[numMoves-1])
		}
	}

	fmt.Println("B", longestScore)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	moves, boards, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("moves", moves)
	fmt.Println("#boards", len(boards))

	solveA(moves, boards)

	for _, board := range boards {
		board.Reset()
	}

	solveB(moves, boards)
}
