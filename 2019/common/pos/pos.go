package pos

type P2 struct {
	X, Y int
}

func (p *P2) Equals(o P2) bool {
	return p.X == o.X && p.Y == o.Y
}
