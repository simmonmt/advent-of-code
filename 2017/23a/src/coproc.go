package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Op struct {
	Name       string
	XReg, YReg string
	XVal, YVal int
}

func parseRegOrVal(str string) (reg string, val int) {
	var err error
	if val, err = strconv.Atoi(str); err != nil {
		reg = str
	}
	return
}

func regOrVal(regFile map[string]int, reg string, val int) int {
	if reg == "" {
		return val
	} else {
		return regFile[reg]
	}
}

func readOps(in io.Reader) []*Op {
	ops := []*Op{}

	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		pieces := strings.Split(line, " ")

		op := &Op{
			Name: pieces[0],
		}

		op.XReg, op.XVal = parseRegOrVal(pieces[1])
		if len(pieces) == 3 {
			op.YReg, op.YVal = parseRegOrVal(pieces[2])
		}

		ops = append(ops, op)
	}

	return ops
}

func runOp(op *Op, regFile map[string]int) (int, bool) {
	rNPC := 1

	switch op.Name {
	case "set":
		regFile[op.XReg] = regOrVal(regFile, op.YReg, op.YVal)
		break
	case "sub":
		regFile[op.XReg] -= regOrVal(regFile, op.YReg, op.YVal)
		break
	case "mul":
		regFile[op.XReg] *= regOrVal(regFile, op.YReg, op.YVal)
		break
	case "jnz":
		if regOrVal(regFile, op.XReg, op.XVal) != 0 {
			rNPC = regOrVal(regFile, op.YReg, op.YVal)
		}
		break
	default:
		panic(fmt.Sprintf("unknown op name %v", op.Name))
	}

	return rNPC, op.Name == "mul"
}

func runProgram(ops []*Op) int {
	regFile := map[string]int{}

	rNPC := 0
	numMul := 0
	for pc := 0; pc >= 0 && pc < len(ops); pc += rNPC {
		op := ops[pc]

		var mulInvoked bool
		rNPC, mulInvoked = runOp(op, regFile)

		if mulInvoked {
			numMul++
		}
	}

	return numMul
}

func main() {
	ops := readOps(os.Stdin)

	numMul := runProgram(ops)
	fmt.Printf("num mul = %d\n", numMul)
}
