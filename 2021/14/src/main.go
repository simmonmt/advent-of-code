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
	"strings"
	"unicode/utf8"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	input    = flag.String("input", "", "input file")
	numSteps = flag.Int("num_steps", 10, "number of steps")
)

func readInput(path string) (string, map[string]rune, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return "", nil, err
	}

	seq := lines[0]
	rewriteMap := map[string]rune{}

	lineNum := 2
	for _, line := range lines[2:] {
		lineNum++

		parts := strings.SplitN(line, " -> ", 2)
		if len(parts) != 2 {
			return "", nil, fmt.Errorf("%d: bad parse", lineNum)
		}

		r, _ := utf8.DecodeRuneInString(parts[1])
		rewriteMap[parts[0]] = r
	}

	return seq, rewriteMap, nil
}

func parseSeq(seq string) *list.List {
	out := list.New()
	for _, r := range seq {
		out.PushBack(r)
	}
	return out
}

func seqToString(seq *list.List) string {
	out := ""
	for e := seq.Front(); e != nil; e = e.Next() {
		out += string(e.Value.(rune))
	}
	return out
}

func expandSeq(seq *list.List, rewriteMap map[string]rune) {
	e := seq.Front()
	for e != nil {
		next := e.Next()
		if next == nil {
			break // we need a pair
		}

		pair := string(e.Value.(rune)) + string(next.Value.(rune))
		if toInsert, found := rewriteMap[pair]; found {
			seq.InsertAfter(toInsert, e)
		}

		e = next
	}
}

func solveA(seqStr string, rewriteMap map[string]rune) {
	seq := parseSeq(seqStr)

	for step := 1; step <= *numSteps; step++ {
		expandSeq(seq, rewriteMap)
		if logger.Enabled() {
			logger.LogF("after step %3d: %v", step,
				seqToString(seq))
		}
		fmt.Println("step", step)

	}

	counts := map[rune]int{}
	for e := seq.Front(); e != nil; e = e.Next() {
		counts[e.Value.(rune)]++
	}

	bySize := []rune{}
	for r := range counts {
		bySize = append(bySize, r)
	}
	sort.Slice(bySize, func(i, j int) bool {
		return counts[bySize[i]] < counts[bySize[j]]
	})

	mostCommon := counts[bySize[len(bySize)-1]]
	leastCommon := counts[bySize[0]]

	fmt.Println("A", mostCommon-leastCommon)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	seq, rewriteMap, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	solveA(seq, rewriteMap)
}
