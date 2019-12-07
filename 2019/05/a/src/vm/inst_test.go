package vm

import "testing"

func MakeRam(vals []int) Ram {
	ram := NewRam(vals...)
	return ram
}

func CheckRam(t *testing.T, ram Ram, vals []int) {
	for i, val := range vals {
		if got := ram.Read(i); got != val {
			t.Errorf("verify mismatch at %v: got %v want %v", i, got, val)
		}
	}
}

func TestAdd(t *testing.T) {
	vals := []int{10, 11, 12, 13, 14}
	ram := MakeRam(vals)

	var inst Instruction = &Add{
		a: &IndirectOperand{1},
		b: &IndirectOperand{2},
		c: &IndirectOperand{4},
	}

	if npc := inst.Execute(ram, 0); npc != 4 {
		t.Errorf("npc = %v, want %v", npc, 4)
	}

	CheckRam(t, ram, []int{10, 11, 12, 13, 11 + 12})
}

func TestMultiply(t *testing.T) {
	vals := []int{10, 11, 12, 13, 14}
	ram := MakeRam(vals)

	var inst Instruction = &Multiply{
		a: &IndirectOperand{1},
		b: &IndirectOperand{2},
		c: &IndirectOperand{4},
	}

	if npc := inst.Execute(ram, 0); npc != 4 {
		t.Errorf("npc = %v, want %v", npc, 4)
	}

	CheckRam(t, ram, []int{10, 11, 12, 13, 11 * 12})
}

func TestHalt(t *testing.T) {
	vals := []int{10, 11, 12, 13, 14}
	ram := MakeRam(vals)

	var inst Instruction = &Halt{}

	if npc := inst.Execute(ram, 0); npc != -1 {
		t.Errorf("npc = %v, want %v", npc, -1)
	}

	CheckRam(t, ram, vals)
}
