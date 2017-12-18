package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func sortWord(word string) string {
	chars := []int{}
	for _, r := range word {
		chars = append(chars, int(r))
	}
	sort.Ints(chars)

	runes := []rune{}
	for _, c := range chars {
		runes = append(runes, rune(c))
	}

	return string(runes)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	numValid := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		wordList := strings.Split(line, " ")
		valid := true
		words := map[string]bool{}
		for _, word := range wordList {
			sortedWord := sortWord(word)

			if _, found := words[sortedWord]; found {
				valid = false
				break
			}
			words[sortedWord] = true
		}

		if valid {
			fmt.Printf("  valid %s\n", line)
			numValid++
		} else {
			fmt.Printf("invalid %s\n", line)
		}

		//fmt.Printf("%v %s\n", valid, wordList)
	}

	fmt.Printf("num valid %d\n", numValid)
}
