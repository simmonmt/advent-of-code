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
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	startStatePattern     = regexp.MustCompile(`^Begin in state ([^.])\.$`)
	numItersPattern       = regexp.MustCompile(`^Perform a diagnostic checksum after ([0-9]+) steps.$`)
	inStatePattern        = regexp.MustCompile(`^In state ([^:]+):$`)
	ifValPattern          = regexp.MustCompile(`^If the current value is ([0-9]+):$`)
	actionWriteValPattern = regexp.MustCompile(`^- Write the value ([0-9]+)\.$`)
	actionMovePattern     = regexp.MustCompile(`^- Move one slot to the (right|left)\.$`)
	actionNewStatePattern = regexp.MustCompile(`^- Continue with state ([^.]+)\.$`)
)

type Condition struct {
	State string
	Val   bool
}

type Action struct {
	WriteVal bool
	Move     int
	NewState string
}

type Machine struct {
	Tape map[int]bool
	Pos  int
}

func newMachine() *Machine {
	return &Machine{
		Tape: map[int]bool{},
		Pos:  0,
	}
}

func (m *Machine) Checksum() int {
	numSet := 0
	for _, v := range m.Tape {
		if v {
			numSet++
		}
	}
	return numSet
}

func (m *Machine) Set(val bool) {
	m.Tape[m.Pos] = val
}

func (m *Machine) Get() bool {
	return m.Tape[m.Pos]
}

func (m *Machine) Move(amt int) {
	m.Pos += amt
}

func readInput(in io.Reader) (string, int, map[Condition]Action, error) {
	reader := bufio.NewReader(in)

	startState := ""
	nIters := -1
	config := map[Condition]Action{}

	var curCond Condition
	var curAction Action

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		if matches := startStatePattern.FindStringSubmatch(line); matches != nil {
			startState = matches[1]
		} else if matches := numItersPattern.FindStringSubmatch(line); matches != nil {
			nIters, err = strconv.Atoi(matches[1])
			if err != nil {
				return "", -1, nil, fmt.Errorf("failed to parse num iters val %v: %v", matches[1], err)
			}
		} else if matches := inStatePattern.FindStringSubmatch(line); matches != nil {
			curCond.State = matches[1]
		} else if matches := ifValPattern.FindStringSubmatch(line); matches != nil {
			curCond.Val = matches[1] == "1"
		} else if matches := actionWriteValPattern.FindStringSubmatch(line); matches != nil {
			curAction.WriteVal = matches[1] == "1"
			config[curCond] = curAction
		} else if matches := actionMovePattern.FindStringSubmatch(line); matches != nil {
			if matches[1] == "left" {
				curAction.Move = -1
			} else {
				curAction.Move = 1
			}
			config[curCond] = curAction
		} else if matches := actionNewStatePattern.FindStringSubmatch(line); matches != nil {
			curAction.NewState = matches[1]
			config[curCond] = curAction
		}
	}

	return startState, nIters, config, nil
}

func main() {
	startState, nIters, config, err := readInput(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	// fmt.Printf("start state %v, nIters %v\n", startState, nIters)
	// fmt.Println(config)

	machine := newMachine()

	state := startState
	for i := 0; i < nIters; i++ {
		action := config[Condition{State: state, Val: machine.Get()}]
		machine.Set(action.WriteVal)
		machine.Move(action.Move)
		state = action.NewState
	}

	fmt.Printf("checksum %v\n", machine.Checksum())
}
