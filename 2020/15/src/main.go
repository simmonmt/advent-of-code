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

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	input    = flag.String("input", "", "input file")
	numTurns = flag.Int("num_turns", 2020, "number of turns")
)

func solve(startingNums []int) int {
	lastByNum := map[int][]int{}
	lastSpoken := 0

	for i := 0; ; i++ {
		turn := i + 1
		num := 0

		if i < len(startingNums) {
			num = startingNums[i]
		} else {
			times := lastByNum[lastSpoken]
			if len(times) == 1 {
				// that was the first time
				num = 0
			} else {
				num = times[0] - times[1]
			}
		}

		logger.LogF("Turn %d: %v", turn, num)

		if turn == *numTurns {
			return num
		}

		lastSpoken = num
		if _, found := lastByNum[num]; !found {
			lastByNum[num] = []int{i}
		} else {
			lastByNum[num] = []int{i, lastByNum[num][0]}
		}

		//logger.LogF("%d: %v last now %v", turn, num, lastByNum[num])
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	nums, err := filereader.OneRowOfNumbers(*input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Answer: %v\n", solve(nums))
}
