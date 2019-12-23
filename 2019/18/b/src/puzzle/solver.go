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

func parseNode(s string) (p pos.P2, keys map[string]bool) {
	parts := strings.Split(s, "_")
	if len(parts) != 2 {
		panic(fmt.Sprintf("bad node '%s'", s))
	}

	p, err := pos.P2FromString(parts[0])
	if err != nil {
		panic(fmt.Sprintf("bad pos '%s'", parts[0]))
	}

	keys = map[string]bool{}
	for _, key := range strings.Split(parts[1], ",") {
		if key != "" {
			keys[key] = true
		}
	}

	return p, keys
}

func nodeToString(p pos.P2, keys map[string]bool) string {
	keyArr := make([]string, len(keys))
	i := 0
	for key := range keys {
		keyArr[i] = key
		i++
	}
	sort.Strings(keyArr)

	return fmt.Sprintf("%s_%s", p.String(), strings.Join(keyArr, ","))
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
	startPos, keys := parseNode(start)

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

	neighbors := make([]string, len(avail))
	for i := 0; i < len(avail); i++ {
		neighbors[i] = nodeToString(a.board.KeyLoc(avail[i].Dest), keys)
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

func (a *astarState) NeighborDistance(n1, n2 string) uint {
	n1Pos, _ := parseNode(n1)
	n2Pos, _ := parseNode(n2)

	n2Key := a.board.KeyAtLoc(n2Pos)

	for _, path := range a.pathsFromPos(n1Pos) {
		if path.Dest == n2Key {
			return uint(path.Dist)
		}
	}

	panic(fmt.Sprintf("%s and %s have no path", n1Pos, n2Pos))
}

func (a *astarState) GoalReached(cand, goal string) bool {
	p, candKeys := parseNode(cand)

	if a.board.Get(p) != TILE_KEY {
		return false
	}

	// Pretend pos is in the keys list because we're standing on
	// that node. A* just doesn't know that means it's been picked
	// up -- from its perspective keys are only picked up when we
	// leave a node.
	candKeys[a.board.KeyAtLoc(p)] = true
	return len(candKeys) == a.numKeys
}

func (a *astarState) findPathCost(path []string) int {
	cost := 0
	for i := len(path) - 1; i >= 1; i-- {
		curPos, _ := parseNode(path[i])
		nextPos, _ := parseNode(path[i-1])

		for _, path := range a.pathsFromPos(curPos) {
			if path.Dest == a.board.KeyAtLoc(nextPos) {
				cost += path.Dist
				break
			}
		}
	}

	return cost
}

func FindShortestPath(board *Board, graph map[string][]Path, numKeys int, start pos.P2) ([]string, int) {
	state := &astarState{
		board:   board,
		graphs:  map[pos.P2]map[string][]Path{start: graph},
		numKeys: numKeys,
	}

	startNode := nodeToString(start, nil)
	path := astar.AStar(startNode, "", state)

	cost := state.findPathCost(path)
	return path, cost
}

func FindShortestPathMultiStart(board *Board, graphs map[pos.P2]map[string][]Path, numKeys int, starts []pos.P2) ([]string, int) {
	return nil, 0
}
