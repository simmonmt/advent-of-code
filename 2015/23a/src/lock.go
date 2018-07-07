package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"instr"
	"reg"
)

var (
	instrPattern = regexp.MustCompile(`^(...) ([^ ]+)(?:, ([^ ]+))?$`)
)

func parseInput(r io.Reader) ([]instr.Instr, error) {
	instrs := []instr.Instr{}

	reader := bufio.NewReader(os.Stdin)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		matches := instrPattern.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("%d: failed to parse line: %v", lineNum, line)
		}

		op := matches[1]
		a := matches[2]
		b := matches[3]

		i, err := instr.Parse(op, a, b)
		if err != nil {
			return nil, fmt.Errorf("%d: %v", lineNum, err)
		}

		instrs = append(instrs, i)
	}

	return instrs, nil
}

func main() {
	instrs, err := parseInput(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	regFile := reg.NewFile()

	pc := 0
	for {
		if pc >= len(instrs) {
			fmt.Println("done")
			regFile.Print()
			break
		}

		i := instrs[pc]
		pc += i.Exec(regFile)
	}
}
