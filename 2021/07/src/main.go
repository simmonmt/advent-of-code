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

func solve(nums []int, costFunc func(dist int) int) int {
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
			fuelNeeded += costFunc(dist)
		}

		if minFuel == -1 || fuelNeeded < minFuel {
			minFuel = fuelNeeded
			minFuelPos = i
		}
	}

	logger.LogF("min fuel %v (at %v)", minFuel, minFuelPos)
	return minFuel
}

func solveA(nums []int) {
	fuel := solve(nums, func(dist int) int { return dist })
	fmt.Println("A", fuel)
}

// I assume there's some clever way to do this but let's try a cache
// for now. It's a map from distance to calculated cost.
var costCache = map[int]int{
	1: 1,
	0: 0,
}

func recursiveCost(dist int) int {
	if cost, found := costCache[dist]; found {
		logger.LogF("reused %d", dist)
		return cost
	}

	cost := recursiveCost(dist-1) + dist
	logger.LogF("calculated %d", dist)
	costCache[dist] = cost
	return cost
}

func solveB(nums []int) {
	fuel := solve(nums, recursiveCost)
	fmt.Println("B", fuel)
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
