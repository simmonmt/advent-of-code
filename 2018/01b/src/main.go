package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	changes := []int64{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		val, err := strconv.ParseInt(strings.TrimPrefix(line, "+"), 0, 32)
		if err != nil {
			log.Fatalf("failed to parse %v: %v", line, err)
		}

		changes = append(changes, val)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read failed: %v", err)
	}

	fmt.Printf("read %v changes\n", len(changes))

	var freq int64
	seen := map[int64]bool{}
	seen[0] = true
	for {
		fmt.Println("loop")
		for _, change := range changes {
			freq += change
			fmt.Printf("with change %v freq now %v\n", change, freq)
			if _, found := seen[freq]; found {
				fmt.Printf("repeat %v\n", freq)
				os.Exit(0)
			}
			seen[freq] = true
		}
	}
}
