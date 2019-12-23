package puzzle

import (
	"fmt"
	"sort"
	"strings"

	"github.com/simmonmt/aoc/2019/common/astar"
	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/pos"
)

type astarState struct {
	board   *Board
	graphs  map[pos.P2]map[string][]Path
	numKeys int
}

func parseNode(s string) (ps []pos.P2, keys map[string]bool) {
	parts := strings.Split(s, "_")

	pStrs := parts[0:(len(parts) - 1)]
	keyStr := parts[len(parts)-1]

	ps = []pos.P2{}
	for _, pStr := range pStrs {
		p, err := pos.P2FromString(parts[0])
		if err != nil {
			panic(fmt.Sprintf("bad pos '%s'", pStr))
		}
		ps = append(ps, p)
	}

	keys = map[string]bool{}
	for _, key := range strings.Split(keyStr, ",") {
		if key != "" {
			keys[key] = true
		}
	}

	return ps, keys
}

func nodeToString(ps []pos.P2, keys map[string]bool) string {
	outs := make([]string, len(ps)+1)
	for i := 0; i < len(ps); i++ {
		outs[i] = ps[i].String()
	}

	keyArr := make([]string, len(keys))
	i := 0
	for key := range keys {
		keyArr[i] = key
		i++
	}
	sort.Strings(keyArr)

	outs[len(ps)] = strings.Join(keyArr, ",")
	return strings.Join(outs, "_")
}

func (a *astarState) pathsFromPos(p pos.P2) []Path {
	t := a.board.Get(p)
	if t == TILE_KEY {
		keyName := a.board.KeyAtLoc(p)
		for _, g := range a.graphs {
			if paths, found := g[keyName]; found {
				return paths
			}
		}
	} else {
		for graphPos, graph := range a.graphs {
			if p.Equals(graphPos) {
				return graph["@"]
			}
		}
	}
	panic(fmt.Sprintf("no graph for %v", p))
}

func (a *astarState) AllNeighbors(start string) []string {
	startPosns, keys := parseNode(start)

	neighbors := []string{}
	for startPosIdx, startPos := range startPosns {
		if t := a.board.Get(startPos); t == TILE_KEY {
			keys[a.board.KeyAtLoc(startPos)] = true
		}

		paths := a.pathsFromPos(startPos)

		avail := []Path{}
		for _, path := range paths {
			allowed := true
			//fmt.Printf("eval path %v\n", path)
			for _, needDoor := range path.Doors {
				needKey := string(needDoor[0] - 'A' + 'a')
				if _, have := keys[needKey]; !have {
					allowed = false
					break
				}
			}

			if !allowed {
				continue
			}

			//fmt.Printf("path allowed\n")

			avail = append(avail, path)
		}

		//fmt.Printf("search: avail: %v\n", avail)

		newPosns := make([]pos.P2, len(startPosns))
		copy(newPosns, startPosns)

		for i := 0; i < len(avail); i++ {
			newPosns[startPosIdx] = a.board.KeyLoc(avail[i].Dest)
			neighbors = append(neighbors, nodeToString(newPosns, keys))
		}
	}

	logger.LogF("neighbors of %s are %v", start, neighbors)
	return neighbors
}

func (a *astarState) EstimateDistance(start, end string) uint {
	_, startKeys := parseNode(start)
	if end == "" {
		return uint(a.numKeys - len(startKeys))
	}

	_, endKeys := parseNode(end)
	return uint(len(endKeys) - len(startKeys))
}

func findChangedPosIdx(a, b []pos.P2) int {
	if len(a) != len(b) {
		panic("mismatch")
	}

	for i := range a {
		if !a[i].Equals(b[i]) {
			return i
		}
	}

	panic("no change")
}

func (a *astarState) findCostFromChange(from, to []pos.P2) int {
	changedIdx := findChangedPosIdx(from, to)

	n2Key := a.board.KeyAtLoc(to[changedIdx])

	for _, path := range a.pathsFromPos(from[changedIdx]) {
		if path.Dest == n2Key {
			return path.Dist
		}
	}

	panic("unable to find change")
}

func (a *astarState) NeighborDistance(n1, n2 string) uint {
	n1Posns, _ := parseNode(n1)
	n2Posns, _ := parseNode(n2)

	return uint(a.findCostFromChange(n1Posns, n2Posns))
}

func (a *astarState) GoalReached(cand, goal string) bool {
	posns, candKeys := parseNode(cand)

	// For each robot that's on a key, pretend its key is in the keys
	// list. A* doesn't know that being on a key means it's picked up --
	// from its perspective keys are only picked up when we leave a node.
	for _, p := range posns {
		if a.board.Get(p) == TILE_KEY {
			candKeys[a.board.KeyAtLoc(p)] = true
		}
	}

	return len(candKeys) == a.numKeys
}

func (a *astarState) findPathCost(path []string) int {
	cost := 0
	for i := len(path) - 1; i >= 1; i-- {
		curPosns, _ := parseNode(path[i])
		nextPosns, _ := parseNode(path[i-1])
		cost += a.findCostFromChange(curPosns, nextPosns)
	}

	return cost
}

func FindShortestPath(board *Board, graph map[string][]Path, numKeys int, start pos.P2) ([]string, int) {
	state := &astarState{
		board:   board,
		graphs:  map[pos.P2]map[string][]Path{start: graph},
		numKeys: numKeys,
	}

	startNode := nodeToString([]pos.P2{start}, nil)
	path := astar.AStar(startNode, "", state)

	cost := state.findPathCost(path)
	return path, cost
}

func FindShortestPathMultiStart(board *Board, graphs map[pos.P2]map[string][]Path, numKeys int, starts []pos.P2) ([]string, int) {
	return nil, 0
}
