package lib

import "fmt"

type Char struct {
	Num   int
	IsElf bool
	P     Pos
	HP    int
	AP    int
}

func (c Char) String() string {
	t := "Elf"
	if !c.IsElf {
		t = "Gob"
	}

	return fmt.Sprintf("#%d: %s HP:%3d AP:%d %+v", c.Num, t, c.HP, c.AP, c.P)
}

func (c Char) Short() rune {
	if c.IsElf {
		return 'E'
	}
	return 'G'
}

func NewChar(num int, isElf bool, pos Pos) *Char {
	return &Char{
		Num:   num,
		IsElf: isElf,
		P:     pos,
		HP:    200,
		AP:    3,
	}
}

type CharByReadingOrder []Char

func (a CharByReadingOrder) Len() int      { return len(a) }
func (a CharByReadingOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a CharByReadingOrder) Less(i, j int) bool {
	if a[i].P.Y != a[j].P.Y {
		return a[i].P.Y < a[j].P.Y
	}
	return a[i].P.X < a[j].P.X
}
