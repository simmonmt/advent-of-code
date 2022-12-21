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
	"container/list"
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

type Operation int

const (
	ADD Operation = 0
	SUB Operation = 1
	MUL Operation = 2
	DIV Operation = 3
	IMM Operation = 4
)

type Monkey struct {
	name string
	op   Operation
	imm  int64
	a, b string
}

func (i *Monkey) Execute(answers map[string]int64) (int64, bool) {
	if i.op == IMM {
		return i.imm, true
	}

	a, found := answers[i.a]
	if !found {
		return 0, false
	}

	b, found := answers[i.b]
	if !found {
		return 0, false
	}

	switch i.op {
	case ADD:
		return a + b, true
	case SUB:
		return a - b, true
	case MUL:
		return a * b, true
	case DIV:
		return a / b, true
	}

	panic("bad op")
}

func parseMonkey(str string) (*Monkey, error) {
	name, rest, ok := strings.Cut(str, ": ")
	if !ok {
		return nil, fmt.Errorf("no ': '")
	}

	num, err := strconv.Atoi(rest)
	if err == nil {
		return &Monkey{name: name, op: IMM, imm: int64(num)}, nil
	}

	parts := strings.Split(rest, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("bad inst split")
	}

	a, b := parts[0], parts[2]

	var op Operation
	switch parts[1] {
	case "+":
		op = ADD
	case "-":
		op = SUB
	case "*":
		op = MUL
	case "/":
		op = DIV
	default:
		return nil, fmt.Errorf("bad op %v", parts[1])
	}

	return &Monkey{
		name: name,
		op:   op,
		a:    a,
		b:    b,
	}, nil
}

func parseInput(lines []string) (map[string]*Monkey, error) {
	monkeys := map[string]*Monkey{}
	for i, line := range lines {
		monkey, err := parseMonkey(line)
		if err != nil {
			return nil, fmt.Errorf("%d: %v", i+1, err)
		}
		monkeys[monkey.name] = monkey
	}
	return monkeys, nil
}

func solveMonkeys(monkeys map[string]*Monkey, start, fail string) (int64, bool) {
	need := list.New()
	need.PushBack(start)

	answers := map[string]int64{}

	for need.Front() != nil {
		name := need.Back().Value.(string)
		if fail == name {
			return 0, false
		}

		monkey, found := monkeys[name]
		if !found {
			panic("bad monkey name '" + name + "'")
		}

		ans, ok := monkey.Execute(answers)
		if ok {
			answers[name] = ans
			need.Remove(need.Back())
			continue
		}

		if _, found := answers[monkey.a]; !found {
			need.PushBack(monkey.a)
		}
		if _, found := answers[monkey.b]; !found {
			need.PushBack(monkey.b)
		}
	}

	return answers[start], true
}

func solveA(monkeys map[string]*Monkey) int64 {
	ans, ok := solveMonkeys(monkeys, "root", "")
	if !ok {
		panic("bad solvea")
	}
	return ans
}

func checkSolveability(monkeys map[string]*Monkey, toCheck *Monkey) (knownName, unknownName string, knownAns int64, knownLeft bool) {
	var ok bool
	knownName, unknownName, knownLeft = toCheck.a, toCheck.b, true
	knownAns, ok = solveMonkeys(monkeys, knownName, "humn")
	if !ok {
		knownName, unknownName, knownLeft = toCheck.b, toCheck.a, false
		knownAns, ok = solveMonkeys(monkeys, knownName, "humn")
	}

	return
}

func unsolve(monkeys map[string]*Monkey, name string, result int64) (string, int64) {
	// given a) a dual-operand monkey, b) a result, and c) knowledge that
	// one side of the monkey (a or b) is solveable, figure out 1) which
	// side is unsolveable and 2) what the value is.
	monkey := monkeys[name]
	knownName, unknownName, knownAns, knownLeft := checkSolveability(monkeys, monkey)
	logger.LogF("unsolve: known %v %v unknown %v", knownName, knownAns, unknownName)

	// We know which side is known/unknown, so now let's reverse the
	// operation.
	return unknownName, unsolveOp(monkey.op, result, knownAns, knownLeft)
}

func unsolveOp(op Operation, result, known int64, knownLeft bool) int64 {
	switch op {
	case ADD:
		return result - known // k + u = result
	case SUB:
		if knownLeft {
			return known - result // k - u = result
		} else {
			return result + known // u - k = result
		}
	case MUL:
		return result / known // k * u = result
	case DIV:
		if knownLeft {
			return known / result // k / u = result
		} else {
			return result * known // u / k = result
		}
	default:
		panic("bad op")
	}
}

func solveB(monkeys map[string]*Monkey) int64 {
	root := monkeys["root"]

	rootKnownName, rootUnknownName, rootKnownAns, _ := checkSolveability(monkeys, root)
	logger.LogF("root is %+v, known %v %v unknown %v",
		root, rootKnownName, rootKnownAns, rootUnknownName)

	name := rootUnknownName
	result := rootKnownAns
	input := int64(-1)
	for {
		logger.LogF("unsolve %v (want result %v)", name, result)
		nextName, nextResult := unsolve(monkeys, name, result)
		logger.LogF("unsolve next %v %v", nextName, nextResult)

		if nextName == "humn" {
			input = nextResult
			break
		}

		name, result = nextName, nextResult
	}

	// Verify
	humn := monkeys["humn"]
	humn.op = IMM
	humn.imm = input

	rootUnknownAns, ok := solveMonkeys(monkeys, rootUnknownName, "")
	if !ok {
		panic("unable to resolve")
	}

	if rootUnknownAns != rootKnownAns {
		panic(fmt.Sprintf("resolve mismatch unk %v known %v",
			rootUnknownAns, rootKnownAns))
	}

	return input
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

	input, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
