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

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]int, error) {
	nums, err := filereader.OneRowOfNumbers(*input)
	if err != nil {
		return nil, err
	}

	return nums, nil
}

func solveA(nums []int) {
	min, max := -1, -1
	for _, num := range nums {
		if min == -1 || num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}

	minFuel := -1
	minFuelPos := -1
	for i := min; i <= max; i++ {
		fuelNeeded := 0

		for _, num := range nums {
			dist := num - i
			if dist < 0 {
				dist = -dist
			}
			fuelNeeded += dist
		}

		if minFuel == -1 || fuelNeeded < minFuel {
			minFuel = fuelNeeded
			minFuelPos = i
		}
	}

	fmt.Printf("A fuel %v (at %v)\n", minFuel, minFuelPos)
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
}
