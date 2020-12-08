package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	topPattern = regexp.MustCompile(
		`^([^ ]+ [^ ]+) bags contain (.*)$`)
	contentsPattern = regexp.MustCompile(
		`^([0-9]+) ([^ ]+ [^ ]+) bags?\.?$`)
)

type Rule struct {
	name     string
	contents map[string]int
}

func (r *Rule) String() string {
	return fmt.Sprintf("%v", *r)
}

func parseRules(path string) ([]*Rule, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	rules := []*Rule{}
	for lineno, line := range lines {
		topParts := topPattern.FindStringSubmatch(line)
		if topParts == nil {
			return nil, fmt.Errorf("%d: bad top parse", lineno)
		}

		topName := topParts[1]
		contents := topParts[2]

		rule := &Rule{name: topName, contents: map[string]int{}}

		if contents != "no other bags." {
			bagStrs := strings.Split(contents, ", ")
			for _, bagStr := range bagStrs {
				parts := contentsPattern.FindStringSubmatch(bagStr)
				if parts == nil {
					return nil, fmt.Errorf(
						"%d: bad bag parse: %v", lineno,
						bagStr)
				}

				num := intmath.AtoiOrDie(parts[1])
				bagName := parts[2]

				rule.contents[bagName] = num
			}
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func solveA(rules []*Rule) {
	containedBy := map[string]map[string]bool{}

	for _, rule := range rules {
		for bag := range rule.contents {
			if _, found := containedBy[bag]; !found {
				containedBy[bag] = map[string]bool{}
			}
			containedBy[bag][rule.name] = true
		}
	}

	start := "shiny gold"
	done := map[string]bool{}
	cands := map[string]bool{}
	for bag := range containedBy[start] {
		cands[bag] = true
	}

	for len(cands) > 0 {
		newCands := map[string]bool{}
		for cand := range cands {
			for container := range containedBy[cand] {
				if _, found := done[container]; !found {
					newCands[container] = true
				}
			}
		}

		for cand := range cands {
			done[cand] = true
		}
		cands = newCands
	}

	fmt.Printf("A: %d bags\n", len(done))
}

func solveB(rules []*Rule) {
	// numByColor[c] = the number of bags contained within a bag of color c
	// including the color c bag itself.
	numByColor := map[string]int{}
	for len(numByColor) != len(rules) {
		for _, rule := range rules {
			allPrereqsFound := true
			total := 1
			for bag, num := range rule.contents {
				if _, found := numByColor[bag]; !found {
					allPrereqsFound = false
					break
				}
				total += num * numByColor[bag]
			}

			if !allPrereqsFound {
				continue
			}

			numByColor[rule.name] = total
		}
	}

	fmt.Printf("B: %d bags\n", numByColor["shiny gold"]-1)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	rules, err := parseRules(*input)
	if err != nil {
		log.Fatal(err)
	}

	solveA(rules)
	solveB(rules)
}
