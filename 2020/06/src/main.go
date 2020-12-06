package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func solveA(groups [][]string) {
	sum := 0
	for i, group := range groups {
		answered := map[rune]bool{}
		for _, person := range group {
			for _, q := range person {
				answered[q] = true
			}
		}

		logger.LogF("group %d answered %d", i+1, len(answered))
		sum += len(answered)
	}

	fmt.Printf("A sum is %d\n", sum)
}

func solveB(groups [][]string) {
	sum := 0
	for i, group := range groups {
		answered := map[rune]int{}
		for _, person := range group {
			for _, q := range person {
				answered[q]++
			}
		}

		numAllAnswered := 0
		for _, num := range answered {
			if num == len(group) {
				numAllAnswered++
			}
		}

		logger.LogF("group %d all answered %d", i+1, numAllAnswered)
		sum += numAllAnswered
	}

	fmt.Printf("B sum is %d\n", sum)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	groups, err := filereader.BlankSeparatedGroups(*input)
	if err != nil {
		log.Fatal(err)
	}

	solveA(groups)
	solveB(groups)
}
