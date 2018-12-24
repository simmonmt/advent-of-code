package xyzpos

import (
	"fmt"
	"strconv"
	"strings"

	"intmath"
)

type Pos struct {
	X, Y, Z int
}

func (p Pos) Eq(o Pos) bool {
	return p.X == o.X && p.Y == o.Y && p.Z == o.Z
}

func (p Pos) Dist(o Pos) int {
	return intmath.Abs(o.X-p.X) + intmath.Abs(o.Y-p.Y) +
		intmath.Abs(o.Z-p.Z)
}

func Parse(str string) (Pos, error) {
	parts := strings.SplitN(str, ",", 3)
	if len(parts) != 3 {
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

	z, err := strconv.Atoi(parts[2])
	if err != nil {
		return Pos{}, fmt.Errorf("bad pos z: %v", err)
	}

	return Pos{x, y, z}, nil
}
