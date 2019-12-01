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

	all := 0
	for _, line := range lines {
		mass, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}

		sum := 0
		for mass > 0 {
			mass = mass/3 - 2
			if mass > 0 {
				sum += mass
			}
			//fmt.Printf("iter: mass=%v, sum=%v\n", mass, sum)
		}

		fmt.Printf("%v: %v\n", line, sum)
		all += sum
	}

	fmt.Println(all)
}
