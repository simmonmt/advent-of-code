// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grid

import (
	"crypto/md5"
	"fmt"
	"node"
)

type Grid struct {
	w, h         uint8
	nodes        []node.Node
	goalX, goalY uint8
}

func New(w, h uint8, goalX, goalY uint8, nodes []node.Node) (*Grid, error) {
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

func (g *Grid) Duplicate() *Grid {
	ng := &Grid{
		w:     g.w,
		h:     g.h,
		nodes: make([]node.Node, len(g.nodes)),
		goalX: g.goalX,
		goalY: g.goalY,
	}
	copy(ng.nodes, g.nodes)
	return ng
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

func Deserialize(w, h uint8, ser []byte) (*Grid, error) {
	if len(ser) < 2 {
		return nil, fmt.Errorf("too little input")
	}

	goalX := ser[0]
	goalY := ser[1]
	serIdx := 2

	nodes := make([]node.Node, int(w)*int(h))
	nodeNum := 0
	for y := 0; y < int(h); y++ {
		for x := 0; x < int(w); x++ {
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
	for y := 0; y < int(g.h); y++ {
		for x := 0; x < int(g.w); x++ {
			n := g.nodes[y*int(g.w)+x]

			bracket := [2]rune{' ', ' '}
			if y == 0 && x == 0 {
				bracket = [2]rune{'(', ')'}
			} else if y == int(g.goalY) && x == int(g.goalX) {
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
			if i%int(g.w) == 0 {
				out += "|"
			} else {
				out += ","
			}
		}
		out += n.String()
	}

	return out + "]"
}

func (g *Grid) SetGoal(x, y uint8) {
	g.goalX = x
	g.goalY = y
}

func (g *Grid) Goal() (uint8, uint8) {
	return g.goalX, g.goalY
}

func (g *Grid) Width() uint8 {
	return g.w
}

func (g *Grid) Height() uint8 {
	return g.h
}

func (g *Grid) Get(x, y uint8) *node.Node {
	var idx int = int(y)*int(g.w) + int(x)
	return &g.nodes[idx]
}

func (g *Grid) Transfer(srcX, srcY, destX, destY uint8) *Grid {
	ng := g.Duplicate()

	src := ng.Get(srcX, srcY)
	dest := ng.Get(destX, destY)
	if dest.Avail() < src.Used {
		panic(fmt.Sprintf("unable to transfer %+v to %+v", src, dest))
	}

	dest.Used += src.Used
	src.Used = 0

	if srcX == ng.goalX && srcY == ng.goalY {
		ng.goalX = destX
		ng.goalY = destY
	}

	return ng
}

func (g *Grid) Hash() string {
	h := md5.Sum(g.Serialize())
	return string(h[:])
}
