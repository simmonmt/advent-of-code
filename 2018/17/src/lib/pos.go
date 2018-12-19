package lib

import (
	"fmt"
	"strconv"
	"strings"
)

type Pos struct {
	X, Y int
}

func (p Pos) Eq(s Pos) bool {
	return p.X == s.X && p.Y == s.Y
}

func PosFromString(str string) (Pos, error) {
	parts := strings.SplitN(str, ",", 2)
	if len(parts) != 2 {
		return Pos{}, fmt.Errorf("invalid pos %v", str)
	}

	x, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return Pos{}, fmt.Errorf("bad x in %v", str)
	}

	y, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return Pos{}, fmt.Errorf("bad y in %v", str)
	}

	return Pos{int(x), int(y)}, nil
}
