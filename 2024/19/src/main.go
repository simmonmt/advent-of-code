package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strings"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/lineio"
	"github.com/simmonmt/aoc/2024/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Input struct {
	Avail map[byte][]string
	Lines []string
}

func parseInput(lines []string) (*Input, error) {
	groups := lineio.BlankSeparatedGroups(lines)
	if len(groups) != 2 {
		return nil, fmt.Errorf("bad num groups: %d", len(groups))
	}

	if len(groups[0]) != 1 {
		return nil, fmt.Errorf("bad avail group size: %d", len(groups[0]))
	}

	avails := map[byte][]string{}
	for _, avail := range strings.Split(groups[0][0], ", ") {
		avails[avail[0]] = append(avails[avail[0]], avail)
	}
	for _, a := range avails {
		sort.Slice(a, func(i, j int) bool { return len(a[i]) < len(a[j]) })
	}

	return &Input{Avail: avails, Lines: groups[1]}, nil
}

func canSolve(line string, avail map[byte][]string, quickStop bool) int {
	if len(line) == 0 {
		return 1
	}

	sum := 0
	for _, cand := range avail[line[0]] {
		if strings.HasPrefix(line, cand) {
			sum += canSolve(line[len(cand):], avail, quickStop)
			if quickStop && sum > 0 {
				return sum
			}
		}
	}
	return sum
}

type Match struct {
	frag  string
	start int
	num   int
}

func canSolveAtIndex(line string, start int, avail map[byte][]string, quickStop bool) []Match {
	candMatches := []Match{}
	for _, arr := range avail {
		for _, cand := range arr {
			for i, r := range cand {
				if byte(r) != line[start] {
					continue
				}

				if strings.HasPrefix(line[start-i:], cand) {
					candMatches = append(candMatches, Match{frag: cand, start: start - i})
				}
			}
		}
	}

	matches := []Match{}
	for _, match := range candMatches {
		if num := canSolve(line[match.start+len(match.frag):], avail, quickStop); num > 0 {
			match.num = num
			matches = append(matches, match)
		}
	}
	return matches
}

func solve(line string, avail map[byte][]string, quickStop bool) int {
	matches := canSolveAtIndex(line, len(line)*2/3, avail, quickStop)
	if len(matches) == 0 {
		return 0
	}

	start := len(line) / 3
	newMatches := []Match{}
	for _, match := range matches {
		if match.start <= start {
			newMatches = append(newMatches, match)
			continue
		}

		for _, newMatch := range canSolveAtIndex(line[0:match.start], start, avail, quickStop) {
			newMatch.num *= match.num
			newMatches = append(newMatches, newMatch)
		}
	}
	if len(newMatches) == 0 {
		return 0
	}
	matches = newMatches

	sum := 0
	for _, match := range matches {
		if num := canSolve(line[0:match.start], avail, quickStop); num > 0 {
			sum += match.num * num
			if quickStop {
				return 1
			}
		}
	}
	return sum
}

func solveA(input *Input) int {
	solveable := 0
	for _, line := range input.Lines {
		fmt.Println(line)
		if solve(line, input.Avail, true) > 0 {
			solveable++
		}
	}
	return solveable
}

func solveB(input *Input) int {
	solutions := 0
	for _, line := range input.Lines {
		num := solve(line, input.Avail, false)
		fmt.Println(line, num)
		solutions += num
	}
	return solutions
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
