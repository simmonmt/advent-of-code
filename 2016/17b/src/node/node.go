package node

import (
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	X, Y     int
	Passcode string
	Path     string
}

func New(x, y int, passcode string) *Node {
	return &Node{
		X:        x,
		Y:        y,
		Passcode: passcode,
		Path:     "",
	}
}

func Deserialize(ser string) (*Node, error) {
	parts := strings.SplitN(ser, ",", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("illegal node name %v", ser)
	}

	x, xErr := strconv.ParseInt(parts[0], 10, 32)
	y, yErr := strconv.ParseInt(parts[1], 10, 32)
	if xErr != nil || yErr != nil {
		return nil, fmt.Errorf("failed to parse x,y from %v", ser)
	}

	return &Node{
		X:        int(x),
		Y:        int(y),
		Passcode: parts[2],
		Path:     parts[3],
	}, nil
}

func (n *Node) Serialize() string {
	return fmt.Sprintf("%d,%d,%v,%v", n.X, n.Y, n.Passcode, n.Path)
}
