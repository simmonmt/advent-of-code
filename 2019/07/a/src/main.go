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

	amp "github.com/simmonmt/aoc/2019/07/a/src/amp"
	vm "github.com/simmonmt/aoc/2019/07/a/src/vm"
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

func tryPhaseCombination(phases []int, ram vm.Ram) int {
	signal := 0
	for i, phase := range phases {
		out, err := amp.Run(phase, signal, ram.Clone())
		if err != nil {
			log.Fatalf("amp %d crashed: %v", i, err)
		}

		signal = out
	}
	return signal
}

func tryAllPhases(ram vm.Ram) ([]int, int) {
	var phases, maxPhases [5]int
	maxResult := 0
	for {
		var phaseCounts [5]bool
		hasRepeats := false
		for _, v := range phases {
			if phaseCounts[v] {
				hasRepeats = true
				break
			} else {
				phaseCounts[v] = true
			}
		}
		if !hasRepeats {
			//fmt.Printf("trying %v\n", phases)
			result := tryPhaseCombination(phases[:], ram)
			//fmt.Printf("phase %v, result %v\n", phases, result)
			if result > maxResult {
				maxResult = result
				maxPhases = phases
			}
		}

		var i int
		for i = 0; i < len(phases); i++ {
			phases[i]++
			if phases[i] == 5 {
				phases[i] = 0
			} else {
				break
			}
		}
		if i == len(phases) {
			break
		}
	}

	return maxPhases[:], maxResult
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *ramPath == "" {
		log.Fatalf("--ram is required")
	}

	ram, err := readRam(*ramPath)
	if err != nil {
		log.Fatal(err)
	}

	if *input != "" {
		inputValues, err := parseInput(*input)
		if err != nil {
			log.Fatal(err)
		}

		in := 0
		for i, phase := range inputValues {
			out, err := amp.Run(phase, in, ram.Clone())
			if err != nil {
				log.Fatalf("amp %d crashed: %v", i, err)
			}

			fmt.Printf("amp %d phase %d in %d out %d\n", i, phase, in, out)
			in = out
		}
	} else {
		seq, result := tryAllPhases(ram)
		fmt.Printf("seq %v result %v\n", seq, result)
	}
}
