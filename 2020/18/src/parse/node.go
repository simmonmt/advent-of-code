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

type NodeType int

const (
	TYPE_IMM NodeType = iota
	TYPE_ADD
	TYPE_MULT
	TYPE_EXPR
)

func (t NodeType) String() string {
	switch t {
	case TYPE_IMM:
		return "imm"
	case TYPE_ADD:
		return "add"
	case TYPE_MULT:
		return "mult"
	case TYPE_EXPR:
		return "expr"
	default:
		return "???"
	}
}

type Node struct {
	Type NodeType
	Imm  int
	Expr []*Node
}
