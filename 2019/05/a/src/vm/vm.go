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

func Run(ram Ram, io IO, pc int) error {
	for {
		var inst Instruction

		val := ram.Read(pc)

		op := val % 100

		ctorA := makeCtor((val / 100) % 10)
		ctorB := makeCtor((val / 1000) % 10)
		ctorC := makeCtor((val / 10000) % 10)

		switch op {
		case 1:
			inst = &Add{
				a: ctorA(ram.Read(pc + 1)),
				b: ctorB(ram.Read(pc + 2)),
				c: ctorC(ram.Read(pc + 3)),
			}
			break
		case 2:
			inst = &Multiply{
				a: ctorA(ram.Read(pc + 1)),
				b: ctorB(ram.Read(pc + 2)),
				c: ctorC(ram.Read(pc + 3)),
			}
			break
		case 3:
			inst = &Input{a: ctorA(ram.Read(pc + 1))}
			break
		case 4:
			inst = &Output{a: ctorA(ram.Read(pc + 1))}
			break
		case 5:
			inst = &JumpIfTrue{
				a: ctorA(ram.Read(pc + 1)),
				b: ctorB(ram.Read(pc + 2)),
			}
			break
		case 6:
			inst = &JumpIfFalse{
				a: ctorA(ram.Read(pc + 1)),
				b: ctorB(ram.Read(pc + 2)),
			}
			break
		case 7:
			inst = &LessThan{
				a: ctorA(ram.Read(pc + 1)),
				b: ctorB(ram.Read(pc + 2)),
				c: ctorC(ram.Read(pc + 3)),
			}
			break
		case 8:
			inst = &Equals{
				a: ctorA(ram.Read(pc + 1)),
				b: ctorB(ram.Read(pc + 2)),
				c: ctorC(ram.Read(pc + 3)),
			}
			break
		case 99:
			inst = &Halt{}
			break

		default:
			return fmt.Errorf("bad opcode %d at %d", op, pc)
		}

		logger.LogF("%d: %s %v", pc, inst.String(), readBytes(ram, pc, inst.Size()))

		npc := inst.Execute(ram, io, pc)
		if npc == -1 {
			return nil
		}

		pc = npc
	}
}
