package vm

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/logger"
)

func immediate(imm int) Operand {
	return &ImmediateOperand{imm}
}

func position(addr int) Operand {
	return &PositionOperand{addr}
}

func makeCtor(mode int) func(int) Operand {
	switch mode {
	case 0:
		return position
	case 1:
		return immediate
	default:
		panic(fmt.Sprintf("unknown mode %d", mode))
	}
}

func readBytes(ram Ram, pc, sz int) []int {
	out := []int{}
	for i := 0; i < sz; i++ {
		out = append(out, ram.Read(pc+i))
	}
	return out
}

type Resources struct {
	ram Ram
	io  IO
}

func decode(r *Resources, pc int) (Instruction, error) {
	var inst Instruction

	val := r.ram.Read(pc)

	op := val % 100

	ctorA := makeCtor((val / 100) % 10)
	ctorB := makeCtor((val / 1000) % 10)
	ctorC := makeCtor((val / 10000) % 10)

	switch op {
	case 1:
		inst = &Add{
			a: ctorA(r.ram.Read(pc + 1)),
			b: ctorB(r.ram.Read(pc + 2)),
			c: ctorC(r.ram.Read(pc + 3)),
		}
		break
	case 2:
		inst = &Multiply{
			a: ctorA(r.ram.Read(pc + 1)),
			b: ctorB(r.ram.Read(pc + 2)),
			c: ctorC(r.ram.Read(pc + 3)),
		}
		break
	case 3:
		inst = &Input{a: ctorA(r.ram.Read(pc + 1))}
		break
	case 4:
		inst = &Output{a: ctorA(r.ram.Read(pc + 1))}
		break
	case 5:
		inst = &JumpIfTrue{
			a: ctorA(r.ram.Read(pc + 1)),
			b: ctorB(r.ram.Read(pc + 2)),
		}
		break
	case 6:
		inst = &JumpIfFalse{
			a: ctorA(r.ram.Read(pc + 1)),
			b: ctorB(r.ram.Read(pc + 2)),
		}
		break
	case 7:
		inst = &LessThan{
			a: ctorA(r.ram.Read(pc + 1)),
			b: ctorB(r.ram.Read(pc + 2)),
			c: ctorC(r.ram.Read(pc + 3)),
		}
		break
	case 8:
		inst = &Equals{
			a: ctorA(r.ram.Read(pc + 1)),
			b: ctorB(r.ram.Read(pc + 2)),
			c: ctorC(r.ram.Read(pc + 3)),
		}
		break
	case 99:
		inst = &Halt{}
		break

	default:
		return nil, fmt.Errorf("bad opcode %d at %d", op, pc)
	}

	logger.LogF("%d: %s %v", pc, inst.String(), readBytes(r.ram, pc, inst.Size()))
	return inst, nil
}

func run(r *Resources, pc int) error {
	for {
		inst, err := decode(r, pc)
		if err != nil {
			return err
		}

		npc := inst.Execute(r.ram, r.io, pc)
		if npc == -1 {
			return nil
		}

		pc = npc
	}
}

func Run(ram Ram, io IO, pc int) error {
	r := &Resources{
		ram: ram,
		io:  io,
	}

	return run(r, pc)
}
