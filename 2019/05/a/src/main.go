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

	vm "github.com/simmonmt/aoc/2019/05/a/src/vm"
	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	ramPath = flag.String("ram", "", "path to file containing ram values")
	input   = flag.String("input", "", "input values")
)

func readRam(path string) (vm.Ram, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var line string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	ram := vm.NewRam()
	for i, str := range strings.Split(line, ",") {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", str, err)
		}
		ram.Write(i, val)
	}

	return ram, nil
}

func parseInput(inputStr string) ([]int, error) {
	out := []int{}
	for _, s := range strings.Split(inputStr, ",") {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse input value %v: %v", s, err)
		}

		out = append(out, v)
	}
	return out, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *ramPath == "" {
		log.Fatalf("--ram is required")
	}

	if *input == "" {
		log.Fatalf("--input is required")
	}

	ram, err := readRam(*ramPath)
	if err != nil {
		log.Fatal(err)
	}

	inputValues, err := parseInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	io := vm.NewIO(inputValues...)

	if err := vm.Run(ram, io, 0); err != nil {
		log.Fatalf("program failed: %v", err)
	}

	fmt.Printf("output: %v\n", io.Written())
}
