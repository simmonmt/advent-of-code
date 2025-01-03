package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/mtsmath"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) (as, bs []int, err error) {
	as, bs = []int{}, []int{}

	for i, line := range lines {
		var a, b int
		if _, err := fmt.Sscanf(line, "%d %d\n", &a, &b); err != nil {
			return nil, nil, fmt.Errorf("%d bad input: %v", i+1, err)
		}

		as = append(as, a)
		bs = append(bs, b)
	}

	return as, bs, nil
}

func dup(in []int) []int {
	out := make([]int, len(in))
	copy(out, in)
	return out
}

func solveA(as, bs []int) int {
	as, bs = dup(as), dup(bs)

	sort.Slice(as, func(i, j int) bool { return as[i] < as[j] })
	sort.Slice(bs, func(i, j int) bool { return bs[i] < bs[j] })

	sum := 0
	for i := 0; i < len(as); i++ {
		diff := mtsmath.Abs(as[i] - bs[i])
		logger.Infof("diff for %v %v is %v", as[i], bs[i], diff)
		sum += diff
	}

	return sum
}

func solveB(as, bs []int) int {
	bMap := map[int]int{}
	for _, b := range bs {
		bMap[b]++
	}

	sum := 0
	for _, a := range as {
		sum += a * bMap[a]
	}

	return sum
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	as, bs, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(as, bs))
	fmt.Println("B", solveB(as, bs))
}
