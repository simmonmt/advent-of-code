package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	DirNorth = iota
	DirSouth
	DirEast
	DirWest

	StateClean = iota
	StateWeakened
	StateInfected
	StateFlagged
)

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func nextDirs(dir int) (left, right, back int) {
	switch dir {
	case DirNorth:
		return DirWest, DirEast, DirSouth
	case DirSouth:
		return DirEast, DirWest, DirNorth
	case DirEast:
		return DirNorth, DirSouth, DirWest
	case DirWest:
		return DirSouth, DirNorth, DirEast
	default:
		panic("unknown direction")
	}
}

func nextPos(pos Point, dir int) Point {
	switch dir {
	case DirNorth:
		return Point{pos.X, pos.Y - 1}
	case DirSouth:
		return Point{pos.X, pos.Y + 1}
	case DirEast:
		return Point{pos.X + 1, pos.Y}
	case DirWest:
		return Point{pos.X - 1, pos.Y}
	default:
		panic("unknown direction")
	}
}

type Point struct {
	X, Y int
}

type Graph struct {
	elems map[Point]int
}

func newGraph() *Graph {
	return &Graph{map[Point]int{}}
}

func (g *Graph) GetState(point Point) int {
	if state, found := g.elems[point]; found {
		return state
	}
	return StateClean
}

func (g *Graph) SetState(point Point, state int) {
	g.elems[point] = state
}

func (g *Graph) Bounds() (Point, Point) {
	var minX, maxX, minY, maxY int
	for point := range g.elems {
		minX, maxX = point.X, point.X
		minY, maxY = point.Y, point.Y
		break
	}

	for point := range g.elems {
		minX = min(minX, point.X)
		minY = min(minY, point.Y)
		maxX = max(maxX, point.X)
		maxY = max(maxY, point.Y)
	}

	return Point{minX, minY}, Point{maxX, maxY}
}

type Carrier struct {
	graph *Graph
	pos   Point
	dir   int
}

func newCarrier(graph *Graph, pos Point) *Carrier {
	return &Carrier{
		graph: graph,
		pos:   pos,
		dir:   DirNorth,
	}
}

func (c *Carrier) Advance() bool {
	dirLeft, dirRight, dirBack := nextDirs(c.dir)
	state := c.graph.GetState(c.pos)

	nextDir := c.dir
	nextState := state
	switch state {
	case StateClean:
		nextDir = dirLeft
		nextState = StateWeakened
		break
	case StateWeakened:
		nextState = StateInfected
		break
	case StateInfected:
		nextDir = dirRight
		nextState = StateFlagged
		break
	case StateFlagged:
		nextDir = dirBack
		nextState = StateClean
		break
	}

	c.graph.SetState(c.pos, nextState)

	c.dir = nextDir
	c.pos = nextPos(c.pos, nextDir)

	return nextState == StateInfected
}

func (c *Carrier) Dump(out io.Writer) {
	//fmt.Println(c.graph)

	minBound, maxBound := c.graph.Bounds()
	// fmt.Printf("min %+v max %+v\n", minBound, maxBound)

	dist := max(maxBound.X-minBound.X+1, maxBound.Y-minBound.Y+1)
	maxBound.X = minBound.X + dist + 3
	minBound.X -= 3
	maxBound.Y = minBound.Y + dist + 3
	minBound.Y -= 3

	// fmt.Printf("dist %v min %+v max %+v\n", dist, minBound, maxBound)

	for y := minBound.Y; y <= maxBound.Y; y++ {
		for x := minBound.X; x <= maxBound.X; x++ {
			point := Point{x, y}

			contents := '?'
			switch c.graph.GetState(point) {
			case StateClean:
				contents = '.'
				break
			case StateWeakened:
				contents = 'W'
				break
			case StateInfected:
				contents = '#'
				break
			case StateFlagged:
				contents = 'F'
				break
			default:
				panic("unknown state")
			}

			if point == c.pos {
				fmt.Fprintf(out, "[%c]", contents)
			} else {
				fmt.Fprintf(out, " %c ", contents)
			}
		}
		fmt.Fprintf(out, "\n")
	}
}

func readInput(in io.Reader) (*Graph, Point) {
	graph := newGraph()

	lineLen := -1
	y := 0
	reader := bufio.NewReader(in)
	for ; ; y++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if lineLen < 0 {
			lineLen = len(line)
		}

		for i, c := range line {
			if c == '#' {
				graph.SetState(Point{i, y}, StateInfected)
			}
		}
	}

	center := Point{lineLen / 2, y / 2}

	return graph, center
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v niters", os.Args[1])
	}

	nIters, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("invalid niters %v: %v", os.Args[1], err)
	}

	graph, center := readInput(os.Stdin)

	carrier := newCarrier(graph, center)
	// carrier.Dump(os.Stdout)

	numInfections := 0
	for i := 1; i <= nIters; i++ {
		causedInfection := carrier.Advance()
		if causedInfection {
			numInfections++
		}

		// fmt.Println()
		// fmt.Printf("iteration %v:\n", i)
		// carrier.Dump(os.Stdout)
		// fmt.Printf("caused infections: %v\n", numInfections)
	}

	// fmt.Println(*graph)
	// fmt.Printf("center is %+v\n", center)
	fmt.Printf("total caused infections: %v\n", numInfections)
}
