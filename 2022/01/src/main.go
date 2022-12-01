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
	"sort"
	"strconv"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([][]int, error) {
	strGroups, err := filereader.BlankSeparatedGroups(path)
	if err != nil {
		return nil, err
	}

	groups := [][]int{}
	for _, strGroup := range strGroups {
		group := []int{}
		for _, str := range strGroup {
			num, err := strconv.ParseInt(str, 0, 32)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to parse %v: %v",
					str, err)
			}
			group = append(group, int(num))
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	groups, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	loads := []int{}
	for _, group := range groups {
		load := 0
		for _, num := range group {
			load += num
		}
		loads = append(loads, load)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(loads)))

	fmt.Println("Part A", loads[0])
	fmt.Println("Part B", loads[0]+loads[1]+loads[2])
}
