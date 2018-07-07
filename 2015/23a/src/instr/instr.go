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
	r   reg.Reg
	off int
}

type hlf struct {
	regBase
}

func newHlf(r reg.Reg) Instr {
	return &hlf{regBase{r: r, off: 0}}
}

func (i *hlf) Exec(file *reg.File) int {
	file.Set(i.r, file.Get(i.r)/2)
	return 1
}

func (i *hlf) String() string {
	return fmt.Sprintf("hlf %s", i.r)
}

type tpl struct {
	regBase
}

func newTpl(r reg.Reg) Instr {
	return &tpl{regBase{r: r, off: 0}}
}

func (i *tpl) Exec(file *reg.File) int {
	file.Set(i.r, file.Get(i.r)*3)
	return 1
}

func (i *tpl) String() string {
	return fmt.Sprintf("tpl %s", i.r)
}

type inc struct {
	regBase
}

func newInc(r reg.Reg) Instr {
	return &inc{regBase{r: r, off: 0}}
}

func (i *inc) Exec(file *reg.File) int {
	file.Set(i.r, file.Get(i.r)+1)
	return 1
}

func (i *inc) String() string {
	return fmt.Sprintf("inc %s", i.r)
}

type jmp struct {
	regBase
}

func newJmp(off int) Instr {
	return &jmp{regBase{r: reg.A, off: off}}
}

func (i *jmp) Exec(file *reg.File) int {
	return i.off
}

func (i *jmp) String() string {
	return fmt.Sprintf("jmp %d", i.off)
}

type jie struct {
	regBase
}

func newJie(r reg.Reg, off int) Instr {
	return &jie{regBase{r: r, off: off}}
}

func (i *jie) Exec(file *reg.File) int {
	if file.Get(i.r)%2 == 0 {
		return i.off
	} else {
		return 1
	}
}

func (i *jie) String() string {
	return fmt.Sprintf("jie %s, %d", i.r, i.off)
}

type jio struct {
	regBase
}

func newJio(r reg.Reg, off int) Instr {
	return &jio{regBase{r: r, off: off}}
}

func (i *jio) Exec(file *reg.File) int {
	if file.Get(i.r) == 1 {
		return i.off
	} else {
		return 1
	}
}

func (i *jio) String() string {
	return fmt.Sprintf("jio %s, %d", i.r, i.off)
}

func parseOffset(str string) (int, error) {
	str = strings.TrimLeft(str, "+")
	val, err := strconv.ParseInt(str, 10, 32)
	return int(val), err
}

func Parse(op, a, b string) (Instr, error) {
	reg, regErr := reg.FromString(a)

	switch op {
	case "hlf":
		fallthrough
	case "tpl":
		fallthrough
	case "inc":
		if regErr != nil {
			return nil, fmt.Errorf("%v reg parse fail: %v", op, regErr)
		}
		if b != "" {
			return nil, fmt.Errorf("%v takes one arg: %v r", op, op)
		}

		switch op {
		case "hlf":
			return newHlf(reg), nil
		case "tpl":
			return newTpl(reg), nil
		case "inc":
			return newInc(reg), nil
		default:
			panic("unreachable")
		}

	case "jmp":
		off, offErr := parseOffset(a)
		if offErr != nil {
			return nil, fmt.Errorf("jmp offset parse fail: %v", offErr)
		}
		if b != "" {
			return nil, fmt.Errorf("%v takes one arg: %v r", op, op)
		}

		return newJmp(off), nil

	case "jie":
		fallthrough
	case "jio":
		if regErr != nil {
			return nil, fmt.Errorf("%v reg parse fail: %v", op, regErr)
		}
		off, offErr := parseOffset(b)
		if offErr != nil {
			return nil, fmt.Errorf("%v offset parse fail: %v", op, offErr)
		}

		switch op {
		case "jie":
			return newJie(reg, off), nil
		case "jio":
			return newJio(reg, off), nil
		default:
			panic("unreachable")
		}

	default:
		return nil, fmt.Errorf("unknown op %v", op)
	}

	panic("unreachable")
}
