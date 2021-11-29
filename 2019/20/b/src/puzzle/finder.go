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

	"github.com/simmonmt/aoc/2019/common/astar"
	"github.com/simmonmt/aoc/2019/common/dir"
	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type finderState struct {
	board *Board
}

func (a *finderState) parseNode(s string) pos.P2 {
	p, err := pos.P2FromString(s)
	if err != nil {
		panic(fmt.Sprintf("bad pos %v", s))
	}
	return p
}

func (a *finderState) nodeToString(p pos.P2) string {
	return p.String()
}

func (a *finderState) AllNeighbors(start string) []string {
	startPos := a.parseNode(start)

	if a.board.Get(startPos) == TILE_GATE {
		return nil
	}

	neighbors := []pos.P2{}
	for _, d := range dir.AllDirs {
		np := d.From(startPos)
		switch t := a.board.Get(np); t {
		case TILE_GATE:
			fallthrough
		case TILE_PATH:
			neighbors = append(neighbors, np)
			break
		case TILE_WALL:
			continue
		case TILE_OPEN:
			continue
		default:
			panic(fmt.Sprintf("unexpected tile type %s at %v", t, np))
		}
	}

	outs := make([]string, len(neighbors))
	for i, n := range neighbors {
		outs[i] = a.nodeToString(n)
	}
	return outs
}

func (a *finderState) NeighborDistance(n1, n2 string) uint {
	return a.EstimateDistance(n1, n2)
}

func (a *finderState) EstimateDistance(start, end string) uint {
	startPos, endPos := a.parseNode(start), a.parseNode(end)
	return uint(startPos.ManhattanDistance(endPos))
}

func (a *finderState) GoalReached(cand, goal string) bool {
	logger.LogF("goal reached? %v %v %v", cand, goal, cand == goal)
	return cand == goal
}

func findShortestPathNoPortals(from pos.P2, to pos.P2, board *Board) (int, bool) {
	state := &finderState{
		board: board,
	}

	rawPath := astar.AStar(state.nodeToString(from), state.nodeToString(to), state)
	if len(rawPath) == 0 {
		return 0, false
	}

	return len(rawPath) - 1, true
}

type Path struct {
	DestPos pos.P2
	Dist    int
}

func FindAllPathsFromPortal(name string, portal pos.P2, board *Board) []Path {
	cands := []pos.P2{}
	for _, gate := range board.Gates() {
		if gate.Name() != name {
			cands = append(cands, gate.GateOut(), gate.GateIn())
		}
	}

	paths := []Path{}
	for _, cand := range cands {
		if cost, found := findShortestPathNoPortals(portal, cand, board); found {
			paths = append(paths, Path{
				DestPos: cand,
				Dist:    cost,
			})
		}
	}

	return paths
}

func FindAllPathsFromAllPortals(board *Board) map[pos.P2][]Path {
	out := map[pos.P2][]Path{}
	for _, gate := range board.Gates() {
		for _, start := range []pos.P2{gate.PortalOut(), gate.PortalIn()} {
			if start.X == -1 && start.Y == -1 {
				continue
			}

			paths := FindAllPathsFromPortal(gate.Name(), start, board)
			if len(paths) > 0 {
				out[start] = paths
			}
		}
	}
	return out
}
