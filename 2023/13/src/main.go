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
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) ([][]string, error) {
	return filereader.BlankSeparatedGroupsFromLines(lines)
}

func btoi(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

func isOneBit(a, b int64) (isOneBit bool, which int) {
	d := a ^ b
	//logger.Infof("isOneBit %x %x d=%x", a, b, d)

	bit := 0
	for d > 0 {
		if d&1 != 0 {
			d = d >> 1
			if d > 0 {
				//logger.Infof("isOneBit multibit %x", d)
				return false, -1
			}
			//logger.Infof("isOneBit true %d", bit)
			return true, bit
		}
		d = d >> 1
		bit++
	}

	//logger.Infof("isOneBit no bit")
	return false, -1
}

func isSplit(a []int64, left int) bool {
	for i := 0; i <= left; i++ {
		l, r := left-i, left+1+i
		if l < 0 || r >= len(a) {
			break
		}

		if a[l] != a[r] {
			return false
		}
	}

	return true
}

func findSplit(a []int64) int {
	for i := range a {
		if i > 0 && a[i] == a[i-1] && isSplit(a, i-1) {
			return i
		}
	}
	return -1
}

func solveAPuzzle(lines []string) (int, int) {
	h := make([]int64, len(lines))
	v := make([]int64, len(lines[0]))

	for i, line := range lines {
		for j, c := range line {
			h[i] = (h[i] << 1) | btoi(c == '#')
			v[j] = (v[j] << 1) | btoi(c == '#')
		}
	}

	return findSplit(h), findSplit(v)
}

func solveA(input [][]string) int {
	sum := 0
	for i, group := range input {
		hSplit, vSplit := solveAPuzzle(group)
		logger.Infof("group %d h %v v %v", i, hSplit, vSplit)

		if hSplit != -1 {
			sum += 100 * hSplit
		} else {
			sum += vSplit
		}
	}
	return sum
}

func isSmudgeSplit(a []int64, left int) (locs [2]pos.P2, found bool) {
	logger.Infof("isSmudgeSplit left %v", left)
	found = false // true if one bit flip found
	for i := 0; i <= left; i++ {
		l, r := left-i, left+1+i
		if l < 0 || r >= len(a) {
			break
		}

		{
			obFound, which := isOneBit(a[l], a[r])
			logger.Infof("l %x r %x - onebit found %v which %v", a[l], a[r], obFound, which)
		}

		if a[l] == a[r] {
			continue
		}
		if found {
			logger.Infof("second flip")
			return locs, false // only one flip; found another
		}

		var which int
		if found, which = isOneBit(a[l], a[r]); found {
			locs[0] = pos.P2{X: which, Y: l}
			locs[1] = pos.P2{X: which, Y: r}
			continue
		}

		return locs, false
	}

	logger.Infof("issmudgesplit bottom %v %v", locs, found)
	return locs, found
}

func findSmudgeSplit(a []int64) (split int, locs [2]pos.P2, found bool) {
	logger.Infof("findsmudgesplit %x", a[0])
	for i := range a {
		if i == 0 {
			continue
		}

		found, _ := isOneBit(a[i], a[i-1])
		if a[i] == a[i-1] || found {
			if locs, found := isSmudgeSplit(a, i-1); found {
				logger.Infof("findsmudgesplit; found at %v, %v", i-1, locs)
				return i - 1, locs, true
			}
		}
	}

	logger.Infof("findsmudgesplit; none found")
	return -1, [2]pos.P2{}, false
}

func verifyNoSplitWithSmudge(a []int64, toFlip [2]pos.P2) {
}

func solveBPuzzle(lines []string) (int, int) {
	logger.Infof("solve b")

	h := make([]int64, len(lines))
	v := make([]int64, len(lines[0]))

	for i, line := range lines {
		for j, c := range line {
			h[i] = (h[i] << 1) | btoi(c == '#')
			v[j] = (v[j] << 1) | btoi(c == '#')
		}
	}

	out := []string{}
	for _, n := range h {
		out = append(out, fmt.Sprintf("%08x ", n))
	}
	logger.Infof("h len=%d %v", len(out), strings.Join(out, " "))

	out = []string{}
	for _, n := range v {
		out = append(out, fmt.Sprintf("%08x ", n))
	}
	logger.Infof("v len=%d %v", len(out), strings.Join(out, " "))

	origHSplit, origVSplit := findSplit(h), findSplit(v)

	hSplit, _, hFound := findSmudgeSplit(h)
	vSplit, _, vFound := findSmudgeSplit(v)

	if hFound && vFound {
		panic("both smudge found")
	}

	if hFound {
		if hSplit+1 == origHSplit {
			panic("same h")
		}
		return hSplit + 1, -1
	}
	if vFound {
		if vSplit+1 == origVSplit {
			panic("same v")
		}
		return -1, vSplit + 1
	}

	panic("neither found")
}

// tried 37136
func solveB(input [][]string) int {
	sum := 0
	for i, group := range input {
		hSplit, vSplit := solveBPuzzle(group)
		logger.Infof("group %d h %v v %v", i, hSplit, vSplit)

		if hSplit != -1 {
			sum += 100 * hSplit
		} else {
			sum += vSplit
		}
	}
	return sum
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
