package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2025/common/filereader"
	"github.com/simmonmt/aoc/2025/common/lineio"
	"github.com/simmonmt/aoc/2025/common/logger"
	"github.com/simmonmt/aoc/2025/common/ranges"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Input struct {
	Ranges []ranges.IncRange
	Ings   []int
}

func parseInput(lines []string) (*Input, error) {
	groups := lineio.BlankSeparatedGroups(lines)
	if len(groups) != 2 {
		return nil, fmt.Errorf("wanted 2 groups, got %d", len(groups))
	}

	rs := []ranges.IncRange{}
	for i, rStr := range groups[0] {
		fStr, tStr, ok := strings.Cut(rStr, "-")
		if !ok {
			return nil, fmt.Errorf("bad range %d: no -", i+1)
		}

		var err error
		var r ranges.IncRange
		if r.From, err = strconv.Atoi(fStr); err != nil {
			return nil, fmt.Errorf("bad range %d: bad from %v: %v", i+1, fStr, err)
		}
		if r.To, err = strconv.Atoi(tStr); err != nil {
			return nil, fmt.Errorf("bad range %d: bad to %v: %v", i+1, tStr, err)
		}

		rs = append(rs, r)
	}

	ings, err := lineio.Numbers(groups[1])
	if err != nil {
		return nil, fmt.Errorf("bad ings: %v", err)
	}

	return &Input{rs, ings}, nil
}

func solveA(input *Input) int {
	tot := 0
	for _, ing := range input.Ings {
		for _, r := range input.Ranges {
			if r.Contains(ing) {
				tot++
				break
			}
		}
	}
	return tot
}

func solveB(input *Input) int {
	rs := []ranges.IncRange{}
	for _, r := range input.Ranges {
		rs = append(rs, r)
	}

	sort.Slice(rs, func(i, j int) bool {
		return rs[i].From < rs[j].From
	})

	for {
		changed := false
		for i := 0; i < len(rs)-1; i++ {
			if rs[i].Overlaps(rs[i+1]) {
				changed = true
				r2 := []ranges.IncRange{}
				if i > 0 {
					r2 = append(r2, rs[0:i]...)
				}

				m, _ := rs[i].Merge(rs[i+1])
				r2 = append(r2, m)

				if i+2 < len(rs) {
					r2 = append(r2, rs[i+2:]...)
				}
				rs = r2
				break
			}
		}
		//fmt.Println(len(rs), changed)
		if !changed {
			break
		}
	}

	tot := 0
	for _, r := range rs {
		tot += r.To - r.From + 1
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
