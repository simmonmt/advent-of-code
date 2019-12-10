package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	vm "github.com/simmonmt/aoc/2019/09/a/src/vm"
	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	ramPath = flag.String("ram", "", "path to file containing ram values")
	input   = flag.String("input", "", "input values")
)

func readRam(path string) (vm.Ram, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var line string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	ram := vm.NewRam()
	for i, str := range strings.Split(line, ",") {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", str, err)
		}
		ram.Write(int64(i), int64(val))
	}

	return ram, nil
}

func parseInput(inputStr string) ([]int64, error) {
	out := []int64{}
	for _, s := range strings.Split(inputStr, ",") {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse input value %v: %v", s, err)
		}

		out = append(out, int64(v))
	}
	return out, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *ramPath == "" {
		log.Fatalf("--ram is required")
	}

	if *input == "" {
		log.Fatalf("--input is required")
	}

	ram, err := readRam(*ramPath)
	if err != nil {
		log.Fatal(err)
	}

	inputValues, err := parseInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	io := vm.NewSaverIO(inputValues...)

	if err := vm.Run(ram, io, 0); err != nil {
		log.Fatalf("program failed: %v", err)
	}

	fmt.Printf("output: %v\n", io.Written())
}
