// Copyright 2023 Google LLC
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
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/mtsmath"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Pulse struct {
	Src, Dest string
	Signal    bool
}

type Node interface {
	Name() string
	Targets() []string
	Init(srcs []string)
	Set(signal bool, src string) []Pulse
	AtRest() bool
	Reset()
	String() string
}

type nodeImpl struct {
	name    string
	targets []string
}

func (n *nodeImpl) Name() string       { return n.name }
func (n *nodeImpl) Targets() []string  { return n.targets }
func (n *nodeImpl) Init(srcs []string) {}
func (n *nodeImpl) AtRest() bool       { return true }
func (n *nodeImpl) Reset()             {}

type Broadcaster struct {
	nodeImpl
}

func (n *Broadcaster) Set(signal bool, _ string) []Pulse {
	out := []Pulse{}
	for _, tgt := range n.targets {
		out = append(out, Pulse{n.name, tgt, signal})
	}
	return out
}

func (n *Broadcaster) String() string {
	return n.name
}

type FlipFlop struct {
	nodeImpl
	state bool
}

func (n *FlipFlop) Set(signal bool, _ string) []Pulse {
	if signal {
		return nil
	}

	n.state = !n.state

	out := []Pulse{}
	for _, tgt := range n.targets {
		out = append(out, Pulse{n.name, tgt, n.state})
	}

	return out
}

func (n *FlipFlop) AtRest() bool {
	return n.state == false
}

func (n *FlipFlop) Reset() {
	n.state = false
}

func (n *FlipFlop) String() string {
	return fmt.Sprintf("%%%s[%v]", n.name, n.state)
}

type Conjunction struct {
	nodeImpl
	srcs  []string
	lasts map[string]bool
}

func (n *Conjunction) Init(srcs []string) {
	n.srcs = make([]string, len(srcs))
	copy(n.srcs, srcs)
	n.Reset()
}

func (n *Conjunction) Set(signal bool, src string) []Pulse {
	n.lasts[src] = signal

	allHigh := true
	for _, v := range n.lasts {
		if !v {
			allHigh = false
			break
		}
	}

	out := []Pulse{}
	for _, tgt := range n.targets {
		out = append(out, Pulse{n.name, tgt, !allHigh})
	}
	return out
}

func (n *Conjunction) AtRest() bool {
	for _, v := range n.lasts {
		if v {
			return false
		}
	}
	return true
}

func (n *Conjunction) Reset() {
	n.lasts = map[string]bool{}
	for _, src := range n.srcs {
		n.lasts[src] = false
	}
}

func (n *Conjunction) String() string {
	return fmt.Sprintf("&%s[%v]", n.name, n.lasts)
}

type Nop struct {
	nodeImpl
}

func (n *Nop) Set(signal bool, src string) []Pulse { return nil }
func (n *Nop) String() string                      { return n.name }

func NewNode(spec string) (Node, error) {
	name, targetList, ok := strings.Cut(spec, " -> ")
	if !ok {
		return nil, fmt.Errorf("missing ->")
	}

	targets := strings.Split(targetList, ", ")

	if name == "broadcaster" {
		return &Broadcaster{nodeImpl{name: name, targets: targets}}, nil
	}

	realName := name[1:]
	switch name[0] {
	case '%':
		return &FlipFlop{
			nodeImpl: nodeImpl{name: realName, targets: targets},
			state:    false,
		}, nil
	case '&':
		return &Conjunction{
			nodeImpl: nodeImpl{name: realName, targets: targets},
			lasts:    map[string]bool{},
		}, nil
	default:
		panic("bad node")
	}
}

func parseInput(lines []string) (map[string]Node, error) {
	ins := map[string][]string{}

	out := map[string]Node{}
	for i, line := range lines {
		n, err := NewNode(line)
		if err != nil {
			return nil, fmt.Errorf("%d: bad node: %v", i+1, err)
		}
		if _, found := out[n.Name()]; found {
			return nil, fmt.Errorf("%d: repeated %v", i+1, n.Name())
		}

		out[n.Name()] = n

		for _, tgt := range n.Targets() {
			ins[tgt] = append(ins[tgt], n.Name())
		}
	}

	for _, node := range out {
		node.Init(ins[node.Name()])
	}

	return out, nil
}

func runIteration(graph map[string]Node) (int, int, bool) {
	pulses := []Pulse{Pulse{"input", "broadcaster", false}}

	numHigh, numLow := 0, 0
	for len(pulses) > 0 {
		next := []Pulse{}

		for _, pulse := range pulses {
			if pulse.Signal {
				numHigh++
			} else {
				numLow++
			}

			node, found := graph[pulse.Dest]
			if !found {
				//logger.Infof("ignoring unknown node %v", pulse.Dest)
				continue
			}

			addl := node.Set(pulse.Signal, pulse.Src)
			//logger.Infof("node %v received %v sent %v", node.Name(), pulse, addl)

			if len(addl) > 0 {
				next = append(next, addl...)
			}
		}

		pulses = next
	}

	atRest := true
	for _, node := range graph {
		if !node.AtRest() {
			atRest = false
			break
		}
	}

	return numHigh, numLow, atRest
}

// 897226215 too high
// 894468890 too high
// 886347020
func solveA(graph map[string]Node) int {
	totHigh, totLow := 0, 0
	var i int
	for i = 0; i < 1000; i++ {
		numHigh, numLow, atRest := runIteration(graph)
		logger.Infof("high %v low %v atRest %v", numHigh, numLow, atRest)
		totHigh += numHigh
		totLow += numLow
		if atRest {
			break
		}
	}
	if i == 1000 {
		return totHigh * totLow
	}

	reps := i + 1

	logger.Infof("done after %d reps, high %v low %v", reps, totHigh, totLow)
	if 1000%reps != 0 {
		panic("uneven")
	}

	return (totHigh * (1000 / reps)) * (totLow * (1000 / reps))
}

func solveB(graph map[string]Node) int64 {
	sizes := []int64{}
	for _, name := range []string{"qs", "sp", "pg", "sv"} {
		old := graph[name]
		ff := FlipFlop{nodeImpl: nodeImpl{name: old.Name(), targets: old.Targets()}, state: false}
		graph[name] = &ff

		for _, node := range graph {
			node.Reset()
		}

		i := 1
		for {
			runIteration(graph)
			if i > 1 && ff.state {
				break
			}

			i++
		}

		sizes = append(sizes, int64(i))
		graph[name] = old
	}

	return mtsmath.LCM(sizes...)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	graph, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}
	fmt.Println("A", solveA(graph))

	graph, err = parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}
	fmt.Println("B", solveB(graph))
}
