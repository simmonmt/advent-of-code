package solver

import (
	"astar"
	"bytes"
	"fmt"

	"grid"
	"intmath"
	"logger"
)

type nodeHash struct {
	w, h uint8
	m    map[string][]byte
}

func newNodeHash(w, h uint8) *nodeHash {
	return &nodeHash{
		w: w,
		h: h,
		m: map[string][]byte{},
	}
}

func (nh *nodeHash) Add(g *grid.Grid) string {
	h := g.Hash()
	ser := g.Serialize()
	if nhSer, found := nh.m[h]; found {
		if nhSer != nil && !bytes.Equal(ser, nhSer) {
			panic("collision in hash")
		}
	} else {
		logger.LogF("nodehash adding %v\n", []byte(h))
		nh.m[h] = ser
		if len(nh.m)%10000 == 0 {
			fmt.Printf("nh sz %v\n", len(nh.m))
		}
	}
	return h
}

func (nh *nodeHash) Get(h string) *grid.Grid {
	ser, found := nh.m[h]
	if !found {
		return nil
	}

	g, err := grid.Deserialize(nh.w, nh.h, ser)
	if err != nil {
		return nil
	}

	return g
}

func (nh *nodeHash) Shed(h string) {
	if _, found := nh.m[h]; !found {
		panic("to shed not found")
	}
	logger.LogF("nodehash shedding %v\n", []byte(h))
	nh.m[h] = nil
}

type aStarHelper struct {
	width, height uint8
	nodes         *nodeHash
}

func newHelper(width, height uint8) *aStarHelper {
	return &aStarHelper{
		width:  width,
		height: height,
		nodes:  newNodeHash(width, height),
	}
}

type transferDesc struct {
	toX, toY uint8
}

func findTransfers(g *grid.Grid, x, y uint8) []transferDesc {
	curNode := g.Get(x, y)

	transfers := make([]transferDesc, 4)
	nTransfers := 0

	if x > 0 {
		if n := g.Get(x-1, y); n != nil && n.Avail() > curNode.Used {
			transfers[nTransfers] = transferDesc{x - 1, y}
			nTransfers++
		}
	}
	if x < g.Width()-1 {
		if n := g.Get(x+1, y); n != nil && n.Avail() > curNode.Used {
			transfers[nTransfers] = transferDesc{x + 1, y}
			nTransfers++
		}
	}
	if y > 0 {
		if n := g.Get(x, y-1); n != nil && n.Avail() > curNode.Used {
			transfers[nTransfers] = transferDesc{x, y - 1}
			nTransfers++
		}
	}
	if y < g.Height()-1 {
		if n := g.Get(x, y+1); n != nil && n.Avail() > curNode.Used {
			transfers[nTransfers] = transferDesc{x, y + 1}
			nTransfers++
		}
	}

	return transfers[0:nTransfers]
}

func (h *aStarHelper) AllNeighbors(gStr string) []string {
	g := h.nodes.Get(gStr)
	if g == nil {
		panic("allneighbors failed to find g")
	}

	//fmt.Printf("g hash %v w %v h %v\n", []byte(gStr), g.Width(), g.Height())

	neighbors := []string{}
	for y := 0; y < int(h.height); y++ {
		for x := 0; x < int(h.width); x++ {
			n := g.Get(uint8(x), uint8(y))
			if n.Used == 0 {
				continue
			}
			transfers := findTransfers(g, uint8(x), uint8(y))

			for _, t := range transfers {
				neighbor := g.Transfer(uint8(x), uint8(y), t.toX, t.toY)
				neighbors = append(neighbors, h.nodes.Add(neighbor))
			}
		}
	}

	return neighbors
}

func (h *aStarHelper) EstimateDistance(startStr, endStr string) uint {
	start := h.nodes.Get(startStr)
	if start == nil {
		panic("failed to find start")
	}
	end := h.nodes.Get(endStr)
	if end == nil {
		panic("failed to find end")
	}

	startGoalX, startGoalY := start.Goal()
	endGoalX, endGoalY := end.Goal()
	return uint(intmath.Abs(int(startGoalX)-int(endGoalX)) +
		intmath.Abs(int(startGoalY)-int(endGoalY)))
}

func (h *aStarHelper) NeighborDistance(n1, n2 string) uint {
	if n1 == n2 {
		return 0
	} else {
		return 1
	}
}

func (h *aStarHelper) GoalReached(candStr, goalStr string) bool {
	cand := h.nodes.Get(candStr)
	if cand == nil {
		panic("failed to find cand")
	}

	x, y := cand.Goal()
	return x == 0 && y == 0
}

func (h *aStarHelper) PrintableKey(key string) string {
	// g, err := grid.Deserialize(h.width, h.height, []byte(key))
	// if err != nil {
	// 	panic("can't deserialize")
	// }

	// // x, y := g.Goal()
	// // return fmt.Sprintf("goal=%v,%v", x, y)
	// return g.String()
	return "key"
}

func (h *aStarHelper) MarkClosed(key string) {
	h.nodes.Shed(key)
}

func Solve(width, height uint8, start *grid.Grid) (found bool, numSteps int) {
	end := start.Duplicate()
	end.SetGoal(0, 0)

	helper := newHelper(width, height)
	helper.nodes.Add(start)
	helper.nodes.Add(end)

	steps := astar.AStar(start.Hash(), end.Hash(), helper)
	if steps == nil {
		return false, 0
	}

	return true, len(steps) - 1
}
