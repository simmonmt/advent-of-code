package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func parseInput(lines []string) ([]int, error) {
	nums := []int{}
	for i, line := range lines {
		num, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("%d: bad num: %v", i+1, err)
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func mix(in, secret int) int {
	return in ^ secret
}

func prune(secret int) int {
	return secret % 16777216
}

func round(secret int) int {
	secret = prune(mix(secret*64, secret))
	secret = prune(mix(secret/32, secret))
	secret = prune(mix(secret*2048, secret))
	return secret
}

func solveA(input []int) int64 {
	sum := 0
	for _, in := range input {
		num := in
		for range 2000 {
			num = round(num)
		}
		sum += num
	}
	return int64(sum)
}

type Elem struct {
	Diff [4]int
}

func makeSeq(secret int) map[Elem]int {
	nums := make([]int, 2000)
	nums[0] = secret
	prices := make([]int, 2000)
	prices[0] = nums[0] % 10
	diffs := make([]int, 2000)

	for i := 1; i < len(nums); i++ {
		secret = round(secret)
		nums[i] = secret
		prices[i] = secret % 10

		diffs[i] = prices[i] - prices[i-1]
	}

	out := map[Elem]int{}
	for i := 4; i < len(diffs); i++ {
		elem := Elem{Diff: [4]int{diffs[i-3], diffs[i-2], diffs[i-1], diffs[i]}}

		if _, found := out[elem]; found {
			continue
		}
		out[elem] = nums[i] % 10
	}
	return out
}

func solveB(input []int) int64 {
	maxNum := -1

	market := map[Elem]int{}
	for _, in := range input {
		for elem, num := range makeSeq(in) {
			num += market[elem]
			if maxNum == -1 || num > maxNum {
				maxNum = num
			}
			market[elem] = num
		}
	}
	return int64(maxNum)
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
