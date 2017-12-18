package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s captcha", os.Args[0])
	}

	digits := []int{}
	for _, c := range os.Args[1] {
		if n, err := strconv.Atoi(string(c)); err != nil {
			log.Fatalf("unable to parse input %v", os.Args[1])
		} else {
			digits = append(digits, n)
		}
	}

	if len(digits) == 0 {
		log.Fatalf("empty input")
	}

	accum := 0
	for i, digit := range digits {
		next := i + 1
		if next == len(digits) {
			next = 0
		}

		fmt.Printf("i=%d, digits[i]=%d, next=%d, digits[next]=%d", i, digits[i], next, digits[next])
		if digit == digits[next] {
			fmt.Printf(" add to accum")
			accum += digit
		}
		fmt.Println("")
	}

	fmt.Printf("out: %d\n", accum)
}
