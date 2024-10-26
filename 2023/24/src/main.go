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
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type StoneSpec struct {
	P, V pos.P3
}

func parseInput(lines []string) ([]StoneSpec, error) {
	out := []StoneSpec{}
	for i, line := range lines {
		parts := strings.Fields(line)
		for j := range parts {
			parts[j] = strings.TrimRight(parts[j], ",")
		}

		nums := []int{}
		for j, part := range parts {
			if j != 3 {
				num, err := strconv.Atoi(part)
				if err != nil {
					return nil, fmt.Errorf(`%d: num %d: "%s": %v`,
						i+1, j+1, part, err)
				}
				nums = append(nums, num)
			}
		}

		p := pos.P3{X: nums[0], Y: nums[1], Z: nums[2]}
		v := pos.P3{X: nums[3], Y: nums[4], Z: nums[5]}

		out = append(out, StoneSpec{p, v})
	}
	return out, nil
}

type Stone2D struct {
	Spec StoneSpec
	M, B float64
}

func Find2DSlope(spec StoneSpec) (m, b float64) {
	// m = spec.V.Y / spec.V.X
	// spec.P.Y = m * spec.P.X + b
	// spec.P.Y - m * spec.P.X = b

	m = float64(spec.V.Y) / float64(spec.V.X)
	b = float64(spec.P.Y) - m*float64(spec.P.X)
	return
}

func Intersects(a, b *Stone2D, lo, hi float64) bool {
	if a.M == b.M {
		return a.B == b.B
	}

	// y=m1x+b1 y=m2x+b2
	// m1x+b1=m2x+b2
	// m1x-m2x=+b2-b1
	// x=(b2-b1)/(m1-m2)

	x := (b.B - a.B) / (a.M - b.M)
	y := a.M*x + a.B

	at := (x - float64(a.Spec.P.X)) / float64(a.Spec.V.X)
	bt := (x - float64(b.Spec.P.X)) / float64(b.Spec.V.X)

	intersects := x >= lo && x <= hi && y >= lo && y <= hi
	logger.Infof("%v and %v %v at %v bt %v pos %v,%v", *a, *b, intersects, at, bt, x, y)

	return intersects && at >= 0 && bt >= 0
}

func solveA(specs []StoneSpec, lo, hi float64) int {
	stones := make([]Stone2D, len(specs))
	for i := range specs {
		stone := &stones[i]
		stone.Spec = specs[i]
		stone.M, stone.B = Find2DSlope(specs[i])
	}

	num := 0
	for i := 0; i < len(stones); i++ {
		a := &stones[i]
		for j := i + 1; j < len(stones); j++ {
			b := &stones[j]
			if Intersects(a, b, lo, hi) {
				num++
			}
		}
	}
	return num
}

func solveB(stones []StoneSpec) int {
	return -1
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

	fmt.Println("A", solveA(input, float64(200000000000000), float64(400000000000000)))
	fmt.Println("B", solveB(input))
}
