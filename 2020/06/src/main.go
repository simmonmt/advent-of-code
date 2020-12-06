package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func readGroups(path string) ([][]string, error) {
	lines, err := readInput(path)
	if err != nil {
		return nil, err
	}
	// So we don't have to special-case the loop end
	lines = append(lines, "")

	groups := [][]string{}
	curGroup := []string{}
	for _, line := range lines {
		if line == "" {
			if len(curGroup) > 0 {
				groups = append(groups, curGroup)
			}
			curGroup = []string{}
			continue
		}

		curGroup = append(curGroup, line)
	}

	return groups, nil
}

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

	groups, err := readGroups(*input)
	if err != nil {
		log.Fatal(err)
	}

	solveA(groups)
	solveB(groups)
}
