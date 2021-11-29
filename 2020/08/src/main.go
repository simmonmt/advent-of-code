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
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2020/08/src/vm"
	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type TermType int

const (
	TERM_INF TermType = iota
	TERM_EXIT
)

func run(resources *vm.Resources, insts []vm.Instruction) TermType {
	seen := map[int64]bool{0: true}
	pc := int64(0)
	for {
		inst := insts[pc]
		npc := inst.Execute(resources, pc)
		logger.LogF("%d: %-20s acc %d => %d", pc, inst,
			resources.Acc, npc)

		if _, found := seen[npc]; found {
			return TERM_INF
		}
		seen[npc] = true

		if npc >= int64(len(insts)) {
			return TERM_EXIT
		}
		pc = npc
	}
}

func solveA(insts []vm.Instruction) {
	resources := &vm.Resources{Acc: 0}
	run(resources, insts)
	fmt.Printf("A: acc %d\n", resources.Acc)
}

func solveB(insts []vm.Instruction) {
	for i, inst := range insts {
		var newInst vm.Instruction
		var err error
		switch inst.Op() {
		case "jmp":
			newInst, err = vm.NewInst("nop", inst.(*vm.Jmp).A())
		case "nop":
			newInst, err = vm.NewInst("jmp", inst.(*vm.Nop).A())
		}

		if err != nil {
			log.Fatalf("failed transmogrify %d: %v: %v",
				i, inst, err)
		}
		if newInst == nil {
			continue
		}

		logger.LogF("trying %v to %v at %v", inst, newInst, i)

		insts[i] = newInst
		resources := &vm.Resources{Acc: 0}
		tt := run(resources, insts)
		insts[i] = inst

		if tt == TERM_EXIT {
			fmt.Printf("B: acc %d\n", resources.Acc)
			return
		}
	}

	fmt.Println("B: not found")
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

	insts := []vm.Instruction{}
	for lineNo, line := range lines {
		inst, err := vm.Decode(line)
		if err != nil {
			log.Fatalf("%d: %v", lineNo, err)
		}

		insts = append(insts, inst)
	}

	solveA(insts)
	solveB(insts)
}
