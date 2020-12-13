package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func solveA(startTS int, busses []int) {
	earliestTS := -1
	earliestNum := 0
	for _, bus := range busses {
		if bus == -1 {
			continue
		}

		next := ((startTS + (bus - 1)) / bus) * bus
		if earliestTS == -1 || next < earliestTS {
			earliestTS = next
			earliestNum = bus
		}
	}

	fmt.Printf("A: %d at %d => %d\n", earliestNum, earliestTS,
		earliestNum*(earliestTS-startTS))
}

type FirstEvery struct {
	first int64
	every int64
}

func findPairAlignment(a, b FirstEvery, off int64) FirstEvery {
	first := int64(-1)

	for ts := a.first; ; ts += a.every {
		if (ts-b.first+off)%b.every == 0 {
			if first == -1 {
				logger.LogF("found first %v", ts)
				first = ts
			} else {
				return FirstEvery{first, ts - first}
			}
		}
	}
}

func findAlignment(busses []int) FirstEvery {
	logger.LogF("busses: %v", busses)

	cum := FirstEvery{0, int64(busses[0])}
	for off := 1; off < len(busses); off++ {
		if busses[off] == -1 {
			continue
		}

		b := FirstEvery{0, int64(busses[off])}
		cum = findPairAlignment(cum, b, int64(off))
	}

	return cum
}

func solveB(busses []int) {
	fmt.Printf("B: %v\n", findAlignment(busses).first)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	startTS := intmath.AtoiOrDie(lines[0])

	busses := []int{}
	for _, str := range strings.Split(lines[1], ",") {
		num := -1
		if str != "x" {
			num = intmath.AtoiOrDie(str)
		}
		busses = append(busses, num)
	}

	solveA(startTS, busses)
	solveB(busses)
}
