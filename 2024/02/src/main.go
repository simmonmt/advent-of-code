package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
	"github.com/simmonmt/aoc/2024/common/mtsmath"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func parseInput(lines []string) ([][]int, error) {
	out := [][]int{}

	for i, line := range lines {
		parts := strings.Split(line, " ")
		nums := make([]int, len(parts))

		for j, s := range parts {
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("%d: bad num: %v", i+1, err)
			}

			nums[j] = n
		}

		out = append(out, nums)
	}

	return out, nil
}

func checkReport(report []int) bool {
	if report[1] == report[0] {
		return false
	}
	inc := report[1] > report[0]

	for i := 1; i < len(report); i++ {
		diff := report[i] - report[i-1]
		if diff == 0 {
			return false
		}
		if diff > 0 != inc {
			return false
		}
		if mtsmath.Abs(diff) > 3 {
			return false
		}
	}
	return true
}

func solveA(reports [][]int) int64 {
	sum := 0
	for _, report := range reports {
		if checkReport(report) {
			sum++
		}
	}

	return int64(sum)
}

// 620 too low
func solveB(reports [][]int) int64 {
	sum := 0
	for _, report := range reports {
		if checkReport(report) {
			sum++
			continue
		}

		foundTrue := false
		for i := 0; i < len(report); i++ {
			sub := make([]int, len(report)-1)
			k := 0
			for j := 0; j < len(report); j++ {
				if j != i {
					sub[k] = report[j]
					k++
				}
			}

			if checkReport(sub) {
				foundTrue = true
				break
			}
		}

		if foundTrue {
			sum++
			continue
		}
	}

	return int64(sum)
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

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
