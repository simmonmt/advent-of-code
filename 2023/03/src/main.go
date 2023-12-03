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
	"strconv"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Finding struct {
	Loc pos.P2
	Num int
}

func parseInput(lines []string) ([]string, error) {
	return lines, nil
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func findNumber(s string, x int) (start, num int) {
	start = x
	end := x
	for start > 0 && isDigit(s[start-1]) {
		start--
	}
	for end < len(s)-1 && isDigit(s[end+1]) {
		end++
	}

	num, _ = strconv.Atoi(string(s[start : end+1]))
	return
}

func findNumbersAt(prev, cur, next string, x, y int) []Finding {
	out := []Finding{}

	saveNumber := func(s string, x, y int) {
		start, num := findNumber(s, x)
		out = append(out, Finding{Loc: pos.P2{X: start, Y: y}, Num: num})
	}

	checkOther := func(s string, x, y int) {
		if isDigit(s[x]) {
			saveNumber(s, x, y)
		} else {
			if x > 0 && isDigit(s[x-1]) {
				saveNumber(s, x-1, y)
			}
			if x < len(cur)-1 && isDigit(s[x+1]) {
				saveNumber(s, x+1, y)
			}
		}
	}

	if prev != "" {
		checkOther(prev, x, y-1)
	}
	if next != "" {
		checkOther(next, x, y+1)
	}

	if x > 0 && isDigit(cur[x-1]) {
		saveNumber(cur, x-1, y)
	}
	if x < len(cur)-1 && isDigit(cur[x+1]) {
		saveNumber(cur, x+1, y)
	}

	return out
}

func findNumbers(prev, cur, next string, curY int) []Finding {
	out := []Finding{}

	for x, r := range cur {
		if r == '.' || isDigit(byte(r)) {
			continue
		}

		out = append(out, findNumbersAt(prev, cur, next, x, curY)...)
	}

	return out
}

func walkBoard(input []string, cb func(prev, cur, next string, y int)) {
	prev, next := "", input[0]
	for i := range input {
		cur := next
		if i < len(input)-1 {
			next = input[i+1]
		} else {
			next = ""
		}

		cb(prev, cur, next, i)

		prev = cur
	}
}

func solveA(input []string) int {
	// keep prev, cur, next lines
	// scan cur. if cur has symbol, look at adjacencies in prev, next.
	// if number found, parse the whole thing, get starting coords
	// put found nums in hash so we don't find them twice

	nums := map[pos.P2]int{}

	walkBoard(input, func(prev, cur, next string, y int) {
		for _, finding := range findNumbers(prev, cur, next, y) {
			if _, found := nums[finding.Loc]; !found {
				logger.Infof("found %+v", finding)
				nums[finding.Loc] = finding.Num
			}
		}
	})

	out := 0
	for _, num := range nums {
		out += num
	}

	return out
}

func solveB(input []string) int {
	out := 0

	walkBoard(input, func(prev, cur, next string, y int) {
		for x, r := range cur {
			if r != '*' {
				continue
			}

			findings := findNumbersAt(prev, cur, next, x, y)
			if len(findings) != 2 {
				continue
			}

			logger.Infof("gear at %v with nums %d %d",
				pos.P2{X: x, Y: y}, findings[0].Num, findings[1].Num)

			out += findings[0].Num * findings[1].Num
		}
	})

	return out
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
