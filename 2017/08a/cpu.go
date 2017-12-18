package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	pattern = regexp.MustCompile(`^([^ ]+) (inc|dec) (-?[0-9]+) if ([^ ]+) ([^ ]+) ([^ ]+)$`)
)

func cmpIsTrue(regFile *map[string]int, reg, op string, val int) bool {
	regVal := (*regFile)[reg]

	switch op {
	case "==":
		return regVal == val
	case "!=":
		return regVal != val
	case "<":
		return regVal < val
	case ">":
		return regVal > val
	case ">=":
		return regVal >= val
	case "<=":
		return regVal <= val
	default:
		log.Fatalf("unknown op %v\n", op)
	}

	return false
}

func maxReg(regFile *map[string]int) int {
	first := true
	maxVal := 0

	for _, val := range *regFile {
		if first || val > maxVal {
			first = false
			maxVal = val
		}
	}

	return maxVal
}

func main() {
	regFile := map[string]int{}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		//fmt.Println(line)

		matches := pattern.FindStringSubmatch(line)
		if len(matches) == 0 {
			log.Fatalf("failed to parse %v\n", line)
		}
		changeReg, op, changeValStr := matches[1], matches[2], matches[3]
		cmpReg, cmpOp, cmpValStr := matches[4], matches[5], matches[6]

		changeVal, err := strconv.Atoi(changeValStr)
		if err != nil {
			log.Fatalf("failed to parse change val %v in %v\n", changeValStr, line)
		}
		cmpVal, err := strconv.Atoi(cmpValStr)
		if err != nil {
			log.Fatalf("failed to parse cmp val %v in %v\n", cmpValStr, line)
		}

		if !cmpIsTrue(&regFile, cmpReg, cmpOp, cmpVal) {
			//fmt.Printf("cmp false cmpReg %v\n", regFile[cmpReg])
			continue
		}

		dir := 1
		if op == "dec" {
			dir = -1
		}

		//fmt.Printf("reg %v was %v ", changeReg, regFile[changeReg])
		regFile[changeReg] += dir * changeVal
		//fmt.Printf("now %v\n", regFile[changeReg])
	}

	fmt.Printf("max in file: %v\n", maxReg(&regFile))
}
