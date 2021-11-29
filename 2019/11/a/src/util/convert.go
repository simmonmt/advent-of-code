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
	"strconv"
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

	code := []int{}

	// Emit preamble
	code = append(code, 11105, 1, 7) // 0: jt 1,7
	code = append(code, 11104, 99)   // 3: out 99 (an invalid value)
	code = append(code, 99)          // 5: hlt
	code = append(code, 0)           // 6: scratch

	scratchAddr := 6
	failAddr := 3
	addr := len(code)

	addInstruction := func(b ...int) {
		code = append(code, b...)
		addr += len(b)
	}

	for _, line := range lines {
		parts := strings.Split(line, " ")
		in := parts[1]
		outs := parts[3:5]

		addInstruction(3, scratchAddr) // in *scratchAddr
		if in == "0" {
			// we expect 0, so fail if non-zero
			// jt *scratchAddr failAddr
			addInstruction(1005, scratchAddr, failAddr)
		} else {
			// we expect non-zero, so fail if zero
			// jf *scratchAddr failAddr
			addInstruction(1006, scratchAddr, failAddr)
		}

		for _, out := range outs {
			val := 0
			if out == "1" {
				val = 1
			}
			addInstruction(11104, val)
		}
	}
	addInstruction(99)

	for i, v := range code {
		if i != 0 {
			fmt.Print(",")
		}
		fmt.Print(strconv.Itoa(v))
	}

	fmt.Println()
}
