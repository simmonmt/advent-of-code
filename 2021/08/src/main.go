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
	"sort"
	"strings"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Line struct {
	Patterns []string
	Outputs  []string
}

func readInput(path string) ([]Line, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	out := []Line{}
	for i, line := range lines {
		parts := strings.SplitN(line, " | ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%d: bad input", i)
		}

		out = append(out, Line{
			Patterns: strings.Split(parts[0], " "),
			Outputs:  strings.Split(parts[1], " "),
		})
	}

	return out, err
}

func solveA(lines []Line) {
	sum := 0
	for _, line := range lines {
		for _, out := range line.Outputs {
			switch len(out) {
			case 2: // digit 1
				fallthrough
			case 3: // digit 7
				fallthrough
			case 4: // digit 4
				fallthrough
			case 7: // digit 8
				sum++
			}
		}
	}
	fmt.Println("A", sum)
}

// Given the constraints in segmap, return all possible mappings
// between signals and segments. The constraints map containts one key
// for each possible signal. The value is a map containing all
// possible segment names for that signal.
//
// Return value: An array of mappings between signals and segments.
func makeSegMapCombinations(constraints map[string]map[string]bool) []map[string]string {
	// A map from signal names to a list of all possible
	// segments. We can't use the constraints map because its
	// value can't be indexed numerically.
	segLists := map[string][]string{}
	for c, tos := range constraints {
		out := []string{}
		for to := range tos {
			out = append(out, to)
		}
		segLists[c] = out
	}

	// lims and nums are used to build a counter where each digit
	// has a different base. nums[i] can have values
	// [0,lims[i]).
	lims := make([]int, len(segLists))
	for i, r := range "abcdefg" {
		lims[i] = len(segLists[string(r)])
	}

	nums := make([]int, len(constraints))

	// Increment the counter stored in nums. Returns true when the
	// counter overflows.
	inc := func() bool {
		for i := 0; i < len(lims); i++ {
			nums[i]++
			if nums[i] < lims[i] {
				return false
			}
			if i == len(lims)-1 {
				return true
			}
			nums[i] = 0
		}
		panic("unreachable")
	}

	// Attempt to create a mapping betwen signals and segments
	// given the current values in nums. nums[i] is used to select
	// the candidate segment name (from segLists) for the i'th
	// signal. Returns the mapping.
	//
	// NOTE: The mapping described by nums may not be legal (it
	// may map multiple signals to the same segment). If that
	// happens, nil is returned.
	fill := func() map[string]string {
		cand := map[string]string{}
		used := map[string]bool{}
		for i, num := range nums {
			from := string("abcdefg"[i])
			to := segLists[from][num]
			if _, found := used[to]; found {
				return nil
			}
			cand[from] = to
			used[to] = true
		}
		return cand
	}

	// Holds all legal mappings.
	out := []map[string]string{}
	for {
		cand := fill()
		if cand != nil {
			out = append(out, cand)
		}

		if inc() {
			break
		}
	}
	return out
}

var (
	// Maps combinations of segment names to the digits they
	// indicate. The segment names must be ordered.
	digits = map[string]int{
		"ABCEFG":  0,
		"CF":      1,
		"ACDEG":   2,
		"ACDFG":   3,
		"BCDF":    4,
		"ABDFG":   5,
		"ABDEFG":  6,
		"ACF":     7,
		"ABCDEFG": 8,
		"ABCDFG":  9,
	}
)

func decodeDigit(str string, comb map[string]string) (int, bool) {
	translated := make([]string, len(str))
	for i, r := range str {
		translated[i] = comb[string(r)]
	}
	sort.Strings(translated)

	digit, found := digits[strings.Join(translated, "")]
	return digit, found
}

func checkComb(pats []string, comb map[string]string) bool {
	mappedDigits := map[int]bool{}
	for _, pat := range pats {
		digit, found := decodeDigit(pat, comb)
		if !found {
			return false
		}

		if _, found := mappedDigits[digit]; found {
			panic("double find")
		}

		mappedDigits[digit] = true
	}
	return true
}

func decodeOutputs(outputs []string, comb map[string]string) int {
	result := 0
	for _, output := range outputs {
		digit, found := decodeDigit(output, comb)
		if !found {
			panic("bad decode")
		}

		result = result*10 + digit
	}
	return result
}

func solveOne(line *Line) int {
	one, four, seven := "", "", ""
	for _, pat := range line.Patterns {
		switch len(pat) {
		case 2:
			one = pat
		case 3:
			seven = pat
		case 4:
			four = pat
		}
	}

	segMap := map[string]map[string]bool{}

	oneChars := map[rune]bool{}
	for _, r := range one {
		oneChars[r] = true
		segMap[string(r)] = map[string]bool{"C": true, "F": true}
	}

	for _, r := range seven {
		if _, found := oneChars[r]; !found {
			segMap[string(r)] = map[string]bool{"A": true}
		}
	}
	for _, r := range four {
		if _, found := oneChars[r]; !found {
			segMap[string(r)] = map[string]bool{
				"B": true,
				"D": true,
			}
		}
	}

	for _, r := range "abcdefg" {
		c := string(r)
		if segMap[c] == nil {
			segMap[c] = map[string]bool{
				"B": true,
				"D": true,
				"E": true,
				"G": true,
			}
		}
	}

	combs := makeSegMapCombinations(segMap)
	if false && logger.Enabled() {
		for _, comb := range combs {
			for _, r := range "abcdefg" {
				c := string(r)
				fmt.Printf("%v=%v ", c, comb[c])
			}
			fmt.Println()
		}
	}

	for _, comb := range combs {
		if checkComb(line.Patterns, comb) {
			return decodeOutputs(line.Outputs, comb)
		}
	}

	return -1
}

func solveB(lines []Line) {
	sum := 0
	for _, line := range lines {
		val := solveOne(&line)
		logger.LogF("mapping for %v: %v", line, val)
		sum += val
	}

	fmt.Println("B", sum)
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
