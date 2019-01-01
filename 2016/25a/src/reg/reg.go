package reg

import (
	"fmt"
	"strings"
)

type Reg int

const (
	A Reg = iota
	B
	C
	D
)

func (r Reg) String() string {
	switch r {
	case A:
		return "a"
	case B:
		return "b"
	case C:
		return "c"
	case D:
		return "d"
	default:
		return "UNKNOWN"
	}
}

func FromString(name string) (Reg, error) {
	switch strings.ToLower(name) {
	case "a":
		return A, nil
	case "b":
		return B, nil
	case "c":
		return C, nil
	case "d":
		return D, nil
	default:
		return A, fmt.Errorf("unknown reg '%v'", name)
	}
}
