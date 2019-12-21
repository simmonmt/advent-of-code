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

func sumDigits(in []uint8, prevSum int, prevStart, prevEnd, start, end int) int {
	if prevStart < 0 {
		// There is no previous sum, so calculate the hard way
		out := 0
		for i := start; i < len(in) && i <= end; i++ {
			out += int(in[i])
		}
		return out
	}

	// There is a previous sum which we can reuse. We'll be asked
	// to sum overlapping ranges, with the new range starting
	// after the first starts and ending after the end starts. So
	// the first call may be for [7,13] while the next is for
	// [8,15]. If that's the case, prevSum is the sum of digits
	// [7,13], with prevStart=7 and prevEnd=13. To calculate
	// [8,15], we'll throw out digit 7 and add 14 and 15. This is
	// silly when the ranges are tight like this, but in practice
	// they'll be much longer so the optimization is worth it.
	out := prevSum
	for i := prevStart; i < start && i < len(in); i++ {
		out -= int(in[i])
	}

	for i := prevEnd + 1; i <= end && i < len(in); i++ {
		out += int(in[i])
	}
	return out

}

func calculate(in []uint8, off int) (out []uint8) {
	out = make([]uint8, len(in))
	lastPos, lastNeg := 0, 0
	lastStartPos, lastStartNeg := -1, -1
	lastEndPos, lastEndNeg := -1, -1
	for digit := off; digit < len(in); digit++ {
		posStart, posEnd := digit, 2*digit-1
		negStart, negEnd := 3*digit, 3*digit-1
		pos := sumDigits(in, lastPos, lastStartPos, lastEndPos, posStart, posEnd)
		neg := sumDigits(in, lastNeg, lastStartNeg, lastEndNeg, negStart, negEnd)

		lastPos, lastStartPos, lastEndPos = pos, posStart, posEnd
		lastNeg, lastStartNeg, lastEndNeg = neg, negStart, negEnd

		d1 := uint8(intmath.Abs(pos-neg) % 10)

		out[digit] = d1
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
