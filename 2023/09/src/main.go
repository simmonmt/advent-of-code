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

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/strutil"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) ([][]int, error) {
	out := [][]int{}
	for i, line := range lines {
		seq, err := strutil.ListOfNumbers(line)
		if err != nil {
			return nil, fmt.Errorf("bad seq line %v", i+1)
		}
		out = append(out, seq)
	}
	return out, nil
}

func dumpSeqs(i int, seqs [][]int) {
	for _, seq := range seqs {
		fmt.Println(i, seq)
	}
}

func makeSeqPyramid(seq []int) [][]int {
	seqs := [][]int{seq}
	curIdx := 0

	loops := 0
	allZeros := false
	for !allZeros {
		loops++
		if loops > 100 {
			panic("out of control")
		}

		allZeros = true

		cur := seqs[curIdx]
		nSeq := make([]int, len(cur)-1)
		for i := 1; i < len(cur); i++ {
			delta := cur[i] - cur[i-1]
			nSeq[i-1] = delta

			if delta != 0 {
				allZeros = false
			}
		}

		seqs = append(seqs, nSeq)
		curIdx += 1
	}

	return seqs
}

func futureSeq(seq []int) int {
	seqs := makeSeqPyramid(seq)
	next := make([]int, len(seqs))

	for i := len(seqs) - 1; i >= 0; i-- {
		diff := 0
		if i != len(seqs)-1 {
			diff = next[i+1]
		}

		next[i] = seqs[i][len(seqs[i])-1] + diff
	}

	return next[0]
}

func pastSeq(seq []int) int {
	seqs := makeSeqPyramid(seq)
	past := make([]int, len(seqs))

	for i := len(seqs) - 1; i >= 0; i-- {
		diff := 0
		if i != len(seqs)-1 {
			diff = past[i+1]
		}

		past[i] = seqs[i][0] - diff
	}

	return past[0]
}

func solveA(input [][]int) int {
	out := 0
	for _, seq := range input {
		out += futureSeq(seq)
	}
	return out
}

func solveB(input [][]int) int {
	out := 0
	for _, seq := range input {
		out += pastSeq(seq)
	}
	return out
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
