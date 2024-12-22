package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/graph"
	"github.com/simmonmt/aoc/2024/common/grid"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/pos"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func parseInput(lines []string) ([]pos.P2, error) {
	out := []pos.P2{}
	for i, line := range lines {
		p, err := pos.P2FromString(line)
		if err != nil {
			return nil, fmt.Errorf("%d: bad pos: %v", i+1, err)
		}
		out = append(out, p)
	}

	return out, nil
}

type GraphHelper struct {
	g *grid.Grid[bool]
}

func (h *GraphHelper) Neighbors(id graph.NodeID) []graph.NodeID {
	p, _ := pos.P2FromString(string(id))

	out := []graph.NodeID{}
	for _, n := range h.g.AllNeighbors(p, false) {
		if v, _ := h.g.Get(n); !v {
			out = append(out, graph.NodeID(n.String()))
		}
	}
	return out
}

func (h *GraphHelper) NeighborDistance(from, to graph.NodeID) int {
	p1, _ := pos.P2FromString(string(from))
	p2, _ := pos.P2FromString(string(to))
	return p1.ManhattanDistance(p2)
}

func solveA(input []pos.P2, num int) int64 {
	if len(input) > num {
		input = input[0:num]
	}

	maxX, maxY := input[0].X, input[0].Y
	for _, p := range input {
		maxX = max(maxX, p.X)
		maxY = max(maxY, p.Y)
	}

	g := grid.New[bool](maxX+1, maxY+1)
	for _, p := range input {
		g.Set(p, true)
	}

	// g.Dump(true, func(p pos.P2, v bool, _ bool) string {
	// 	if v {
	// 		return "#"
	// 	}
	// 	return "."
	// })

	start, end := pos.P2{X: 0, Y: 0}, pos.P2{X: maxX, Y: maxY}
	helper := &GraphHelper{g}

	path := graph.ShortestPath(
		graph.NodeID(start.String()),
		graph.NodeID(end.String()), helper)

	return int64(len(path))
}

func solveB(input []pos.P2, num int) string {
	maxX, maxY := input[0].X, input[0].Y
	for _, p := range input {
		maxX = max(maxX, p.X)
		maxY = max(maxY, p.Y)
	}

	start, end := pos.P2{X: 0, Y: 0}, pos.P2{X: maxX, Y: maxY}

	g := grid.New[bool](maxX+1, maxY+1)
	for i := 0; i < num; i++ {
		g.Set(input[i], true)
	}

	path := graph.ShortestPath(
		graph.NodeID(start.String()),
		graph.NodeID(end.String()), &GraphHelper{g})

	makePathNodes := func(path []graph.NodeID) map[pos.P2]bool {
		out := map[pos.P2]bool{}
		for _, ps := range path {
			p, _ := pos.P2FromString(string(ps))
			out[p] = true
		}
		return out
	}

	pathNodes := makePathNodes(path)

	for i := num; i < len(input); i++ {
		p := input[i]
		g.Set(p, true)

		if _, found := pathNodes[p]; !found {
			continue
		}

		// this node interrupted the existing path. See if there's another one.
		path = graph.ShortestPath(
			graph.NodeID(start.String()),
			graph.NodeID(end.String()), &GraphHelper{g})

		if len(path) == 0 {
			return p.String()
		}

		pathNodes = makePathNodes(path)
	}

	return "no"
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input, 1024))
	fmt.Println("B", solveB(input, 1024))
}
