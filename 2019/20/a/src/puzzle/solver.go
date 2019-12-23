package puzzle

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/astar"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type solverState struct {
	board *Board
	graph map[pos.P2][]Path
}

func (a *solverState) parseNode(s string) pos.P2 {
	p, err := pos.P2FromString(s)
	if err != nil {
		panic(fmt.Sprintf("bad pos %v", s))
	}
	return p
}

func (a *solverState) nodeToString(p pos.P2) string {
	return p.String()
}

func (a *solverState) AllNeighbors(start string) []string {
	startPos := a.parseNode(start)

	out := []string{}
	for _, path := range a.graph[startPos] {
		out = append(out, a.nodeToString(path.DestPos))
	}
	return out
}

func (a *solverState) NeighborDistance(n1, n2 string) uint {
	n1Pos, n2Pos := a.parseNode(n1), a.parseNode(n2)

	for _, path := range a.graph[n1Pos] {
		if path.DestPos.Equals(n2Pos) {
			return uint(path.Dist)
		}
	}

	panic("no neighbor")
}

func (a *solverState) EstimateDistance(start, end string) uint {
	// It can't be Manhattan distance between start and end because the
	// portals let us go shorter than Manhattan distance. A* only works if
	// we *underestimate* the distance between start and end.
	return 1
	// startPos, endPos := a.parseNode(start), a.parseNode(end)
	// return uint(startPos.ManhattanDistance(endPos))
}

func (a *solverState) GoalReached(cand, goal string) bool {
	return cand == goal
}

func Solve(board *Board, graph map[pos.P2][]Path, start, end pos.P2) (int, bool) {
	opposite := func(board *Board, p pos.P2) pos.P2 {
		gate := board.GateByGateLoc(p)
		var other pos.P2
		if gp := gate.Gate1(); gp.Equals(p) {
			other = gate.Portal2()
		} else {
			other = gate.Portal1()
		}
		if other.X == -1 && other.Y == -1 {
			return p
		}
		return other
	}

	warpGraph := map[pos.P2][]Path{}
	for from, paths := range graph {
		warp := []Path{}
		for _, p := range paths {
			warp = append(warp, Path{opposite(board, p.DestPos), p.Dist})
		}
		warpGraph[from] = warp
	}

	state := &solverState{
		board: board,
		graph: warpGraph,
	}

	rawPath := astar.AStar(state.nodeToString(start), state.nodeToString(end), state)
	if len(rawPath) == 0 {
		return 0, false
	}

	cost := uint(0)
	for i := len(rawPath) - 1; i >= 1; i-- {
		from, to := rawPath[i], rawPath[i-1]
		cost += state.NeighborDistance(from, to)
	}

	fmt.Printf("raw path %v\n", rawPath)

	return int(cost), true
}
