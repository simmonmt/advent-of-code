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
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	pattern = regexp.MustCompile(`^([0-9]+) <-> (.*)$`)
)

type Graph struct {
	Edges map[int]map[int]bool
}

func NewGraph() *Graph {
	return &Graph{
		Edges: map[int]map[int]bool{},
	}
}

func (g *Graph) Add(src, dest int) {
	if _, found := g.Edges[src]; !found {
		g.Edges[src] = map[int]bool{}
	}
	g.Edges[src][dest] = true
}

func readGraph(in io.Reader) (*Graph, error) {
	graph := NewGraph()
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		matches := pattern.FindStringSubmatch(line)
		srcStr, destStrs := matches[1], strings.Split(matches[2], ", ")

		src, err := strconv.Atoi(srcStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse src %v in %v", srcStr, line)
		}

		for _, destStr := range destStrs {
			dest, err := strconv.Atoi(destStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse dest %v in %v", destStr, line)
			}

			graph.Add(src, dest)
			graph.Add(dest, src)
		}
	}

	return graph, nil
}

func findConnected(graph *Graph, start int, foundNodes map[int]bool) {
	var candidates map[int]bool
	var found bool
	candidates, found = graph.Edges[start]
	if !found {
		panic("can't find start")
	}

	for candidate, _ := range candidates {
		if _, found := foundNodes[candidate]; found {
			continue
		}

		foundNodes[candidate] = true
		findConnected(graph, candidate, foundNodes)
	}
}

func main() {
	// if len(os.Args) != 2 {
	// 	log.Fatalf("Usage: %v start", os.Args[0])
	// }
	// start, err := strconv.Atoi(os.Args[1])
	// if err != nil {
	// 	log.Fatalf("failed to parse start %v: %v", os.Args[1], err)
	// }

	graph, err := readGraph(os.Stdin)
	if err != nil {
		log.Fatalf("failed to build graph: %v\n", err)
	}

	fmt.Printf("graph edges %v\n", graph.Edges)

	allNodes := []int{}
	for node, _ := range graph.Edges {
		allNodes = append(allNodes, node)
	}
	sort.Ints(allNodes)

	sets := map[string]bool{}
	for _, start := range allNodes {
		found := map[int]bool{}
		findConnected(graph, start, found)

		foundNodes := []string{}
		for node, _ := range found {
			foundNodes = append(foundNodes, strconv.Itoa(node))
		}
		sort.Strings(foundNodes)
		//fmt.Printf("found: %v nodes: %v\n", len(foundNodes), foundNodes)

		sets[strings.Join(foundNodes, " ")] = true
		fmt.Printf("sets: %d\n", len(sets))
	}
}
