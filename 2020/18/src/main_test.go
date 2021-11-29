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
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2020/18/src/parse"
	"github.com/simmonmt/aoc/2020/common/logger"
)

func TestEvalA(t *testing.T) {
	type TestCase struct {
		in   string
		want int
	}

	testCases := []TestCase{
		TestCase{"1 + 2 * 3 + 4 * 5 + 6", 71},
		TestCase{"1 + (2 * 3) + (4 * (5 + 6))", 51},
		TestCase{"2 * 3 + (4 * 5)", 26},
		TestCase{"5 + (8 * 3 + 9 + 3 * 4 * 3)", 437},
		TestCase{"5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))", 12240},
		TestCase{"((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2", 13632},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := EvalA(parse.Parse(tc.in)); got != tc.want {
				t.Errorf("Eval(_) = %d, want %d", got, tc.want)
			}
		})
	}
}

func immNode(imm int) *parse.Node {
	return &parse.Node{
		Type: parse.TYPE_IMM,
		Imm:  imm,
	}
}

func immNodes(imms ...int) []*parse.Node {
	out := []*parse.Node{}
	for _, imm := range imms {
		out = append(out, immNode(imm))
	}
	return out
}

func immsFromNodes(nodes []*parse.Node) []int {
	out := []int{}
	for _, node := range nodes {
		if node.Type != parse.TYPE_IMM {
			panic("bad type")
		}
		out = append(out, node.Imm)
	}
	return out
}

func TestReplace(t *testing.T) {
	in := immNodes(0, 1, 2, 3, 4, 5)
	rep := immNode(99)

	type TestCase struct {
		start, end int
		want       []int
	}

	testCases := []TestCase{
		TestCase{0, 2, []int{99, 3, 4, 5}},
		TestCase{1, 3, []int{0, 99, 4, 5}},
		TestCase{3, 5, []int{0, 1, 2, 99}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := immsFromNodes(replace(in, tc.start, tc.end, rep))
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("replace(_, %d, %d, _) = %v, want %v",
					tc.start, tc.end, got, tc.want)
			}
		})
	}
}

func TestEvalB(t *testing.T) {
	type TestCase struct {
		in   string
		want int
	}

	testCases := []TestCase{
		TestCase{"1 + 2 * 3 + 4 * 5 + 6", 231},
		TestCase{"1 + (2 * 3) + (4 * (5 + 6))", 51},
		TestCase{"2 * 3 + (4 * 5)", 46},
		TestCase{"5 + (8 * 3 + 9 + 3 * 4 * 3)", 1445},
		TestCase{"5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))", 669060},
		TestCase{"((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2", 23340},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := EvalB(parse.Parse(tc.in)); got != tc.want {
				t.Errorf("Eval(_) = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
