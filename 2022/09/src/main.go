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
	"strconv"
	"strings"

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

type Move struct {
	Dir dir.Dir
	Num int
}

func readInput(path string) ([]Move, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	dirMap := map[string]dir.Dir{
		"R": dir.DIR_EAST,
		"L": dir.DIR_WEST,
		"U": dir.DIR_NORTH,
		"D": dir.DIR_SOUTH,
	}

	moves := []Move{}
	for _, line := range lines {
		ds, ns, ok := strings.Cut(line, " ")
		if !ok {
			return nil, fmt.Errorf("bad cut: %v", line)
		}

		d, found := dirMap[ds]
		if !found {
			return nil, fmt.Errorf("bad dir: %v", line)
		}

		n, err := strconv.Atoi(ns)
		if err != nil {
			return nil, fmt.Errorf("bad num: %v: %v", err, line)
		}

		moves = append(moves, Move{d, n})
	}

	return moves, nil
}

func moveTail(head pos.P2, tail pos.P2) pos.P2 {
	diff := pos.P2{head.X - tail.X, head.Y - tail.Y}
	if mtsmath.Abs(diff.X) < 2 && mtsmath.Abs(diff.Y) < 2 {
		return tail
	}

	if diff.X == 0 {
		tail.Y += diff.Y / mtsmath.Abs(diff.Y)
	} else if diff.Y == 0 {
		tail.X += diff.X / mtsmath.Abs(diff.X)
	} else {
		tail.Y += diff.Y / mtsmath.Abs(diff.Y)
		tail.X += diff.X / mtsmath.Abs(diff.X)
	}
	return tail
}

func solveA(moves []Move) int {
	head := pos.P2{0, 0}
	tail := head
	seen := map[pos.P2]bool{head: true}

	for _, move := range moves {
		for i := 1; i <= move.Num; i++ {
			head = move.Dir.From(head)
			tail = moveTail(head, tail)
			seen[tail] = true
		}
	}

	return len(seen)
}

func solveB(moves []Move) int {
	elems := [10]pos.P2{}
	seen := map[pos.P2]bool{elems[0]: true}

	for _, move := range moves {
		for i := 1; i <= move.Num; i++ {
			elems[0] = move.Dir.From(elems[0])
			for j := 1; j < 10; j++ {
				elems[j] = moveTail(elems[j-1], elems[j])
			}
			seen[elems[9]] = true
		}
	}

	return len(seen)
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

	fmt.Println("A", solveA(moves))
	fmt.Println("B", solveB(moves))
}
