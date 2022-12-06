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

func startMarker(line string) bool {
	//	fmt.Println(line)
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			//	fmt.Println(i,line[i],j,line[j])
			if line[i] == line[j] {
				return false
			}
		}
	}
	return true
}

func solveA(lines []string) int {
	line := lines[0]

	for i := range line {
		if i < 3 {
			continue
		}
		if startMarker(line[i-3 : i+1]) {
			return i + 1
		}
	}

	return -1
}

func solveB(lines []string) int {
	return -1
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
