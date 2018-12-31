package main

import "intmath"

type Pos struct {
	X, Y int
}

func (p Pos) Eq(o Pos) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p Pos) Dist(o Pos) int {
	return intmath.Abs(p.X-o.X) + intmath.Abs(p.Y-o.Y)
}
