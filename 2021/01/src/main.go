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

func solveA(nums []int) {
	numInc := 0
	for i := 1; i < len(nums); i++ {
		if nums[i] > nums[i-1] {
			numInc++
		}
	}

	fmt.Println("A: num increases", numInc)
}

func sumList(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func solveB(nums []int) {
	numInc := 0
	prev := -1
	for i := 2; i < len(nums); i++ {
		this := sumList(nums[i-2 : i+1])
		if prev > 0 && this > prev {
			numInc++
		}
		prev = this
	}

	fmt.Println("B: num increases", numInc)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	nums, err := filereader.Numbers(*input)
	if err != nil {
		log.Fatal(err)
	}

	solveA(nums)
	solveB(nums)
}
