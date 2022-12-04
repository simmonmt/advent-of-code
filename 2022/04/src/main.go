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

	"github.com/simmonmt/aoc/2022/common/area"
	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([][2]area.Area1D, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	input := [][2]area.Area1D{}
	for _, line := range lines {
		partOne, partTwo, ok := strings.Cut(line, ",")
		if !ok {
			return nil, fmt.Errorf("bad cut on line: %v", line)
		}

		one, err := area.ParseArea1D(partOne)
		if err != nil {
			return nil, fmt.Errorf("bad interval 1 on line: %v: %v", line, err)
		}
		two, err := area.ParseArea1D(partTwo)
		if err != nil {
			return nil, fmt.Errorf("bad interval 2 on line: %v: %v", line, err)
		}

		input = append(input, [2]area.Area1D{one, two})
	}

	return input, nil
}

func solveA(input [][2]area.Area1D) int {
	num := 0
	for _, ranges := range input {
		if ranges[0].Contains(ranges[1]) || ranges[1].Contains(ranges[0]) {
			num++
		}
	}

	return num
}

func solveB(input [][2]area.Area1D) int {
	num := 0
	for _, ranges := range input {
		if ranges[0].Overlaps(ranges[1]) || ranges[1].Overlaps(ranges[0]) {
			logger.LogF("overlap: %v %v", ranges[0], ranges[1])
			num++
		}
	}

	return num
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
