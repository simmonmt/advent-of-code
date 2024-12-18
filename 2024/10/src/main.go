package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

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

func parseInput(lines []string) (*grid.Grid[byte], error) {
	return grid.NewFromLines(lines, func(p pos.P2, r rune) (byte, error) {
		if r == '.' {
			return 99, nil
		}
		return byte(r - '0'), nil
	})
}

func scoreTrailheadDFS(input *grid.Grid[byte], p pos.P2, peaks map[pos.P2]bool, soFar map[pos.P2]bool) int {
	soFar[p] = true
	defer delete(soFar, p)

	pv, _ := input.Get(p)
	if pv == 9 {
		peaks[p] = true
		return 1
	}

	sum := 0
	for _, n := range input.AllNeighbors(p, false) {
		if _, found := soFar[n]; found {
			continue
		}

		nv, _ := input.Get(n)

		if pv+1 != nv {
			continue
		}

		sum += scoreTrailheadDFS(input, n, peaks, soFar)
	}
	return sum
}

func solveA(input *grid.Grid[byte]) int64 {
	sum := 0
	input.Walk(func(p pos.P2, b byte) {
		if b == 0 {
			peaks := map[pos.P2]bool{}
			scoreTrailheadDFS(input, p, peaks, map[pos.P2]bool{})
			sum += len(peaks)
		}
	})
	return int64(sum)
}

func solveB(input *grid.Grid[byte]) int64 {
	sum := 0
	input.Walk(func(p pos.P2, b byte) {
		if b == 0 {
			sum += scoreTrailheadDFS(input, p,
				map[pos.P2]bool{}, map[pos.P2]bool{})
		}
	})
	return int64(sum)
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
