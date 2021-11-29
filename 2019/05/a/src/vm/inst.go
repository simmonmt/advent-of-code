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

	"github.com/simmonmt/aoc/2019/common/logger"
)

type Operand interface {
	Read(ram Ram, pc int) int
	Write(raw Ram, pc, val int)
	String() string
}

type ImmediateOperand struct {
	imm int
}

func (o *ImmediateOperand) Read(ram Ram, pc int) int {
	return o.imm
}

func (o *ImmediateOperand) Write(ram Ram, pc, val int) {
	panic("attempt to write immediate")
}

func (o *ImmediateOperand) String() string {
	return fmt.Sprintf("%v", o.imm)
}

type PositionOperand struct {
	loc int
}

func (o *PositionOperand) Read(ram Ram, pc int) int {
	return ram.Read(o.loc)
}

func (o *PositionOperand) Write(ram Ram, pc, val int) {
	ram.Write(o.loc, val)
}

func (o *PositionOperand) String() string {
	return fmt.Sprintf("*%v", o.loc)
}

type Instruction interface {
	Size() int
	Execute(ram Ram, io IO, pc int) (npc int)
	String() string
}

type Add struct {
	a, b, c Operand
}

func (i *Add) Size() int {
	return 4
}

func (i *Add) Execute(ram Ram, io IO, pc int) (npc int) {
	a := i.a.Read(ram, pc)
	b := i.b.Read(ram, pc)
	out := a + b
	logger.LogF("add exec: %d + %d (=%d) => %s", a, b, out, i.c)
	i.c.Write(ram, pc, out)
	npc = pc + i.Size()
	return
}

func (i *Add) String() string {
	return fmt.Sprintf("add %s, %s => %s", i.a, i.b, i.c)
}

type Multiply struct {
	a, b, c Operand
}

func (i *Multiply) Size() int {
	return 4
}

func (i *Multiply) Execute(ram Ram, io IO, pc int) (npc int) {
	i.c.Write(ram, pc, i.a.Read(ram, pc)*i.b.Read(ram, pc))
	npc = pc + i.Size()
	return
}

func (i *Multiply) String() string {
	return fmt.Sprintf("mul %s, %s => %s", i.a, i.b, i.c)
}

type Input struct {
	a Operand
}

func (i *Input) Size() int {
	return 2
}

func (i *Input) Execute(ram Ram, io IO, pc int) (npc int) {
	in := io.Read()
	logger.LogF("in exec: %d => %s", in, i.a)
	i.a.Write(ram, pc, in)
	return pc + i.Size()
}

func (i *Input) String() string {
	return fmt.Sprintf("in => %s", i.a)
}

type Output struct {
	a Operand
}

func (i *Output) Size() int {
	return 2
}

func (i *Output) Execute(ram Ram, io IO, pc int) (npc int) {
	io.Write(i.a.Read(ram, pc))
	return pc + i.Size()
}

func (i *Output) String() string {
	return fmt.Sprintf("out %s", i.a)
}

type Halt struct{}

func (i *Halt) Size() int {
	return 1
}

func (i *Halt) Execute(ram Ram, io IO, pc int) (npc int) {
	return -1
}

func (i *Halt) String() string {
	return "hlt"
}

type JumpIfTrue struct {
	a, b Operand
}

func (i *JumpIfTrue) Size() int {
	return 3
}

func (i *JumpIfTrue) Execute(ram Ram, io IO, pc int) (npc int) {
	a := i.a.Read(ram, pc)
	b := i.b.Read(ram, pc)
	logger.LogF("jit exec: %d ? goto %v", a, b)
	if a > 0 {
		npc = b
	} else {
		npc = pc + i.Size()
	}
	return
}

func (i *JumpIfTrue) String() string {
	return fmt.Sprintf("jit %s? to %s", i.a, i.b)
}

type JumpIfFalse struct {
	a, b Operand
}

func (i *JumpIfFalse) Size() int {
	return 3
}

func (i *JumpIfFalse) Execute(ram Ram, io IO, pc int) (npc int) {
	a := i.a.Read(ram, pc)
	b := i.b.Read(ram, pc)
	logger.LogF("jif exec: %d =0? goto %v", a, b)
	if a == 0 {
		npc = b
	} else {
		npc = pc + i.Size()
	}
	return
}

func (i *JumpIfFalse) String() string {
	return fmt.Sprintf("jif %s=0? to %s", i.a, i.b)
}

type LessThan struct {
	a, b, c Operand
}

func (i *LessThan) Size() int {
	return 4
}

func (i *LessThan) Execute(ram Ram, io IO, pc int) (npc int) {
	a := i.a.Read(ram, pc)
	b := i.b.Read(ram, pc)

	out := 0
	if a < b {
		out = 1
	}

	logger.LogF("lt exec: %d<%d? %d => %d", a, b, out, i.c)
	i.c.Write(ram, pc, out)
	npc = pc + i.Size()
	return
}

func (i *LessThan) String() string {
	return fmt.Sprintf("lt %s<%s => %s", i.a, i.b, i.c)
}

type Equals struct {
	a, b, c Operand
}

func (i *Equals) Size() int {
	return 4
}

func (i *Equals) Execute(ram Ram, io IO, pc int) (npc int) {
	a := i.a.Read(ram, pc)
	b := i.b.Read(ram, pc)

	out := 0
	if a == b {
		out = 1
	}

	logger.LogF("eq exec: %d==%d? %d => %d", a, b, out, i.c)
	i.c.Write(ram, pc, out)
	npc = pc + i.Size()
	return
}

func (i *Equals) String() string {
	return fmt.Sprintf("eq %s==%s => %s", i.a, i.b, i.c)
}
