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

func sumDigits(in []uint8, prevSum int, prevStart, start, end int) int {
	out := 0
	for i := start; i < len(in) && i <= end; i++ {
		out += int(in[i])
	}
	return out
}

func calculate(in []uint8, off int) (out []uint8) {
	out = make([]uint8, len(in))
	// lastSumPos, lastSumNeg := 0, 0
	// lastStartPos, lastStartNeg := -1, -1
	for digit := off; digit < len(in); digit++ {
		// pos := sumDigits(in, lastSumPos, lastStartPos, digit, 2*digit-1)
		// neg := sumDigits(in, lastSumNeg, lastStartNeg, 3*digit, 3*digit-1)

		// out[digit] = uint8(intmath.Abs(pos-neg) * 10)

		out[digit] = uint8(intmath.Abs(sumDigits(in, 0, 0, digit, 2*digit-1)-
			sumDigits(in, 0, 0, 3*digit, 4*digit-1)) % 10)
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

	vals := make([]uint8, len(line))
	for i, r := range line {
		vals[i] = uint8(r - '0')
	}

	signal := make([]uint8, len(line)*10000)
	for i := 0; i < 10000; i++ {
		off := i * len(line)
		copy(signal[off:off+len(line)], vals)
	}

	fmt.Printf("signal len %d\n", len(signal))

	offStr := line[0:7]
	off, err := strconv.Atoi(offStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("offset is %d\n", off)

	if off*4 < len(signal) {
		panic("needs repeat")
	}

	for phase := 1; phase <= 100; phase++ {
		out := calculate(signal, off)
		fmt.Printf("phase %4d\n", phase)
		signal = out
	}

	fmt.Printf("after 100 phases, selected 8 at %d: ", off)
	for _, d := range signal[off : off+8] {
		fmt.Print(d)
	}
	fmt.Println()
}
