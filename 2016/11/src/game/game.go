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

package game

import (
	"fmt"

	"astar"
	"board"
)

type aStarHelper struct{}

func (h *aStarHelper) AllNeighbors(start string) []string {
	b, err := board.Deserialize(start)
	if err != nil {
		panic(fmt.Sprintf("failed to deserialize %v", start))
	}
	moves := b.AllMoves()

	neighbors := []string{}
	for _, move := range moves {
		nb := b.Apply(move)
		neighbors = append(neighbors, nb.Serialize())
	}
	return neighbors
}

func (h *aStarHelper) Estimate(start, end string) uint {
	if start == end {
		return 0
	}
	return 1
}

func (h *aStarHelper) NeighborDistance(n1, n2 string) uint {
	if n1 == n2 {
		return 0
	}
	return 1
}

func Play(b *board.Board) []string {
	helper := &aStarHelper{}
	return astar.AStar(b.Serialize(), b.SuccessBoard().Serialize(),
		helper, helper)
}
