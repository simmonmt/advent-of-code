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
	"container/list"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/graph"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	inputPattern = regexp.MustCompile(
		`^Valve ([^ ]+) has flow rate=(\d+); tunnels? leads? to valves? ([^ ,]+(?:, [^ ,]+)*)`)
)

type InputNode struct {
	Name  string
	Rate  int
	Dests []string
}

func parseInput(lines []string) ([]*InputNode, error) {
	out := []*InputNode{}
	for i, line := range lines {
		parts := inputPattern.FindStringSubmatch(line)
		if len(parts) == 0 {
			return nil, fmt.Errorf("%d: bad match", i+1)
		}

		name := parts[1]
		rate, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("%d: bad rate: %v", i+1, err)
		}

		dests := strings.Split(parts[3], ", ")

		out = append(out, &InputNode{name, rate, dests})
	}
	return out, nil
}

type InputGraph struct {
	nodes map[graph.NodeID]*InputNode
}

func (g *InputGraph) NeighborDistance(from, to graph.NodeID) int {
	return 1
}

func (g *InputGraph) Neighbors(id graph.NodeID) []graph.NodeID {
	inputNode, found := g.nodes[id]
	if !found {
		panic("unknown")
	}

	//logger.LogF("%v: input node %+v", id, inputNode)

	out := []graph.NodeID{}
	for _, dest := range inputNode.Dests {
		out = append(out, graph.NodeID(dest))
	}

	//logger.LogF("neighbors of %v: %v", id, out)
	return out
}

type SimpleGraph struct {
	edges map[graph.NodeID]map[graph.NodeID]int
	rates map[graph.NodeID]int
}

func NewSimpleGraph(nodes []*InputNode) *SimpleGraph {
	g := &SimpleGraph{
		edges: map[graph.NodeID]map[graph.NodeID]int{},
		rates: map[graph.NodeID]int{},
	}

	for _, node := range nodes {
		if node.Rate > 0 {
			g.rates[graph.NodeID(node.Name)] = node.Rate
		}
	}

	return g
}

func (g *SimpleGraph) AddEdge(from, to graph.NodeID, cost int) {
	if _, found := g.edges[from]; !found {
		g.edges[from] = map[graph.NodeID]int{}
	}
	g.edges[from][to] = cost
}

func (g *SimpleGraph) HasEdge(from, to graph.NodeID) bool {
	sub, found := g.edges[from]
	if !found {
		return false
	}
	_, found = sub[to]
	return found
}

func (g *SimpleGraph) EdgeCost(from, to graph.NodeID) int {
	sub, found := g.edges[from]
	if !found {
		panic("bad from")
	}
	cost, found := sub[to]
	if !found {
		panic("bad to")
	}
	return cost
}

func (g *SimpleGraph) AllEdges(from graph.NodeID) map[graph.NodeID]int {
	sub, found := g.edges[from]
	if !found {
		panic("bad from")
	}
	return sub
}

func (g *SimpleGraph) Rate(id graph.NodeID) int {
	return g.rates[id]
}

func simplifyInputGraph(nodes []*InputNode) *SimpleGraph {
	allNodes := map[graph.NodeID]*InputNode{}
	usableNodes := []graph.NodeID{"AA"}
	for _, node := range nodes {
		allNodes[graph.NodeID(node.Name)] = node
		if node.Rate > 0 {
			usableNodes = append(usableNodes, graph.NodeID(node.Name))
		}
	}

	inputGraph := &InputGraph{allNodes}
	simpleGraph := NewSimpleGraph(nodes)
	for _, from := range usableNodes {
		for _, to := range usableNodes {
			if from == to || to == "AA" {
				continue
			}
			if simpleGraph.HasEdge(from, to) {
				continue
			}

			//logger.LogF("%v to %v", from, to)

			path := graph.ShortestPath(graph.NodeID(from), graph.NodeID(to), inputGraph)
			if path == nil {
				panic(fmt.Sprintf("disconnect; no path from %v to %v", from, to))
			}

			cost := len(path) + 1 // + 1 to turn on dest valve

			simpleGraph.AddEdge(from, to, cost)
			if from != "AA" {
				simpleGraph.AddEdge(to, from, cost)
			}
		}
	}

	return simpleGraph
}

type PathManager struct {
	g     *SimpleGraph
	start graph.NodeID

	visited map[graph.NodeID]bool
	path    *list.List

	maxRelease     int
	maxReleasePath []graph.NodeID
}

func NewPathManager(g *SimpleGraph, start graph.NodeID) *PathManager {
	return &PathManager{
		g:       g,
		start:   start,
		visited: map[graph.NodeID]bool{},
		path:    list.New(),
	}
}

func pathListToSlice(path *list.List) []graph.NodeID {
	out := []graph.NodeID{}
	for elem := path.Front(); elem != nil; elem = elem.Next() {
		out = append(out, elem.Value.(graph.NodeID))
	}
	return out
}

func pathListToString(path *list.List) string {
	out := []string{}
	for elem := path.Front(); elem != nil; elem = elem.Next() {
		out = append(out, string(elem.Value.(graph.NodeID)))
	}
	return strings.Join(out, ",")
}

func (m *PathManager) computeRelease(left int) int {
	cur := m.start
	totalRelease := 0
	releasePerMin := 0
	for elem := m.path.Front().Next(); elem != nil; elem = elem.Next() {
		stepID := elem.Value.(graph.NodeID)
		stepCost := m.g.EdgeCost(cur, stepID)

		// cost includes the minutes to get to stepID and 1min to turn
		// stepID on.
		totalRelease += releasePerMin * stepCost

		// Future minutes will include this rate too now that it's on.
		// If we turn a valve on in t=1 it doesn't count until t=2.
		releasePerMin += m.g.Rate(stepID)

		// We've now accounted for t=0 through the time this valve was
		// turned on.

		cur = stepID
	}

	// We've accounted for t=0 through the time the last valve was turned
	// on. We'll let everything run for the remaining `left` minutes.
	totalRelease += releasePerMin * left

	return totalRelease
}

func (m *PathManager) Visit(id graph.NodeID, left int) {
	m.path.PushBack(id)
	m.visited[id] = true

	release := m.computeRelease(left)
	if release > m.maxRelease {
		m.maxRelease = release
		m.maxReleasePath = pathListToSlice(m.path)
	}
}

func (m *PathManager) UnvisitLast() {
	elem := m.path.Back()
	last := elem.Value.(graph.NodeID)
	m.path.Remove(elem)
	delete(m.visited, last)
}

func (m *PathManager) Seen(id graph.NodeID) bool {
	return m.visited[id]
}

func (m *PathManager) MaxReleased() int {
	return m.maxRelease
}

func findMaxPath(g *SimpleGraph, cur graph.NodeID, left int, pathManager *PathManager) {
	pathManager.Visit(cur, left)
	defer pathManager.UnvisitLast()

	for neighbor, cost := range g.AllEdges(cur) {
		if cost > left || pathManager.Seen(neighbor) {
			continue
		}

		findMaxPath(g, neighbor, left-cost, pathManager)
	}
}

func solveA(nodes []*InputNode) int {
	simpleGraph := simplifyInputGraph(nodes)
	pathManager := NewPathManager(simpleGraph, "AA")
	findMaxPath(simpleGraph, "AA", 30, pathManager)

	return pathManager.MaxReleased()
}

func solveB(nodes []*InputNode) int {
	return -1
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	nodes, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(nodes))
	fmt.Println("B", solveB(nodes))
}
