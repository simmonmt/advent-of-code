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
	var freq int64

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		val, err := strconv.ParseInt(strings.TrimPrefix(line, "+"), 0, 32)
		if err != nil {
			log.Fatalf("failed to parse %v: %v", line, err)
		}

		freq += val
		fmt.Printf("%v %v %v\n", line, val, freq)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read failed: %v", err)
	}
}
