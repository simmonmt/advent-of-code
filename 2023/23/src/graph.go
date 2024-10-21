package main

import (
	"fmt"

	"github.com/simmonmt/aoc/2023/common/dir"
	"github.com/simmonmt/aoc/2023/common/grid"
	"github.com/simmonmt/aoc/2023/common/pos"
)

// 1. find all intersections. those are the nodes
// 2. for each intersection
//      for each cardinal direction
//        can i get to another intersection?
//          add directed edge
//
// result:
//   intersections []pos.P2
//   edges map[pos.P2]pos.P2

type Graph struct {
	Start, End    pos.P2
	Intersections map[pos.P2]bool
	Edges         map[pos.P2][]Edge
}

type Edge struct {
	End  pos.P2
	Dist int
}

func allowedDir(r rune) dir.Dir {
	switch r {
	case '>':
		return dir.DIR_EAST
	case '<':
		return dir.DIR_WEST
	case 'v':
		return dir.DIR_SOUTH
	default:
		panic(fmt.Sprintf("bad r %c", r))
	}
}

func findEdge(isecs map[pos.P2]bool, board *grid.Grid[rune], cur pos.P2, curDir dir.Dir, restrict bool) (Edge, bool) {
	dist := 1
	for {
		next := curDir.From(cur)
		if _, found := isecs[next]; found {
			return Edge{End: next, Dist: dist}, true
		}

		r, ok := board.Get(next)
		if !ok {
			// N is off the board
			return Edge{}, false
		}

		if r == '#' {
			return Edge{}, false // it's a wall
		}

		// next isn't an intersection, and it's on the board. It can't
		// be '#', so it's either a direction restriction or a
		// free-travel space.
		if r != '.' && restrict {
			// It's a travel restriction
			if allowedDir(r) != curDir {
				return Edge{}, false // .. that we can't pass through
			}
		}

		// Advance to next
		cur = next

		// Update curDir
		nextDir := dir.DIR_UNKNOWN
		for _, d := range dir.AllDirs {
			if d == curDir.Reverse() {
				continue
			}

			if r, ok := board.Get(d.From(next)); !ok || r == '#' {
				continue // off the board or a wall
			}

			if nextDir != dir.DIR_UNKNOWN {
				panic(fmt.Sprintf("too many options at %v: had %v, also %v",
					next, nextDir, d))
			}
			nextDir = d
		}

		if nextDir == dir.DIR_UNKNOWN {
			return Edge{}, false // nowhere to go
		}

		curDir = nextDir
		dist++
	}
}

func findEdges(isecs map[pos.P2]bool, board *grid.Grid[rune], start pos.P2, restrict bool) []Edge {
	out := []Edge{}
	for _, d := range dir.AllDirs {
		if edge, ok := findEdge(isecs, board, start, d, restrict); ok {
			out = append(out, edge)
		}
	}
	return out
}

func BuildGraph(board *grid.Grid[rune], restrict bool) (*Graph, error) {
	g := &Graph{
		Intersections: map[pos.P2]bool{},
		Edges:         map[pos.P2][]Edge{},
	}

	board.Walk(func(p pos.P2, v rune) {
		if v == '#' {
			return
		}

		if p.Y == 0 {
			g.Start = p
			g.Intersections[p] = true
			return
		}
		if p.Y == board.Height()-1 {
			g.End = p
			g.Intersections[p] = true
			return
		}

		dirs := 0
		for _, n := range board.AllNeighbors(p, false) {
			if r, ok := board.Get(n); ok && r != '#' {
				dirs++
			}
		}
		if dirs > 2 {
			g.Intersections[p] = true
		}
	})

	for start := range g.Intersections {
		for _, edge := range findEdges(g.Intersections, board, start, restrict) {
			g.Edges[start] = append(g.Edges[start], edge)
		}
	}

	return g, nil
}
