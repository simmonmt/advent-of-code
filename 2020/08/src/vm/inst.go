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

package vm

import (
	"fmt"
	"strings"

	"github.com/simmonmt/aoc/2020/common/intmath"
)

type Operand interface {
	Read(r *Resources, pc int64) int64
	Write(r *Resources, pc, val int64)
	String() string
}

type ImmediateOperand struct {
	imm int64
}

func (o *ImmediateOperand) Read(r *Resources, pc int64) int64 {
	return o.imm
}

func (o *ImmediateOperand) Write(r *Resources, pc, val int64) {
	panic("attempt to write immediate")
}

func (o *ImmediateOperand) String() string {
	return fmt.Sprintf("%v", o.imm)
}

type Instruction interface {
	Size() int64
	Execute(r *Resources, pc int64) (npc int64)
	Op() string
	String() string
}

type Acc struct {
	a Operand
}

func (i *Acc) A() Operand {
	return i.a
}

func (i *Acc) Size() int64 {
	return 1
}

func (i *Acc) Execute(r *Resources, pc int64) (npc int64) {
	r.Acc += i.a.Read(r, pc)
	return pc + i.Size()
}

func (i *Acc) Op() string {
	return "acc"
}

func (i *Acc) String() string {
	return fmt.Sprintf("acc %s", i.a)
}

type Nop struct {
	a Operand
}

func (i *Nop) A() Operand {
	return i.a
}

func (i *Nop) Size() int64 {
	return 1
}

func (i *Nop) Execute(r *Resources, pc int64) (npc int64) {
	return pc + i.Size()
}

func (i *Nop) Op() string {
	return "nop"
}

func (i *Nop) String() string {
	return fmt.Sprintf("nop %s", i.a)
}

type Jmp struct {
	a Operand
}

func (i *Jmp) A() Operand {
	return i.a
}

func (i *Jmp) Size() int64 {
	return 1
}

func (i *Jmp) Execute(r *Resources, pc int64) (npc int64) {
	npc = pc + i.a.Read(r, pc)
	return
}

func (i *Jmp) Op() string {
	return "jmp"
}

func (i *Jmp) String() string {
	return fmt.Sprintf("jmp %s", i.a)
}

var (
	factories = map[string]func(a Operand) Instruction{
		"acc": func(a Operand) Instruction { return &Acc{a: a} },
		"jmp": func(a Operand) Instruction { return &Jmp{a: a} },
		"nop": func(a Operand) Instruction { return &Nop{a: a} },
	}
)

func Decode(str string) (Instruction, error) {
	parts := strings.SplitN(str, " ", 2)
	op := parts[0]
	a := &ImmediateOperand{imm: int64(intmath.AtoiOrDie(parts[1]))}

	factory, found := factories[op]
	if !found {
		return nil, fmt.Errorf("invalid op %v", op)
	}

	return factory(a), nil
}

func NewInst(op string, a Operand) (Instruction, error) {
	factory, found := factories[op]
	if !found {
		return nil, fmt.Errorf("invalid op %v", op)
	}

	return factory(a), nil
}
