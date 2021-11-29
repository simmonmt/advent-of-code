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
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	OpSet = iota
	OpAnd
	OpOr
	OpLShift
	OpRShift
	OpNot
)

type Gate struct {
	Deps       []string
	AReg, BReg string
	AVal, BVal uint
	Op         int
	Result     string
}

func makeDeps(a, b string) []string {
	deps := []string{}
	if a != "" {
		deps = append(deps, a)
	}
	if b != "" {
		deps = append(deps, b)
	}
	return deps
}

func makeDep(a string) []string {
	return makeDeps(a, "")
}

func parseRegOrVal(str string) (string, uint) {
	val, err := strconv.ParseUint(str, 0, 16)
	if err != nil {
		return str, 0
	}
	return "", uint(val)
}

func regOrVal(reg string, val uint, registers map[string]uint) uint {
	if reg == "" {
		return val
	}
	return registers[reg]
}

func parseGate(line string) (*Gate, error) {
	sides := strings.SplitN(line, " -> ", 2)
	left, result := sides[0], sides[1]

	parts := strings.Split(left, " ")
	switch len(parts) {
	case 1:
		reg, val := parseRegOrVal(parts[0])
		return &Gate{
			Deps:   makeDep(reg),
			AReg:   reg,
			AVal:   val,
			Op:     OpSet,
			Result: result,
		}, nil
		break

	case 2:
		if parts[0] != "NOT" {
			return nil, fmt.Errorf("unexpected 2-part op %v", parts[0])
		}

		reg, val := parseRegOrVal(parts[1])
		return &Gate{
			Deps:   makeDep(reg),
			AReg:   reg,
			AVal:   val,
			Op:     OpNot,
			Result: result,
		}, nil
		break

	case 3:
		opStr := parts[1]

		areg, aval := parseRegOrVal(parts[0])
		breg, bval := parseRegOrVal(parts[2])

		var op int
		switch opStr {
		case "AND":
			op = OpAnd
			break
		case "OR":
			op = OpOr
			break
		case "LSHIFT":
			op = OpLShift
			break
		case "RSHIFT":
			op = OpRShift
			break
		default:
			return nil, fmt.Errorf("unknown op %v", opStr)
		}

		return &Gate{
			Deps:   makeDeps(areg, breg),
			AReg:   areg,
			BReg:   breg,
			AVal:   aval,
			BVal:   bval,
			Op:     op,
			Result: result,
		}, nil
		break

	default:
		return nil, fmt.Errorf("unexpected number of parts: %v", len(parts))
	}

	panic("unreachable")
}

func processGate(gate *Gate, results map[string]uint) uint {
	aval := regOrVal(gate.AReg, gate.AVal, results)
	bval := regOrVal(gate.BReg, gate.BVal, results)

	switch gate.Op {
	case OpSet:
		return aval
	case OpAnd:
		return aval & bval
	case OpOr:
		return aval | bval
	case OpLShift:
		return (aval << bval) & 0xffff
	case OpRShift:
		return aval >> bval
	case OpNot:
		return (^aval) & 0xffff
	default:
		panic(fmt.Sprintf("unknown op %d", gate.Op))
	}
}

func main() {
	gates := []*Gate{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		gate, err := parseGate(line)
		if err != nil {
			log.Fatalf("failed to parse gate %v: %v", line, err)
		}

		gates = append(gates, gate)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading stdin: %v", err)
	}

	results := map[string]uint{}

	lastNumUnresolved := 0
	for i := 0; ; i++ {
		numUnresolved := 0

		for _, gate := range gates {
			if _, found := results[gate.Result]; found {
				continue
			}

			unfoundDeps := 0
			for _, dep := range gate.Deps {
				if _, found := results[dep]; !found {
					unfoundDeps++
					break
				}
			}

			if unfoundDeps > 0 {
				numUnresolved++
				continue
			}

			fmt.Printf("round %d found %v\n", i, gate.Result)
			results[gate.Result] = processGate(gate, results)
		}

		fmt.Printf("round %d ends; %d unresolved, %d results\n", i, numUnresolved, len(results))
		if numUnresolved == 0 {
			break
		}

		if numUnresolved == lastNumUnresolved {
			log.Fatalf("no convergence; results has %v", results)
		}
		lastNumUnresolved = numUnresolved
	}

	resultNames := []string{}
	for name, _ := range results {
		resultNames = append(resultNames, name)
	}
	sort.Strings(resultNames)

	for _, name := range resultNames {
		fmt.Printf("%v: %v\n", name, results[name])
	}
}
