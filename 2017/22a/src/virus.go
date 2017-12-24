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

func dirLeftRight(dir int) (left, right int) {
	switch dir {
	case DirNorth:
		return DirWest, DirEast
	case DirSouth:
		return DirEast, DirWest
	case DirEast:
		return DirNorth, DirSouth
	case DirWest:
		return DirSouth, DirNorth
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
	elems map[Point]bool
}

func newGraph() *Graph {
	return &Graph{map[Point]bool{}}
}

func (g *Graph) IsInfected(point Point) bool {
	if infected, found := g.elems[point]; found {
		return infected
	}
	return false
}

func (g *Graph) SetInfected(point Point, infected bool) {
	g.elems[point] = infected
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
	dirLeft, dirRight := dirLeftRight(c.dir)
	isInfected := c.graph.IsInfected(c.pos)

	var nextDir int
	if isInfected {
		nextDir = dirRight
	} else {
		nextDir = dirLeft
	}

	c.graph.SetInfected(c.pos, !isInfected)

	c.dir = nextDir
	c.pos = nextPos(c.pos, nextDir)

	return !isInfected
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

			contents := '.'
			if c.graph.IsInfected(point) {
				contents = '#'
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
				graph.SetInfected(Point{i, y}, true)
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
