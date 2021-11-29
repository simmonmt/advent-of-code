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

package parse

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2020/common/logger"
)

func Flatten(node *Node) string {
	switch node.Type {
	case TYPE_IMM:
		return strconv.Itoa(node.Imm)
	case TYPE_ADD:
		return "+"
	case TYPE_MULT:
		return "*"
	case TYPE_EXPR:
		subs := []string{}
		for _, expr := range node.Expr {
			subs = append(subs, Flatten(expr))
		}
		return fmt.Sprintf("(%s)", strings.Join(subs, " "))
	default:
		panic("bad node type")
	}
}

func TestFlatten(t *testing.T) {
	type TestCase struct {
		in   *Node
		want string
	}

	testCases := []TestCase{
		TestCase{
			in:   &Node{Type: TYPE_IMM, Imm: 1},
			want: "1",
		},
		TestCase{
			in:   &Node{Type: TYPE_ADD},
			want: "+",
		},
		TestCase{
			in:   &Node{Type: TYPE_MULT},
			want: "*",
		},
		TestCase{
			in: &Node{
				Type: TYPE_EXPR,
				Expr: []*Node{
					&Node{Type: TYPE_IMM, Imm: 1},
					&Node{Type: TYPE_ADD},
					&Node{Type: TYPE_IMM, Imm: 3},
				},
			},
			want: "(1 + 3)",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := Flatten(tc.in); got != tc.want {
				t.Errorf("case %d: got %v, want %v",
					i, got, tc.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []string{
		"1 + 2 * 3 + 4 * 5 + 6",
		"1 + (2 * 3) + (4 * (5 + 6))",
		"((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2",
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			node, sz := doParse([]byte(tc))
			flat := Flatten(node)

			want := "(" + tc + ")"

			if flat != want {
				t.Errorf("flat got %v, want %v", flat, want)
			}
			if sz != len(tc) {
				t.Errorf("len got %v, want %v", sz, len(want))
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
