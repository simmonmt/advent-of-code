package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
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
	case "mod":
		regFile[op.XReg] = regOrVal(regFile, op.XReg, op.XVal) %
			regOrVal(regFile, op.YReg, op.YVal)
		break
	case "jnz":
		//fmt.Printf("jnz %+v x %v y %v\n", op, regOrVal(regFile, op.XReg, op.XVal), regOrVal(regFile, op.YReg, op.YVal))
		if regOrVal(regFile, op.XReg, op.XVal) != 0 {
			rNPC = regOrVal(regFile, op.YReg, op.YVal)
		}
		break
	case "nop":
		break
	default:
		panic(fmt.Sprintf("unknown op name %v", op.Name))
	}

	return rNPC, op.Name == "mul"
}

func isPrime(value int) bool {
	for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func runProgram(ops []*Op) int {
	regFile := map[string]int{"a": 1}

	rNPC := 0
	numMul := 0
	// lastH := regFile["h"]
	// lastF := regFile["f"]
	i := 0
	for pc := 0; pc >= 0 && pc < len(ops); pc += rNPC {
		i++

		op := ops[pc]

		var mulInvoked bool
		rNPC, mulInvoked = runOp(op, regFile)

		if mulInvoked {
			numMul++
		}

		if pc == 28 {
			fmt.Println(regFile["b"], regFile["h"], isPrime(regFile["b"]))
		}

		// if regFile["h"] != lastH || regFile["f"] != lastF || (i%10000000) == 0 {
		// 	fmt.Printf("i %v f %v h %v\n", i, regFile["f"], regFile["h"])
		// 	lastF = regFile["f"]
		// 	lastH = regFile["h"]
		// }
	}

	return numMul
}

func main() {
	ops := readOps(os.Stdin)

	numMul := runProgram(ops)
	fmt.Printf("num mul = %d\n", numMul)
}
