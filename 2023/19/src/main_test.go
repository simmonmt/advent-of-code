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
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/simmonmt/aoc/2023/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestParseInput(t *testing.T) {
	workflows, parts, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	wantWorkflows := map[string]*Workflow{
		"crn": &Workflow{
			Name:  "crn",
			Rules: []Rule{{Condition: "x>2662", Action: "A"}, {Action: "R"}},
		},
		"gd": &Workflow{
			Name:  "gd",
			Rules: []Rule{{Condition: "a>3333", Action: "R"}, {Action: "R"}},
		},
		"hdj": &Workflow{
			Name:  "hdj",
			Rules: []Rule{{Condition: "m>838", Action: "A"}, {Action: "pv"}},
		},
		"in": &Workflow{
			Name:  "in",
			Rules: []Rule{{Condition: "s<1351", Action: "px"}, {Action: "qqz"}},
		},
		"lnx": &Workflow{
			Name:  "lnx",
			Rules: []Rule{{Condition: "m>1548", Action: "A"}, {Action: "A"}},
		},
		"pv": &Workflow{
			Name:  "pv",
			Rules: []Rule{{Condition: "a>1716", Action: "R"}, {Action: "A"}},
		},
		"px": &Workflow{
			Name:  "px",
			Rules: []Rule{{Condition: "a<2006", Action: "qkq"}, {Condition: "m>2090", Action: "A"}, {Action: "rfg"}},
		},
		"qkq": &Workflow{
			Name:  "qkq",
			Rules: []Rule{{Condition: "x<1416", Action: "A"}, {Action: "crn"}},
		},
		"qqz": &Workflow{
			Name:  "qqz",
			Rules: []Rule{{Condition: "s>2770", Action: "qs"}, {Condition: "m<1801", Action: "hdj"}, {Action: "R"}},
		},
		"qs": &Workflow{
			Name:  "qs",
			Rules: []Rule{{Condition: "s>3448", Action: "A"}, {Action: "lnx"}},
		},
		"rfg": &Workflow{
			Name:  "rfg",
			Rules: []Rule{{Condition: "s<537", Action: "gd"}, {Condition: "x>2440", Action: "R"}, {Action: "A"}},
		},
	}

	wantParts := []map[string]int{
		{"a": 1222, "m": 2655, "s": 2876, "x": 787},
		{"a": 2067, "m": 44, "s": 496, "x": 1679},
		{"a": 79, "m": 264, "s": 2244, "x": 2036},
		{"a": 466, "m": 1339, "s": 291, "x": 2461},
		{"a": 2188, "m": 1623, "s": 1013, "x": 2127},
	}

	if diff := cmp.Diff(wantWorkflows, workflows); diff != "" {
		t.Errorf("workflows mismatch; -want,+got:\n%s\n", diff)
	}

	if diff := cmp.Diff(wantParts, parts); diff != "" {
		t.Errorf("parts mismatch; -want,+got:\n%s\n", diff)
	}
}

func TestSolveA(t *testing.T) {
	workflows, parts, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(workflows, parts), 19114; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	workflows, parts, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(workflows, parts), int64(167409079868000); got != want {
		t.Logf("want %d", want)
		t.Logf("got  %d", got)
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
