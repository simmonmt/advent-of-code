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
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2020/18/src/parse"
	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func EvalA(node *parse.Node) int {
	if node.Type == parse.TYPE_IMM {
		return node.Imm
	}
	if node.Type != parse.TYPE_EXPR {
		panic("bad type")
	}

	subs := node.Expr
	if len(subs)%2 != 1 {
		panic("bad subs num")
	}

	num := EvalA(subs[0])

	for i := 1; i < len(subs); i += 2 {
		imm := EvalA(subs[i+1])
		if subs[i].Type == parse.TYPE_ADD {
			num += imm
		} else {
			num *= imm
		}
	}

	return num
}

func replace(arr []*parse.Node, start, end int, rep *parse.Node) []*parse.Node {
	out := []*parse.Node{}

	if start > 0 {
		out = append(out, arr[0:start]...)
	}

	out = append(out, rep)

	if end < len(arr)-1 {
		out = append(out, arr[end+1:]...)
	}

	return out
}

func EvalB(node *parse.Node) int {
	if node.Type == parse.TYPE_IMM {
		return node.Imm
	}
	if node.Type != parse.TYPE_EXPR {
		panic("bad type")
	}

	flat := []*parse.Node{}
	for _, sub := range node.Expr {
		switch sub.Type {
		case parse.TYPE_IMM:
			fallthrough
		case parse.TYPE_ADD:
			fallthrough
		case parse.TYPE_MULT:
			flat = append(flat, sub)
		case parse.TYPE_EXPR:
			imm := &parse.Node{Type: parse.TYPE_IMM, Imm: EvalB(sub)}
			flat = append(flat, imm)
		default:
			panic("bad type")
		}
	}

	collapse := func(nodeType parse.NodeType, combine func(a, b int) int) {
		for len(flat) > 1 {
			opIdx := -1
			for i := 1; i < len(flat); i += 2 {
				if flat[i].Type == nodeType {
					opIdx = i
					break
				}
			}

			if opIdx == -1 {
				break // no more to combine
			}

			a := flat[opIdx-1].Imm
			b := flat[opIdx+1].Imm

			combined := &parse.Node{
				Type: parse.TYPE_IMM,
				Imm:  combine(a, b),
			}

			flat = replace(flat, opIdx-1, opIdx+1, combined)
		}
	}

	collapse(parse.TYPE_ADD, func(a, b int) int { return a + b })
	collapse(parse.TYPE_MULT, func(a, b int) int { return a * b })

	return flat[0].Imm
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	sumA, sumB := 0, 0
	for _, line := range lines {
		sumA += EvalA(parse.Parse(line))
		sumB += EvalB(parse.Parse(line))
	}

	fmt.Printf("A: %d\n", sumA)
	fmt.Printf("B: %d\n", sumB)
}
