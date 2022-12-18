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
	_ "embed"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2022/common/graph"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestParseInput(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	want := []*InputNode{
		&InputNode{"AA", 0, []string{"DD", "II", "BB"}},
		&InputNode{"BB", 13, []string{"CC", "AA"}},
		&InputNode{"CC", 2, []string{"DD", "BB"}},
		&InputNode{"DD", 20, []string{"CC", "AA", "EE"}},
		&InputNode{"EE", 3, []string{"FF", "DD"}},
		&InputNode{"FF", 0, []string{"EE", "GG"}},
		&InputNode{"GG", 0, []string{"FF", "HH"}},
		&InputNode{"HH", 22, []string{"GG"}},
		&InputNode{"II", 0, []string{"AA", "JJ"}},
		&InputNode{"JJ", 21, []string{"II"}},
	}

	if !reflect.DeepEqual(input, want) {
		t.Errorf("parseInput(sampleLines) = %v, want %v",
			input, want)
	}
}

func TestSimplifyInputGraph(t *testing.T) {
	nodes, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	g := simplifyInputGraph(nodes)

	wants := map[graph.NodeID]map[graph.NodeID]int{
		"AA": map[graph.NodeID]int{
			"BB": 2, "CC": 3, "DD": 2, "EE": 3, "HH": 6, "JJ": 3,
		},
		"BB": map[graph.NodeID]int{
			"CC": 2, "DD": 3, "EE": 4, "HH": 7, "JJ": 4,
		},
		"CC": map[graph.NodeID]int{
			"BB": 2, "DD": 2, "EE": 3, "HH": 6, "JJ": 5,
		},
		"DD": map[graph.NodeID]int{
			"BB": 3, "CC": 2, "EE": 2, "HH": 5, "JJ": 4,
		},
		"EE": map[graph.NodeID]int{
			"BB": 4, "CC": 3, "DD": 2, "HH": 4, "JJ": 5,
		},
		"HH": map[graph.NodeID]int{
			"BB": 7, "CC": 6, "DD": 5, "EE": 4, "JJ": 8,
		},
		"JJ": map[graph.NodeID]int{
			"BB": 4, "CC": 5, "DD": 4, "EE": 5, "HH": 8,
		},
	}

	for from, want := range wants {
		got := g.AllEdges(from)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("edges for %v = %v, want %v", from, got, want)
		}
	}
}

func TestPathManagerRelease(t *testing.T) {
	nodes, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	g := simplifyInputGraph(nodes)

	pathManager := NewPathManager(g, "AA")
	left := 30

	path := []graph.NodeID{"AA", "DD", "BB", "JJ", "HH", "EE", "CC"}
	for i, id := range path {
		if i > 0 {
			left -= g.EdgeCost(path[i-1], id)
		}
		pathManager.Visit(id, left)

	}

	if got, want := pathManager.MaxReleased(), 1651; got != want {
		t.Errorf("maxreleased = %v, want %v", got, want)
	}
}

func TestSolveA(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(input), 1651; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), -1; got != want {
		t.Errorf("solveB(sample) = %v, want %v", got, want)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleLines = strings.Split(rawSample, "\n")
	if len(sampleLines) > 0 && sampleLines[len(sampleLines)-1] == "" {
		sampleLines = sampleLines[0 : len(sampleLines)-1]
	}

	os.Exit(m.Run())
}
