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
	"strings"

	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/vm"
)

var (
	verbose   = flag.Bool("verbose", false, "verbose")
	ramPath   = flag.String("ram", "", "path to file containing ram values")
	inputPath = flag.String("input", "", "input file")
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
		line = strings.TrimSpace(strings.Split(line, "#")[0])
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func makeIO(lines []string) *vm.SaverIO {
	arr := []int64{}
	for _, line := range lines {
		for _, r := range line {
			arr = append(arr, int64(r))
		}
		arr = append(arr, 10) // \n
	}

	return vm.NewSaverIO(arr...)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *ramPath == "" {
		log.Fatalf("--ram is required")
	}

	ram, err := vm.NewRamFromFile(*ramPath)
	if err != nil {
		log.Fatal(err)
	}

	if *inputPath == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*inputPath)
	if err != nil {
		log.Fatal(err)
	}

	io := makeIO(lines)

	if err := vm.Run(ram, io); err != nil {
		log.Fatal(err)
	}

	for _, v := range io.Written() {
		if v > 255 {
			fmt.Printf("\n\noutput value %v\n\n", v)
		} else {
			fmt.Printf("%c", rune(v))
		}
	}
}
