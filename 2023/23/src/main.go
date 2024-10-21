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

	"github.com/simmonmt/aoc/2023/common/dir"
	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/grid"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) (*grid.Grid[rune], error) {
	return grid.NewFromLines(lines, func(p pos.P2, r rune) (rune, error) {
		if r != '#' && r != '.' && r != '^' && r != '<' && r != '>' && r != 'v' {
			return '?', fmt.Errorf(fmt.Sprintf("bad rune %c", r))
		}
		return r, nil
	})
}

// type Graph struct {
// 	g *grid.Grid[rune]
// }

// func (g *Graph) Deserialize(id graph.NodeID) pos.P2 {
// 	p, err := pos.P2FromString(string(id))
// 	if err != nil {
// 		panic("bad pos")
// 	}
// 	return p
// }

// func (g *Graph) Serialize(p pos.P2) graph.NodeID {
// 	return graph.NodeID(p.String())
// }

// func (g *Graph) Neighbors(id graph.NodeID) []graph.NodeID {
// 	p := g.Deserialize(id)

// 	out := []graph.NodeID{}
// 	for _, d := range dir.AllDirs {
// 		n := d.From(p)
// 		r, found := g.g.Get(n)
// 		if !found || r == '#' {
// 			continue
// 		}
// 		if r == '.' {
// 			out = append(out, g.Serialize(n))
// 			continue
// 		}

// 		ad, found := allowedDirs[r]
// 		if !found {
// 			panic("bad rune")
// 		}

// 		logger.Infof("at %v, considering %v which is %c to the %s ad %s",
// 			p, n, r, d, ad)

// 		if d == ad {
// 			out = append(out, g.Serialize(n))
// 			continue
// 		}
// 	}
// 	return out
// }

// func (g *Graph) NeighborDistance(from, to graph.NodeID) int {
// 	return -1
// }

var (
	allowedDirs = map[rune]dir.Dir{
		'^': dir.DIR_NORTH,
		'v': dir.DIR_SOUTH,
		'<': dir.DIR_WEST,
		'>': dir.DIR_EAST,
	}
)

func Neighbors(g *grid.Grid[rune], p pos.P2, disallowSlopes bool) []pos.P2 {
	out := []pos.P2{}
	for _, d := range dir.AllDirs {
		n := d.From(p)
		r, found := g.Get(n)
		if !found || r == '#' {
			continue
		}
		if r == '.' {
			out = append(out, n)
			continue
		}

		if disallowSlopes {
			ad, found := allowedDirs[r]
			if !found {
				panic("bad rune")
			}
			if d != ad {
				continue
			}
		}

		//logger.Infof("at %v, considering %v which is %c to the %s ad %s",
		//p, n, r, d, ad)

		out = append(out, n)
	}
	return out
}

func dfs(g *Graph, start, end pos.P2, path map[pos.P2]bool, dist int, foundEnd func(map[pos.P2]bool, int)) {
	//logger.Infof("start %v end %v", start, end)

	if start.Equals(end) {
		foundEnd(path, dist)
		return
	}

	edges := []Edge{}
	for _, edge := range g.Edges[start] {
		if found := path[edge.End]; found {
			continue
		}

		edges = append(edges, edge)
	}

	if len(edges) == 0 {
		return
	}

	for _, edge := range edges {
		path[edge.End] = true
		dfs(g, edge.End, end, path, dist+edge.Dist, foundEnd)
		path[edge.End] = false
	}
}

func solve(board *grid.Grid[rune], restrict bool) int {
	g, err := BuildGraph(board, restrict)
	if err != nil {
		panic("bad graph")
	}

	sofar := map[pos.P2]bool{g.Start: true}
	longest := 0
	dfs(g, g.Start, g.End, sofar, 0, func(path map[pos.P2]bool, dist int) {
		longest = max(longest, dist)
	})

	return longest
}

func solveA(board *grid.Grid[rune]) int {
	return solve(board, true)
}

func solveB(board *grid.Grid[rune]) int {
	return solve(board, false)
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
