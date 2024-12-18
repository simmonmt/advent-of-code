// Copyright 2024 Google LLC
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

package astar

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2024/common/logger"
)

type helperNode struct {
	distances map[string]uint
}

type aStarHelper struct {
	nodes map[string]helperNode
}

func (h *aStarHelper) AllNeighbors(start string) []string {
	node, found := h.nodes[start]
	if !found {
		return nil
	}

	neighbors := []string{}
	for neighbor := range node.distances {
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}

func (h *aStarHelper) EstimateDistance(start, end string) uint {
	if start == end {
		return 0
	} else {
		return 1
	}
}

func (h *aStarHelper) NeighborDistance(n1, n2 string) uint {
	if n1 == n2 {
		return 0
	}

	node, found := h.nodes[n1]
	if !found {
		return 0
	}

	for neighbor, dist := range node.distances {
		if neighbor == n2 {
			return dist
		}
	}

	panic(fmt.Sprintf("no distance for %v to %v", n1, n2))
}

func (h *aStarHelper) GoalReached(cand, goal string) bool {
	return cand == goal
}

func (h *aStarHelper) Serialize(node string) string {
	return "ZZ" + node
}

func (h *aStarHelper) Deserialize(val string) (string, error) {
	if !strings.HasPrefix(val, "ZZ") {
		return "", fmt.Errorf("no prefix")
	}
	return val[2:], nil
}

func TestAStar(t *testing.T) {
	helper := aStarHelper{
		nodes: map[string]helperNode{
			"start": helperNode{distances: map[string]uint{"a": 15, "d": 20}},
			"a":     helperNode{distances: map[string]uint{"start": 15, "b": 20}},
			"b":     helperNode{distances: map[string]uint{"a": 20, "c": 30}},
			"c":     helperNode{distances: map[string]uint{"b": 30, "end": 40}},
			"d":     helperNode{distances: map[string]uint{"start": 20, "d1": 2, "e": 20}},
			"d1":    helperNode{distances: map[string]uint{"d": 2, "d2": 1}},
			"d2":    helperNode{distances: map[string]uint{"d1": 1, "d3": 1}},
			"d3":    helperNode{distances: map[string]uint{"d2": 1, "e": 1}},
			"e":     helperNode{distances: map[string]uint{"d": 20, "d3": 1, "end": 20}},
			"end":   helperNode{distances: map[string]uint{"c": 40, "e": 20}},
		},
	}

	solver := New("start", "end", &helper)
	result := solver.Solve()

	expected := []string{"end", "e", "d3", "d2", "d1", "d", "start"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("start->end, got %v, want %v", result, expected)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	logger.Init(true)

	os.Exit(m.Run())
}
