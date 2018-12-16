package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
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

type TestCase struct {
	Before, After *reg.File
	Op            int
	A, B, C       int
}

func (tc TestCase) String() string {
	return fmt.Sprintf("%v %d %d %d %d %v", *tc.Before, tc.Op, tc.A, tc.B, tc.C, *tc.After)
}

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

type InstrLine struct {
	Op, A, B, C int
}

func readInstrs() ([]InstrLine, error) {
	lines, err := readLines()
	if err != nil {
		return nil, err
	}

	instrs := []InstrLine{}
	for _, line := range lines {
		op, a, b, c := parseInstr(line)
		instrs = append(instrs, InstrLine{op, a, b, c})
	}
	return instrs, nil
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

func findCands(testCases []TestCase) {
	testCases, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	cands := [16]map[string]bool{}

	goodCases := []TestCase{}
	for _, testCase := range testCases {
		matches := map[string]bool{}
		for _, desc := range instr.All {
			regFile := reg.File{}
			regFile = *testCase.Before

			desc.F(&regFile, testCase.Op, testCase.A, testCase.B, testCase.C)

			if regFile == *testCase.After {
				matches[desc.Name] = true
			}
		}

		if cands[testCase.Op] == nil {
			cands[testCase.Op] = map[string]bool{}
			for m := range matches {
				cands[testCase.Op][m] = true
			}
		} else {
			// Remove from cands anything that's not in matches
			bad := []string{}
			for c := range cands[testCase.Op] {
				if _, found := matches[c]; !found {
					bad = append(bad, c)
				}
			}

			for _, b := range bad {
				delete(cands[testCase.Op], b)
			}
		}
	}

	for i := range cands {
		names := []string{}
		for name := range cands[i] {
			names = append(names, name)
		}
		sort.Strings(names)
		fmt.Printf("%2d: %s\n", i, strings.Join(names, ", "))
	}

	fmt.Println(len(goodCases))
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	instrs, err := readInstrs()
	if err != nil {
		log.Fatal(err)
	}

	descByOp := map[int]instr.Desc{}
	for _, desc := range instr.All {
		descByOp[desc.Op] = desc
	}

	regFile := reg.File{}
	for _, inst := range instrs {
		desc := descByOp[inst.Op]
		fmt.Printf("%s %d %d %d\n", desc.Name, inst.A, inst.B, inst.C)
		desc.F(&regFile, inst.Op, inst.A, inst.B, inst.C)
		fmt.Println(regFile)
	}

}
