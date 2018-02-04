package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func readLines(r io.Reader) ([]string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

const (
	stateNormal    = iota
	stateEscape    // seen \
	stateHexEscape // seen x in \x27
	stateHexDigit  // seen 2 in \x27
)

func countInMemory(line string) int {
	inMem := 0
	state := stateNormal
	for _, c := range line {
		switch state {
		case stateNormal:
			if c == '\\' {
				state = stateEscape
			} else {
				inMem++
			}
			break

		case stateEscape:
			if c == 'x' {
				state = stateHexEscape
			} else {
				inMem++
				state = stateNormal
			}
			break

		case stateHexEscape:
			state = stateHexDigit
			break

		case stateHexDigit:
			inMem++
			state = stateNormal
		}

		//fmt.Printf("on %c res state %v inMem %v\n", c, state, inMem)
	}

	//arr := ([]byte(line))[1 : len(line)-2]

	inMem -= 2
	//fmt.Printf("line %v inMem %v\n", line, inMem)
	return inMem
}

func encode(line string) string {
	out := ""
	for _, c := range line {
		switch c {
		case '"':
			out += "\\\""
			break
		case '\\':
			out += "\\\\"
			break
		default:
			out += string(c)
		}
	}

	out = "\"" + out + "\""
	fmt.Printf("%v encodes as %v\n", line, out)
	return out
}

func main() {
	lines, err := readLines(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	inMemory, inCode, inEncoded := 0, 0, 0
	for _, line := range lines {
		inCode += len(line)
		inMemory += countInMemory(line)
		inEncoded += len(encode(line))
	}

	fmt.Printf("in code: %v\n", inCode)
	fmt.Printf("in mem : %v\n", inMemory)
	fmt.Printf("in enc : %v\n", inEncoded)
	fmt.Printf("code-mem: %v\n", inCode-inMemory)
	fmt.Printf("enc-code: %v\n", inEncoded-inCode)
}
