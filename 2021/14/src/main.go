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

func solve(seq string, stepNum, numSteps int, rewriteMap map[string]rune) map[string]int {
	logger.LogF("entering step %v", stepNum)

	if stepNum == numSteps {
		out := map[string]int{}
		for _, r := range seq {
			out[string(r)]++
		}
		return out
	}

	seen := map[string]map[string]int{}

	for i := 0; i < len(seq)-1; i++ {
		a, b := string(seq[i]), string(seq[i+1])
		pair := a + b

		if _, found := seen[pair]; found {
			continue
		}
		c, found := rewriteMap[pair]
		if !found {
			continue
		}

		newSeq := a + string(c) + b
		seen[pair] = solve(newSeq, stepNum+1, numSteps, rewriteMap)
	}

	totals := map[string]int{}
	for i := 0; i < len(seq)-1; i++ {
		a, b := string(seq[i]), string(seq[i+1])
		pair := a + b

		pairCounts, found := seen[pair]
		if !found {
			totals[a]++
			totals[b]++
			continue
		}

		for s, n := range pairCounts {
			totals[s] += n
		}

		if i != 0 {
			totals[a]--
		}
	}

	return totals
}

func solvePart(seqStr string, numSteps int, rewriteMap map[string]rune) {
	counts := solve(seqStr, 0, numSteps, rewriteMap)

	bySize := []string{}
	for s := range counts {
		bySize = append(bySize, s)
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

	solvePart(seq, *numSteps, rewriteMap)
}
