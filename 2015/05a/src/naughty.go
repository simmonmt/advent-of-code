package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func isNice(line string) bool {
	numVowels, numDups, numBad := 0, 0, 0

	var last rune
	for _, c := range line {
		if c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' {
			numVowels++
		}

		if last == c {
			numDups++
		}

		if (last == 'a' && c == 'b') ||
			(last == 'c' && c == 'd') ||
			(last == 'p' && c == 'q') ||
			(last == 'x' && c == 'y') {
			numBad++
		}

		last = c
	}

	return numVowels >= 3 && numDups > 0 && numBad == 0
}

func main() {
	numNice := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if isNice(line) {
			fmt.Printf("nice: %v\n", line)
			numNice++
		} else {
			fmt.Printf("naughty: %v\n", line)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading stdin: %v", err)
	}

	fmt.Println(numNice)
}
