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

func solvePairs(first string, pairs map[string]bool, stepNum, numSteps int, rewriteMap map[string]rune) map[string]map[string]int {
	logger.LogF("%v: (max %v) pairs %v", stepNum, numSteps, pairs)

	if stepNum >= numSteps {
		out := map[string]map[string]int{}
		for p := range pairs {
			counts := map[string]int{}
			counts[string(p[0])]++
			counts[string(p[1])]++
			out[p] = counts
		}
		logger.LogF("%v: returning %v", stepNum, out)
		return out
	}

	expanded := map[string]string{}
	for p := range pairs {
		if mid, found := rewriteMap[p]; found {
			expanded[p] = string(p[0]) + string(mid) + string(p[1])
		} else {
			panic("no rewrite")
		}
	}

	subPairs := map[string]bool{}
	for _, t := range expanded {
		subPairs[string(t[0])+string(t[1])] = true
		subPairs[string(t[1])+string(t[2])] = true
	}

	subPairCounts := solvePairs(first, subPairs, stepNum+1, numSteps,
		rewriteMap)

	counts := map[string]map[string]int{}
	for p := range pairs {
		exp := expanded[p]
		left, right := exp[0:2], exp[1:3]

		pairCounts := map[string]int{}
		for c, n := range subPairCounts[left] {
			pairCounts[c] += n
		}
		for c, n := range subPairCounts[right] {
			pairCounts[c] += n
		}
		pairCounts[string(exp[1])]-- // it was double-counted

		logger.LogF("%v: pair %v left %v=%v right %v=%v => %v",
			stepNum, p, left, subPairCounts[left],
			right, subPairCounts[right], pairCounts)

		counts[p] = pairCounts
	}

	logger.LogF("%d: returning %v", stepNum, counts)
	return counts

}

func solveSeq(seq string, numSteps int, rewriteMap map[string]rune) map[string]int {
	pairs := map[string]bool{}
	for i := 0; i < len(seq)-1; i++ {
		pairs[seq[i:i+2]] = true
	}

	countsByPair := solvePairs(string(seq[0]), pairs, 0, numSteps, rewriteMap)

	logger.LogF("totalizing %v", seq)

	totals := map[string]int{}
	for i := 0; i < len(seq)-1; i++ {
		pair := seq[i : i+2]
		counts := countsByPair[pair]

		for c, n := range counts {
			totals[c] += n
		}

		totals[seq[i:i+1]]-- // it will be double-counted

		logger.LogF("total: pair %v counts %v => %v",
			pair, counts, totals)
	}

	totals[seq[0:1]]++ // not double-counted

	logger.LogF("total: totals %v", totals)

	return totals
}

func solve(seq string, numSteps int, rewriteMap map[string]rune) {
	counts := solveSeq(seq, numSteps, rewriteMap)

	bySize := []string{}
	for s := range counts {
		bySize = append(bySize, s)
	}
	sort.Slice(bySize, func(i, j int) bool {
		return counts[bySize[i]] < counts[bySize[j]]
	})

	mostCommon := counts[bySize[len(bySize)-1]]
	leastCommon := counts[bySize[0]]

	fmt.Println("Result", mostCommon-leastCommon)
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

	solve(seq, *numSteps, rewriteMap)
}
