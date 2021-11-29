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
	"unicode"
	"unicode/utf8"
)

func eatNum(in []byte) (num, numSz int) {
	num = 0
	numSz = 0
	for len(in) > 0 {
		r, sz := utf8.DecodeRune(in)
		if r == utf8.RuneError {
			panic("rune error")
		}

		if !unicode.IsDigit(r) {
			break
		}

		num = num*10 + int(r-'0')
		numSz += sz
		in = in[sz:]
	}
	return
}

func Parse(in string) *Node {
	n, _ := doParse([]byte(in))
	return n
}

func doParse(in []byte) (*Node, int) {
	//logger.LogF("parse %v", string(in))
	consumed := 0

	node := &Node{Type: TYPE_EXPR, Expr: []*Node{}}

	for len(in) > 0 {
		r, sz := utf8.DecodeRune(in)
		if r == utf8.RuneError {
			panic("rune error")
		}

		//logger.LogF("decoded rune %v (consumed %d)",
		//strconv.QuoteRune(r), consumed)

		if unicode.IsDigit(r) {
			num, numSz := eatNum(in)
			imm := &Node{Type: TYPE_IMM, Imm: num}
			node.Expr = append(node.Expr, imm)

			sz = numSz
		} else if unicode.IsSpace(r) {
			// Nothing; just eat it

		} else if r == '+' {
			node.Expr = append(node.Expr, &Node{Type: TYPE_ADD})
		} else if r == '*' {
			node.Expr = append(node.Expr, &Node{Type: TYPE_MULT})
		} else if r == '(' {
			sub, subSz := doParse(in[sz:])
			//logger.LogF("back from parse, sz %v", subSz)
			node.Expr = append(node.Expr, sub)
			sz += subSz
		} else if r == ')' {
			return node, consumed + sz
		} else {
			panic("bad char")
		}

		in = in[sz:]
		consumed += sz
	}

	return node, consumed
}
