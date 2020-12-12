package pos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2020/common/intmath"
)

func fromString(str string, wantParts int) ([]int, error) {
	parts := strings.Split(str, ",")
	if len(parts) != wantParts {
		return nil, fmt.Errorf("invalid input")
	}
	vs := make([]int, wantParts)
	for i := range vs {
		v, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, fmt.Errorf("invalid coord %v", parts[i])
		}
		vs[i] = v
	}
	return vs, nil
}

type P2 struct {
	X, Y int
}

func P2FromString(str string) (P2, error) {
	vs, err := fromString(str, 2)
	if err != nil {
		return P2{}, err
	}
	return P2{vs[0], vs[1]}, nil
}

func (p *P2) Equals(o P2) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p *P2) Add(o P2) {
	p.X += o.X
	p.Y += o.Y
}

func (p *P2) LessThan(o P2) bool {
	if p.X < o.X {
		return true
	} else if p.X > o.X {
		return false
	} else {
		return p.Y < o.Y
	}
}

func (p P2) ManhattanDistance(o P2) int {
	return intmath.Abs(o.X-p.X) + intmath.Abs(o.Y-p.Y)
}

func (p P2) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func (p P2) AllNeighbors(includeDiag bool) []P2 {
	num := 4
	if includeDiag {
		num = 8
	}

	out := make([]P2, num)
	out[0] = P2{p.X - 1, p.Y}
	out[1] = P2{p.X + 1, p.Y}
	out[2] = P2{p.X, p.Y - 1}
	out[3] = P2{p.X, p.Y + 1}

	if includeDiag {
		out[4] = P2{p.X - 1, p.Y - 1}
		out[5] = P2{p.X + 1, p.Y - 1}
		out[6] = P2{p.X - 1, p.Y + 1}
		out[7] = P2{p.X + 1, p.Y + 1}
	}

	return out
}

type P3 struct {
	X, Y, Z int
}

func (p *P3) Equals(o P3) bool {
	return p.X == o.X && p.Y == o.Y && p.Z == o.Z
}

func P3FromString(str string) (P3, error) {
	vs, err := fromString(str, 3)
	if err != nil {
		return P3{}, err
	}
	return P3{vs[0], vs[1], vs[2]}, nil
}

func (p P3) String() string {
	return fmt.Sprintf("%d,%d,%d", p.X, p.Y, p.Z)
}
