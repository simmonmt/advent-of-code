package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"instr"
	"intmath"
	"logger"
	"reg"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")

	regFilePattern = regexp.MustCompile(`^[^ ]+ +\[(\d+), (\d+), (\d+), (\d+)\]$`)
)

func parseInstr(str string) (op, a, b, c int) {
	parts := strings.Split(str, " ")
	if len(parts) != 4 {
		panic("bad parse: " + str)
	}

	vals := [4]int{}
	for i, s := range parts {
		vals[i] = intmath.AtoiOrDie(s)
	}

	return vals[0], vals[1], vals[2], vals[3]
}

func readLines() ([]string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

type TestCase struct {
	Before, After *reg.File
	Op            int
	A, B, C       int
}

func (tc TestCase) String() string {
	return fmt.Sprintf("%v %d %d %d %d %v", *tc.Before, tc.Op, tc.A, tc.B, tc.C, *tc.After)
}

func readInput() ([]TestCase, error) {
	lines, err := readLines()
	if err != nil {
		return nil, err
	}

	testCases := []TestCase{}
	for i := 0; i < len(lines); i += 4 {
		beforeStr := lines[i]
		instStr := lines[i+1]
		afterStr := lines[i+2]

		before := reg.ParseFile(beforeStr)
		after := reg.ParseFile(afterStr)
		op, a, b, c := parseInstr(instStr)

		testCases = append(testCases, TestCase{
			Before: before,
			After:  after,
			Op:     op,
			A:      a,
			B:      b,
			C:      c,
		})
	}

	return testCases, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	testCases, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	goodCases := []TestCase{}
	for i, testCase := range testCases {
		numMatches := 0
		for _, desc := range instr.All {
			regFile := reg.File{}
			regFile = *testCase.Before

			desc.F(&regFile, testCase.Op, testCase.A, testCase.B, testCase.C)

			if regFile == *testCase.After {
				numMatches++
				logger.LogF("case %d %s match for %s", i, testCase, desc.Name)
			} else {
				logger.LogF("case %d %s fail for %s", i, testCase, desc.Name)
			}
		}

		if numMatches >= 3 {
			goodCases = append(goodCases, testCase)
		}
	}

	fmt.Println(len(goodCases))
}
