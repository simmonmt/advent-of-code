package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/lineio"
	"github.com/simmonmt/aoc/2024/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func parseInput(lines []string) ([]int, error) {
	nums, err := lineio.NumbersFromLine(lines[0], " ")
	if err != nil {
		return nil, err
	}

	out := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		out[i] = int(nums[i])
	}
	return out, nil
}

func splitDigits(in int) (num, left, right int) {
	if in == 0 {
		panic("shouldn't happen")
	}

	v := in
	num, left, right = 0, 0, 0
	div := 1
	for v > 0 {
		v /= 10
		num++
		if num%2 == 0 {
			div *= 10
		}
	}

	left = in / div
	right = in % div
	return
}

func transform(in []int) []int {
	out := []int{}
	for _, n := range in {
		if n == 0 {
			out = append(out, 1)
		} else if num, left, right := splitDigits(n); num%2 == 0 {
			out = append(out, left, right)
		} else {
			new := n * 2024
			if new < n {
				panic("overflow")
			}
			out = append(out, new)
		}
	}
	return out
}

func solveA(input []int) int {
	in := make([]int, len(input))
	copy(in, input)

	for i := 0; i < 25; i++ {
		new := transform(in)
		in = new
	}
	return len(in)
}

type Pair struct {
	Num   int
	Level int
}

func solveMemo(in int, level, maxLevel int, memo map[Pair]int) int {
	p := Pair{in, level}
	if v, found := memo[p]; found {
		return v
	}

	if level == maxLevel {
		memo[p] = 1
		return 1
	}

	desc := []int{}
	if in == 0 {
		desc = append(desc, 1)
	} else if num, left, right := splitDigits(in); num%2 == 0 {
		desc = append(desc, left, right)
	} else {
		desc = append(desc, in*2024)
	}

	sum := 0
	for _, d := range desc {
		sum += solveMemo(d, level+1, maxLevel, memo)
	}
	memo[p] = sum
	return sum
}

func solveB(input []int) int {
	sum := 0
	memo := map[Pair]int{}
	for _, n := range input {
		sum += solveMemo(n, 0, 75, memo)
	}
	return sum
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
