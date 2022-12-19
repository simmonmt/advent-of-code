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

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
	"github.com/simmonmt/aoc/2022/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) ([]pos.P3, error) {
	out := []pos.P3{}
	for _, line := range lines {
		p, err := pos.P3FromString(line)
		if err != nil {
			return nil, fmt.Errorf("bad pos %v: %v", line, err)
		}
		out = append(out, p)
	}
	return out, nil
}

func solveA(locList []pos.P3) int {
	sum := 0

	locs := map[pos.P3]bool{}
	for _, loc := range locList {
		locs[loc] = true
	}

	for loc := range locs {
		closed := 0
		for _, n := range loc.AllNeighbors(false) {
			if _, found := locs[n]; found {
				closed++
			}
		}

		logger.LogF("%v closed %v open %v", loc, closed, 6-closed)

		sum += 6 - closed
	}

	return sum
}

func solveB(locs []pos.P3) int {
	return -1
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

	input, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
