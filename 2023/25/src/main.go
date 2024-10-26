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
	"strings"

	"github.com/simmonmt/aoc/2023/common/collections"
	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Graph struct {
	g map[string][]string
}

func NewGraph() *Graph {
	return &Graph{
		g: map[string][]string{},
	}
}

func (g *Graph) Insert(a, b string) {
	g.g[a] = append(g.g[a], b)
}

func (g *Graph) Neighbors(n string) []string {
	return g.g[n]
}

func parseInput(lines []string) (*Graph, error) {
	g := NewGraph()

	for _, line := range lines {
		from, toList, _ := strings.Cut(line, ": ")
		tos := strings.Split(toList, " ")

		for _, to := range tos {
			g.Insert(from, to)
			g.Insert(to, from)
		}
	}

	return g, nil
}

func shortestPath(start string, ends map[string]bool, skipEdges map[string]bool, graph *Graph) []string {
	visited := map[string]bool{}

	queue := collections.NewPriorityQueue[string](collections.LessThan)
	queue.Insert(start, 0)

	distances := map[string]int{}
	distances[start] = 0

	froms := map[string]string{}

	for !queue.IsEmpty() {
		cur, _ := queue.Next()
		curDist := distances[cur]

		for _, neighbor := range graph.Neighbors(cur) {
			if _, found := visited[neighbor]; found {
				continue
			}

			if _, found := skipEdges[cur+":"+neighbor]; found {
				continue
			}

			throughCurDist := curDist + 1
			neighborDist, found := distances[neighbor]
			if !found || throughCurDist < neighborDist {
				distances[neighbor] = throughCurDist
				queue.Insert(neighbor, throughCurDist)
				froms[neighbor] = cur
			}

		}

		if _, found := ends[cur]; found {
			revPath := []string{}
			for id := cur; id != start; id = froms[id] {
				revPath = append(revPath, id)
			}
			return revPath
		}

		visited[cur] = true
	}

	return nil // no path found
}

// We start by picking the node with the highest number of edges because a) it
// can't be removed from the graph with 3 cuts and b) it seems like that node
// should be the easiest to find paths to, which is important when the size of
// the component is small. This node is the initial member of the maximal subgraph.
//
// We then iterate through its neighbors, trying to find ones that have at least
// four unique paths to any node in the subgraph. A node with at least four
// unique paths can't be severed from the subgraph. The iteration ends when
// we've run out of nodes that a) we haven't visited b) aren't part of the
// subgraph and c) don't have at least four connections.
//
// The end resultis lots and lots of (next-)shortest-path searches, each
// starting from scratch. There's probably a more elegant way to do it, but this
// one works quickly enough given the input.
func solveA(g *Graph) int {
	start := ""
	startCard := -1
	for node, neighbors := range g.g {
		if l := len(neighbors); l > 4 && l > startCard {
			start = node
			startCard = l
		}
	}

	ends := map[string]bool{start: true}

	todo := []string{}
	for _, neighbors := range g.g[start] {
		todo = append(todo, neighbors)
	}

	visited := map[string]bool{}

	for len(todo) > 0 {
		cur := todo[0]
		todo = todo[1:]

		if _, found := visited[cur]; found {
			continue
		}
		visited[cur] = true

		skipEdges := map[string]bool{}
		numPaths := 0
		for limit := 0; numPaths < 4 && limit < 1000; limit++ {
			path := shortestPath(cur, ends, skipEdges, g)
			if len(path) == 0 {
				break
			}
			numPaths += 1

			path = append(path, cur)
			for i := len(path) - 1; i > 0; i-- {
				a, b := path[i], path[i-1]
				skipEdges[a+":"+b] = true
				skipEdges[b+":"+a] = true
			}
		}

		if numPaths >= 4 {
			ends[cur] = true
			for _, n := range g.g[cur] {
				if _, found := ends[n]; !found {
					todo = append(todo, n)
				}
			}
		}

	}

	return len(ends) * (len(g.g) - len(ends))
}

func solveB(g *Graph) int {
	return -1
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

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
