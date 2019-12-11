package puzzle

import (
	"github.com/simmonmt/aoc/2019/common/intmath"
)

type Pos struct {
	X, Y int
}

type ByManhattanOriginDistance []Pos

func (a ByManhattanOriginDistance) Len() int      { return len(a) }
func (a ByManhattanOriginDistance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByManhattanOriginDistance) Less(i, j int) bool {
	iDist := intmath.Abs(a[i].X) + intmath.Abs(a[i].Y)
	jDist := intmath.Abs(a[j].X) + intmath.Abs(a[j].Y)
	return iDist < jDist
}
