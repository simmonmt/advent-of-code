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

package instr

import (
	"fmt"
	"reg"
	"strconv"
	"strings"
)

type Instr interface {
	Exec(file *reg.File) int
	String() string
}

type regBase struct {
	r1, r2 *reg.Reg
	imm    uint32
	off    int
}

type cpy struct {
	regBase
}

func newCpyReg(r1, r2 reg.Reg) Instr {
	return &cpy{regBase{r1: &r1, r2: &r2}}
}

func newCpyImm(imm uint32, r2 reg.Reg) Instr {
	return &cpy{regBase{imm: imm, r2: &r2}}
}

func (i *cpy) Exec(file *reg.File) int {
	var imm uint32 = i.imm
	if i.r1 != nil {
		imm = file.Get(*i.r1)
	}

	file.Set(*i.r2, imm)
	return 1
}

func (i *cpy) String() string {
	if i.r1 == nil {
		return fmt.Sprintf("cpy %d, %s", i.imm, *i.r2)
	} else {
		return fmt.Sprintf("cpy %v, %s", i.r1, *i.r2)
	}
}

type inc struct {
	regBase
}

func newInc(r reg.Reg) Instr {
	return &inc{regBase{r1: &r}}
}

func (i *inc) Exec(file *reg.File) int {
	file.Set(*i.r1, file.Get(*i.r1)+1)
	return 1
}

func (i *inc) String() string {
	return fmt.Sprintf("inc %s", *i.r1)
}

type dec struct {
	regBase
}

func newDec(r reg.Reg) Instr {
	return &dec{regBase{r1: &r}}
}

func (i *dec) Exec(file *reg.File) int {
	file.Set(*i.r1, file.Get(*i.r1)-1)
	return 1
}

func (i *dec) String() string {
	return fmt.Sprintf("dec %s", *i.r1)
}

type jnz struct {
	regBase
}

func newJnzReg(r1 reg.Reg, off int) Instr {
	return &jnz{regBase{r1: &r1, off: off}}
}

func newJnzImm(imm uint32, off int) Instr {
	return &jnz{regBase{imm: imm, off: off}}
}

func (i *jnz) Exec(file *reg.File) int {
	var val uint32 = i.imm
	if i.r1 != nil {
		val = file.Get(*i.r1)
	}

	if val != 0 {
		return i.off
	}
	return 1
}

func (i *jnz) String() string {
	if i.r1 == nil {
		return fmt.Sprintf("jnz %d, %d", i.imm, i.off)
	} else {
		return fmt.Sprintf("jnz %s, %d", *i.r1, i.off)
	}
}

func parseImm(str string) (uint32, error) {
	str = strings.TrimLeft(str, "+")
	val, err := strconv.ParseUint(str, 10, 32)
	return uint32(val), err
}

func parseOffset(str string) (int, error) {
	str = strings.TrimLeft(str, "+")
	val, err := strconv.ParseInt(str, 10, 32)
	return int(val), err
}

func Parse(op, a, b string) (Instr, error) {
	reg1, reg1Err := reg.FromString(a)
	reg2, reg2Err := reg.FromString(b)

	switch op {
	case "cpy":
		if reg2Err != nil {
			return nil, fmt.Errorf("%v reg2 parse fail: %v", op, reg2Err)
		}

		if reg1Err != nil {
			imm, immErr := parseImm(a)
			if immErr != nil {
				return nil, fmt.Errorf("%v reg1/imm parse fail: %v", op, immErr)
			}
			return newCpyImm(imm, reg2), nil
		} else {
			return newCpyReg(reg1, reg2), nil
		}
		break

	case "inc":
		fallthrough
	case "dec":
		if reg1Err != nil {
			return nil, fmt.Errorf("%v reg1 parse fail: %v", op, reg1Err)
		}
		if b != "" {
			return nil, fmt.Errorf("%v takes one arg: %v r1", op, op)
		}

		switch op {
		case "inc":
			return newInc(reg1), nil
		case "dec":
			return newDec(reg1), nil
		default:
			panic("unreachable")
		}

	case "jnz":
		off, offErr := parseOffset(b)
		if offErr != nil {
			return nil, fmt.Errorf("jmp offset parse fail: %v", offErr)
		}

		if reg1Err != nil {
			imm, immErr := parseImm(a)
			if immErr != nil {
				return nil, fmt.Errorf("%v reg1/imm parse fail: %v", op, immErr)
			}
			return newJnzImm(imm, off), nil
		} else {
			return newJnzReg(reg1, off), nil
		}

	default:
		return nil, fmt.Errorf("unknown op %v", op)
	}

	panic("unreachable")
}
