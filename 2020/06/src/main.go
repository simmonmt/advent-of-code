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
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func solveA(groups [][]string) {
	sum := 0
	for i, group := range groups {
		answered := map[rune]bool{}
		for _, person := range group {
			for _, q := range person {
				answered[q] = true
			}
		}

		logger.LogF("group %d answered %d", i+1, len(answered))
		sum += len(answered)
	}

	fmt.Printf("A sum is %d\n", sum)
}

func solveB(groups [][]string) {
	sum := 0
	for i, group := range groups {
		answered := map[rune]int{}
		for _, person := range group {
			for _, q := range person {
				answered[q]++
			}
		}

		numAllAnswered := 0
		for _, num := range answered {
			if num == len(group) {
				numAllAnswered++
			}
		}

		logger.LogF("group %d all answered %d", i+1, numAllAnswered)
		sum += numAllAnswered
	}

	fmt.Printf("B sum is %d\n", sum)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	groups, err := filereader.BlankSeparatedGroups(*input)
	if err != nil {
		log.Fatal(err)
	}

	solveA(groups)
	solveB(groups)
}
