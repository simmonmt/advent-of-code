package main

import "fmt"

type OpCode int

const (
	OP_ADV OpCode = 0
	OP_BXL OpCode = 1
	OP_BST OpCode = 2
	OP_JNZ OpCode = 3
	OP_BXC OpCode = 4
	OP_OUT OpCode = 5
	OP_BDV OpCode = 6
	OP_CDV OpCode = 7
)

func comboToString(opnd byte) string {
	switch opnd {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "A"
	case 5:
		return "B"
	case 6:
		return "C"
	case 7:
		return "ERR"
	default:
		panic("bad combo")
	}
}

func comboToValue(opnd byte, regs map[string]int) int {
	switch opnd {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return regs["A"]
	case 5:
		return regs["B"]
	case 6:
		return regs["C"]
	case 7:
		panic("combo 7")
	default:
		panic("bad combo")
	}
}

type Inst interface {
	Execute(regs map[string]int, pc int) (npc int, out []byte)
	String() string
}

func ParseInst(x, y byte) (Inst, error) {
	switch OpCode(x) {
	case OP_ADV:
		return &instDiv{name: "ADV", dest: "A", opnd: y}, nil
	case OP_BDV:
		return &instDiv{name: "BDV", dest: "B", opnd: y}, nil
	case OP_CDV:
		return &instDiv{name: "CDV", dest: "C", opnd: y}, nil
	case OP_BXL:
		return &instBXL{opnd: y}, nil
	case OP_BST:
		return &instBST{opnd: y}, nil
	case OP_BXC:
		return &instBXC{opnd: y}, nil
	case OP_JNZ:
		return &instJNZ{opnd: y}, nil
	case OP_OUT:
		return &instOUT{opnd: y}, nil
	default:
		return nil, fmt.Errorf("bad instruction %d %d", x, y)
	}
}

type instDiv struct {
	name string
	dest string
	opnd byte
}

func (inst *instDiv) Execute(regs map[string]int, pc int) (npc int, out []byte) {
	num := regs["A"]

	den := 1
	for i := 0; i < comboToValue(inst.opnd, regs); i++ {
		den *= 2
	}

	regs[inst.dest] = num / den
	return pc + 2, nil
}

func (inst *instDiv) String() string {
	return fmt.Sprintf("%s %s", inst.name, comboToString(inst.opnd))
}

type instBXL struct {
	opnd byte
}

func (inst *instBXL) Execute(regs map[string]int, pc int) (npc int, out []byte) {
	regs["B"] = regs["B"] ^ int(inst.opnd)
	return pc + 2, nil
}

func (inst *instBXL) String() string {
	return fmt.Sprintf("BXL %d", inst.opnd)
}

type instBST struct {
	opnd byte
}

func (inst *instBST) Execute(regs map[string]int, pc int) (npc int, out []byte) {
	regs["B"] = comboToValue(inst.opnd, regs) % 8
	return pc + 2, nil
}

func (inst *instBST) String() string {
	return fmt.Sprintf("BST %s", comboToString(inst.opnd))
}

type instBXC struct {
	opnd byte
}

func (inst *instBXC) Execute(regs map[string]int, pc int) (npc int, out []byte) {
	regs["B"] = regs["B"] ^ regs["C"]
	return pc + 2, nil
}

func (inst *instBXC) String() string {
	return fmt.Sprintf("BXC %s", comboToString(inst.opnd))
}

type instJNZ struct {
	opnd byte
}

func (inst *instJNZ) Execute(regs map[string]int, pc int) (npc int, out []byte) {
	a := regs["A"]
	if a == 0 {
		return pc + 2, nil
	}
	return int(inst.opnd), nil
}

func (inst *instJNZ) String() string {
	return fmt.Sprintf("JNZ %d", inst.opnd)
}

type instOUT struct {
	opnd byte
}

func (inst *instOUT) Execute(regs map[string]int, pc int) (npc int, out []byte) {
	return pc + 2, []byte{byte(comboToValue(inst.opnd, regs) % 8)}
}

func (inst *instOUT) String() string {
	return fmt.Sprintf("OUT %s", comboToString(inst.opnd))
}
