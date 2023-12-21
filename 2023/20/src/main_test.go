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
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string

	//go:embed sample2.txt
	rawSample2   string
	sample2Lines []string

	leftLines = []string{
		"broadcaster -> jp",
		"%jp -> pz, vr",
		"%vr -> hp",
		"%hp -> tx",
		"%tx -> dx",
		"%dx -> pz, ph",
		"%ph -> mb, pz",
		"%mb -> jc",
		"%jc -> pz, kt",
		"%kt -> ct",
		"%ct -> kd, pz",
		"%kd -> pz, pp",
		"%pp -> pz",
		"&pz -> kt, pg, mb, vr, hp, jp, tx",
	}
)

func dumpStates(i int, graph map[string]Node, order []string) {
	fmt.Printf("%4d:", i)
	for _, name := range order {
		node := graph[name]
		fmt.Printf(" %-10s", node)
	}
	fmt.Println()
}

func TestSolveA(t *testing.T) {
	ins := [][]string{sampleLines, sample2Lines}
	wants := []int{8000 * 4000, 4250 * 2750}

	for i := range ins {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			input, err := parseInput(ins[i])
			if err != nil {
				t.Fatal(err)
			}

			if got, want := solveA(input), wants[i]; got != want {
				t.Errorf("solveA(sample %d) = %v, want %v", i, got, want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleLines = strings.Split(rawSample, "\n")
	if len(sampleLines) > 0 && sampleLines[len(sampleLines)-1] == "" {
		sampleLines = sampleLines[0 : len(sampleLines)-1]
	}

	sample2Lines = strings.Split(rawSample2, "\n")
	if len(sample2Lines) > 0 && sample2Lines[len(sample2Lines)-1] == "" {
		sample2Lines = sample2Lines[0 : len(sample2Lines)-1]
	}

	os.Exit(m.Run())
}
