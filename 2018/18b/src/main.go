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
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"lib"
	"logger"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	numSteps = flag.Int("num_steps", 10, "num steps")
)

func readInput() ([]string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func runLoop(board *lib.Board, results map[int]int, start, end int) {
	for t := start; t <= end; t++ {
		board.Step()

		if *verbose {
			logger.LogF("\nAfter %d minute(s)", t)
			board.Dump()
		}

		numWoods, numLumber := board.Score()
		result := numWoods * numLumber

		fmt.Printf("%d: %d woods, %d lumber = %d", t, numWoods, numLumber, result)

		if last, found := results[result]; found {
			fmt.Printf(" (delta=%d)\n", t-last)
		} else {
			fmt.Println()
		}
		results[result] = t
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	lines, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	board := lib.NewBoardFromString(lines)

	if *verbose {
		logger.LogF("Initial state:")
		board.Dump()
	}

	results := map[int]int{}

	runLoop(board, results, 1, 2000)

	//goal := 10000
	goal := 1000000000
	t := 2000 + ((goal-2000)/28)*28 + 1
	runLoop(board, results, t, goal)
}
