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
	"regexp"
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/mtsmath"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	mapPattern = regexp.MustCompile(`^([^ ]+) = \(([^,]+), ([^,]+)\)$`)
)

type Node struct {
	Left, Right string
}

func parseInput(lines []string) (string, map[string]Node, error) {
	instructions := lines[0]
	nodeMap := map[string]Node{}

	for i, line := range lines[2:] {
		parts := mapPattern.FindStringSubmatch(line)
		if parts == nil {
			return "", nil, fmt.Errorf("bad match line %d", i+2)
		}

		nodeMap[parts[1]] = Node{parts[2], parts[3]}
	}

	return instructions, nodeMap, nil
}

func solveA(instructions string, nodeMap map[string]Node) int {
	cur := "AAA"
	steps := 0
	for {
		for _, r := range instructions {
			node := nodeMap[cur]
			if r == 'L' {
				cur = node.Left
			} else {
				cur = node.Right
			}
			steps++
		}

		if cur == "ZZZ" {
			break
		}
	}
	return steps
}

func solveBForNode(instructions string, start string, nodeMap map[string]Node) int {
	cur := start
	steps := 0
	for {
		for _, r := range instructions {
			node := nodeMap[cur]
			if r == 'L' {
				cur = node.Left
			} else {
				cur = node.Right
			}
			steps++
		}

		if strings.HasSuffix(cur, "Z") {
			break
		}
	}
	return steps
}

func solveB(instructions string, nodeMap map[string]Node) int64 {
	loops := []int64{}
	for name := range nodeMap {
		if strings.HasSuffix(name, "A") {
			loops = append(loops, int64(solveBForNode(instructions, name, nodeMap)))
		}
	}
	return mtsmath.LCM(loops...)
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

	instructions, nodeMap, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(instructions, nodeMap))
	fmt.Println("B", solveB(instructions, nodeMap))
}
