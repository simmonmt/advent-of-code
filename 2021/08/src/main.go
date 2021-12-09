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
	"strings"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Line struct {
	Patterns []string
	Outputs  []string
}

func readInput(path string) ([]Line, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	out := []Line{}
	for i, line := range lines {
		parts := strings.SplitN(line, " | ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%d: bad input", i)
		}

		out = append(out, Line{
			Patterns: strings.Split(parts[0], " "),
			Outputs:  strings.Split(parts[1], " "),
		})
	}

	return out, err
}

func solveA(lines []Line) {
	sum := 0
	for _, line := range lines {
		for _, out := range line.Outputs {
			switch len(out) {
			case 2: // digit 1
				fallthrough
			case 3: // digit 7
				fallthrough
			case 4: // digit 4
				fallthrough
			case 7: // digit 8
				sum++
			}
		}
	}
	fmt.Println("A", sum)
}

func solveOne(line *Line) int {
	return 0
}

func solveB(lines []Line) {
	sum := 0
	for _, line := range lines {
		sum += solveOne(&line)
	}

	fmt.Println("B", sum)
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

	solveA(lines)
	solveB(lines)
}
