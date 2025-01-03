package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2024/common/dir"
	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/grid"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/pos"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func parseInput(lines []string) (*grid.Grid[rune], error) {
	return grid.NewFromLines[rune](lines, func(p pos.P2, r rune) (rune, error) {
		return r, nil
	})
}

type Border struct {
	P pos.P2
	D dir.Dir
}

type Region struct {
	R       rune
	P       []pos.P2
	Borders map[Border]bool
}

func findRegion(g *grid.Grid[rune], start pos.P2, r rune) *Region {
	region := &Region{R: r, P: []pos.P2{}, Borders: map[Border]bool{}}

	seen := map[pos.P2]bool{}
	todo := []pos.P2{start}

	for len(todo) > 0 {
		next := []pos.P2{}

		for _, cur := range todo {
			if _, found := seen[cur]; found {
				continue
			}
			seen[cur] = true
			region.P = append(region.P, cur)

			for _, n := range g.AllNeighbors(cur, false) {
				nr, _ := g.Get(n)
				if nr == r {
					next = append(next, n)
				}
			}
		}

		todo = next
	}

	for _, p := range region.P {
		for _, d := range dir.AllDirs {
			if _, found := seen[d.From(p)]; !found {
				region.Borders[Border{P: p, D: d}] = true
			}
		}
	}

	return region
}

func findRegions(g *grid.Grid[rune]) []*Region {
	regions := []*Region{}
	seen := map[pos.P2]*Region{}

	g.Walk(func(p pos.P2, r rune) {
		if _, found := seen[p]; found {
			return
		}

		region := findRegion(g, p, r)
		for _, p := range region.P {
			seen[p] = region
		}
		regions = append(regions, region)
	})

	return regions
}

func solveA(g *grid.Grid[rune]) int {
	regions := findRegions(g)

	score := 0
	for _, region := range regions {
		area := len(region.P)
		perimeter := len(region.Borders)
		score += area * perimeter
	}
	return score
}

func clearBorders(start Border, travel dir.Dir, left map[Border]bool) {
	for p := travel.From(start.P); ; p = travel.From(p) {
		pb := Border{p, start.D}
		if _, found := left[pb]; !found {
			break
		}
		delete(left, pb)
	}
}

func solveB(g *grid.Grid[rune]) int {
	regions := findRegions(g)

	score := 0
	for _, region := range regions {
		area := len(region.P)

		left := map[Border]bool{}
		for b := range region.Borders {
			left[b] = true
		}

		perimeter := 0

		for len(left) > 0 {
			var b Border
			for b = range left {
				break
			}

			perimeter++
			delete(left, b)

			clearBorders(b, b.D.Left(), left)
			clearBorders(b, b.D.Right(), left)
		}

		logger.Infof("area %v perimeter %v score %v", area, perimeter, area*perimeter)

		score += area * perimeter
	}

	return score
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
