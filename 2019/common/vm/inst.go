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

type PositionOperand struct {
	loc int64
}

func (o *PositionOperand) Read(r *Resources, pc int64) int64 {
	return r.ram.Read(o.loc)
}

func (o *PositionOperand) Write(r *Resources, pc, val int64) {
	r.ram.Write(o.loc, val)
}

func (o *PositionOperand) String() string {
	return fmt.Sprintf("*%v", o.loc)
}

type RelativeOperand struct {
	imm int64
}

func (o *RelativeOperand) Read(r *Resources, pc int64) int64 {
	return r.ram.Read(r.relBase + o.imm)
}

func (o *RelativeOperand) Write(r *Resources, pc, val int64) {
	r.ram.Write(r.relBase+o.imm, val)
}

func (o *RelativeOperand) String() string {
	return fmt.Sprintf("*R%v", o.imm)
}

type Instruction interface {
	Size() int64
	Execute(r *Resources, pc int64) (npc int64)
	String() string
}

type Add struct {
	a, b, c Operand
}

func (i *Add) Size() int64 {
	return 4
}

func (i *Add) Execute(r *Resources, pc int64) (npc int64) {
	a := i.a.Read(r, pc)
	b := i.b.Read(r, pc)
	out := a + b
	logger.LogF("add exec: %d + %d (=%d) => %s", a, b, out, i.c)
	i.c.Write(r, pc, out)
	npc = pc + i.Size()
	return
}

func (i *Add) String() string {
	return fmt.Sprintf("add %s, %s => %s", i.a, i.b, i.c)
}

type Multiply struct {
	a, b, c Operand
}

func (i *Multiply) Size() int64 {
	return 4
}

func (i *Multiply) Execute(r *Resources, pc int64) (npc int64) {
	i.c.Write(r, pc, i.a.Read(r, pc)*i.b.Read(r, pc))
	npc = pc + i.Size()
	return
}

func (i *Multiply) String() string {
	return fmt.Sprintf("mul %s, %s => %s", i.a, i.b, i.c)
}

type Input struct {
	a Operand
}

func (i *Input) Size() int64 {
	return 2
}

func (i *Input) Execute(r *Resources, pc int64) (npc int64) {
	in := r.io.Read()
	logger.LogF("in exec: %d => %s", in, i.a)
	i.a.Write(r, pc, in)
	return pc + i.Size()
}

func (i *Input) String() string {
	return fmt.Sprintf("in => %s", i.a)
}

type Output struct {
	a Operand
}

func (i *Output) Size() int64 {
	return 2
}

func (i *Output) Execute(r *Resources, pc int64) (npc int64) {
	out := i.a.Read(r, pc)
	logger.LogF("out exec: write %d", out)
	r.io.Write(i.a.Read(r, pc))
	return pc + i.Size()
}

func (i *Output) String() string {
	return fmt.Sprintf("out %s", i.a)
}

type Halt struct{}

func (i *Halt) Size() int64 {
	return 1
}

func (i *Halt) Execute(r *Resources, pc int64) (npc int64) {
	return -1
}

func (i *Halt) String() string {
	return "hlt"
}

type JumpIfTrue struct {
	a, b Operand
}

func (i *JumpIfTrue) Size() int64 {
	return 3
}

func (i *JumpIfTrue) Execute(r *Resources, pc int64) (npc int64) {
	a := i.a.Read(r, pc)
	b := i.b.Read(r, pc)
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

func (i *JumpIfFalse) Size() int64 {
	return 3
}

func (i *JumpIfFalse) Execute(r *Resources, pc int64) (npc int64) {
	a := i.a.Read(r, pc)
	b := i.b.Read(r, pc)
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

func (i *LessThan) Size() int64 {
	return 4
}

func (i *LessThan) Execute(r *Resources, pc int64) (npc int64) {
	a := i.a.Read(r, pc)
	b := i.b.Read(r, pc)

	var out int64 = 0
	if a < b {
		out = 1
	}

	logger.LogF("lt exec: %d<%d? %d => %d", a, b, out, i.c)
	i.c.Write(r, pc, out)
	npc = pc + i.Size()
	return
}

func (i *LessThan) String() string {
	return fmt.Sprintf("lt %s<%s => %s", i.a, i.b, i.c)
}

type Equals struct {
	a, b, c Operand
}

func (i *Equals) Size() int64 {
	return 4
}

func (i *Equals) Execute(r *Resources, pc int64) (npc int64) {
	a := i.a.Read(r, pc)
	b := i.b.Read(r, pc)

	var out int64 = 0
	if a == b {
		out = 1
	}

	logger.LogF("eq exec: %d==%d? %d => %d", a, b, out, i.c)
	i.c.Write(r, pc, out)
	npc = pc + i.Size()
	return
}

func (i *Equals) String() string {
	return fmt.Sprintf("eq %s==%s => %s", i.a, i.b, i.c)
}

type SetRelBase struct {
	a Operand
}

func (i *SetRelBase) Size() int64 {
	return 2
}

func (i *SetRelBase) Execute(r *Resources, pc int64) (npc int64) {
	a := i.a.Read(r, pc)
	logger.LogF("setrelbase exec: old %v + %d", r.relBase, a)
	r.relBase += a
	npc = pc + i.Size()
	return
}

func (i *SetRelBase) String() string {
	return fmt.Sprintf("setrelbase %s", i.a)
}
