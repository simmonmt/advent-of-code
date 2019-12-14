package pos

import (
	"fmt"
	"strconv"
	"strings"
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
