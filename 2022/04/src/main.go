// Copyright 2022 Google LLC
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
	"strings"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Range struct {
	From, To int
}

func (r *Range) Contains(o Range) bool {
	return r.From <= o.From && r.To >= o.To
}

func (r *Range) Overlaps(o Range) bool {
	if r.From <= o.From {
		return r.To >= o.From
	} else {
		return r.From <= o.To
	}
}

func parseRange(s string) (Range, error) {
	left, right, ok := strings.Cut(s, "-")
	if !ok {
		return Range{}, fmt.Errorf("bad range cut")
	}

	parseNum := func(s string) (int, error) {
		num, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return 0, err
		}
		if num <= 0 {
			return 0, fmt.Errorf("num out of range")
		}
		return int(num), nil
	}

	var r Range
	var err error
	r.From, err = parseNum(left)
	if err != nil {
		return Range{}, err
	}
	r.To, err = parseNum(right)
	if err != nil {
		return Range{}, err
	}

	return r, nil
}

func readInput(path string) ([][2]Range, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	ranges := [][2]Range{}
	for _, line := range lines {
		partOne, partTwo, ok := strings.Cut(line, ",")
		if !ok {
			return nil, fmt.Errorf("bad cut on line: %v", line)
		}

		one, err := parseRange(partOne)
		if err != nil {
			return nil, fmt.Errorf("bad range 1 on line: %v: %v", line, err)
		}
		two, err := parseRange(partTwo)
		if err != nil {
			return nil, fmt.Errorf("bad range 1 on line: %v: %v", line, err)
		}

		ranges = append(ranges, [2]Range{one, two})
	}

	return ranges, nil
}

func solveA(ranges [][2]Range) int {
	num := 0
	for _, r := range ranges {
		if r[0].Contains(r[1]) || r[1].Contains(r[0]) {
			num++
		}
	}

	return num
}

func solveB(ranges [][2]Range) int {
	num := 0
	for _, r := range ranges {
		if r[0].Overlaps(r[1]) || r[1].Overlaps(r[0]) {
			logger.LogF("overlap: %v %v", r[0], r[1])
			num++
		}
	}

	return num
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

	fmt.Println("A", solveA(lines))
	fmt.Println("B", solveB(lines))
}
