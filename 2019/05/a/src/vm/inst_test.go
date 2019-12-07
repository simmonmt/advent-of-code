package vm

import (
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2019/05/a/src/testutils"
)

func CheckRam(t *testing.T, ram Ram, vals []int) {
	for i, val := range vals {
		if got := ram.Read(i); got != val {
			t.Errorf("verify mismatch at %v: got %v want %v", i, got, val)
		}
	}
}

func CheckEmptyOutput(t *testing.T, io IO) {
	if got := io.Written(); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("output = %v, want []", got)
	}
}

func TestImmediateOperand(t *testing.T) {
	ramVals := []int{10, 11, 12, 13, 14}
	ram := NewRam(ramVals...)

	var op Operand = &ImmediateOperand{imm: 2}

	if got := op.Read(ram, 0); got != 2 {
		t.Errorf("Read(ram, 0) = %d, want %d", got, 2)
	}

	testutils.AssertPanic(t, "write failed to panic",
		func() { op.Write(ram, 0, 99) })

	CheckRam(t, ram, ramVals)
}

func TestPositionOperand(t *testing.T) {
	ram := NewRam(10, 11, 12, 13, 14)

	var op Operand = &PositionOperand{loc: 2}

	if got := op.Read(ram, 0); got != 12 {
		t.Errorf("Read(ram, 0) = %d, want %d", got, 12)
	}

	op.Write(ram, 0, 99)
	if got := op.Read(ram, 0); got != 99 {
		t.Errorf("Read(ram, 0) = %d, want %d", got, 99)
	}

	CheckRam(t, ram, []int{10, 11, 99, 13, 14})
}

func TestAddInstruction(t *testing.T) {
	ram := NewRam(10, 11, 12, 13, 14)
	io := NewIO()

	var inst Instruction = &Add{
		a: &PositionOperand{1},
		b: &PositionOperand{2},
		c: &PositionOperand{4},
	}

	if npc := inst.Execute(ram, io, 1); npc != 5 {
		t.Errorf("npc = %v, want %v", npc, 5)
	}

	CheckEmptyOutput(t, io)
	CheckRam(t, ram, []int{10, 11, 12, 13, 11 + 12})
}

func TestMultiplyInstruction(t *testing.T) {
	ram := NewRam(10, 11, 12, 13, 14)
	io := NewIO()

	var inst Instruction = &Multiply{
		a: &PositionOperand{1},
		b: &PositionOperand{2},
		c: &PositionOperand{4},
	}

	if npc := inst.Execute(ram, io, 1); npc != 5 {
		t.Errorf("npc = %v, want %v", npc, 5)
	}

	CheckEmptyOutput(t, io)
	CheckRam(t, ram, []int{10, 11, 12, 13, 11 * 12})
}

func TestInputInstruction(t *testing.T) {
	ram := NewRam(10, 11, 12)
	io := NewIO(5)

	var inst Instruction = &Input{&PositionOperand{1}}

	if npc := inst.Execute(ram, io, 1); npc != 3 {
		t.Errorf("npc = %v, want %v", npc, 3)
	}

	CheckEmptyOutput(t, io)
	CheckRam(t, ram, []int{10, 5, 12})
}

func TestOutputInstruction(t *testing.T) {
	ram := NewRam(10, 11, 12)
	io := NewIO()

	var inst Instruction = &Output{&PositionOperand{1}}

	if npc := inst.Execute(ram, io, 1); npc != 3 {
		t.Errorf("npc = %v, want %v", npc, 3)
	}

	if got := io.Written(); !reflect.DeepEqual(got, []int{11}) {
		t.Errorf("Written() = %v, want [11]")
	}

	CheckRam(t, ram, []int{10, 11, 12})
}

func TestHaltInstruction(t *testing.T) {
	vals := []int{10, 11, 12, 13, 14}
	ram := NewRam(vals...)
	io := NewIO()

	var inst Instruction = &Halt{}

	if npc := inst.Execute(ram, io, 0); npc != -1 {
		t.Errorf("npc = %v, want %v", npc, -1)
	}

	CheckEmptyOutput(t, io)
	CheckRam(t, ram, vals)
}
