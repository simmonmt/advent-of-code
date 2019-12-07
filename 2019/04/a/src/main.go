package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	from    = flag.Int("from", -1, "from")
	to      = flag.Int("to", -1, "to")
)

func readInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func matches(num string) bool {
	arr := []byte(num)
	hasDouble := false
	for i := 1; i < len(arr); i++ {
		if arr[i-1] == arr[i] {
			hasDouble = true
		} else if arr[i-1] > arr[i] {
			return false
		}
	}

	return hasDouble
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *from == -1 {
		log.Fatalf("--from is required")
	}
	if *to == -1 {
		log.Fatalf("--to is required")
	}

	numMatches := 0
	for i := *from; i <= *to; i++ {
		if matches(strconv.Itoa(i)) {
			numMatches++
		}
	}

	fmt.Printf("matches %d\n", numMatches)
}
