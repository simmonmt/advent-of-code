package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func count(line string) (int, int) {
	found := map[rune]int{}
	for _, c := range line {
		found[c]++
	}

	numTwo := 0
	numThree := 0
	for _, num := range found {
		if num == 2 {
			numTwo++
		}
		if num == 3 {
			numThree++
		}
	}

	return numTwo, numThree
}

func main() {
	hasTwice := 0
	hasThrice := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		numTwo, numThree := count(line)
		fmt.Printf("line: %v #2: %v #3: %v\n", line, numTwo, numThree)

		if numTwo > 0 {
			hasTwice++
		}
		if numThree > 0 {
			hasThrice++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read failed: %v", err)
	}

	fmt.Println(hasTwice * hasThrice)
}
