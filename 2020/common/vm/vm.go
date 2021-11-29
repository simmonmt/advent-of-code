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

	"github.com/simmonmt/aoc/2020/common/logger"
)

func immediate(imm int64) Operand {
	return &ImmediateOperand{imm}
}

func position(addr int64) Operand {
	return &PositionOperand{addr}
}

func relative(addr int64) Operand {
	return &RelativeOperand{addr}
}

func makeCtor(mode int64) func(int64) Operand {
	switch mode {
	case 0:
		return position
	case 1:
		return immediate
	case 2:
		return relative
	default:
		panic(fmt.Sprintf("unknown mode %d", mode))
	}
}

func readBytes(ram Ram, pc, sz int64) []int64 {
	out := []int64{}
	for i := int64(0); i < sz; i++ {
		out = append(out, ram.Read(pc+i))
	}
	return out
}

type Resources struct {
	ram     Ram
	io      IO
	relBase int64
}

func decode(r *Resources, pc int64) (Instruction, error) {
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
	case 9:
		inst = &SetRelBase{
			a: ctorA(r.ram.Read(pc + 1)),
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

func run(r *Resources, pc int64) error {
	for {
		inst, err := decode(r, pc)
		if err != nil {
			return err
		}

		npc := inst.Execute(r, pc)
		if npc == -1 {
			return nil
		}

		pc = npc
	}
}

func Run(ram Ram, io IO) error {
	r := &Resources{
		ram:     ram,
		io:      io,
		relBase: 0,
	}

	return run(r, 0)
}

type Async struct {
	In  chan *ChanIOMessage
	Out chan *ChanIOMessage
}

func RunAsync(id string, ram Ram) *Async {
	async := &Async{
		In:  make(chan *ChanIOMessage, 2),
		Out: make(chan *ChanIOMessage, 2),
	}

	io := NewChanIO(async.In, async.Out)

	go func() {
		if err := Run(ram, io); err != nil {
			async.Out <- &ChanIOMessage{
				Err: fmt.Errorf("vm %s failed: %v", id, err),
			}
		}

		close(async.Out)
		fmt.Sprintf("vm %s terminating\n", id)
	}()

	return async
}
