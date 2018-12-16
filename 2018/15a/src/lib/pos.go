package lib

type Pos struct {
	X, Y int
}

func PosLess(a, b Pos) bool {
	if a.Y != b.Y {
		return a.Y < b.Y
	}
	return a.X < b.X
}

type PosByReadingOrder []Pos

func (a PosByReadingOrder) Len() int           { return len(a) }
func (a PosByReadingOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PosByReadingOrder) Less(i, j int) bool { return PosLess(a[i], a[j]) }
