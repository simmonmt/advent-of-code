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
	"sort"

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

type Result int

const (
	RES_OK Result = iota
	RES_CORRUPT
	RES_INCOMPLETE
)

func validate(line string) (rune, *list.List, Result) {
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
				return r, stack, RES_CORRUPT
			}

			opener := stack.Remove(back).(rune)
			if opener != openerFor(r) {
				return r, stack, RES_CORRUPT
			}

		default:
			panic("bad char")
		}
	}

	if stack.Back() == nil {
		return '_', nil, RES_OK
	}

	return '_', stack, RES_INCOMPLETE
}

func isCorrupt(line string) (rune, bool) {
	last, _, result := validate(line)

	if result == RES_CORRUPT {
		return last, true
	}

	return '_', false
}

func solveA(lines []string) {
	score := 0
	for _, line := range lines {

		last, _, result := validate(line)
		if result == RES_CORRUPT {
			logger.LogF("corrupt: %v at %v", line, string(last))

			switch last {
			case ')':
				score += 3
			case ']':
				score += 57
			case '}':
				score += 1197
			case '>':
				score += 25137
			default:
				panic(fmt.Sprintf("bad corrupt: %v", last))
			}
		}
	}

	fmt.Println("A", score)
}

func completeSequence(stack *list.List) int {
	score := 0
	for stack.Back() != nil {
		last := stack.Remove(stack.Back()).(rune)

		inc := 0
		switch last {
		case '(':
			inc = 1
		case '[':
			inc = 2
		case '{':
			inc = 3
		case '<':
			inc = 4
		default:
			panic("bad last")
		}

		score = score*5 + inc
	}

	return score
}

func solveB(lines []string) {
	scores := []int{}
	for _, line := range lines {
		_, stack, result := validate(line)
		if result == RES_CORRUPT {
			continue
		}
		if result == RES_OK {
			panic("unexpected ok")
		}

		score := completeSequence(stack)
		logger.LogF("score for %v is %v", line, score)
		scores = append(scores, score)
	}

	sort.Ints(scores)
	logger.LogF("scores %v", scores)
	if len(scores)%2 == 0 {
		panic("even scores")
	}

	score := scores[len(scores)/2]

	fmt.Println("B", score)
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
