package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func diffByOne(a, b string) bool {
	aArr := []rune(a)
	bArr := []rune(b)

	numDiff := 0
	for i := range aArr {
		if aArr[i] != bArr[i] {
			numDiff++
			if numDiff > 1 {
				return false
			}
		}
	}
	return numDiff == 1
}

func common(a, b string) string {
	comm := []rune{}

	aChars := []rune(a)
	bChars := []rune(b)

	for i := range a {
		if aChars[i] == bChars[i] {
			comm = append(comm, aChars[i])
		}
	}

	return string(comm)
}

func main() {
	lines := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read failed: %v", err)
	}

	for i := 0; i < len(lines); i++ {
		for j := i + 1; j < len(lines); j++ {
			if diffByOne(lines[i], lines[j]) {
				fmt.Printf("%v\n", common(lines[i], lines[j]))
			}
		}
	}
}
