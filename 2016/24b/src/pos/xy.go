package pos

import "intmath"

type XY struct {
	X, Y int
}

func (p XY) Eq(o XY) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p XY) Dist(o XY) int {
	return intmath.Abs(p.X-o.X) + intmath.Abs(p.Y-o.Y)
}
