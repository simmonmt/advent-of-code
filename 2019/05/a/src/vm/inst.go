package vm

import "fmt"

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
	Execute(ram Ram, pc int) (npc int)
	String() string
}

type Add struct {
	a, b, c Operand
}

func (i *Add) Execute(ram Ram, pc int) (npc int) {
	i.c.Write(ram, pc, i.a.Read(ram, pc)+i.b.Read(ram, pc))
	npc = pc + 4
	return
}

func (i *Add) String() string {
	return fmt.Sprintf("add %s, %s => %s", i.a, i.b, i.c)
}

type Multiply struct {
	a, b, c Operand
}

func (i *Multiply) Execute(ram Ram, pc int) (npc int) {
	i.c.Write(ram, pc, i.a.Read(ram, pc)*i.b.Read(ram, pc))
	npc = pc + 4
	return
}

func (i *Multiply) String() string {
	return fmt.Sprintf("mul %s, %s => %s", i.a, i.b, i.c)
}

type Halt struct{}

func (i *Halt) Execute(ram Ram, pc int) (npc int) {
	return -1
}

func (i *Halt) String() string {
	return "hlt"
}
