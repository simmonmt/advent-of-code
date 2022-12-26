// Copyright 2022 Google LLC
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
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func fromSnafu(in string) int {
	out := 0
	for _, r := range in {
		var n int
		switch r {
		case '2':
			n = 2
		case '1':
			n = 1
		case '0':
			n = 0
		case '-':
			n = -1
		case '=':
			n = -2
		}
		out = out*5 + n
	}
	return out
}

func toSnafu(in int) string {
	out := ""
	var carry int
	for in > 0 {
		r := in%5 + carry
		var d string
		switch r {
		case 0:
			d, carry = "0", 0
		case 1:
			d, carry = "1", 0
		case 2:
			d, carry = "2", 0
		case 3:
			d, carry = "=", 1
		case 4:
			d, carry = "-", 1
		case 5:
			d, carry = "0", 1
		}

		out = d + out
		in /= 5
	}
	if carry == 1 {
		out = "1" + out
	}

	return out
}

func solveA(lines []string) string {
	sum := 0
	for _, line := range lines {
		sum += fromSnafu(line)
	}

	logger.LogF("translating %v", sum)
	return toSnafu(sum)
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

	fmt.Println("A", solveA(lines))
}
