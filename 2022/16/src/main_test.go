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
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
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

func sortFutures(in [][]Action) []string {
	out := []string{}
	for _, group := range in {
		gs := []string{}
		for _, action := range group {
			gs = append(gs, fmt.Sprintf("%v", action))
		}
		out = append(out, strings.Join(gs, ","))
	}
	sort.Strings(out)
	return out
}

func makeClaimed(ids ...graph.NodeID) map[graph.NodeID]bool {
	out := map[graph.NodeID]bool{}
	for _, id := range ids {
		out[id] = true
	}
	return out
}

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

func TestFindAvailableNeighbors(t *testing.T) {
	g := simplifyInputGraph([]*InputNode{
		&InputNode{"AA", 0, []string{"BB"}},
		&InputNode{"BB", 1, []string{"AA", "CC"}},
		&InputNode{"CC", 2, []string{"BB"}},
	})

	type TestCase struct {
		id      graph.NodeID
		left    int
		claimed map[graph.NodeID]bool
		want    []graph.NodeID
	}

	testCases := []TestCase{
		TestCase{"AA", 30, makeClaimed("AA"), []graph.NodeID{"BB", "CC"}},
		TestCase{"BB", 30, makeClaimed("AA"), []graph.NodeID{"CC"}},
		TestCase{"BB", 30, makeClaimed("AA", "CC"), []graph.NodeID{}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := findAvailableNeighbors(g, tc.id, tc.left, tc.claimed)
			sort.Slice(got, func(i, j int) bool {
				return got[i] < got[j]
			})

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("findAvailableNeighbors(_, %v, %v, %v) = %v, want %v",
					tc.id, tc.left, tc.claimed, got, tc.want)
			}
		})
	}
}

func TestExecuteMinute(t *testing.T) {
	g := simplifyInputGraph([]*InputNode{
		&InputNode{"AA", 0, []string{"BB"}},
		&InputNode{"BB", 1, []string{"AA", "CC"}},
		&InputNode{"CC", 2, []string{"BB"}},
	})

	type TestCase struct {
		release Release
		claimed map[graph.NodeID]bool
		players []PlayerState

		wantRelease Release
		wantPlayers []PlayerState
		wantFutures [][]Action
	}

	testCases := []TestCase{
		TestCase{
			release: Release{},
			claimed: makeClaimed("AA"),
			players: []PlayerState{PlayerState{"AA", 0}},

			wantRelease: Release{total: 0},
			wantPlayers: []PlayerState{PlayerState{"AA", 0}},
			wantFutures: [][]Action{
				[]Action{Action{0, "BB"}},
				[]Action{Action{0, "CC"}},
				[]Action{}, // no action
			},
		},

		TestCase{ // Getting to BB bumps rate
			release: Release{},
			claimed: makeClaimed("AA"),
			players: []PlayerState{PlayerState{"BB", 1}},

			wantRelease: Release{total: 0, rate: 1},
			wantPlayers: []PlayerState{PlayerState{"BB", 0}},
			wantFutures: [][]Action{[]Action{}},
		},

		TestCase{ // Increment release but nowhere to go
			release: Release{total: 3, rate: 4},
			claimed: makeClaimed("AA", "BB"),
			players: []PlayerState{PlayerState{"CC", 0}},

			wantRelease: Release{total: 7, rate: 4},
			wantPlayers: []PlayerState{PlayerState{"CC", 0}},
			wantFutures: [][]Action{[]Action{}},
		},

		TestCase{
			release: Release{},
			claimed: makeClaimed("AA"),
			players: []PlayerState{
				PlayerState{"AA", 0},
				PlayerState{"AA", 0},
			},

			wantRelease: Release{total: 0},
			wantPlayers: []PlayerState{
				PlayerState{"AA", 0},
				PlayerState{"AA", 0},
			},
			wantFutures: [][]Action{
				[]Action{}, // no action
				[]Action{Action{0, "BB"}},
				[]Action{Action{0, "CC"}},
				[]Action{Action{1, "BB"}},
				[]Action{Action{1, "CC"}},
				[]Action{Action{0, "BB"}, Action{1, "CC"}},
				[]Action{Action{0, "CC"}, Action{1, "BB"}},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			logger.LogF("start test case %v", i)

			release := tc.release
			players := make([]PlayerState, len(tc.players))
			copy(players, tc.players)
			futures := executeMinute(g, 1, 30, &release, players, tc.claimed)

			sortedFutures := sortFutures(futures)
			sortedWantFutures := sortFutures(tc.wantFutures)

			if !reflect.DeepEqual(release, tc.wantRelease) {
				t.Errorf("got release %v, want %v", release, tc.wantRelease)
			}
			if !reflect.DeepEqual(players, tc.wantPlayers) {
				t.Errorf("got players %v, want %v", players, tc.wantPlayers)
			}

			if !reflect.DeepEqual(sortedFutures, sortedWantFutures) {
				t.Errorf("got futures %v, want %v", futures,
					tc.wantFutures)
			}
		})
	}

}

func TestRunWorldSinglePlayer(t *testing.T) {
	g := simplifyInputGraph([]*InputNode{
		&InputNode{"AA", 0, []string{"BB"}},
		&InputNode{"BB", 1, []string{"AA", "CC"}},
		&InputNode{"CC", 2, []string{"BB"}},
	})

	if got, want := solve(g, 5, 1, "AA"), 5; got != want {
		t.Errorf("solve = %v, want %v", got, want)
	}
}

func TestRunWorldTwoPlayer(t *testing.T) {
	g := simplifyInputGraph([]*InputNode{
		&InputNode{"AA", 0, []string{"BB"}},
		&InputNode{"BB", 1, []string{"AA", "CC"}},
		&InputNode{"CC", 2, []string{"BB"}},
	})

	if got, want := solve(g, 5, 2, "AA"), 7; got != want {
		t.Errorf("solve = %v, want %v", got, want)
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

func NoTestSolveB(t *testing.T) {
	input, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(input), 1707; got != want {
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
