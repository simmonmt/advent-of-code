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
	"regexp"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	workflowPattern = regexp.MustCompile(`^([^\}]+)\{([^\}]+)\}$$`)
)

type Rule struct {
	Condition string
	Action    string
}

type Workflow struct {
	Name  string
	Rules []Rule
}

func parseWorkflow(line string) (*Workflow, error) {
	parts := workflowPattern.FindStringSubmatch(line)
	if parts == nil {
		return nil, fmt.Errorf("regexp mismatch")
	}

	name := parts[1]

	rules := []Rule{}
	ruleStrs := strings.Split(parts[2], ",")
	for _, ruleStr := range ruleStrs {
		condition, action, ok := strings.Cut(ruleStr, ":")
		if !ok {
			rules = append(rules, Rule{Action: ruleStr})
		} else {
			rules = append(rules, Rule{
				Condition: condition,
				Action:    action,
			})
		}
	}

	return &Workflow{Name: name, Rules: rules}, nil
}

func parseWorkflows(lines []string) (map[string]*Workflow, error) {
	workflows := map[string]*Workflow{}
	for i, line := range lines {
		workflow, err := parseWorkflow(line)
		if err != nil {
			return nil, fmt.Errorf("bad workflow %d: %v", i+1, err)
		}
		workflows[workflow.Name] = workflow
	}
	return workflows, nil
}

func parsePart(line string) (map[string]int, error) {
	if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
		return nil, fmt.Errorf("bad {}")
	}

	line = line[1 : len(line)-1]

	part := map[string]int{}
	kvs := strings.Split(line, ",")
	for _, kv := range kvs {
		n, vStr, ok := strings.Cut(kv, "=")
		if !ok {
			return nil, fmt.Errorf("bad kv")
		}

		v, err := strconv.Atoi(vStr)
		if err != nil {
			return nil, fmt.Errorf("bad val: %v", err)
		}

		part[n] = v
	}

	return part, nil
}

func parseParts(lines []string) ([]map[string]int, error) {
	parts := []map[string]int{}
	for i, line := range lines {
		part, err := parsePart(line)
		if err != nil {
			return nil, fmt.Errorf("bad part %d: %v", i+1, err)
		}

		parts = append(parts, part)
	}
	return parts, nil
}

func parseInput(lines []string) (map[string]*Workflow, []map[string]int, error) {
	groups, err := filereader.BlankSeparatedGroupsFromLines(lines)
	if err != nil || len(groups) != 2 {
		return nil, nil, fmt.Errorf("bad groups")
	}

	workflows, err := parseWorkflows(groups[0])
	if err != nil {
		return nil, nil, fmt.Errorf("bad workflows: %v", err)
	}

	parts, err := parseParts(groups[1])
	if err != nil {
		return nil, nil, fmt.Errorf("bad parts: %v", err)
	}

	return workflows, parts, nil
}

func evaluateCondition(condition string, part map[string]int) bool {
	if condition == "" {
		return true
	}

	key := condition[0]
	op := condition[1]
	num, err := strconv.Atoi(condition[2:])
	if err != nil {
		panic("bad num")
	}

	v := part[string(key)]
	switch op {
	case '<':
		return v < num
	case '>':
		return v > num
	default:
		panic(fmt.Sprintf("bad op %c", op))
	}
}

func handlePart(workflows map[string]*Workflow, part map[string]int) bool {
	cur := "in"
	for {
		workflow := workflows[cur]
		for _, rule := range workflow.Rules {
			if evaluateCondition(rule.Condition, part) {
				if rule.Action == "A" {
					return true
				} else if rule.Action == "R" {
					return false
				} else {
					cur = rule.Action
					break
				}
			}
		}
	}
}

func solveA(workflows map[string]*Workflow, parts []map[string]int) int {
	sum := 0
	for _, part := range parts {
		if handlePart(workflows, part) {
			for _, v := range part {
				sum += v
			}
		}
	}
	return sum
}

type Acceptable struct {
	Mins map[rune]int
	Maxs map[rune]int
	Path []string
}

func (a *Acceptable) String() string {
	out := ""
	for _, r := range []rune{'x', 'm', 'a', 's'} {
		if out != "" {
			out += ","
		}
		out += fmt.Sprintf("%c=[%d,%d]", r, a.Mins[r], a.Maxs[r])
	}
	out += "<=" + strings.Join(a.Path, ",")

	return out
}

func (a *Acceptable) Clone() *Acceptable {
	out := &Acceptable{
		Mins: map[rune]int{},
		Maxs: map[rune]int{},
	}

	for k, v := range a.Mins {
		out.Mins[k] = v
	}
	for k, v := range a.Maxs {
		out.Maxs[k] = v
	}

	out.Path = make([]string, len(a.Path))
	copy(out.Path, a.Path)

	return out
}

func applyCondition(condition string, cur *Acceptable) *Acceptable {
	if condition == "" {
		return cur.Clone()
	}

	key := rune(condition[0])
	op := condition[1]
	num, err := strconv.Atoi(condition[2:])
	if err != nil {
		panic("bad num")
	}

	with := cur.Clone()

	switch op {
	case '<':
		with.Maxs[key] = num - 1
		cur.Mins[key] = num
	case '>':
		with.Mins[key] = num + 1
		cur.Maxs[key] = num
	default:
		panic("bad op")
	}

	return with
}

func solveB(workflows map[string]*Workflow, parts []map[string]int) int64 {
	curs := map[string]*Acceptable{
		"in": &Acceptable{
			Mins: map[rune]int{'x': 1, 'm': 1, 'a': 1, 's': 1},
			Maxs: map[rune]int{'x': 4000, 'm': 4000, 'a': 4000, 's': 4000},
		},
	}

	accepted := []*Acceptable{}

	for len(curs) > 0 {
		out := map[string]*Acceptable{}

		for name, cur := range curs {
			cur.Path = append(cur.Path, name)
			for _, rule := range workflows[name].Rules {
				nc := applyCondition(rule.Condition, cur)

				if rule.Action == "A" {
					accepted = append(accepted, nc)
				} else if rule.Action != "R" {
					out[rule.Action] = nc
				}
			}
		}

		curs = out
	}

	var sum int64
	for _, acc := range accepted {
		//fmt.Println(acc)

		tot := int64(1)
		for _, r := range []rune{'x', 'm', 'a', 's'} {
			tot *= int64(acc.Maxs[r] - acc.Mins[r] + 1)
		}
		sum += tot
	}
	return sum
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

	workflows, parts, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(workflows, parts))
	fmt.Println("B", solveB(workflows, parts))
}
