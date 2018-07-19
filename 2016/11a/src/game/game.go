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
