package vm

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/logger"
)

func indirect(addr int) Operand {
	return &IndirectOperand{addr}
}

func Run(pc int, ram Ram) error {
	for {
		var inst Instruction

		op := ram.Read(pc)
		switch op {
		case 1:
			inst = &Add{
				a: indirect(ram.Read(pc + 1)),
				b: indirect(ram.Read(pc + 2)),
				c: indirect(ram.Read(pc + 3)),
			}
			break
		case 2:
			inst = &Multiply{
				a: indirect(ram.Read(pc + 1)),
				b: indirect(ram.Read(pc + 2)),
				c: indirect(ram.Read(pc + 3)),
			}
			break
		case 99:
			inst = &Halt{}
			break

		default:
			return fmt.Errorf("bad opcode %d at %d", op, pc)
		}

		logger.LogF("%d: %s", pc, inst.String())

		npc := inst.Execute(ram, pc)
		if npc == -1 {
			return nil
		}

		pc = npc
	}
}
