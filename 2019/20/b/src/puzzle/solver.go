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

package puzzle

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2019/common/astar"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type warpPath struct {
	nearPos pos.P2
	farPos  pos.P2
	gate    Gate
	dist    int
}

type solverState struct {
	board *Board
	graph map[pos.P2][]warpPath
}

func gateOpen(level int, goingIn bool, gate *Gate) bool {
	if gate.name == "AA" {
		return false
	}
	if gate.name == "ZZ" {
		return level == 0
	}

	if !goingIn && level == 0 {
		// AA is closed always, ZZ is available outbound only
		// from level 0, and they were checked above.
		return false
	}

	return true
}

func (a *solverState) parseNode(s string) (level int, p pos.P2) {
	parts := strings.Split(s, ":")

	level, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(fmt.Sprintf("bad level %v", s))
	}

	p, err = pos.P2FromString(parts[1])
	if err != nil {
		panic(fmt.Sprintf("bad pos %v", s))
	}
	return level, p
}

func (a *solverState) nodeToString(level int, p pos.P2) string {
	return fmt.Sprintf("%d:%s", level, p.String())
}

func (a *solverState) AllNeighbors(start string) []string {
	startLevel, startPos := a.parseNode(start)

	out := []string{}
	for _, path := range a.graph[startPos] {
		//fmt.Printf("an of %v eval %+v\n", start, path)

		goingIn := path.nearPos.Equals(path.gate.GateIn())
		//fmt.Printf("goingIn %v\n", goingIn)

		farLevel := startLevel
		if goingIn {
			farLevel++
		} else {
			farLevel--
		}

		if gateOpen(startLevel, goingIn, &path.gate) {
			out = append(out, a.nodeToString(farLevel, path.farPos))
		}
	}

	return out

}

func (a *solverState) NeighborDistance(n1, n2 string) uint {
	_, n1Pos := a.parseNode(n1)
	_, n2Pos := a.parseNode(n2)

	for _, path := range a.graph[n1Pos] {
		if path.farPos.Equals(n2Pos) {
			return uint(path.dist)
		}
	}

	panic("no neighbor")
}

func (a *solverState) EstimateDistance(start, end string) uint {
	// We can't use Manhattan distance within a level between start and end
	// because the portals let us go shorter than Manhattan distance. A*
	// only works if we *underestimate* the distance between start and end.

	startLevel, _ := a.parseNode(start)
	endLevel, _ := a.parseNode(end)
	return uint(endLevel - startLevel + 1)
}

func (a *solverState) GoalReached(cand, goal string) bool {
	return cand == goal
}

func Solve(board *Board, graph map[pos.P2][]Path, start, end pos.P2) (int, bool) {
	opposite := func(board *Board, p pos.P2) pos.P2 {
		gate := board.GateByGateLoc(p)
		var other pos.P2
		if gp := gate.GateOut(); gp.Equals(p) {
			other = gate.PortalIn()
		} else {
			other = gate.PortalOut()
		}
		if other.X == -1 && other.Y == -1 {
			return p
		}
		return other
	}

	warpGraph := map[pos.P2][]warpPath{}
	for from, paths := range graph {
		warps := []warpPath{}
		for _, p := range paths {
			warps = append(warps, warpPath{
				nearPos: p.DestPos,
				farPos:  opposite(board, p.DestPos),
				gate:    board.GateByGateLoc(p.DestPos),
				dist:    p.Dist,
			})
		}
		warpGraph[from] = warps
	}

	state := &solverState{
		board: board,
		graph: warpGraph,
	}

	rawPath := astar.AStar(state.nodeToString(0, start), state.nodeToString(-1, end), state)
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
