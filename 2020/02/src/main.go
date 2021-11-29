// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	inputPattern = regexp.MustCompile(`^([0-9]+)-([0-9]+) ([a-z]): (.*)$`)
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

func matchesOld(from, to int, char rune, password string) bool {
	numMatched := 0
	for _, c := range password {
		if c == char {
			numMatched++
		}
	}

	return numMatched >= from && numMatched <= to
}

func matchesNew(from, to int, char rune, password string) bool {
	pr := []rune(password)
	return (pr[from-1] == char) != (pr[to-1] == char)
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

	numValidOld, numValidNew := 0, 0
	for _, line := range lines {
		parts := inputPattern.FindStringSubmatch(line)
		if parts == nil {
			log.Fatalf("bad input %v", line)
		}

		from := intmath.AtoiOrDie(parts[1])
		to := intmath.AtoiOrDie(parts[2])
		char := []rune(parts[3])[0]
		password := parts[4]

		if matchesOld(from, to, char, password) {
			numValidOld++
		}
		if matchesNew(from, to, char, password) {
			numValidNew++
		}
	}

	fmt.Printf("num valid old: %v\n", numValidOld)
	fmt.Printf("num valid new: %v\n", numValidNew)
}
