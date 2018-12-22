package main

type Pos struct {
	X, Y int
}

func (p Pos) Eq(o Pos) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p Pos) Before(o Pos) bool {
	if p.Y < o.Y {
		return true
	} else if p.Y > o.Y {
		return false
	} else {
		return p.X < o.X
	}
}
