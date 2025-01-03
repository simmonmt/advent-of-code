package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Line struct {
	Total int
	Nums  []int
}

func parseInput(lines []string) ([]Line, error) {
	out := []Line{}
	for i, line := range lines {
		totStr, rest, ok := strings.Cut(line, ": ")
		if !ok {
			return nil, fmt.Errorf("%d: no total", i+1)
		}

		total, err := strconv.ParseInt(totStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%d: bad total: %v", i+1, err)
		}

		nums, err := filereader.ParseNumbersFromLine(rest, " ")
		if err != nil {
			return nil, fmt.Errorf("%d: bad nums: %v", i+1, err)
		}

		out = append(out, Line{Total: int(total), Nums: nums})
	}

	return out, nil
}

type Op int

const (
	OP_ADD Op = iota
	OP_MUL
	OP_CONCAT
)

func canSolve(want, cur int, ops []Op, left []int) bool {
	for _, op := range ops {
		nextCur := -1
		switch op {
		case OP_ADD:
			nextCur = cur + left[0]
		case OP_MUL:
			nextCur = cur * left[0]
		case OP_CONCAT:
			nextCur = cur
			shift := 0
			for v := left[0]; v > 0; v /= 10 {
				shift++
			}
			for i := 0; i < shift; i++ {
				nextCur *= 10
			}
			nextCur += left[0]
		default:
			panic("bad op")
		}

		rest := left[1:]
		if len(rest) == 0 {
			if nextCur == want {
				return true
			}
		} else if nextCur > want {
			continue
		} else if canSolve(want, nextCur, ops, rest) {
			return true
		}

	}

	return false
}

func solve(input []Line, ops []Op) int {
	sum := 0
	for _, line := range input {
		if canSolve(line.Total, line.Nums[0], ops, line.Nums[1:]) {
			sum += line.Total
		}
	}

	return sum
}

func solveA(input []Line) int {
	return solve(input, []Op{OP_ADD, OP_MUL})
}

func solveB(input []Line) int {
	return solve(input, []Op{OP_ADD, OP_MUL, OP_CONCAT})
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
