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
	"log"
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	digits = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
)

func parseInput(lines []string) ([]string, error) {
	return lines, nil
}

func solveALine(input string) int {
	first, last := -1, -1

	for _, c := range input {
		if c >= '0' && c <= '9' {
			d := int(c - '0')
			if first == -1 {
				first = d
			}
			last = d
		}
	}

	return first*10 + last
}

func solveA(input []string) int {
	total := 0
	for _, line := range input {
		num := solveALine(line)
		logger.Infof("input %v output %v", line, num)
		total += num
	}

	return total
}

func solveBLine(input string) int {
	first, last := -1, -1

	for i, c := range input {
		d := -1
		if c >= '0' && c <= '9' {
			d = int(c - '0')
		} else {
			for j, name := range digits {
				if strings.HasPrefix(string(input[i:]), name) {
					d = j
					break
				}
			}
		}

		if d == -1 {
			continue
		}

		if first == -1 {
			first = d
		}
		last = d
	}

	return first*10 + last
}

func solveB(input []string) int {
	total := 0
	for _, line := range input {
		num := solveBLine(line)
		logger.Infof("input %v output %v", line, num)
		total += num
	}

	return total
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
