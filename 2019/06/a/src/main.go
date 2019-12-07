package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/simmonmt/aoc/2019/common/logger"
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

func walkDepthFirst(orbitedBy map[string][]string, cur string, level int, callback func(cur string, children []string, level int)) {
	children := orbitedBy[cur]
	callback(cur, children, level)
	for _, child := range children {
		walkDepthFirst(orbitedBy, child, level+1, callback)
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	orbitedBy := map[string][]string{}
	orbits := map[string]string{}

	for _, line := range lines {
		parts := strings.Split(line, ")")
		orbitee := parts[0]
		orbiter := parts[1]

		orbits[orbiter] = orbitee
		if _, found := orbitedBy[orbitee]; !found {
			orbitedBy[orbitee] = []string{}
		}
		orbitedBy[orbitee] = append(orbitedBy[orbitee], orbiter)
	}

	numDirect := 0
	numIndirect := 0

	walkDepthFirst(orbitedBy, "COM", 0,
		func(cur string, children []string, level int) {
			numDirect += len(children)
			if level > 1 {
				numIndirect += level - 1
			}
			//fmt.Printf("%v %v\n", cur, level)
		})

	fmt.Printf("direct %d\n", numDirect)
	fmt.Printf("indirect %d\n", numIndirect)
	fmt.Printf("sum %d\n", numDirect+numIndirect)
}
