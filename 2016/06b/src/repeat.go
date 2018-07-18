package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func readInput(r io.Reader) ([]string, error) {
	lines := []string{}

	lineLen := -1
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if lineLen != -1 && lineLen != len(line) {
			return nil, fmt.Errorf("uneven lines; found %v and %v", lineLen, len(line))
		}

		lines = append(lines, strings.TrimSpace(line))
	}

	return lines, nil
}

func main() {
	lines, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err.Error())
	}

	for i := 0; i < len(lines[0]); i++ {
		freqs := map[rune]int{}

		for _, line := range lines {
			freqs[rune(line[i])]++
		}
		minFreq := -1
		minFreqChar := ' '
		for r, f := range freqs {
			if minFreq == -1 || f < minFreq {
				minFreq = f
				minFreqChar = r
			}
		}

		fmt.Printf("%c", minFreqChar)
	}

	fmt.Println()
}
