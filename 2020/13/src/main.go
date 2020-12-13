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

func findCommon(a, b int64, off int64) FirstEvery {
	first := int64(-1)

	for ts := a; ; ts += a {
		if (ts+off)%b == 0 {
			if first == -1 {
				logger.LogF("found first %v", ts)
				first = ts
			} else {
				return FirstEvery{first, ts - first}
			}
		}
	}
}

func findAlignment(a, b FirstEvery) FirstEvery {
	first := int64(-1)

	for ts := a.first; ; ts += a.every {
		if (ts-b.first)%b.every == 0 {
			if first == -1 {
				logger.LogF("found first %v", ts)
				first = ts
			} else {
				return FirstEvery{first, ts - first}
			}
		}
	}
}

func findCumAlignment(busses []int) FirstEvery {
	logger.LogF("busses: %v", busses)

	alignments := []FirstEvery{}

	a := int64(busses[0])
	for off := 1; off < len(busses); off++ {
		if busses[off] == -1 {
			continue
		}

		b := int64(busses[off])

		logger.LogF("findCommon %v %v %v", a, b, off)
		alignment := findCommon(a, b, int64(off))

		logger.LogF("a %v b %v off %v first %v every %v",
			a, b, off, alignment.first, alignment.every)

		alignments = append(alignments, alignment)
	}

	logger.LogF("alignments %v", alignments)
	cumAlignment := alignments[0]
	for i := 1; i < len(alignments); i++ {
		newAlignment := findAlignment(cumAlignment, alignments[i])
		logger.LogF("alignment %v %v = %v",
			cumAlignment, alignments[i], newAlignment)
		cumAlignment = newAlignment
	}

	return cumAlignment
}

func solveB(busses []int) {
	fmt.Printf("B: %v\n", findCumAlignment(busses).first)
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
