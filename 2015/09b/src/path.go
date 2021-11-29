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
	"strconv"
	"strings"
)

var (
	edgePattern = regexp.MustCompile(`^([^ ]+) to ([^ ]+) = ([0-9]+)$`)
)

type EdgeDest struct {
	Dest string
	Dist int
}

type Graph struct {
	Nodes map[string]bool
	Edges map[string][]EdgeDest
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: map[string]bool{},
		Edges: map[string][]EdgeDest{},
	}
}

func (g *Graph) addOneEdge(src, dest string, dist int) {
	if _, found := g.Edges[src]; !found {
		g.Edges[src] = []EdgeDest{}
	}
	g.Edges[src] = append(g.Edges[src], EdgeDest{dest, dist})
	g.Nodes[src] = true
}

func (g *Graph) AddEdge(src, dest string, dist int) {
	g.addOneEdge(src, dest, dist)
	g.addOneEdge(dest, src, dist)
}

func (g *Graph) doWalk(src string, cb func(string, map[string]int), seen map[string]int) {
	for _, edgeDest := range g.Edges[src] {
		if _, found := seen[edgeDest.Dest]; found {
			continue
		}

		seen[edgeDest.Dest] = len(seen)
		cb(edgeDest.Dest, seen)
		g.doWalk(edgeDest.Dest, cb, seen)
		delete(seen, edgeDest.Dest)
	}
}

func (g *Graph) Walk(src string, cb func(string, map[string]int)) {
	seen := map[string]int{src: 0}
	g.doWalk(src, cb, seen)
}

type Path struct {
	Nodes []string
	Dist  int
}

func NewPath() *Path {
	return &Path{Nodes: []string{}, Dist: 0}
}

func (p *Path) ToString() string {
	out := fmt.Sprintf("%d:", p.Dist)
	out += strings.Join(p.Nodes, "->")
	return out
}

func isBest(ref, cand *Path) bool {
	if ref == nil {
		return true
	}
	return cand.Dist > ref.Dist
}

func readGraph(r io.Reader) (*Graph, error) {
	graph := NewGraph()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		matches := edgePattern.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("failed to parse %v", line)
		}

		from, to := matches[1], matches[2]

		dist, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("failed to parse distance %v in %v: %v",
				matches[3], line, err)
		}

		graph.AddEdge(from, to, dist)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading stdin: %v", err)
	}

	return graph, nil
}

func makePath(graph *Graph, nodes []string) *Path {
	path := NewPath()

	for i, node := range nodes {
		path.Nodes = append(path.Nodes, node)
		if i == 0 {
			continue
		}

		src := nodes[i-1]
		for _, edge := range graph.Edges[src] {
			if edge.Dest == node {
				path.Dist += edge.Dist
				break
			}
		}
	}

	return path
}

func shortestPath(src, dest string, graph *Graph) *Path {
	var best *Path

	fmt.Printf("-- looking for best path %v to %v\n", src, dest)

	graph.Walk(src, func(node string, seen map[string]int) {
		// Ignore the path if it doesn't end at dest or if it
		// hasn't seen all nodes.
		if node != dest || len(seen) != len(graph.Nodes) {
			//fmt.Printf("rejecting %v (len(seen) %v len(nodes) %v)\n",
			//seen, len(seen), len(graph.Nodes))
			return
		}

		pathNodes := make([]string, len(seen))
		for node, pos := range seen {
			pathNodes[pos] = node
		}

		path := makePath(graph, pathNodes)
		//fmt.Printf("   made path %v\n", path.ToString())

		if isBest(best, path) {
			//fmt.Printf("   path is new best\n")
			best = path
		}
	})

	if best == nil {
		fmt.Printf("   no path %v to %v\n", src, dest)
		return nil
	}

	fmt.Printf("-- best path %v to %v is %v\n", src, dest, best.ToString())
	return best
}

func shortestPathToAny(src string, graph *Graph) *Path {
	var best *Path

	fmt.Printf("- looking for best path %v to all\n", src)

	for dest := range graph.Nodes {
		if dest == src {
			continue
		}

		path := shortestPath(src, dest, graph)
		if path != nil && isBest(best, path) {
			best = path
		}
	}

	if best == nil {
		fmt.Printf("- no path %v to all\n", src)
		return nil
	}

	fmt.Printf("- best path from %v is %v\n", src, best.ToString())
	return best
}

func main() {
	graph, err := readGraph(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("read graph %v\n", *graph)

	paths := map[string]*Path{}

	for node := range graph.Edges {
		if _, found := paths[node]; found {
			continue
		}

		path := shortestPathToAny(node, graph)
		if path == nil {
			continue
		}
		paths[node] = path
	}

	var best *Path
	for _, path := range paths {
		if isBest(best, path) {
			best = path
		}
	}

	fmt.Printf("best: %v\n", best.ToString())
}
