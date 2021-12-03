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

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func solveA(lines []string) {
	numLines := len(lines)

	counts := make([]int, len(lines[0]))
	for _, line := range lines {
		for i, rune := range line {
			if rune == '1' {
				counts[i]++
			}
		}
	}

	calc := func(counts []int, computeBit func(num int) bool) int {
		out := 0

		for _, count := range counts {
			if numLines%2 == 0 && count*2 == numLines {
				panic("even")
			}

			bit := 0
			if computeBit(count) {
				bit = 1
			}

			out = (out << 1) | bit
		}

		return out
	}

	gamma := calc(counts, func(num int) bool {
		fmt.Println(num)
		return num*2 > numLines
	})
	epsilon := calc(counts, func(num int) bool {
		return num*2 < numLines
	})

	fmt.Println("gamma", gamma)
	fmt.Println("epsilon", epsilon)
	fmt.Println("A", gamma*epsilon)
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

	solveA(lines)
}
