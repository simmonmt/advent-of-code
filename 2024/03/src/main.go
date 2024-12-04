package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	mulPattern = regexp.MustCompile(`mul\(([0-9]+),([0-9]+)\)`)
	doPattern  = regexp.MustCompile(`do(?:n't)?\(\)`)
)

func parseInput(lines []string) (string, error) {
	out := ""
	for _, line := range lines {
		if len(out) > 0 {
			out += " "
		}
		out += line
	}

	return out, nil
}

func solveA(input string) int64 {
	results := mulPattern.FindAllStringSubmatch(input, -1)

	sum := 0
	for _, result := range results {
		a, _ := strconv.Atoi(result[1])
		b, _ := strconv.Atoi(result[2])
		sum += a * b
	}

	return int64(sum)
}

func solveB(input string) int64 {
	in := []byte(input)
	do := true

	sum := 0

	consumeMul := func(match []int) {
		logger.Infof("consumeMul %s", string(in[match[0]:match[1]]))

		if do {
			a, _ := strconv.Atoi(string(in[match[2]:match[3]]))
			b, _ := strconv.Atoi(string(in[match[4]:match[5]]))
			sum += a * b
			logger.Infof("sum %d", sum)
		}

		in = in[match[1]:]
	}

	consumeDo := func(match []int) {
		logger.Infof("consumeDo %s", string(in[match[0]:match[1]]))

		verb := string(in[match[0]:match[1]])
		do = verb == "do()"
		logger.Infof("do %v", do)
		in = in[match[1]:]
	}

	for i := 0; ; i++ {
		if i > 10000 {
			panic("runaway")
		}

		nextMul := mulPattern.FindSubmatchIndex(in)
		nextDo := doPattern.FindSubmatchIndex(in)

		if nextMul != nil && nextDo != nil {
			if nextMul[0] < nextDo[0] {
				consumeMul(nextMul)
			} else {
				consumeDo(nextDo)
			}
		} else if nextMul != nil {
			consumeMul(nextMul)
		} else if nextDo != nil {
			consumeDo(nextDo)
		} else {
			// no more
			break
		}
	}

	return int64(sum)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input))
	fmt.Println("B", solveB(input))
}
