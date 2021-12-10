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
	"container/list"
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

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	return lines, err
}

func openerFor(r rune) rune {
	switch r {
	case ')':
		return '('
	case ']':
		return '['
	case '}':
		return '{'
	case '>':
		return '<'
	default:
		panic("bad closer")
	}
}

func isCorrupt(line string) (rune, bool) {
	stack := list.New()

	for _, r := range line {
		switch r {
		case '(':
			fallthrough
		case '[':
			fallthrough
		case '{':
			fallthrough
		case '<':
			stack.PushBack(r)

		case ')':
			fallthrough
		case ']':
			fallthrough
		case '}':
			fallthrough
		case '>':
			back := stack.Back()
			if back == nil {
				return '_', true
			}

			opener := stack.Remove(back).(rune)
			if opener != openerFor(r) {
				return r, true
			}

		default:
			panic("bad char")
		}
	}

	return '0', false
}

func solveA(lines []string) {
	score := 0
	for _, line := range lines {
		c, corrupt := isCorrupt(line)
		if corrupt {
			logger.LogF("corrupt: %v at %v", line, string(c))

			switch c {
			case ')':
				score += 3
			case ']':
				score += 57
			case '}':
				score += 1197
			case '>':
				score += 25137
			default:
				panic(fmt.Sprintf("bad corrupt: %v", c))
			}
		}
	}

	fmt.Println("A", score)
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
}
