package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2025/common/filereader"
	"github.com/simmonmt/aoc/2025/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Range struct {
	From, To int
}

type Input struct {
	Ranges []Range
}

func parseInput(lines []string) (*Input, error) {
	out := &Input{}

	for _, str := range strings.Split(lines[0], ",") {
		fs, ts, _ := strings.Cut(str, "-")
		var err error

		var r Range
		r.From, err = strconv.Atoi(fs)
		if err != nil {
			return nil, fmt.Errorf("bad from %v: %v", fs, err)
		}
		r.To, err = strconv.Atoi(ts)
		if err != nil {
			return nil, fmt.Errorf("bad to %v: %v", ts, err)
		}

		out.Ranges = append(out.Ranges, r)
	}

	return out, nil
}

func isAValid(n int) bool {
	s := strconv.Itoa(n)
	if len(s)%2 != 0 {
		return true
	}

	if string(s[0:len(s)/2]) != string(s[len(s)/2:]) {
		return true
	}

	return false
}

func isBValid(n int) bool {
	s := []byte(strconv.Itoa(n))

	for i := 1; i <= len(s)/2; i++ {
		if len(s)%i != 0 {
			continue
		}

		first := s[0:i]
		all := true
		for j := 1; j < len(s)/i; j++ {
			this := s[i*j : i*(j+1)]
			if !slices.Equal(first, this) {
				all = false
				break
			}
		}

		if all {
			return false
		}
	}

	return true
}

func invalidsInRange(r Range, validator func(n int) bool) []int {
	out := []int{}

	for i := r.From; i <= r.To; i++ {
		if !validator(i) {
			out = append(out, i)
		}
	}

	return out
}

func solveA(input *Input) int {
	out := 0
	for _, r := range input.Ranges {
		for _, inv := range invalidsInRange(r, isAValid) {
			out += inv
		}
	}
	return out
}

func solveB(input *Input) int {
	out := 0
	for _, r := range input.Ranges {
		for _, inv := range invalidsInRange(r, isBValid) {
			out += inv
		}
	}
	return out
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
