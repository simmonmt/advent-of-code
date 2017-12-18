package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		firstBasement := -1
		floor := 0
		for i, c := range line {
			switch c {
			case '(':
				floor++
				break
			case ')':
				floor--
				break
			default:
				log.Fatalf("unknown char %v\n", c)
			}

			if floor == -1 && firstBasement == -1 {
				firstBasement = i + 1
			}
		}

		fmt.Printf("%d first basement: %v\n", lineNum, firstBasement)
		fmt.Printf("%d end floor: %v\n", lineNum, floor)
	}
}
