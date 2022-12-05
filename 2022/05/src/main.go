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
	"regexp"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Stack struct {
	id    int
	elems *list.List
}

func NewStack(id int) *Stack {
	return &Stack{id: id, elems: list.New()}
}

func (s *Stack) Clone() *Stack {
	n := NewStack(s.id)
	for elem := s.elems.Front(); elem != nil; elem = elem.Next() {
		n.elems.PushBack(elem.Value)
	}
	return n
}

func (s *Stack) ID() int {
	return s.id
}

func (s *Stack) PushFront(r rune) {
	s.elems.PushFront(r)
}

func (s *Stack) Push(r rune) {
	s.elems.PushBack(r)
}

func (s *Stack) Pop() rune {
	if s.elems.Front() == nil {
		panic(fmt.Sprintf("empty stack %d", s.id))
	}

	r := s.elems.Back().Value.(rune)
	s.elems.Remove(s.elems.Back())
	return r
}

func (s *Stack) Peek() rune {
	if s.elems.Front() == nil {
		panic(fmt.Sprintf("empty stack %d", s.id))
	}

	return s.elems.Back().Value.(rune)
}

func (s *Stack) String() string {
	out := []string{}
	for elem := s.elems.Front(); elem != nil; elem = elem.Next() {
		out = append(out, string(elem.Value.(rune)))
	}
	return strings.Join(out, ",")
}

func parseStacks(lines []string, num int) ([]*Stack, error) {
	stacks := make([]*Stack, num+1)
	for i := 0; i <= num; i++ {
		stacks[i] = NewStack(i)
	}

	for _, line := range lines {
		off := 1
		for i := 1; i <= num; i++ {
			if off >= len(line) {
				continue
			}
			r := rune(line[off])
			if r != ' ' {
				stacks[i].PushFront(rune(line[off]))
			}
			off += 4
		}
	}

	return stacks, nil
}

var (
	instPattern = regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)
)

func parseInstructions(lines []string) ([]Instruction, error) {
	insts := []Instruction{}

	for _, line := range lines {
		parts := instPattern.FindStringSubmatch(line)
		if parts == nil {
			return nil, fmt.Errorf("no regexp match")
		}

		var inst Instruction
		var err error
		inst.Num, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("bad num")
		}
		inst.From, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("bad from")
		}
		inst.To, err = strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("bad to")
		}

		insts = append(insts, inst)
	}

	return insts, nil
}

type Instruction struct {
	From int
	Num  int
	To   int
}

func readInput(path string) ([]*Stack, []Instruction, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, nil, err
	}

	var stackLines, instLines []string
	var numStacks int
	for i, line := range lines {
		if line != "" {
			continue
		}

		stackLines = lines[0 : i-1]
		instLines = lines[i+1:]
		numStacks = len(strings.Fields(lines[i-1]))
		break
	}

	stacks, err := parseStacks(stackLines, numStacks)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse stacks: %v", err)
	}
	insts, err := parseInstructions(instLines)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse instructions: %v", err)
	}

	return stacks, insts, nil
}

func solveA(stacks []*Stack, insts []Instruction) string {
	for _, inst := range insts {
		for i := 0; i < inst.Num; i++ {
			r := stacks[inst.From].Pop()
			stacks[inst.To].Push(r)
		}
	}

	out := []string{}
	for i := 1; i < len(stacks); i++ {
		out = append(out, string(stacks[i].Peek()))
	}

	return strings.Join(out, "")
}

func solveB(stacks []*Stack, insts []Instruction) string {
	for _, inst := range insts {
		accum := make([]rune, inst.Num)
		for i := 0; i < inst.Num; i++ {
			accum[i] = stacks[inst.From].Pop()
		}

		for i := len(accum) - 1; i >= 0; i-- {
			stacks[inst.To].Push(accum[i])
		}
	}

	out := []string{}
	for i := 1; i < len(stacks); i++ {
		out = append(out, string(stacks[i].Peek()))
	}

	return strings.Join(out, "")
}

func dumpStacks(stacks []*Stack) {
	fmt.Println("stacks:")
	for i := 1; i < len(stacks); i++ {
		fmt.Printf("  %2d: %s\n", stacks[i].ID(), stacks[i])
	}
}

func cloneStacks(stacks []*Stack) []*Stack {
	n := make([]*Stack, len(stacks))
	for i, stack := range stacks {
		n[i] = stack.Clone()
	}
	return n
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	stacks, insts, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	if logger.Enabled() {
		dumpStacks(stacks)
		fmt.Println(insts)
	}

	fmt.Println("A", solveA(cloneStacks(stacks), insts))
	fmt.Println("B", solveB(cloneStacks(stacks), insts))
}
