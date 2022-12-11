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
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Operation interface {
	Execute(old int) int
}

type MulOp struct {
	num int
}

func (o *MulOp) Execute(old int) int {
	return old * o.num
}

type AddOp struct {
	num int
}

func (o *AddOp) Execute(old int) int {
	return old + o.num
}

type SquareOp struct{}

func (o *SquareOp) Execute(old int) int {
	return old * old
}

type Monkey struct {
	id                  int
	items               []int
	numInspections      int
	op                  Operation
	testDivisor         int
	trueDest, falseDest int
}

func NewMonkey(id int, items []int, op Operation, testDivisor int, trueDest, falseDest int) *Monkey {
	return &Monkey{
		id:             id,
		items:          items[:],
		numInspections: 0,
		op:             op,
		testDivisor:    testDivisor,
		trueDest:       trueDest,
		falseDest:      falseDest,
	}
}

func (m *Monkey) Clone() *Monkey {
	return &Monkey{
		id:             m.id,
		items:          m.items[:],
		numInspections: m.numInspections,
		op:             m.op,
		testDivisor:    m.testDivisor,
		trueDest:       m.trueDest,
		falseDest:      m.falseDest,
	}
}

func (m *Monkey) ID() int {
	return m.id
}

func (m *Monkey) TestDivisor() int {
	return m.testDivisor
}

func (m *Monkey) Items() []int {
	return m.items[:]
}

func (m *Monkey) NumInspections() int {
	return m.numInspections
}

func (m *Monkey) AddItem(item int) {
	m.items = append(m.items, item)
}

func (m *Monkey) DumpItems() {
	fmt.Printf("%2d: ", m.id)
	for i, item := range m.items {
		if i != 0 {
			fmt.Printf(", ")
		}
		fmt.Print(item)
	}
	fmt.Printf(" (%d)\n", m.numInspections)
}

// Step processes a single round, iterating over each item and deciding what to
// do with it. `div` handles the differences between the part A and B
// algorithms, and is the value by which the worry level is divided for each
// item. `mod` should be the product of the (known prime) testDivisor values for
// all monkeys, and is used to keep the item values from getting out of
// hand. Math using those values will be done using modular arithmetic.
func (m *Monkey) Step(div, mod int) map[int][]int {
	out := map[int][]int{}
	if len(m.items) == 0 {
		return out
	}

	for _, level := range m.items {
		m.numInspections++
		level = m.op.Execute(level)

		level /= div
		if level > 0 {
			level = level % mod
		}

		var dest int
		if level%m.testDivisor == 0 {
			dest = m.trueDest
		} else {
			dest = m.falseDest
		}

		if _, found := out[dest]; !found {
			out[dest] = []int{}
		}
		out[dest] = append(out[dest], level)
	}

	m.items = []int{}
	return out
}

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

var (
	monkeyPattern = regexp.MustCompile(`` +
		`^Monkey (\d+):.*` +
		`Starting items: ([^|]+).*` +
		`Operation: new = old (.) ([^|]+).*` +
		`Test: divisible by (\d+).*` +
		`If true: throw to monkey (\d+).*` +
		`If false: throw to monkey (\d+)`)
)

func parseOperation(opChar string, opBStr string) (Operation, error) {
	opB, err := strconv.Atoi(opBStr)

	if opChar == "+" {
		if err != nil {
			return nil, err
		}
		return &AddOp{opB}, nil
	} else if opChar == "*" {
		if opBStr == "old" {
			return &SquareOp{}, nil
		}
		if err != nil {
			return nil, err
		}
		return &MulOp{opB}, nil
	}
	return nil, fmt.Errorf("bad operation %v", opChar)
}

func parseItems(all string) ([]int, error) {
	items := strings.Split(all, ", ")
	out := make([]int, len(items))
	for i, item := range items {
		num, err := strconv.Atoi(item)
		if err != nil {
			return nil, err
		}
		out[i] = num
	}
	return out, nil
}

func parseMonkey(lines []string) (*Monkey, error) {
	parts := monkeyPattern.FindStringSubmatch(strings.Join(lines, "|"))
	if len(parts) != 8 {
		return nil, fmt.Errorf("failed to regexp match monkey")
	}

	idStr, itemsAll, opChar, opB, divStr, trueStr, falseStr :=
		parts[1], parts[2], parts[3], parts[4], parts[5], parts[6],
		parts[7]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("bad monkey id: %v", err)
	}

	items, err := parseItems(itemsAll)
	if err != nil {
		return nil, fmt.Errorf("bad items: %v", err)
	}

	op, err := parseOperation(opChar, opB)
	if err != nil {
		return nil, fmt.Errorf("bad operation; %v", err)
	}

	div, err := strconv.Atoi(divStr)
	if err != nil {
		return nil, fmt.Errorf("bad test divisor: %v", err)
	}

	trueNum, err := strconv.Atoi(trueStr)
	if err != nil {
		return nil, fmt.Errorf("bad true monkey: %v", err)
	}

	falseNum, err := strconv.Atoi(falseStr)
	if err != nil {
		return nil, fmt.Errorf("bad false monkey: %v", err)
	}

	return NewMonkey(id, items, op, div, trueNum, falseNum), nil
}

func parseMonkeys(lines []string) ([]*Monkey, error) {
	monkeys := []*Monkey{}
	for i := 0; i < len(lines); i += 7 {
		if i+6 > len(lines) {
			return nil, fmt.Errorf("short input")
		}

		monkey, err := parseMonkey(lines[i : i+6])
		if err != nil {
			return nil, fmt.Errorf("bad monkey at line %v: %v",
				i+1, err)
		}
		monkeys = append(monkeys, monkey)
	}
	return monkeys, nil
}

func playGame(rounds int, div int, monkeys []*Monkey) int {
	mod := 1
	monkeysById := map[int]*Monkey{}
	for _, monkey := range monkeys {
		monkeysById[monkey.ID()] = monkey
		mod *= monkey.TestDivisor()
	}

	for round := 1; round <= rounds; round++ {
		for _, monkey := range monkeys {
			for id, items := range monkey.Step(div, mod) {
				for _, item := range items {
					monkeysById[id].AddItem(item)
				}
			}
		}

		// if logger.Enabled() {
		// 	logger.LogF("end of round %d", round)
		// 	for _, monkey := range monkeys {
		// 		monkey.DumpItems()
		// 	}
		// }
	}

	inspections := make([]int, len(monkeys))
	for i, monkey := range monkeys {
		inspections[i] = monkey.NumInspections()
	}

	sort.Ints(inspections)
	return inspections[len(inspections)-2] * inspections[len(inspections)-1]
}

func solveA(monkeys []*Monkey) int {
	return playGame(20, 3, monkeys)
}

func solveB(monkeys []*Monkey) int {
	return playGame(10000, 1, monkeys)
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

	monkeys, err := parseMonkeys(lines)
	if err != nil {
		log.Fatal(err)
	}

	// The items lists change as the game is played so we need to present
	// each solver with a new set of monkeys.
	cloneMonkeys := func(in []*Monkey) []*Monkey {
		out := make([]*Monkey, len(in))
		for i, m := range in {
			out[i] = m.Clone()
		}
		return out
	}

	fmt.Println("A", solveA(cloneMonkeys(monkeys)))
	fmt.Println("B", solveB(cloneMonkeys(monkeys)))
}
