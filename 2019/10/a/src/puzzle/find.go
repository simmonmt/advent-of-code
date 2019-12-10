package puzzle

import (
	"github.com/simmonmt/aoc/2019/common/logger"
)

func FindNumVisible(ctr Pos, all map[Pos]bool) int {
	visible := map[Pos]bool{}

	for p := range all {
		if p.X == ctr.X && p.Y == ctr.Y {
			continue
		}

		slope := Pos{
			X: p.X - ctr.X,
			Y: p.Y - ctr.Y,
		}

		var simplified Pos
		simplified.Y, simplified.X = Factor(slope.Y, slope.X)

		logger.LogF("considering %+v from %+v slope %d/%d simp %d/%d",
			p, ctr, slope.Y, slope.X, simplified.Y, simplified.X)
		visible[simplified] = true
	}

	return len(visible)
}

func FindBest(all map[Pos]bool) (Pos, int) {
	var bestPos Pos
	bestVisible := -1
	for ctr := range all {
		numVisible := FindNumVisible(ctr, all)
		logger.LogF("%+v visible %v\n", ctr, numVisible)
		if bestVisible == -1 || bestVisible < numVisible {
			bestPos = ctr
			bestVisible = numVisible
		}
	}
	return bestPos, bestVisible
}
