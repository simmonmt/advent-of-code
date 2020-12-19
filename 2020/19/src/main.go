package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/simmonmt/aoc/2020/19/src/parse"
	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

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

	rules := []string{}
	for _, line := range lines {
		if line == "" {
			break
		}
		rules = append(rules, line)
	}

	patStr, err := parse.Parse(rules, 0)
	if err != nil {
		log.Fatal(err)
	}

	pat, err := regexp.Compile(patStr)
	if err != nil {
		log.Fatal(err)
	}

	messages := lines[len(rules)+1:]
	numMatches := 0
	for _, message := range messages {
		if pat.MatchString(message) {
			numMatches++
		}
	}

	fmt.Printf("result: %v matches\n", numMatches)
}
