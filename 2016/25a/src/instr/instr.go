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
	IsTgl() bool
}

type regBase struct {
	r1, r2 *reg.Reg
	imm    int
	off    int
}

func (r *regBase) IsTgl() bool {
	return false
}

type add struct {
	regBase
}

func newAdd(r1, r2 reg.Reg) Instr {
	return &add{regBase{r1: &r1, r2: &r2}}
}

func (i *add) Exec(file *reg.File) int {
	file.Set(*i.r2, file.Get(*i.r1)+file.Get(*i.r2))
	return 1
}

func (i *add) String() string {
	return fmt.Sprintf("add %s %s", *i.r1, *i.r2)
}

type mul struct {
	regBase
}

func newMul(r1, r2 reg.Reg) Instr {
	return &mul{regBase{r1: &r1, r2: &r2}}
}

func (i *mul) Exec(file *reg.File) int {
	file.Set(*i.r2, file.Get(*i.r1)*file.Get(*i.r2))
	return 1
}

func (i *mul) String() string {
	return fmt.Sprintf("mul %s %s", *i.r1, *i.r2)
}

type cpy struct {
	regBase
}

func newCpyReg(r1, r2 reg.Reg) Instr {
	return &cpy{regBase{r1: &r1, r2: &r2}}
}

func newCpyImm(imm int, r2 reg.Reg) Instr {
	return &cpy{regBase{imm: imm, r2: &r2}}
}

func (i *cpy) Exec(file *reg.File) int {
	imm := i.imm
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

type out struct {
	regBase
}

func newOutReg(r1 reg.Reg) Instr {
	return &out{regBase{r1: &r1}}
}

func newOutImm(imm int) Instr {
	return &out{regBase{imm: imm}}
}

func (i *out) Exec(file *reg.File) int {
	val := i.imm
	if i.r1 != nil {
		val = file.Get(*i.r1)
	}
	fmt.Println(val)
	return 1
}

func (i *out) String() string {
	if i.r1 != nil {
		return fmt.Sprintf("out %s", *i.r1)
	} else {
		return fmt.Sprintf("out %d", i.imm)
	}
}

type jnz struct {
	regBase
	offReg *reg.Reg
}

func newJnzReg(r1 reg.Reg, off int, offReg *reg.Reg) Instr {
	return &jnz{
		regBase: regBase{r1: &r1, off: off},
		offReg:  offReg,
	}
}

func newJnzImm(imm int, off int, offReg *reg.Reg) Instr {
	return &jnz{
		regBase: regBase{imm: imm, off: off},
		offReg:  offReg,
	}
}

func (i *jnz) Exec(file *reg.File) int {
	val := i.imm
	if i.r1 != nil {
		val = file.Get(*i.r1)
	}

	if val != 0 {
		if i.offReg != nil {
			return file.Get(*i.offReg)
		} else {
			return i.off
		}
	}
	return 1
}

func (i *jnz) String() string {
	out := "jnz "

	if i.r1 != nil {
		out += fmt.Sprintf("%s, ", *i.r1)
	} else {
		out += fmt.Sprintf("%d, ", i.imm)
	}

	if i.offReg != nil {
		out += fmt.Sprint(*i.offReg)
	} else {
		out += fmt.Sprint(i.off)
	}

	return out
}

type nop struct {
	regBase
}

func newNop() Instr {
	return &nop{}
}

func (i *nop) Exec(file *reg.File) int {
	return 1
}

func (i *nop) String() string {
	return "nop"
}

type tgl struct {
	r1 *reg.Reg
}

func newTgl(r1 reg.Reg) Instr {
	return &tgl{r1: &r1}
}

func (i *tgl) Exec(file *reg.File) int {
	return int(int32(file.Get(*i.r1)))
}

func (i *tgl) String() string {
	return fmt.Sprintf("tgl %s", *i.r1)
}

func (i *tgl) IsTgl() bool {
	return true
}

func parseImm(str string) (int, error) {
	str = strings.TrimLeft(str, "+")
	val, err := strconv.ParseInt(str, 10, 32)
	return int(val), err
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
	case "add":
		if reg1Err != nil {
			return nil, fmt.Errorf("%v reg1 parse fail: %v", op, reg1Err)
		}
		if reg2Err != nil {
			return nil, fmt.Errorf("%v reg2 parse fail: %v", op, reg2Err)
		}
		return newAdd(reg1, reg2), nil

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

	case "dec":
		fallthrough
	case "inc":
		if reg1Err != nil {
			return nil, fmt.Errorf("%v reg1 parse fail: %v", op, reg1Err)
		}
		if b != "" {
			return nil, fmt.Errorf("%v takes one arg: %v r1", op, op)
		}

		switch op {
		case "dec":
			return newDec(reg1), nil
		case "inc":
			return newInc(reg1), nil
		default:
			panic("unreachable")
		}

	case "jnz":
		off, offErr := parseOffset(b)
		var offReg *reg.Reg
		if reg2Err == nil {
			offReg = &reg2
		}

		if offErr != nil && reg2Err != nil {
			return nil, fmt.Errorf("jmp offset parse fail: %v", offErr)
		}
		if offErr == nil && reg2Err == nil {
			return nil, fmt.Errorf("jmp offset parse fail: both nil")
		}

		if reg1Err != nil {
			imm, immErr := parseImm(a)
			if immErr != nil {
				return nil, fmt.Errorf("%v reg1/imm parse fail: %v", op, immErr)
			}
			return newJnzImm(imm, off, offReg), nil
		} else {
			return newJnzReg(reg1, off, offReg), nil
		}

	case "mul":
		if reg1Err != nil {
			return nil, fmt.Errorf("%v reg1 parse fail: %v", op, reg1Err)
		}
		if reg2Err != nil {
			return nil, fmt.Errorf("%v reg2 parse fail: %v", op, reg2Err)
		}
		return newMul(reg1, reg2), nil

	case "nop":
		if a != "" || b != "" {
			return nil, fmt.Errorf("%v takes no args: %v", op, op)
		}

		return newNop(), nil

	case "out":
		if b != "" {
			return nil, fmt.Errorf("%v takes one arg: %v r1", op, op)
		}

		if reg1Err != nil {
			imm, immErr := parseImm(a)
			if immErr != nil {
				return nil, fmt.Errorf("%v reg1/imm parse fail: %v", op, immErr)
			}
			return newOutImm(imm), nil
		} else {
			return newOutReg(reg1), nil
		}

	case "tgl":
		if reg1Err != nil {
			return nil, fmt.Errorf("%v reg1 parse fail: %v", op, reg1Err)
		}
		if b != "" {
			return nil, fmt.Errorf("%v takes one arg: %v r1", op, op)
		}

		return newTgl(reg1), nil

	default:
		return nil, fmt.Errorf("unknown op %v", op)
	}

	panic("unreachable")
}
