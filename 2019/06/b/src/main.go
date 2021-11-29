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
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func walkDepthFirst(orbitedBy map[string][]string, cur string, level int, callback func(cur string, children []string, level int)) {
	children := orbitedBy[cur]
	callback(cur, children, level)
	for _, child := range children {
		walkDepthFirst(orbitedBy, child, level+1, callback)
	}
}

func findParents(orbits map[string]string, node string) []string {
	out := []string{}
	for node != "COM" {
		node = orbits[node]
		out = append(out, node)
	}
	return out
}

func findParentIndex(parents []string, parent string) int {
	for i, v := range parents {
		if v == parent {
			return i
		}
	}
	panic("not found")
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

	orbitedBy := map[string][]string{}
	orbits := map[string]string{}

	for _, line := range lines {
		parts := strings.Split(line, ")")
		orbitee := parts[0]
		orbiter := parts[1]

		orbits[orbiter] = orbitee
		if _, found := orbitedBy[orbitee]; !found {
			orbitedBy[orbitee] = []string{}
		}
		orbitedBy[orbitee] = append(orbitedBy[orbitee], orbiter)
	}

	sanParents := findParents(orbits, "SAN")
	youParents := findParents(orbits, "YOU")

	pmap := map[string]bool{}
	for _, p := range sanParents {
		pmap[p] = true
	}
	var commonParent string
	for _, p := range youParents {
		if _, found := pmap[p]; found {
			commonParent = p
			break
		}
	}

	if commonParent == "" {
		panic("no common parent")
	}

	sanIdx := findParentIndex(sanParents, commonParent)
	youIdx := findParentIndex(youParents, commonParent)

	// fmt.Println(sanParents)
	// fmt.Println(youParents)
	// fmt.Println(commonParent)
	// fmt.Println(sanIdx)
	// fmt.Println(youIdx)
	fmt.Println(sanIdx + youIdx)
}
