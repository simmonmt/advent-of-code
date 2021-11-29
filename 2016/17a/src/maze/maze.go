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

package maze

import (
	"crypto/md5"
	"fmt"

	"astar"
	"intmath"
	"node"
)

type aStarHelper struct {
	width, height int
}

func NewHelper(width, height int) *aStarHelper {
	return &aStarHelper{
		width:  width,
		height: height,
	}
}

func (h *aStarHelper) tryNeighbor(start *node.Node, dir string, xDelta, yDelta int, val byte, neighbors *[]string) {
	if val < 0xb {
		return
	}

	n := &node.Node{
		X:        start.X + xDelta,
		Y:        start.Y + yDelta,
		Passcode: start.Passcode,
		Path:     start.Path + dir,
	}

	if n.X < 0 || n.Y < 0 {
		return
	}
	if n.X >= h.width || n.Y >= h.height {
		return
	}

	*neighbors = append(*neighbors, n.Serialize())
}

func (h *aStarHelper) AllNeighbors(start string) []string {
	neighbors := []string{}

	n, err := node.Deserialize(start)
	if err != nil {
		panic(fmt.Sprintf(`failed to deserialize "%v": %v`, start, err))
	}

	hash := md5.Sum([]byte(n.Passcode + n.Path))
	h.tryNeighbor(n, "U", 0, -1, hash[0]>>4, &neighbors)
	h.tryNeighbor(n, "D", 0, 1, hash[0]&0xf, &neighbors)
	h.tryNeighbor(n, "L", -1, 0, hash[1]>>4, &neighbors)
	h.tryNeighbor(n, "R", 1, 0, hash[1]&0xf, &neighbors)

	return neighbors
}

func (h *aStarHelper) EstimateDistance(start, end string) uint {
	sNode, err := node.Deserialize(start)
	if err != nil {
		panic(fmt.Sprintf("failed to deserialize start %v: %v", start, err))
	}

	eNode, err := node.Deserialize(end)
	if err != nil {
		panic(fmt.Sprintf("failed to deserialize end %v: %v", end, err))
	}

	return uint(intmath.Abs(sNode.X-eNode.X) + intmath.Abs(sNode.Y-eNode.Y))
}

func (h *aStarHelper) NeighborDistance(n1, n2 string) uint {
	if n1 == n2 {
		return 0
	} else {
		return 1
	}
}

func (h *aStarHelper) GoalReached(cand, goal string) bool {
	cNode, err := node.Deserialize(cand)
	if err != nil {
		panic(fmt.Sprintf("failed to deserialize cand %v: %v", cand, err))
	}

	return cNode.X == h.width-1 && cNode.Y == h.height-1
}

func RunMaze(width, height int, passcode string) (found bool, lastStep string) {
	helper := NewHelper(width, height)

	start := node.New(0, 0, passcode)
	end := node.New(width-1, height-1, passcode)

	steps := astar.AStar(start.Serialize(), end.Serialize(), helper)
	if steps == nil {
		return false, ""
	}

	lastNode, err := node.Deserialize(steps[0])
	if err != nil {
		panic(fmt.Sprintf("failed to deserialize last step %v", steps[0]))
	}

	return true, lastNode.Path
}
