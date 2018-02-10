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

func foo(amt, goal int, fixed, containers []int) (int, [][]int) {
	// fmt.Printf("foo(amt=%d,goal=%d,fixed=%v,containers=%v)\n", amt, goal, fixed, containers)

	if amt == goal {
		return 1, [][]int{fixed}
	} else if amt > goal || len(containers) == 0 {
		return 0, nil
	}

	// fmt.Printf("without %d: ", containers[0])
	found, winners := foo(amt, goal, fixed, containers[1:])

	fixedWith := make([]int, len(fixed))
	copy(fixedWith, fixed)
	fixedWith = append(fixedWith, containers[0])

	// fmt.Printf("with %d: ", containers[0])
	foundWith, winnersWith := foo(amt+containers[0], goal, fixedWith, containers[1:])
	found += foundWith

	for _, winner := range winnersWith {
		winners = append(winners, winner)
	}

	// fmt.Println(winners)
	return found, winners
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

	found, winners := foo(0, goal, []int{}, containers)
	fmt.Printf("found=%d\n", found)
	//fmt.Printf("winners=%v\n", winners)

	minWinnerLen := 0
	numMinWinnerLen := 0
	for _, winner := range winners {
		if minWinnerLen == 0 || len(winner) < minWinnerLen {
			minWinnerLen = len(winner)
			numMinWinnerLen = 1
		} else if len(winner) == minWinnerLen {
			numMinWinnerLen++
		}
	}

	fmt.Printf("min winner len %d (%d)\n", minWinnerLen, numMinWinnerLen)

}
