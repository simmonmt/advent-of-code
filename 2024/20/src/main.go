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

type Input struct {
	Start, End pos.P2
	Grid       *grid.Grid[rune]
}

func parseInput(lines []string) (*Input, error) {
	var start, end pos.P2

	g, _ := grid.NewFromLines[rune](lines, func(p pos.P2, r rune) (rune, error) {
		if r == 'S' {
			start = p
			r = '.'
		} else if r == 'E' {
			end = p
			r = '.'
		}

		return r, nil
	})

	return &Input{Start: start, End: end, Grid: g}, nil
}

type GraphHelper struct {
	g *grid.Grid[rune]
}

func (h *GraphHelper) Neighbors(id graph.NodeID) []graph.NodeID {
	p, _ := pos.P2FromString(string(id))

	out := []graph.NodeID{}
	for _, n := range h.g.AllNeighbors(p, false) {
		if v, _ := h.g.Get(n); v == '.' {
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

func nodeToPos(n graph.NodeID) pos.P2 {
	p, _ := pos.P2FromString(string(n))
	return p
}

func nodesToPos(ns []graph.NodeID) []pos.P2 {
	out := []pos.P2{}
	for _, n := range ns {
		p, _ := pos.P2FromString(string(n))
		out = append(out, p)
	}
	return out
}

func findBestPath(input *Input) []pos.P2 {
	helper := &GraphHelper{g: input.Grid}

	startID := graph.NodeID(input.Start.String())
	endID := graph.NodeID(input.End.String())

	path := []pos.P2{input.Start}
	path = append(path, nodesToPos(graph.ShortestPath(startID, endID, helper))...)
	return path
}

func findSolutions(input *Input, maxDist int) map[int]int {
	path := findBestPath(input)
	noCheatLen := len(path) - 1

	dists := map[pos.P2]int{}
	for i, p := range path {
		dists[p] = len(path) - 1 - i
	}

	diffs := map[int]int{}
	for i := 0; i < len(path)-1; i++ {
		for j := i + 1; j < len(path); j++ {
			p, n := path[i], path[j]

			pnDist := p.ManhattanDistance(n)
			if pnDist > maxDist {
				continue
			}

			withCheatLen := i + pnDist + dists[n]
			if withCheatLen < noCheatLen {
				diffs[noCheatLen-withCheatLen]++
			}
		}
	}
	return diffs
}

func solve(input *Input, maxDist int) int64 {
	sum := 0
	for savings, num := range findSolutions(input, maxDist) {
		if savings >= 100 {
			sum += num
		}
	}
	return int64(sum)
}

func solveA(input *Input) int64 {
	return solve(input, 2)
}

func solveB(input *Input) int64 {
	return solve(input, 20)
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

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
