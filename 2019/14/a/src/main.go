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
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	input    = flag.String("input", "", "input file")
	wantFuel = flag.Int("wantfuel", 1, "desired fuel")
)

type Ingredient struct {
	Name string
	Qty  int
}

type Formula struct {
	Out Ingredient
	In  []Ingredient
}

func readInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func makeIngredient(s string) (Ingredient, error) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return Ingredient{}, fmt.Errorf("bad #parts in %v", s)
	}

	name := parts[1]
	qty, err := strconv.Atoi(parts[0])
	if err != nil {
		return Ingredient{}, err
	}

	return Ingredient{name, qty}, nil
}

func parseReactions(lines []string) (map[string]Formula, error) {
	reactions := map[string]Formula{}

	for _, line := range lines {
		parts := strings.Split(line, " => ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid =>")
		}

		output, err := makeIngredient(parts[1])
		if err != nil {
			return nil, fmt.Errorf("bad result: %v", err)
		}

		inputs := []Ingredient{}
		for i, s := range strings.Split(parts[0], ", ") {
			input, err := makeIngredient(s)
			if err != nil {
				return nil, fmt.Errorf("bad input %d in %s: %v", i, line, err)
			}

			inputs = append(inputs, input)
		}

		if _, found := reactions[output.Name]; found {
			panic(fmt.Sprintf("double for %s", output.Name))
		}
		reactions[output.Name] = Formula{
			Out: output,
			In:  inputs,
		}
	}

	return reactions, nil
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

	reactions, err := parseReactions(lines)
	if err != nil {
		log.Fatal(err)
	}

	have := map[string]int{"FUEL": *wantFuel}
	extra := map[string]int{}
	for {
		names := []string{}
		for name := range have {
			names = append(names, name)
		}
		sort.Strings(names)

		if len(names) == 1 && names[0] == "ORE" {
			break
		}

		var name string
		for _, name = range names {
			if name != "ORE" {
				break
			}
		}
		if name == "" {
			panic("no name")
		}
		qty := have[name]

		if extraQty, found := extra[name]; found {
			qty -= extraQty
			extra[name] = 0
			//fmt.Printf("found %d extra %s, qty now %d\n", extraQty, name, qty)
		}

		formula, found := reactions[name]
		if !found {
			panic(fmt.Sprintf("no reaction for %v", name))
		}

		formulaTimes := (qty + formula.Out.Qty - 1) / formula.Out.Qty
		left := (formulaTimes * formula.Out.Qty) - qty
		//fmt.Printf("have %d %s, formula %v making it %dx (left %d)\n",
		//qty, name, formula, formulaTimes, left)

		if left > 0 {
			extra[formula.Out.Name] += left
		}

		for _, ing := range formula.In {
			addQty := ing.Qty * formulaTimes
			have[ing.Name] += addQty
		}

		delete(have, name)

		//fmt.Printf("have: %v, extra: %v\n", have, extra)
	}

	fmt.Printf("%d FUEL needs %d ORE (delta from 1T: %d)\n",
		*wantFuel, have["ORE"], (1e12 - have["ORE"]))
}
