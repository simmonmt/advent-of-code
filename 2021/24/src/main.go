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
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Instruction struct {
	Op, A, B string
}

func parseInstruction(line string) (*Instruction, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 || len(parts) > 3 {
		return nil, fmt.Errorf("bad instruction")
	}

	inst := &Instruction{
		Op: parts[0],
		A:  parts[1],
	}

	if len(parts) == 3 {
		inst.B = parts[2]
	}

	return inst, nil
}

func readBVal(bStr string, regs map[string]int) (int, error) {
	if bStr == "w" || bStr == "x" || bStr == "y" || bStr == "z" {
		return regs[bStr], nil
	}

	return strconv.Atoi(bStr)
}

func decodeZ(z int) string {
	out := []string{}
	for z > 0 {
		out = append([]string{strconv.Itoa(z % 26)}, out...)
		z /= 26
	}
	return strings.Join(out, " ")
}

func executeInstruction(inst *Instruction, regs map[string]int, input *string) error {
	switch inst.Op {
	case "inp":
		logger.LogF("before inp regs: %v (%v)", regs, decodeZ(regs["z"]))
		if len(*input) == 0 {
			return fmt.Errorf("ran out of input")
		}

		inCh := string((*input)[0])
		*input = (*input)[1:]

		num, err := strconv.Atoi(inCh)
		if err != nil || num < 0 || num > 9 {
			return fmt.Errorf("bad input %v", inCh)
		}

		regs[inst.A] = num

	case "add":
		bVal, err := readBVal(inst.B, regs)
		if err != nil {
			return err
		}

		regs[inst.A] = regs[inst.A] + bVal

	case "mul":
		bVal, err := readBVal(inst.B, regs)
		if err != nil {
			return err
		}

		regs[inst.A] = regs[inst.A] * bVal

	case "div":
		bVal, err := readBVal(inst.B, regs)
		if err != nil {
			return err
		}

		regs[inst.A] = regs[inst.A] / bVal

	case "mod":
		bVal, err := readBVal(inst.B, regs)
		if err != nil {
			return err
		}

		regs[inst.A] = regs[inst.A] % bVal

	case "eql":
		bVal, err := readBVal(inst.B, regs)
		if err != nil {
			return err
		}

		if regs[inst.A] == bVal {
			regs[inst.A] = 1
		} else {
			regs[inst.A] = 0
		}

	default:
		return fmt.Errorf("unknown opcode")
	}

	return nil
}

func runProgram(lines []string, regs map[string]int, input string) error {
	for lineNo := 1; lineNo < len(lines); lineNo++ {
		inst, err := parseInstruction(lines[lineNo-1])
		if err != nil {
			return fmt.Errorf("%d: parse failure: %v", lineNo, err)
		}

		if err := executeInstruction(inst, regs, &input); err != nil {
			return fmt.Errorf("%d: exec failure: %v", lineNo, err)
		}
	}

	logger.LogF("end regs: %v (%v)", regs, decodeZ(regs["z"]))
	return nil
}

func solve(lines []string, part, input string) {
	regs := map[string]int{"w": 0, "x": 0, "y": 0, "z": 0}
	if err := runProgram(lines, regs, input); err != nil {
		log.Fatal(err)
	}

	if regs["z"] == 0 {
		fmt.Println(part, input)
	} else {
		fmt.Println(part, "bad, z=", regs["z"])
	}
}

func solveA(lines []string) {
	input := "93499629698999" // see scratchpad file
	solve(lines, "A", input)
}

func solveB(lines []string) {
	input := "11164118121471"
	solve(lines, "B", input)
}

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	return lines, err
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

	solveA(lines)
	solveB(lines)
}
