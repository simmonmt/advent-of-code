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
	"sort"
	"strings"

	"github.com/simmonmt/aoc/2019/common/astar"
	"github.com/simmonmt/aoc/2019/common/logger"
)

type astarState struct {
	graph   map[string][]Path
	numKeys int
}

func parseNode(s string) (pos string, keys map[string]bool) {
	parts := strings.Split(s, "_")
	if len(parts) != 2 {
		panic(fmt.Sprintf("bad node '%s'", s))
	}

	keys = map[string]bool{}
	for _, key := range strings.Split(parts[1], ",") {
		if key != "" {
			keys[key] = true
		}
	}

	return parts[0], keys
}

func nodeToString(pos string, keys map[string]bool) string {
	keyArr := make([]string, len(keys))
	i := 0
	for key := range keys {
		keyArr[i] = key
		i++
	}
	sort.Strings(keyArr)

	return fmt.Sprintf("%s_%s", pos, strings.Join(keyArr, ","))
}

func (a *astarState) AllNeighbors(start string) []string {
	pos, keys := parseNode(start)

	//fmt.Printf("allneigbors %v => %v,%v\n", start, pos, keys)

	if pos != "@" {
		keys[pos] = true
	}

	//fmt.Printf("keys now %s\n", keys)

	avail := []Path{}
	for _, path := range a.graph[pos] {
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
		neighbors[i] = nodeToString(avail[i].Dest, keys)
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

	for _, path := range a.graph[n1Pos] {
		if path.Dest == n2Pos {
			return uint(path.Dist)
		}
	}

	panic(fmt.Sprintf("%s and %s have no path", n1Pos, n2Pos))
}

func (a *astarState) GoalReached(cand, goal string) bool {
	pos, candKeys := parseNode(cand)
	if pos == "@" {
		return false
	}

	// Pretend pos is in the keys list because we're standing on
	// that node. A* just doesn't know that means it's been picked
	// up -- from its perspective keys are only picked up when we
	// leave a node.
	candKeys[pos] = true
	return len(candKeys) == a.numKeys
}

func findPathCost(graph map[string][]Path, path []string) int {
	cost := 0
	for i := len(path) - 1; i >= 1; i-- {
		curPos, _ := parseNode(path[i])
		nextPos, _ := parseNode(path[i-1])

		for _, path := range graph[curPos] {
			if path.Dest == nextPos {
				cost += path.Dist
				break
			}
		}
	}

	return cost
}

func FindShortestPath(graph map[string][]Path, numKeys int, start string) ([]string, int) {
	state := &astarState{
		graph:   graph,
		numKeys: numKeys,
	}

	startNode := nodeToString(start, nil)
	path := astar.AStar(startNode, "", state)

	cost := findPathCost(graph, path)
	return path, cost
}
