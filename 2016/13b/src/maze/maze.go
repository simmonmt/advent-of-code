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
	"fmt"
	"strconv"
	"strings"

	"astar"
	"intmath"
	"logger"
)

type aStarHelper struct {
	magicNumber int
}

func parsePosition(str string) (int, int) {
	parts := strings.SplitN(str, ",", 2)

	x, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		panic(fmt.Sprintf("failed to parse x in %v", str))
	}

	y, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		panic(fmt.Sprintf("failed to parse y in %v", str))
	}

	return int(x), int(y)
}

func serializePosition(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func countBits(val uint32) int {
	num := 0
	for i := 0; i < 32; i++ {
		if val&1 == 1 {
			num++
		}
		val = val >> 1
	}
	return num
}

func IsOpenSpace(magicNumber int, x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}

	val := x*x + 3*x + 2*x*y + y + y*y
	val += magicNumber
	if val < 0 {
		panic(fmt.Sprintf("unexpected negative for %v,%v", x, y))
	}

	return countBits(uint32(val))%2 == 0
}

func tryAddMove(magicNumber int, x, y int, moves *[][2]int) {
	if x < 0 || y < 0 {
		return
	}

	if IsOpenSpace(magicNumber, x, y) {
		*moves = append(*moves, [2]int{x, y})
	}
}

func (h *aStarHelper) AllNeighbors(start string) []string {
	x, y := parsePosition(start)

	moves := [][2]int{}
	tryAddMove(h.magicNumber, x, y-1, &moves)
	tryAddMove(h.magicNumber, x-1, y, &moves)
	tryAddMove(h.magicNumber, x+1, y, &moves)
	tryAddMove(h.magicNumber, x, y+1, &moves)

	moveStrs := make([]string, len(moves))
	for i, move := range moves {
		moveStrs[i] = serializePosition(move[0], move[1])
	}

	logger.LogF("neighbors for %v: %v\n", start, moveStrs)
	return moveStrs
}

func (h *aStarHelper) Estimate(start, end string) uint {
	startX, startY := parsePosition(start)
	endX, endY := parsePosition(end)

	return uint(intmath.Abs(startX-endX) + intmath.Abs(startY-endY))
}

func (h *aStarHelper) NeighborDistance(n1, n2 string) uint {
	if n1 == n2 {
		return 0
	}
	return 1
}

func WalkMaze(magicNumber int, startX, startY int, goalX, goalY int) []string {
	helper := &aStarHelper{magicNumber: magicNumber}
	positions := astar.AStar(
		serializePosition(startX, startY),
		serializePosition(goalX, goalY),
		helper, helper)

	return positions
}
