package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	lastSndRegName = "_lastSnd"
	recvRegName    = "_recv"
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

func runOp(op *Op, regFile map[string]int, sndCh, rcvCh chan int) (int, bool) {
	rNPC := 1
	sent := false

	switch op.Name {
	case "snd":
		sndCh <- regOrVal(regFile, op.XReg, op.XVal)
		sent = true
		break
	case "set":
		regFile[op.XReg] = regOrVal(regFile, op.YReg, op.YVal)
		break
	case "add":
		regFile[op.XReg] += regOrVal(regFile, op.YReg, op.YVal)
		break
	case "mul":
		regFile[op.XReg] *= regOrVal(regFile, op.YReg, op.YVal)
		break
	case "mod":
		regFile[op.XReg] = regOrVal(regFile, op.XReg, op.XVal) %
			regOrVal(regFile, op.YReg, op.YVal)
		break
	case "rcv":
		regFile[op.XReg] = <-rcvCh
		break
	case "jgz":
		if regOrVal(regFile, op.XReg, op.XVal) > 0 {
			rNPC = regOrVal(regFile, op.YReg, op.YVal)
		}
		break
	default:
		panic(fmt.Sprintf("unknown op name %v", op.Name))
	}

	return rNPC, sent
}

func runProgram(ops []*Op, programId int, sndCh, rcvCh chan int) {
	regFile := map[string]int{"p": programId}
	numSent := 0

	rNPC := 0
	for pc := 0; pc >= 0 && pc < len(ops); pc += rNPC {
		// fmt.Printf("%d %d %d\n", programId, pc, len(ops))

		op := ops[pc]
		var sent bool
		if rNPC, sent = runOp(op, regFile, sndCh, rcvCh); sent {
			numSent++
			if programId == 1 {
				fmt.Printf("program %d numSent %d\n",
					programId, numSent)
			}
		}
		// fmt.Printf("%d pc:%d rnpc:%d npc:%d %+v %v\n",
		// 	programId, pc, rNPC, pc+rNPC, *op, regFile)
	}
	fmt.Printf("program %v exiting\n", programId)
}

func main() {
	ops := readOps(os.Stdin)

	ch1 := make(chan int, 100000)
	ch2 := make(chan int, 100000)

	go runProgram(ops, 0, ch1, ch2)
	go runProgram(ops, 1, ch2, ch1)

	time.Sleep(10 * time.Second)
}
