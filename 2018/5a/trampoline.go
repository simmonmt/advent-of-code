package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func renderOffsets(offsets []int, addr int) string {
	out := "["
	for i, offset := range offsets {
		if i != 0 {
			out += fmt.Sprintf(" ")
		}
		if i == addr {
			out += fmt.Sprintf("(%d %d)", offset-1, offset)
		} else {
			out += fmt.Sprintf("%d", offset)
		}
	}
	out += "]"
	return out
}

func main() {
	offsets := []int{}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		offset, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("failed to parse input line %v", line)
		}

		offsets = append(offsets, offset)
	}

	addr := 0
	var numSteps int
	for numSteps = 0; true; numSteps++ {
		if addr < 0 || addr >= len(offsets) {
			break
		}

		naddr := addr + offsets[addr]
		offsets[addr]++

		//fmt.Printf("addr %v naddr %v %v\n", addr, naddr, renderOffsets(offsets, addr))
		addr = naddr
	}

	fmt.Println(numSteps)
}
