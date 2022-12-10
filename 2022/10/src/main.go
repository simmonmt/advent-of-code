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
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Inst struct {
	Name string
	Arg  int
}

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func parseInstructions(lines []string) ([]Inst, error) {
	insts := []Inst{}
	for i, line := range lines {
		parts := strings.SplitN(line, " ", 2)

		inst := Inst{Name: parts[0]}
		if len(parts) == 2 {
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("%d: bad arg: %v", i, err)
			}
			inst.Arg = num
		}

		insts = append(insts, inst)
	}
	return insts, nil
}

type VM struct {
	insts         []Inst
	curInst       *Inst
	curInstCycles int
	reg           int
}

func NewVM(insts []Inst) *VM {
	return &VM{
		insts:         insts,
		curInst:       nil,
		curInstCycles: 0,
		reg:           1,
	}
}

func (vm *VM) Next() (int, bool) {
	if len(vm.insts) == 0 {
		return vm.reg, true
	}

	regToReturn := vm.reg

	// Fetch
	if vm.curInst == nil {
		vm.curInst = &vm.insts[0]
		vm.insts = vm.insts[1:]
		vm.curInstCycles = 1
	}

	// Execute
	switch vm.curInst.Name {
	case "noop":
		vm.curInst = nil
	case "addx":
		if vm.curInstCycles == 2 {
			vm.reg += vm.curInst.Arg
			vm.curInst = nil
		}
	default:
		panic(fmt.Sprintf("bad instruction %v", vm.curInst.Name))
	}

	vm.curInstCycles++
	return regToReturn, false
}

func solveA(insts []Inst) int {
	vm := NewVM(insts)

	score := 0
	for cycle := 1; ; cycle++ {
		reg, done := vm.Next()
		if done {
			break
		}

		if cycle == 20 || (cycle-20)%40 == 0 {
			score += cycle * reg
		}

		if cycle > 10000 {
			panic("too many cycles")
		}
	}
	return score
}

func solveB(insts []Inst) int {
	return -1
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

	insts, err := parseInstructions(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(insts))
	fmt.Println("B", solveB(insts))
}
