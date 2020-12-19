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

func countMatches(pat *regexp.Regexp, messages []string) int {
	numMatches := 0
	for _, message := range messages {
		out := pat.FindStringIndex(message)
		if len(out) > 0 && out[0] == 0 && out[1] == len(message) {
			numMatches++
		}
	}
	return numMatches
}

func solveA(rules, messages []string) int {
	patStr, err := parse.Parse(rules, 0)
	if err != nil {
		log.Fatal(err)
	}

	pat, err := regexp.Compile(patStr)
	if err != nil {
		log.Fatal(err)
	}

	return countMatches(pat, messages)
}

// The only reference to rules 8 and 11 are in rule 0. Rule 0 is 8 11
//
// The rewritten rule 8 is: 42 | 42 8
//    This means "one or more 42's"
// The rewritten rule 11 is: 42 31 | 42 11 31
//    This means "42s followed by 31s, with as many 42s as 31s"
//
// Putting an 8 before an 11 modifies the constraint that the number of
// 42s and 31s be equal, turning it into "more 42s than 31s". So we
// have this constraint:
//
//     42s followed by 31s, with more 42s than 31s
//
// If we ignore rules 0 8 and 11, concentrating on verifying the above
// constraint, we don't have to deal with the rule looping in code. We
// can just verify "42+ 31+" and make sure the counts match the
// constraint.
func solveB(rules, messages []string) int {
	rule42, err := parse.Parse(rules, 42)
	if err != nil {
		log.Fatal(err)
	}
	rule42pat := regexp.MustCompile(rule42)

	rule31, err := parse.Parse(rules, 31)
	if err != nil {
		log.Fatal(err)
	}
	rule31pat := regexp.MustCompile(rule31)

	combined := fmt.Sprintf("^((?:%s)+)((?:%s)+)$", rule42, rule31)
	pat := regexp.MustCompile(combined)

	numMatches := 0
	for _, message := range messages {
		// 42+ 31+ where #42s > #31s

		parts := pat.FindStringSubmatch(message)
		if parts == nil {
			continue
		}

		num42s := len(rule42pat.FindAllString(parts[1], -1))
		num31s := len(rule31pat.FindAllString(parts[2], -1))
		logger.LogF("%v %d %d\n", message, num42s, num31s)

		if num42s <= num31s {
			continue
		}

		numMatches++
	}

	return numMatches
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

	rules := []string{}
	for _, line := range lines {
		if line == "" {
			break
		}
		rules = append(rules, line)
	}
	messages := lines[len(rules)+1:]

	fmt.Printf("A: %v\n", solveA(rules, messages))
	fmt.Printf("B: %v\n", solveB(rules, messages))
}
