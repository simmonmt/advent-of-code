package maze

import (
	"fmt"
	"sort"
	"unicode"

	"pos"

	"github.com/soniakeys/graph"
)

var (
	dirs = []pos.XY{pos.XY{-1, 0}, pos.XY{1, 0}, pos.XY{0, -1}, pos.XY{0, 1}}
)

type Board struct {
	cells     [][]bool
	h, w      int
	al        graph.LabeledAdjacencyList
	numsByPos map[pos.XY]int
	numPosns  map[int]pos.XY
}

func NewBoard(lines []string) *Board {
	h := len(lines)
	w := len(lines[0])
	cells := make([][]bool, h)
	for y := range lines {
		cells[y] = make([]bool, w)
	}

	board := &Board{
		cells:     cells,
		h:         h,
		w:         w,
		al:        make([][]graph.Half, h*w),
		numsByPos: map[pos.XY]int{},
		numPosns:  map[int]pos.XY{},
	}

	for y := range lines {
		for x, c := range lines[y] {
			p := pos.XY{x, y}
			if c == '.' {
				board.Set(p, true)
			} else if unicode.IsDigit(c) {
				digit := int(c - '0')
				board.numsByPos[p] = digit
				board.numPosns[digit] = p
				board.Set(p, true)
			}
		}
	}

	return board
}

func (b *Board) graph() graph.LabeledAdjacencyList {
	return graph.LabeledAdjacencyList(b.al)
}

func (b *Board) Get(p pos.XY) bool {
	return b.cells[p.Y][p.X]
}

func (b *Board) Set(p pos.XY, open bool) {
	if p.X < 0 || p.X >= b.w || p.Y < 0 || p.Y >= b.h {
		panic("out of range")
	}

	if b.cells[p.Y][p.X] == open {
		panic("already there")
	}

	for _, dir := range dirs {
		cand := pos.XY{p.X + dir.X, p.Y + dir.Y}
		if cand.X < 0 || cand.X >= b.w || cand.Y < 0 || cand.Y >= b.h {
			continue
		}

		if !b.Get(cand) {
			continue // can't route to closed
		}

		if !open {
			b.removeEdge(p, cand)
			b.removeEdge(cand, p)
		} else {
			b.addEdge(p, cand)
			b.addEdge(cand, p)
		}
	}

	b.cells[p.Y][p.X] = open
}

func (b *Board) addEdge(from, to pos.XY) {
	fromIdx := int(b.posToNI(from))
	toNI := b.posToNI(to)

	for _, h := range b.al[fromIdx] {
		if h.To == toNI {
			return // already there
		}
	}

	b.al[fromIdx] = append(b.al[fromIdx], graph.Half{To: toNI})
}

func (b *Board) removeEdge(from, to pos.XY) {
	fromIdx := int(b.posToNI(from))
	toNI := b.posToNI(to)

	out := []graph.Half{}
	for _, h := range b.al[fromIdx] {
		if h.To != toNI {
			out = append(out, h)
		}
	}
	b.al[fromIdx] = out
}

func (b *Board) posToNI(p pos.XY) graph.NI {
	return graph.NI(p.Y*b.w + p.X)
}

func (b *Board) niToPos(ni graph.NI) pos.XY {
	val := int(ni)
	y := val / b.w
	x := val % b.w
	return pos.XY{x, y}
}

func (b *Board) Nums() []int {
	out := []int{}
	for num := range b.numPosns {
		out = append(out, num)
	}
	sort.Ints(out)
	return out
}

func (b *Board) ShortestPath(from, to int) (int, bool) {
	var found bool
	var fromPos, toPos pos.XY
	if fromPos, found = b.numPosns[from]; !found {
		panic("from unfound")
	}
	if toPos, found = b.numPosns[to]; !found {
		panic("to unfound")
	}

	// if b.Get(fromPos) {
	// 	panic("from is routable")
	// }
	// if b.Get(toPos) {
	// 	panic("to is routable")
	// }

	// b.Set(fromPos, true)
	// b.Set(toPos, true)

	weight := func(cand graph.LI) float64 { return 1 }
	heuristic := func(cand graph.NI) float64 {
		return float64(b.niToPos(cand).Dist(toPos))
	}

	fullpath, _ := b.graph().AStarAPath(
		b.posToNI(fromPos), b.posToNI(toPos), heuristic, weight)

	// b.Set(fromPos, false)
	// b.Set(toPos, false)

	if fullpath.Path == nil {
		return 0, false
	}

	return len(fullpath.Path), true
}

func (b *Board) Dump() {
	for y := range b.cells {
		for x, c := range b.cells[y] {
			p := pos.XY{x, y}

			if num, found := b.numsByPos[p]; found {
				fmt.Print(num)
			} else if c {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}
