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
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func runeToPriority(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r) - 'a' + 1
	} else if r >= 'A' && r <= 'Z' {
		return int(r) - 'A' + 27
	} else {
		panic("bad rune")
	}
}

func solveA(lines []string) int {
	sum := 0

	for _, line := range lines {
		left, right := line[0:len(line)/2], line[len(line)/2:]
		if len(left) != len(right) {
			panic("not even")
		}

		leftItems := map[rune]bool{}
		for _, r := range left {
			leftItems[r] = true
		}

		for _, r := range right {
			if _, found := leftItems[r]; !found {
				continue
			}

			sum += runeToPriority(r)
			break
		}
	}

	return sum
}

func findBadge(group []string) rune {
	// Tracks items that have been found in an unbroken string from 0. If a
	// rune is found any number of times in group[0], foundItems[r]=0. If it
	// was found in group[0] and group[1], foundItems[r]=1. If it was found
	// in group[1] but not in group[0], it won't be present in the map.
	foundItems := map[rune]int{}

	for j, sack := range group {
		for _, r := range sack {
			if n, found := foundItems[r]; !found {
				if j == 0 {
					// Every item found in group[0] gets
					// added.
					foundItems[r] = j
				}
			} else {
				// Only update the value for items that were
				// found for the previous sack (i.e. group[j-1])
				if j > 0 && n == j-1 {
					foundItems[r] = j
					if j == 2 {
						// We only expect to find one
						// item in all sacks so we can
						// quit immediately upon finding
						// it.
						return r
					}
				}
			}
		}
	}

	panic("not found")
}

func solveB(lines []string) int {
	sum := 0

	for i := 0; i < len(lines); i += 3 {
		group := lines[i : i+3]
		badge := findBadge(group)
		logger.LogF("group %d badge %s", i/3, string(badge))
		sum += runeToPriority(badge)
	}

	return sum
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(lines))
	fmt.Println("B", solveB(lines))
}
