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

func initRam(line string) (map[int]int, error) {
	ram := map[int]int{}
	for i, str := range strings.Split(line, ",") {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", str, err)
		}
		ram[i] = val
	}

	return ram, nil
}

func runProgram(ram map[int]int) {
	for pc := 0; ; pc += 4 {
		op := ram[pc]
		logger.LogF("pc=%d op %d %d %d %d",
			pc, op, ram[pc+1], ram[pc+2], ram[pc+3])
		switch op {
		case 1:
			ram[ram[pc+3]] = ram[ram[pc+1]] + ram[ram[pc+2]]
			break
		case 2:
			ram[ram[pc+3]] = ram[ram[pc+1]] * ram[ram[pc+2]]
			break
		case 99:
			logger.LogLn("exiting")
			return

		default:
			panic(fmt.Sprintf("bad opcode %d at %d", op, pc))
		}
	}
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

	for i, line := range lines {
		logger.LogF("program %d", i)
		ram, err := initRam(line)
		if err != nil {
			log.Fatalf("program %d ram init fail: %v", i, err)
		}

		logger.LogLn("read ram")

		for noun := 0; noun <= 99; noun++ {
			fmt.Printf("noun %d\n", noun)
			for verb := 0; verb <= 99; verb++ {
				ramCopy := map[int]int{}
				for k, v := range ram {
					ramCopy[k] = v
				}

				ramCopy[1] = noun
				ramCopy[2] = verb

				runProgram(ramCopy)

				if ramCopy[0] == 19690720 {
					fmt.Println(100*noun + verb)
					return
				}
			}
		}
	}
}
