package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"sort"
	"strings"

	"github.com/simmonmt/aoc/2024/common/collections"
	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/lineio"
	"github.com/simmonmt/aoc/2024/common/logger"
)

var (
	verbose      = flag.Bool("verbose", false, "verbose")
	input        = flag.String("input", "", "input file")
	cpuprofile   = flag.String("cpuprofile", "", "write cpu profile to file")
	labelWires   = flag.Bool("label_wires", false, "label wires")
	minimizeOuts = flag.String("minimize_outs", "", "comma-separated list of outs")

	gatePattern = regexp.MustCompile(`^([^ ]+) (AND|XOR|OR) ([^ ]+) -> ([^ ]+)$`)
)

type State int
type GateType int

const (
	ST_UNK State = iota
	ST_OFF
	ST_ON

	GT_OR GateType = iota
	GT_AND
	GT_XOR
)

func (s State) String() string {
	switch s {
	case ST_UNK:
		return "?"
	case ST_OFF:
		return "0"
	case ST_ON:
		return "1"
	default:
		panic("bad state")
	}
}

func (gt GateType) String() string {
	switch gt {
	case GT_OR:
		return "OR"
	case GT_AND:
		return "AND"
	case GT_XOR:
		return "XOR"
	default:
		panic("bad gate")
	}
}

type Gate struct {
	ID       int
	In1, In2 string
	Out      string
	Type     GateType
}

type Input struct {
	Wires map[string]State
	Gates []Gate
}

func parseWires(lines []string) (map[string]State, error) {
	out := map[string]State{}
	for i, line := range lines {
		name, val, ok := strings.Cut(line, ": ")
		if !ok {
			return nil, fmt.Errorf("wire %d: bad cut", i+1)
		}

		if val == "1" {
			out[name] = ST_ON
		} else {
			out[name] = ST_OFF
		}
	}
	return out, nil
}

func parseGates(lines []string) ([]Gate, error) {
	gates := make([]Gate, len(lines))
	for i, line := range lines {
		parts := gatePattern.FindStringSubmatch(line)
		if len(parts) != 5 {
			return nil, fmt.Errorf("gate %d: bad match", i+1)
		}

		var gateType GateType
		if t := parts[2]; t == "AND" {
			gateType = GT_AND
		} else if t == "OR" {
			gateType = GT_OR
		} else if t == "XOR" {
			gateType = GT_XOR
		} else {
			panic("bad gate type")
		}

		gate := &gates[i]
		gate.ID = i
		gate.In1 = parts[1]
		gate.In2 = parts[3]
		gate.Out = parts[4]
		gate.Type = gateType
	}
	return gates, nil
}

func parseInput(lines []string) (*Input, error) {
	groups := lineio.BlankSeparatedGroups(lines)
	if len(groups) != 2 {
		return nil, fmt.Errorf("bad groups, want 2, got %d", len(groups))
	}

	wires, err := parseWires(groups[0])
	if err != nil {
		return nil, fmt.Errorf("bad wires: %v", err)
	}

	gates, err := parseGates(groups[1])
	if err != nil {
		return nil, fmt.Errorf("bad gates: %v", err)
	}

	return &Input{
		Wires: wires,
		Gates: gates,
	}, nil
}

func executeGate(gate *Gate, in1, in2 State) State {
	toState := func(in bool) State {
		if in {
			return ST_ON
		}
		return ST_OFF
	}

	if in1 == ST_UNK || in2 == ST_UNK {
		return ST_UNK
	}

	switch gate.Type {
	case GT_OR:
		return toState(in1 == ST_ON || in2 == ST_ON)
	case GT_AND:
		return toState(in1 == ST_ON && in2 == ST_ON)
	case GT_XOR:
		return toState(in1 != in2)
	default:
		panic("bad type")
	}
}

func calculateSolution(wires map[string]State) int {
	zWires := collections.FilterMap(wires, func(name string, _ State) bool {
		return name[0] == 'z'
	})

	zNames := collections.MapKeys(zWires)
	sort.Strings(zNames)

	var out int
	for i := len(zNames) - 1; i >= 0; i-- {
		var bit int
		if zWires[zNames[i]] == ST_ON {
			bit = 1
		}

		out = (out << 1) | bit
	}
	return out
}

func solve(wiresIn map[string]State, gates []Gate) int {
	solvedGates := map[*Gate]bool{}
	wires := collections.CloneMap(wiresIn)

	for len(solvedGates) != len(gates) {
		for i := range len(gates) {
			gate := &gates[i]
			if solvedGates[gate] {
				continue
			}

			out := executeGate(gate, wires[gate.In1], wires[gate.In2])
			if out == ST_UNK {
				continue
			}

			wires[gate.Out] = out
			solvedGates[gate] = true
		}

		logger.Infof("solved gates: %d of %d", len(solvedGates), len(gates))
	}
	return calculateSolution(wires)
}

func solveA(input *Input) int {
	return solve(input.Wires, input.Gates)
}

func makeDot(gates []Gate, input *Input) {
	fmt.Println("digraph G {")

	needInWires := map[string]bool{}
	for _, gate := range gates {
		fmt.Printf("gate%d_ [label=\"%s %d\"];\n", gate.ID, gate.Type, gate.ID)

		for _, g2 := range gates {
			if gate.Out == g2.In1 || gate.Out == g2.In2 {
				label := ""
				if *labelWires {
					label = fmt.Sprintf("[label=\"%s\"]", gate.Out)
				}

				fmt.Printf("gate%d_ -> gate%d_ %s\n", gate.ID, g2.ID, label)
			}
		}

		if in := gate.In1[0]; in == 'x' || in == 'y' || in == 'z' {
			needInWires[gate.In1] = true
			fmt.Printf("in_%s -> gate%d_\n", gate.In1, gate.ID)
		}
		if in := gate.In2[0]; in == 'x' || in == 'y' || in == 'z' {
			needInWires[gate.In2] = true
			fmt.Printf("in_%s -> gate%d_\n", gate.In2, gate.ID)
		}
		if out := gate.Out[0]; out == 'x' || out == 'y' || out == 'z' {
			fmt.Printf("gate%d_ -> out_%s\n", gate.ID, gate.Out)
		}
	}

	for wire := range needInWires {
		state := input.Wires[wire]

		var name string
		if wire[0] == 'x' || wire[0] == 'y' {
			name = fmt.Sprintf("in_%s", wire)
		} else if wire[0] == 'z' {
			name = fmt.Sprintf("out_%s", wire)
		} else {
			continue
		}
		label := fmt.Sprintf("%s = %s", name, state)

		fmt.Printf("%s [label=\"%s\"]\n", name, label)
	}

	fmt.Println("}")
}

func minimize(gates []Gate, outs []string) []Gate {
	needWires := map[string]bool{}
	needGates := map[int]bool{}
	for _, out := range outs {
		needWires[out] = true
	}

	changed := true
	for changed {
		changed = false
		for _, gate := range gates {
			if !needWires[gate.Out] {
				continue // doesn't help us
			}
			if needGates[gate.ID] {
				continue // already processed
			}

			needGates[gate.ID] = true
			needWires[gate.In1] = true
			needWires[gate.In2] = true
			changed = true
		}
	}

	out := []Gate{}
	for _, gate := range gates {
		if _, found := needGates[gate.ID]; found {
			out = append(out, gate)
		}
	}

	return out
}

func findGate(gates []Gate, id int) *Gate {
	for i := range gates {
		if gates[i].ID == id {
			return &gates[i]
		}
	}

	panic(fmt.Sprintf("failed to find gate %d", id))
}

func swap(gates []Gate, a, b int) ([]Gate, []string) {
	out := make([]Gate, len(gates))
	copy(out, gates)

	aGate := findGate(out, a)
	bGate := findGate(out, b)

	aGate.Out, bGate.Out = bGate.Out, aGate.Out
	return out, []string{aGate.Out, bGate.Out}
}

func solveB(input *Input) string {
	swaps := [][2]int{[2]int{83, 50}, [2]int{45, 81}, [2]int{23, 46}, [2]int{183, 63}}

	gates := input.Gates

	swappedWires := []string{}
	for _, s := range swaps {
		var wires []string
		gates, wires = swap(gates, s[0], s[1])
		swappedWires = append(swappedWires, wires...)
	}

	if *minimizeOuts != "" {
		gates = minimize(gates, strings.Split(*minimizeOuts, ","))
	}

	// wires := map[string]State{}
	// for i := 0; i < 45; i++ {
	// 	wires[fmt.Sprintf("x%02d", i)] = ST_ON
	// 	wires[fmt.Sprintf("y%02d", i)] = ST_OFF
	// }

	// fmt.Printf("%x\n", solve(wires, gates))

	//makeDot(gates, input)

	sort.Strings(swappedWires)
	return strings.Join(swappedWires, ",")
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
