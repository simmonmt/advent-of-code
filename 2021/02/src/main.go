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
	"strings"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/intmath"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]pos.P2, error) {
	out := []pos.P2{}

	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		parts := strings.Split(line, " ")
		dir, numStr := parts[0], parts[1]

		num := intmath.AtoiOrDie(numStr)

		var p pos.P2
		switch dir {
		case "forward":
			p.X = num
		case "up":
			p.Y = num
		case "down":
			p.Y = -num
		default:
			panic(fmt.Sprint("bad dir", dir))
		}

		out = append(out, p)
	}
	return out, nil
}

func solveA(ps []pos.P2) {
	var cur pos.P2

	for _, p := range ps {
		cur.Add(p)
	}

	fmt.Println("A", cur.X*(-cur.Y))
}

func solveB(ps []pos.P2) {
	var cur pos.P2
	var aim int

	for _, p := range ps {
		if p.Y != 0 {
			aim += -p.Y
		}
		if p.X != 0 {
			cur.X += p.X
			cur.Y += -aim * p.X
		}
	}

	fmt.Println("B", cur.X*(-cur.Y))
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	input, err := readInput(*input)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	solveA(input)
	solveB(input)
}
