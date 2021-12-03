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

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func binToInt(in string) int {
	val, err := strconv.ParseInt(in, 2, 32)
	if err != nil {
		panic("bad bin")
	}
	return int(val)
}

func countBit(lines []string, idx int) int {
	num := 0
	for _, line := range lines {
		if line[idx] == '1' {
			num++
		}
	}
	return num
}

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
		return num*2 > numLines
	})
	epsilon := calc(counts, func(num int) bool {
		return num*2 < numLines
	})

	fmt.Println("gamma", gamma)
	fmt.Println("epsilon", epsilon)
	fmt.Println("A", gamma*epsilon)
}

func filter(lines []string, idx int, shouldKeepOnes func(numOnes, numLeft int) bool) []string {
	keepVal := '0'
	numOnes := countBit(lines, idx)
	if shouldKeepOnes(numOnes, len(lines)) {
		keepVal = '1'
	}

	out := []string{}
	for _, line := range lines {
		if rune(line[idx]) == keepVal {
			out = append(out, line)
		}
	}

	//fmt.Printf("numones %v keep val %v => %v\n",
	//numOnes, string(keepVal), out)

	return out
}

func filterLoop(lines []string, shouldKeepOnes func(numOnes, numLeft int) bool) string {
	idx := 0
	for len(lines) > 1 {
		lines = filter(lines, idx, shouldKeepOnes)
		idx++
	}

	if len(lines) == 0 {
		panic("no lines")
	}
	return lines[0]
}

func solveB(lines []string) {
	oxygen := binToInt(filterLoop(lines, func(numOnes, numLeft int) bool {
		return numOnes*2 >= numLeft
	}))
	fmt.Println("oxygen", oxygen)

	co2 := binToInt(filterLoop(lines, func(numOnes, numLeft int) bool {
		return numOnes*2 < numLeft
	}))
	fmt.Println("co2", co2)

	fmt.Println("B", oxygen*co2)
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
	solveB(lines)
}
