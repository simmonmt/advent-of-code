package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func foo(amt, goal int, containers []int) int {
	//fmt.Printf("foo(amt=%d,goal=%d,containers=%d)\n", amt, goal, containers)

	if amt == goal {
		return 1
	} else if amt > goal || len(containers) == 0 {
		return 0
	}

	found := foo(amt+containers[0], goal, containers[1:])
	found += foo(amt, goal, containers[1:])
	return found
}

func readInput(r io.Reader) ([]int, error) {
	vals := []int{}

	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		val, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", line, err)
		}

		vals = append(vals, val)
	}

	return vals, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v goal", os.Args[0])
	}
	goal, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("invalid goal %v: %v", os.Args[1], err)
	}

	containers, err := readInput(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read containers: %v", err)
	}
	fmt.Println(containers)

	fmt.Println(foo(0, goal, containers))
}
