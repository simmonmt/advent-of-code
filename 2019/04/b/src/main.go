package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
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
	for i := 1; i < len(arr); i++ {
		if arr[i-1] > arr[i] {
			return false
		}
	}

	hasDouble := false
	last := arr[0]
	streak := 0
	for i := 1; i < len(arr); i++ {
		if arr[i] == last {
			streak++
		} else {
			if streak == 1 {
				hasDouble = true
			}
			streak = 0
		}
		last = arr[i]
	}
	if streak == 1 {
		hasDouble = true
	}

	return hasDouble
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range lines {
		parts := strings.Split(line, "-")

		from, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatalf("bad from %v: %v", parts[0], err)
		}

		to, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatalf("bad to %v: %v", parts[1], err)
		}

		numMatches := 0
		for i := from; i <= to; i++ {
			if matches(strconv.Itoa(i)) {
				numMatches++
			}
		}

		fmt.Printf("%d-%d matches %d\n", from, to, numMatches)
	}
}
