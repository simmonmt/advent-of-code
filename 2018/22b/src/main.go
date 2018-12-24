// 978 too low

package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"intmath"
	"logger"
	"xypos"

	"github.com/soniakeys/graph"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	depth    = flag.Int("depth", -1, "depth")
	target   = flag.String("target", "", "target")
	size     = flag.String("size", "", "size")
	dumpFill = flag.Bool("dump_fill", false, "dump fill")
)

type Fill int

const (
	FILL_ROCKY Fill = iota
	FILL_WET
	FILL_NARROW
)

func (e Fill) String() string {
	switch e {
	case FILL_ROCKY:
		return "rocky"
	case FILL_WET:
		return "wet"
	case FILL_NARROW:
		return "narrow"
	default:
		panic(strconv.Itoa(int(e)))
	}
}

type Equip int

const (
	EQUIP_NEITHER Equip = 0
	EQUIP_CLIMB   Equip = 1
	EQUIP_TORCH   Equip = 2
)

func (e Equip) String() string {
	switch e {
	case EQUIP_NEITHER:
		return "neither"
	case EQUIP_CLIMB:
		return "climb"
	case EQUIP_TORCH:
		return "torch"
	default:
		panic(strconv.Itoa(int(e)))
	}
}

func erosion(geo int) int {
	return (geo + *depth) % 20183
}

func erosionToFill(geo int) Fill {
	switch geo % 3 {
	case 0:
		return FILL_ROCKY
	case 1:
		return FILL_WET
	case 2:
		return FILL_NARROW
	default:
		panic("unknown")
	}
}

func posToNodeId(p xypos.Pos, w int, equip Equip) graph.NI {
	ni := graph.NI(((p.Y*w + p.X) << 2) | int(equip))

	outPos, outEquip := nodeIdToPos(ni, w)
	if !outPos.Eq(p) || equip != outEquip {
		panic(fmt.Sprintf("in pos %v equip %s %d out (%v) pos %v equip %s %d", p, equip, int(equip), int(ni), outPos, outEquip, int(outEquip)))
	}

	return ni
}

func nodeIdToPos(ni graph.NI, w int) (xypos.Pos, Equip) {
	val := int(ni) >> 2
	y := val / w
	x := val % w
	return xypos.Pos{x, y}, Equip(int(ni) & 3)
}

var (
	dirs = []xypos.Pos{xypos.Pos{-1, 0}, xypos.Pos{1, 0}, xypos.Pos{0, -1}, xypos.Pos{0, 1}}

	rockyTools  = []Equip{EQUIP_CLIMB, EQUIP_TORCH}
	wetTools    = []Equip{EQUIP_CLIMB, EQUIP_NEITHER}
	narrowTools = []Equip{EQUIP_TORCH, EQUIP_NEITHER}
)

func allowedTools(fill Fill) []Equip {
	switch fill {
	case FILL_ROCKY:
		return rockyTools
	case FILL_WET:
		return wetTools
	case FILL_NARROW:
		return narrowTools
	default:
		panic("unknown")
	}
}

func addEdge(g [][]graph.Half, w int, fromPos xypos.Pos, fromTool Equip, toPos xypos.Pos, toTool Equip, cost int) graph.LabeledAdjacencyList {
	fromNI := posToNodeId(fromPos, w, fromTool)
	toNI := posToNodeId(toPos, w, toTool)

	h := graph.Half{To: toNI, Label: graph.LI(cost)}

	if g[int(fromNI)] == nil {
		g[int(fromNI)] = []graph.Half{h}
	} else {
		g[int(fromNI)] = append(g[int(fromNI)], h)
	}

	return g
}

func makeGraph(fill [][]Fill, w, h int) graph.LabeledAdjacencyList {
	g := make([][]graph.Half, w*h*4)

	for y := range fill {
		for x := range fill[y] {
			p := xypos.Pos{x, y}
			pType := fill[p.Y][p.X]
			pTools := allowedTools(pType)

			for _, dir := range dirs {
				op := xypos.Pos{p.X + dir.X, p.Y + dir.Y}
				if op.X < 0 || op.Y < 0 {
					continue
				} else if op.X >= len(fill[y]) || op.Y >= len(fill) {
					continue
				}

				opType := fill[op.Y][op.X]
				opTools := allowedTools(opType)

				for _, pTool := range pTools {
					for _, opTool := range opTools {
						cost := 1
						if pTool != opTool {
							cost += 7
						}
						g = addEdge(g, w, p, pTool, op, opTool, cost)
					}
				}
			}
		}
	}

	return graph.LabeledAdjacencyList(g)
}

func dump(start, target xypos.Pos, fill [][]Fill) {
	for y := range fill {
		for x := range fill[y] {
			pos := xypos.Pos{x, y}
			f := fill[y][x]

			var char string
			switch {
			case pos.Eq(start):
				char = "M"
			case pos.Eq(target):
				char = "T"
			case f == FILL_ROCKY:
				char = "."
			case f == FILL_WET:
				char = "="
			case f == FILL_NARROW:
				char = "|"
			default:
				panic("unknown")
			}

			fmt.Print(char)
		}
		fmt.Println()
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *target == "" {
		log.Fatal("--target is required")
	}
	if *depth == -1 {
		log.Fatal("--depth is required")
	}

	size, err := xypos.Parse(*size)
	if err != nil {
		log.Fatal(err)
	}

	target, err := xypos.Parse(*target)
	if err != nil {
		log.Fatal(err)
	}

	if size.X <= target.X || size.Y <= target.Y {
		panic("size too small")
	}

	w := size.X
	h := size.Y

	geo := make([][]int, h)
	for y := range geo {
		geo[y] = make([]int, w)
		for x := range geo[y] {
			pos := xypos.Pos{x, y}
			val := -1

			if target.Eq(pos) {
				val = 0
			} else if x == 0 && y == 0 {
				val = 0
			} else if y == 0 {
				val = 16807 * x
			} else if x == 0 {
				val = 48271 * y
			} else {
				val = erosion(geo[y][x-1]) *
					erosion(geo[y-1][x])
			}

			geo[y][x] = val
		}
	}

	fill := make([][]Fill, h)
	for y := range geo {
		fill[y] = make([]Fill, w)
		for x, g := range geo[y] {
			fill[y][x] = erosionToFill(erosion(g))
		}
	}

	if fill[target.Y][target.X] != FILL_ROCKY {
		panic("not rocky")
	}

	srcPos := xypos.Pos{0, 0}
	srcTool := EQUIP_TORCH

	destPos := target
	destTool := EQUIP_TORCH

	if *dumpFill {
		dump(srcPos, destPos, fill)
	}

	heuristic := func(cand graph.NI) float64 {
		candPos, _ := nodeIdToPos(cand, w)
		return float64(intmath.Abs(candPos.X-destPos.X) + intmath.Abs(candPos.Y-destPos.Y))
	}

	weight := func(cand graph.LI) float64 { return float64(cand) }

	g := makeGraph(fill, w, h)
	path, dist := g.AStarAPath(posToNodeId(srcPos, w, srcTool),
		posToNodeId(destPos, w, destTool),
		heuristic, weight)

	elems := []graph.Half{graph.Half{To: path.Start, Label: graph.LI(float64(0))}}
	for _, h := range path.Path {
		elems = append(elems, h)
	}

	cum := 0
	for _, h := range elems {
		pos, tool := nodeIdToPos(h.To, w)
		cum += int(h.Label)
		fmt.Printf("%-10v %-10s %-10s %d %d\n", fmt.Sprint(pos), fill[pos.Y][pos.X], tool, int(h.Label), cum)
	}

	fmt.Printf("dist=%v\n", dist)
}
