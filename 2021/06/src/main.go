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
	numDays = flag.Int("num_days", 10, "number of days")
)

func readInput(path string) ([]int, error) {
	nums, err := filereader.OneRowOfNumbers(*input)
	if err != nil {
		return nil, err
	}

	return nums, nil
}

func solve(seed []int, numDays int) {
	fish := map[int]int{}
	for _, f := range seed {
		fish[f]++
	}

	logger.LogF("seed: %v", fish)

	for i := 1; i <= numDays; i++ {
		out := map[int]int{}
		for j := 8; j >= 0; j-- {
			curNum := fish[j]
			addFish := false
			newIdx := j - 1
			if newIdx == -1 {
				addFish = true
				newIdx = 6
			}

			out[newIdx] += curNum
			if curNum > 0 && addFish {
				out[8] += curNum
			}
		}
		fish = out

		logger.LogF("After %2d days: %v", i, fish)
	}

	sum := 0
	for _, num := range fish {
		sum += num
	}

	fmt.Println(sum)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	seed, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	solve(seed, *numDays)
}
