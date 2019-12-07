package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	vm "github.com/simmonmt/aoc/2019/05/a/src/vm"
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

func initRam(line string) (vm.Ram, error) {
	ram := vm.NewRam()
	for i, str := range strings.Split(line, ",") {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", str, err)
		}
		ram.Write(i, val)
	}

	return ram, nil
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

	ram, err := initRam(lines[0])
	if err != nil {
		log.Fatalf("ram init fail: %v", err)
	}

	ram.Write(1, 12)
	ram.Write(2, 2)
	if err := vm.Run(0, ram); err != nil {
		log.Fatalf("program failed: %v", err)
	}
	fmt.Printf("ram[0] = %v\n", ram.Read(0))
}
