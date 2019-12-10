package vm

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2019/07/a/src/testutils"
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

type InstructionTestCase struct {
	inst        Instruction
	expectedNPC int
	expectedRam []int
}

func CheckInstruction(t *testing.T, startRam Ram, startPC int, testCases []InstructionTestCase) {
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			ram := startRam.Clone()
			io := NewIO()

			var wantNPC int
			if tc.expectedNPC != 0 {
				wantNPC = tc.expectedNPC
			} else {
				wantNPC = startPC + tc.inst.Size()
			}

			if npc := tc.inst.Execute(ram, io, startPC); npc != wantNPC {
				t.Errorf("Execute, npc=%v, want %v", npc, wantNPC)
			}

			if tc.expectedRam != nil {
				CheckRam(t, ram, tc.expectedRam)
			}

			CheckEmptyOutput(t, io)
		})

	}
}

func TestInstructions(t *testing.T) {
	initialRamValues := []int{10, 11, 12, 13, 14}
	ram := NewRam(initialRamValues...)
	startPC := 1

	testCases := []InstructionTestCase{
		// Simple instructions
		InstructionTestCase{
			inst: &Add{
				a: &PositionOperand{1},
				b: &PositionOperand{2},
				c: &PositionOperand{4},
			},
			expectedRam: []int{10, 11, 12, 13, 11 + 12},
		},
		InstructionTestCase{
			inst: &Multiply{
				a: &PositionOperand{1},
				b: &PositionOperand{2},
				c: &PositionOperand{4},
			},
			expectedRam: []int{10, 11, 12, 13, 11 * 12},
		},
		InstructionTestCase{
			inst:        &Halt{},
			expectedRam: initialRamValues,
			expectedNPC: -1,
		},

		// JumpIfTrue
		InstructionTestCase{
			inst:        &JumpIfTrue{&ImmediateOperand{1}, &ImmediateOperand{99}},
			expectedNPC: 99,
		},
		InstructionTestCase{
			inst:        &JumpIfTrue{&ImmediateOperand{0}, &ImmediateOperand{99}},
			expectedNPC: startPC + (&JumpIfTrue{}).Size(),
		},

		// JumpIfFalse
		InstructionTestCase{
			inst:        &JumpIfFalse{&ImmediateOperand{0}, &ImmediateOperand{99}},
			expectedNPC: 99,
		},
		InstructionTestCase{
			inst:        &JumpIfFalse{&ImmediateOperand{1}, &ImmediateOperand{99}},
			expectedNPC: startPC + (&JumpIfFalse{}).Size(),
		},

		// LessThan
		InstructionTestCase{
			inst: &LessThan{
				a: &ImmediateOperand{1},
				b: &ImmediateOperand{2},
				c: &PositionOperand{1},
			},
			expectedRam: []int{10, 1, 12, 13, 14},
		},
		InstructionTestCase{
			inst: &LessThan{
				a: &ImmediateOperand{2},
				b: &ImmediateOperand{1},
				c: &PositionOperand{1},
			},
			expectedRam: []int{10, 0, 12, 13, 14},
		},

		// Equals
		InstructionTestCase{
			inst: &Equals{
				a: &ImmediateOperand{1},
				b: &ImmediateOperand{1},
				c: &PositionOperand{1},
			},
			expectedRam: []int{10, 1, 12, 13, 14},
		},
		InstructionTestCase{
			inst: &Equals{
				a: &ImmediateOperand{1},
				b: &ImmediateOperand{2},
				c: &PositionOperand{1},
			},
			expectedRam: []int{10, 0, 12, 13, 14},
		},
	}

	CheckInstruction(t, ram, startPC, testCases)
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
