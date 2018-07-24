package grid

import (
	"fmt"
	"node"
)

type Grid struct {
	w, h         int
	nodes        []node.Node
	goalX, goalY uint8
}

func New(w, h int, goalX, goalY uint8, nodes []node.Node) (*Grid, error) {
	if int(w)*int(h) != len(nodes) {
		return nil, fmt.Errorf("w=%v h=%v expected %v nodes, got %v",
			w, h, w*h, len(nodes))
	}

	return &Grid{
		w:     w,
		h:     h,
		nodes: nodes,
		goalX: goalX,
		goalY: goalY,
	}, nil
}

func deserializeNode(out []byte) (*node.Node, int) {
	decode := func(out []byte) (uint16, int) {
		if len(out) < 1 {
			return 0, 0
		}
		val := uint16(out[0])
		if val&0x80 != 0 {
			if len(out) < 2 {
				return 0, 0
			}
			return ((val & 0x7f) << 8) | uint16(out[1]), 2
		} else {
			return val, 1
		}
	}

	outIdx := 0
	size, consumed := decode(out[outIdx:])
	if consumed == 0 {
		return nil, 0
	}
	outIdx += consumed
	used, consumed := decode(out[outIdx:])
	if consumed == 0 {
		return nil, 0
	}
	outIdx += consumed

	return node.New(size, used), outIdx
}

func Deserialize(w, h int, ser []byte) (*Grid, error) {
	if len(ser) < 2 {
		return nil, fmt.Errorf("too little input")
	}

	goalX := ser[0]
	goalY := ser[1]
	serIdx := 2

	nodes := make([]node.Node, w*h)
	nodeNum := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			n, consumed := deserializeNode(ser[serIdx:])
			if n == nil {
				panic("can't deserialize")
			}
			nodes[nodeNum] = *n
			serIdx += consumed
			nodeNum++
		}
	}

	g, err := New(w, h, goalX, goalY, nodes)
	return g, err
}

func (g *Grid) Print() {
	for row := 0; row < g.h; row++ {
		for col := 0; col < g.w; col++ {
			n := g.nodes[row*g.w+col]

			bracket := [2]rune{' ', ' '}
			if row == 0 && col == 0 {
				bracket = [2]rune{'(', ')'}
			} else if row == int(g.goalY) && col == int(g.goalX) {
				bracket = [2]rune{'[', ']'}
			}

			fmt.Printf(" %c%3d/%3d%c", bracket[0], n.Used, n.Size, bracket[1])
		}
		fmt.Println()
	}
}

func serializeNode(n *node.Node, out []byte) int {
	need := 2
	if len(out) < need {
		return 0
	}
	outIdx := 0

	encode := func(val uint16) bool {
		if val < 128 {
			out[outIdx] = uint8(val)
			outIdx++
		} else {
			need++
			if len(out) < need {
				return false
			}

			out[outIdx] = byte((val >> 8) | 0x80)
			out[outIdx+1] = byte(val & 0xff)
			outIdx += 2
		}

		return true
	}

	if !encode(n.Size) || !encode(n.Used) {
		return 0
	}

	return need
}

func (g *Grid) Serialize() []byte {
	out := make([]byte, 2+6*len(g.nodes))
	out[0] = g.goalX
	out[1] = g.goalY
	outNext := 2

	for _, n := range g.nodes {
		if used := serializeNode(&n, out[outNext:]); used == 0 {
			panic("didn't fit")
		} else {
			outNext += used
		}
	}

	return append([]byte(nil), out[0:outNext]...)
}

func (g Grid) String() string {
	out := fmt.Sprintf("%dx%d,g:%dx%d,[", g.w, g.h, g.goalX, g.goalY)

	for i, n := range g.nodes {
		if i > 0 {
			out += ","
		}
		out += n.String()
	}

	return out + "]"
}
