package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/simmonmt/aoc/2025/common/filereader"
	"github.com/simmonmt/aoc/2025/common/lineio"
	"github.com/simmonmt/aoc/2025/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Input struct {
	Lines []string
}

func parseInput(lines []string) (*Input, error) {
	expLen := len(lines[0])
	for i, line := range lines {
		if got := len(line); got != expLen {
			return nil, fmt.Errorf("%d: bad len, want %d got %d", i+1, expLen, got)
		}
	}

	return &Input{Lines: lines}, nil
}

func solveA(input *Input) int {
	groups := [][]int{}
	for i := 0; i < len(input.Lines)-1; i++ {
		grp, err := lineio.NumbersFromLine(input.Lines[i], " ")
		if err != nil {
			panic(fmt.Sprintf("%d: %v", i+1, err))
		}
		groups = append(groups, grp)
	}

	ops := []rune{}
	for _, str := range strings.Split(input.Lines[len(input.Lines)-1], " ") {
		if str == "" {
			continue
		}
		ops = append(ops, rune(str[0]))
	}

	tot := 0
	for i, op := range ops {
		sub := 0
		for j := range groups {
			n := groups[j][i]
			if j == 0 {
				sub = n
				continue
			}
			if op == '+' {
				sub += n
			} else {
				sub *= n
			}
		}
		tot += sub
	}
	return tot
}

func extractArgs(lines []string, start, lim int) []int {
	out := []int{}
	for i := lim - 1; i >= start; i-- {
		accum := 0
		for j := range lines {
			if r := lines[j][i]; r != ' ' {
				accum = 10*accum + int(r-'0')
			}
		}
		out = append(out, accum)
	}
	return out
}

func processOp(op rune, args []int) int {
	accum := args[0]
	for i := 1; i < len(args); i++ {
		if op == '+' {
			accum += args[i]
		} else {
			accum *= args[i]
		}
	}
	return accum
}

func solveB(input *Input) int {
	opsStr := input.Lines[len(input.Lines)-1]
	ops := map[int]rune{}
	opLocs := []int{}
	for i, r := range opsStr {
		if r != ' ' {
			opLocs = append(opLocs, i)
			ops[i] = r
		}
	}

	tot := 0
	argLines := input.Lines[0 : len(input.Lines)-1]
	for i, loc := range opLocs {
		lim := len(opsStr)
		if i != len(opLocs)-1 {
			lim = opLocs[i+1] - 1
		}

		args := extractArgs(argLines, loc, lim)
		sub := processOp(ops[loc], args)
		//fmt.Println("op", loc, lim, args, sub)
		tot += sub
	}
	return tot
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
