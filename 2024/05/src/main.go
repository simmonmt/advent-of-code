package main

import (
	"flag"
	"fmt"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Input struct {
	Befores map[int]map[int]bool
	Updates [][]int
}

func parseBefores(lines []string) (befores map[int]map[int]bool, err error) {
	befores = map[int]map[int]bool{}

	for _, line := range lines {
		var a, b int
		if _, err = fmt.Sscanf(line, "%d|%d", &a, &b); err != nil {
			return
		}

		if _, ok := befores[a]; !ok {
			befores[a] = map[int]bool{}
		}
		befores[a][b] = true
	}

	return
}

func parseUpdates(lines []string) (updates [][]int, err error) {
	updates = [][]int{}

	for _, line := range lines {
		nums, err := filereader.ParseNumbersFromLine(line, ",")
		if err != nil {
			return nil, err
		}

		updates = append(updates, nums)
	}

	return
}

func parseInput(lines []string) (input *Input, err error) {
	groups := filereader.BlankSeparatedGroupsFromLines(lines)
	if len(groups) != 2 {
		return nil, fmt.Errorf("want two groups; found %d", len(groups))
	}

	input = &Input{}
	if input.Befores, err = parseBefores(groups[0]); err != nil {
		return nil, err
	}
	if input.Updates, err = parseUpdates(groups[1]); err != nil {
		return nil, err
	}

	return
}

func isSorted(in []int, befores map[int]map[int]bool) bool {
	idxs := map[int]int{}
	for i, n := range in {
		idxs[n] = i
	}

	for i, n := range in {
		for before := range befores[n] {
			if bi, found := idxs[before]; found && i > bi {
				logger.Infof("bad: %v because %v after %v",
					in, n, before)
				return false
			}
		}
	}

	return true
}

func solveA(input *Input) int {
	sum := 0
	for _, update := range input.Updates {
		if isSorted(update, input.Befores) {
			logger.Infof("ok: %v", update)
			sum += update[len(update)/2]
		}
	}

	return sum
}

func resort(in []int, befores map[int]map[int]bool) []int {
	logger.Infof("resorting %v", in)

	out := []int{in[0]}

	for inIdx := 1; inIdx < len(in); inIdx++ {
		toInsert := in[inIdx]
		logger.Infof("inserting %v", toInsert)

		for outIdx := 0; outIdx <= len(out); outIdx++ {
			cand := []int{}
			cand = append(cand, out[:outIdx]...)
			cand = append(cand, toInsert)
			cand = append(cand, out[outIdx:]...)

			logger.Infof("trying %v", cand)

			if isSorted(cand, befores) {
				logger.Infof("success")
				out = cand
				break
			}
		}
	}

	return out
}

func solveB(input *Input) int {
	sum := 0
	for _, update := range input.Updates {
		if isSorted(update, input.Befores) {
			continue
		}

		sorted := resort(update, input.Befores)
		sum += sorted[len(sorted)/2]
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

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
