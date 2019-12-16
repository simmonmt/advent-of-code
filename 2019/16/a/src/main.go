package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/simmonmt/aoc/2019/common/intmath"
	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read failed: %v", err)
	}

	panic("unreachable")
}

func makePattern(base []int, rep int) []int {
	out := make([]int, len(base)*rep)

	off := 0
	for _, b := range base {
		for j := 0; j < rep; j++ {
			out[off+j] = b
		}
		off += rep
	}

	return out
}

func calculate(in string) (out string) {
	basePattern := []int{0, 1, 0, -1}
	for digitIdx := range in {
		//fmt.Printf("digit %d\n", digitIdx)

		pattern := makePattern(basePattern, digitIdx+1)
		//fmt.Printf("using pattern %v\n", pattern)

		accum := 0
		for i, r := range in {
			d := int(r - '0')
			m := pattern[(i+1)%len(pattern)]
			accum += d * m
			//fmt.Printf("term %d: %d * %d = %d => %d\n", i, d, m, d*m, accum)
		}
		accum = intmath.Abs(accum % 10)
		//fmt.Printf("result: %d\n", accum)

		out += strconv.Itoa(accum)
	}

	return out
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	line, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	signal := line
	for phase := 1; phase <= 100; phase++ {
		out := calculate(signal)
		//fmt.Printf("phase %4d %20s => %20s\n", phase, signal, out)
		signal = out
	}

	fmt.Printf("after 100 phases, first 8 are: %s\n", signal[0:8])
}
