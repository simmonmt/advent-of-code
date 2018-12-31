package main

import "logger"

// heuristic
//
// moving left: manhattan distance from goal to 0,0
// moving to goal: (ml + manhattan distance empty to goal)

type aStarHelper struct {
	board *Board
}

func NewAStarHelper(board *Board) *aStarHelper {
	return &aStarHelper{board}
}

var (
	dirs = []Pos{
		Pos{-1, 0},
		Pos{1, 0},
		Pos{0, -1},
		Pos{0, 1},
	}
)

func (h *aStarHelper) AllNeighbors(start string) []string {
	ps := Decode(start)

	width, height := h.board.Size()

	outs := []*PlayState{}

	for _, dir := range dirs {
		cand := Pos{
			X: ps.Empty.X + dir.X,
			Y: ps.Empty.Y + dir.Y,
		}

		if cand.X < 0 || cand.Y < 0 || cand.X >= width || cand.Y >= height {
			continue
		}

		if !h.board.IsMoveable(cand) {
			continue
		}

		var out *PlayState
		if cand.Eq(ps.Goal) {
			out = &PlayState{Empty: ps.Goal, Goal: ps.Empty}
		} else {
			out = &PlayState{Empty: cand, Goal: ps.Goal}
		}
		outs = append(outs, out)
	}

	if *verbose {
		logger.LogLn("AllNeighbors in:")
		h.board.Dump(ps)
		logger.LogLn("AllNeighbors out:")
		for i, out := range outs {
			logger.LogF("out %d:", i)
			h.board.Dump(out)
		}
	}

	neighbors := make([]string, len(outs))
	for i, out := range outs {
		neighbors[i] = out.Encode()
	}
	return neighbors
}

func (h *aStarHelper) EstimateDistance(start, goal string) uint {
	startPs := Decode(start)
	goalPs := Decode(goal)

	return uint(startPs.Empty.Dist(startPs.Goal) +
		startPs.Goal.Dist(goalPs.Goal))
}

func (h *aStarHelper) NeighborDistance(n1, n2 string) uint {
	return 1
}

func (h *aStarHelper) GoalReached(cand, goal string) bool {
	ps := Decode(cand)
	return ps.Goal.X == 0 && ps.Goal.Y == 0
}

func (h *aStarHelper) PrintableKey(key string) string {
	return key
}

func (h *aStarHelper) MarkClosed(key string) {
}
