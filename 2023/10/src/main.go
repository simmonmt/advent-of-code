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
	"github.com/simmonmt/aoc/2023/common/graph"
	"github.com/simmonmt/aoc/2023/common/grid"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func nodeType(p pos.P2, n Node) rune {
	hasNorth := n.A.Equals(dir.DIR_NORTH.From(p)) || n.B.Equals(dir.DIR_NORTH.From(p))
	hasSouth := n.A.Equals(dir.DIR_SOUTH.From(p)) || n.B.Equals(dir.DIR_SOUTH.From(p))
	hasEast := n.A.Equals(dir.DIR_EAST.From(p)) || n.B.Equals(dir.DIR_EAST.From(p))
	hasWest := n.A.Equals(dir.DIR_WEST.From(p)) || n.B.Equals(dir.DIR_WEST.From(p))

	if hasNorth && hasSouth {
		return '|'
	} else if hasEast && hasWest {
		return '-'
	} else if hasNorth {
		if hasEast {
			return 'L'
		} else if hasWest {
			return 'J'
		}
	} else if hasSouth {
		if hasEast {
			return 'F'
		} else if hasWest {
			return '7'
		}
	}

	panic("unknown type")
}

type Node struct {
	R    rune
	A, B pos.P2
}

func (n Node) ConnectsTo(p pos.P2) bool {
	return n.A.Equals(p) || n.B.Equals(p)
}

func NewBoard(lines []string) (*grid.Grid[Node], pos.P2, error) {
	var start pos.P2
	g, err := grid.NewFromLines(lines, func(p pos.P2, r rune) (Node, error) {
		node := Node{R: r}

		switch r {
		case '|':
			node.A = dir.DIR_NORTH.From(p)
			node.B = dir.DIR_SOUTH.From(p)
		case '-':
			node.A = dir.DIR_EAST.From(p)
			node.B = dir.DIR_WEST.From(p)
		case 'L':
			node.A = dir.DIR_NORTH.From(p)
			node.B = dir.DIR_EAST.From(p)
		case 'J':
			node.A = dir.DIR_NORTH.From(p)
			node.B = dir.DIR_WEST.From(p)
		case '7':
			node.A = dir.DIR_SOUTH.From(p)
			node.B = dir.DIR_WEST.From(p)
		case 'F':
			node.A = dir.DIR_SOUTH.From(p)
			node.B = dir.DIR_EAST.From(p)
		case 'S':
			start = p
		case '.':
		default:
			return node, fmt.Errorf("bad char %c at %s", r, p)
		}

		return node, nil
	})
	if err != nil {
		return nil, start, err
	}

	// We weren't given connections for S, so figure them out.
	numNeighbors := 0
	c1, c2 := pos.P2{}, pos.P2{}
	for _, neighbor := range start.AllNeighbors(false) {
		if nn, found := g.Get(neighbor); !found || !nn.ConnectsTo(start) {
			continue
		}
		if numNeighbors == 0 {
			c1 = neighbor
		} else {
			c2 = neighbor
		}

		numNeighbors += 1
	}
	if numNeighbors != 2 {
		return nil, start, fmt.Errorf("start has %d neighbors", numNeighbors)
	}

	sNode := Node{A: c1, B: c2}
	sNode.R = nodeType(start, sNode)
	g.Set(start, sNode)

	return g, start, nil
}

// PipeGraph is used to calculate the loop within the grid. If we sever one edge
// within the loop, we're really looking for the shortest path. sever1, sever2
// hold the ends of the edge that we're pretending doesn't exist. We fabricate
// the severed node like this because altering the grid to make it so would be
// bad form. Among other things, it would prevent board reuse in part B.
type PipeGraph struct {
	g              *grid.Grid[Node]
	sever1, sever2 pos.P2
}

func (pg *PipeGraph) Neighbors(id graph.NodeID) []graph.NodeID {
	cur, err := pos.P2FromString(string(id))
	if err != nil {
		panic("bad cur")
	}

	out := []graph.NodeID{}
	for _, neighbor := range cur.AllNeighbors(false) {
		if nn, found := pg.g.Get(neighbor); !found || !nn.ConnectsTo(cur) {
			continue
		}
		sever := (cur.Equals(pg.sever1) && neighbor.Equals(pg.sever2)) ||
			(cur.Equals(pg.sever2) && neighbor.Equals(pg.sever1))
		if sever {
			continue
		}

		out = append(out, graph.NodeID(neighbor.String()))
	}
	return out
}

func (pg *PipeGraph) NeighborDistance(_, _ graph.NodeID) int { return 1 }

func solveA(g *grid.Grid[Node], start pos.P2) int {
	startNode, _ := g.Get(start)

	path := graph.ShortestPath(graph.NodeID(startNode.A.String()),
		graph.NodeID(start.String()),
		&PipeGraph{g: g, sever1: startNode.A, sever2: start})

	// The full path is [start, startNode.A, path...], so the full path length is len(path) + 2

	return (len(path) + 2) / 2
}

func walkBorders(g *grid.Grid[Node], cb func(p pos.P2)) {
	for x := 0; x < g.Width(); x++ {
		cb(pos.P2{X: x, Y: 0})
		cb(pos.P2{X: x, Y: g.Height() - 1})
	}

	for y := 1; y < g.Height()-1; y++ {
		cb(pos.P2{X: 0, Y: y})
		cb(pos.P2{X: g.Width() - 1, Y: y})
	}
}

func findStartingPosition(g *grid.Grid[Node], pathPositions map[pos.P2]int) (borderPos pos.P2, insideDir dir.Dir) {
	todo := []pos.P2{}
	walkBorders(g, func(p pos.P2) {
		if _, found := pathPositions[p]; !found {
			todo = append(todo, p)
		}
	})

	if len(todo) == 0 {
		// There's no border space, so any edge will do
		return pos.P2{X: 0, Y: 1}, dir.DIR_EAST
	}

	seen := map[pos.P2]bool{}
	for len(todo) > 0 {
		next := []pos.P2{}

		for _, cur := range todo {
			if _, found := seen[cur]; found {
				continue
			}
			seen[cur] = true

			for _, n := range cur.AllNeighbors(false) {
				if _, found := seen[n]; found {
					continue
				}

				if _, found := pathPositions[n]; found {
					return n, inferDir(cur, n)
				}

				next = append(next, n)
			}
		}

		todo = next
	}

	panic("no starting found")
}

func inferDir(old, new pos.P2) dir.Dir {
	if new.Y == old.Y {
		if new.X > old.X {
			return dir.DIR_EAST
		} else {
			return dir.DIR_WEST
		}
	} else if new.Y > old.Y {
		return dir.DIR_SOUTH
	} else {
		return dir.DIR_NORTH
	}
}

type TurnAttrs struct {
	LeftTurnDir dir.Dir
	Ignore      [2]dir.Dir
	Decorate    []pos.P2
}

var (
	turnAttrs = map[rune]TurnAttrs{
		'L': TurnAttrs{
			LeftTurnDir: dir.DIR_SOUTH,
			Ignore:      [2]dir.Dir{dir.DIR_NORTH, dir.DIR_EAST},
			Decorate:    []pos.P2{pos.P2{X: -1, Y: 0}, pos.P2{X: -1, Y: 1}, pos.P2{X: 0, Y: 1}},
		},
		'J': TurnAttrs{
			LeftTurnDir: dir.DIR_EAST,
			Ignore:      [2]dir.Dir{dir.DIR_NORTH, dir.DIR_WEST},
			Decorate:    []pos.P2{pos.P2{X: 0, Y: 1}, pos.P2{X: 1, Y: 1}, pos.P2{X: 1, Y: 0}},
		},
		'7': TurnAttrs{
			LeftTurnDir: dir.DIR_NORTH,
			Ignore:      [2]dir.Dir{dir.DIR_SOUTH, dir.DIR_WEST},
			Decorate:    []pos.P2{pos.P2{X: 0, Y: -1}, pos.P2{X: 1, Y: -1}, pos.P2{X: 1, Y: 0}},
		},
		'F': TurnAttrs{
			LeftTurnDir: dir.DIR_WEST,
			Ignore:      [2]dir.Dir{dir.DIR_SOUTH, dir.DIR_EAST},
			Decorate:    []pos.P2{pos.P2{X: -1, Y: 0}, pos.P2{X: -1, Y: -1}, pos.P2{X: 0, Y: -1}},
		},
	}
)

func decorate(g *grid.Grid[Node], pathPositions map[pos.P2]int, p pos.P2, insideDir dir.Dir, insides map[pos.P2]bool) {
	cands := []pos.P2{}
	node, _ := g.Get(p)

	if node.R == '|' || node.R == '-' {
		cands = append(cands, insideDir.From(p))
	} else if node.R == 'L' || node.R == 'J' || node.R == '7' || node.R == 'F' {
		attrs := turnAttrs[node.R]
		if insideDir != attrs.Ignore[0] && insideDir != attrs.Ignore[1] {
			for _, off := range attrs.Decorate {
				dp := p
				dp.Add(off)
				cands = append(cands, dp)
			}
		}
	} else {
		panic("bad node type")
	}

	for _, cand := range cands {
		if _, found := insides[cand]; found {
			continue
		}
		if _, found := pathPositions[cand]; found {
			continue
		}

		insides[cand] = true
	}
}

func advance(g *grid.Grid[Node], p pos.P2, travelDir, insideDir dir.Dir) (newTravelDir, newInsideDir dir.Dir) {
	node, _ := g.Get(p)

	if node.R == '|' || node.R == '-' {
		return travelDir, insideDir // no change
	}

	attrs := turnAttrs[node.R]
	if attrs.LeftTurnDir == travelDir {
		return travelDir.Left(), insideDir.Left()
	} else {
		return travelDir.Right(), insideDir.Right()
	}
}

func fillInsideFrom(p pos.P2, g *grid.Grid[Node], pathPositions map[pos.P2]int, insides map[pos.P2]bool) {
	todo := []pos.P2{p}
	seen := map[pos.P2]bool{}

	for len(todo) > 0 {
		next := []pos.P2{}

		for _, cand := range todo {
			if seen[cand] {
				continue
			}
			seen[cand] = true

			for _, n := range g.AllNeighbors(cand, false) {
				if _, found := pathPositions[n]; found {
					continue
				}
				if _, found := insides[n]; found {
					continue
				}

				insides[n] = true
				next = append(next, n)
			}
		}

		todo = next
	}
}

func fillInsides(g *grid.Grid[Node], pathPositions map[pos.P2]int, insides map[pos.P2]bool) {
	g.Walk(func(p pos.P2, node Node) {
		for _, neighbor := range p.AllNeighbors(false) {
			if _, found := pathPositions[neighbor]; found {
				continue
			}
			if _, found := insides[neighbor]; found {
				fillInsideFrom(neighbor, g, pathPositions, insides)
			}
		}
	})
}

func solveB(g *grid.Grid[Node], start pos.P2) int {
	startNode, _ := g.Get(start)

	nodeIDPath := graph.ShortestPath(graph.NodeID(startNode.A.String()),
		graph.NodeID(start.String()),
		&PipeGraph{g: g, sever1: startNode.A, sever2: start})

	// nodeIDPath started at A, so the first element is the position chosen
	// *after* A. We therefore need to pre-add it.
	pathPositions := map[pos.P2]int{startNode.A: 0}
	path := []pos.P2{startNode.A}
	loc := 1
	for _, nid := range nodeIDPath {
		p, _ := pos.P2FromString(string(nid))
		pathPositions[p] = loc
		path = append(path, p)
		loc++
	}

	insides := map[pos.P2]bool{}
	borderPos, insideDir := findStartingPosition(g, pathPositions)
	next := path[(pathPositions[borderPos]+1)%len(path)]
	travelDir := inferDir(borderPos, next)

	logger.Infof("starting with border pos %s inside %s", borderPos, insideDir)

	// We find the starting position by breadth first search from the
	// outside. Consider the F in this example:
	//
	//                ......
	//                ...F7.
	//                ...J|.
	//
	// From the perspective of the starting position finder, which doesn't
	// know about the path, the inside direction for F could be EAST or
	// SOUTH. The rest of the code, however, really wants one or the other
	// depending on the direction of travel. It wants what advance()
	// would've returned. If we were going clockwise, advance() would've
	// returned SOUTH.  If counterclockwise, EAST.
	//
	// The bad cases happen when the direction of travel (which we can
	// figure out because we have the path) is the same as the inside
	// direction. In the example, this happens when we discover F *from* the
	// north when the direction of travel is CCW, or when we discover it
	// from the west when the direction of travel is CW. In the other two
	// cases we'll guess the correct direction.
	//
	// As it happens, in the bad cases the direction is what the inside
	// direction would've been on the previous node (the J for CW, 7 for
	// CCW). So to fix the inside direction we just need to advance it as if
	// we were going from the previous node to the discovered border node,
	// leaving the inferred travel direction alone.
	if travelDir == insideDir {
		prev := path[(len(path)+pathPositions[borderPos]-1)%len(path)]
		_, insideDir = advance(g, borderPos, inferDir(prev, borderPos), insideDir)
	}

	for cur := borderPos; !next.Equals(borderPos); cur = next {
		logger.Infof("cur %s travel %s inside %s", cur, travelDir, insideDir)

		// decorate
		decorate(g, pathPositions, cur, insideDir, insides)

		// advance
		next = travelDir.From(cur)
		travelDir, insideDir = advance(g, next, travelDir, insideDir)
		if travelDir == insideDir {
			panic("dirs unexpectedly match")
		}
	}

	logger.Infof("filling insides")
	fillInsides(g, pathPositions, insides)

	if false {
		g.Dump(true, func(p pos.P2, n Node, _ bool) string {
			if _, found := pathPositions[p]; found {
				return string(n.R)
			}
			if _, found := insides[p]; found {
				return "o"
			}
			return "."
		})
	}

	return len(insides)
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

	g, start, err := NewBoard(lines)
	if err != nil {
		logger.Fatalf("bad board: %v", err)
	}

	fmt.Println("A", solveA(g, start))
	fmt.Println("B", solveB(g, start))
}
