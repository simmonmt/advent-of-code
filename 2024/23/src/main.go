package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"sort"
	"strings"

	"github.com/simmonmt/aoc/2024/common/collections"
	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Input struct {
	Map map[string]map[string]bool
}

func parseInput(lines []string) (*Input, error) {
	out := map[string]map[string]bool{}

	add := func(a, b string) {
		if _, found := out[a]; !found {
			out[a] = map[string]bool{}
		}
		out[a][b] = true
	}

	for i, line := range lines {
		a, b, ok := strings.Cut(line, "-")
		if !ok {
			return nil, fmt.Errorf("%d: bad cut", i+1)
		}
		add(a, b)
		add(b, a)
	}
	return &Input{out}, nil
}

func solveA(input *Input) int {
	triples := map[string]bool{}

	for host1, connMap := range input.Map {
		conns := collections.MapKeys(connMap)
		for i := 0; i < len(conns)-1; i++ {
			host2 := conns[i]
			h2Conns := input.Map[host2]

			for j := i + 1; j < len(conns); j++ {
				host3 := conns[j]
				if _, found := h2Conns[host3]; found {
					hosts := []string{host1, host2, host3}
					sort.Strings(hosts)
					triples[strings.Join(hosts, ",")] = true
				}
			}
		}
	}

	num := 0
	for triple := range triples {
		if triple[0] == 't' || strings.Contains(triple, ",t") {
			num++
		}
	}
	return num
}

func solveB(input *Input) string {
	curSet := map[string]bool{}
	for name := range input.Map {
		curSet[name] = true
	}
	curSetSize := 1

	for {
		nextSet := map[string]bool{}
		for setStr := range curSet {
			set := strings.Split(setStr, ",")

			for name, conns := range input.Map {
				if slices.Contains(set, name) {
					continue
				}

				// can name be added to set?
				connected := true
				for _, member := range set {
					if !conns[member] {
						connected = false
						break
					}
				}

				if !connected {
					continue
				}

				set = append(set, name)
				sort.Strings(set)
				nextSet[strings.Join(set, ",")] = true
				break
			}

		}

		logger.Infof("processed %d size %d clusters; found %d size %d",
			len(curSet), curSetSize, len(nextSet), curSetSize+1)
		if len(nextSet) == 0 {
			return collections.OneMapKey(curSet, "")
		}

		curSet, curSetSize = nextSet, curSetSize+1
	}
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
