package xypos

import (
	"fmt"
	"strconv"
	"strings"
)

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

func Parse(str string) (Pos, error) {
	parts := strings.SplitN(str, ",", 2)
	if len(parts) != 2 {
		return Pos{}, fmt.Errorf("bad pos %v", str)
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return Pos{}, fmt.Errorf("bad pos x: %v", err)
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return Pos{}, fmt.Errorf("bad pos y: %v", err)
	}

	return Pos{x, y}, nil
}
