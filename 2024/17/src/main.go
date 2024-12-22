package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2024/common/filereader"
	"github.com/simmonmt/aoc/2024/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

	regPattern = regexp.MustCompile(`^Register ([ABC]): ([-0-9]+)$`)
)

type Input struct {
	A, B, C int
	Mem     []byte
}

func parseInput(lines []string) (*Input, error) {
	groups := filereader.BlankSeparatedGroupsFromLines(lines)
	if len(groups) != 2 {
		return nil, fmt.Errorf("bad num groups: %d", len(groups))
	}

	regs := map[string]int{}
	for i, line := range groups[0] {
		parts := regPattern.FindStringSubmatch(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("%d: bad register line", i+1)
		}

		name := parts[1]
		val, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("%d: bad register value: %v", i+1, err)
		}

		regs[name] = val
	}

	if len(groups[1]) != 1 {
		return nil, fmt.Errorf("bad program line num lines; found %d",
			len(groups[1]))
	}

	_, nums, found := strings.Cut(groups[1][0], " ")
	if !found {
		return nil, fmt.Errorf("bad program line")
	}

	memInt, err := filereader.ParseNumbersFromLine(nums, ",")
	if err != nil {
		return nil, fmt.Errorf("bad program line: %v", err)
	}
	mem := make([]byte, len(memInt))
	for i, v := range memInt {
		mem[i] = byte(v)
	}

	return &Input{
		A: regs["A"], B: regs["B"], C: regs["C"],
		Mem: mem,
	}, nil
}

func runProgram(mem []byte, regs map[string]int) []byte {
	pc := 0

	out := []byte{}
	for num := 1; ; num++ {
		if pc >= len(mem) {
			break
		}

		inst, err := ParseInst(mem[pc], mem[pc+1])
		if err != nil {
			panic(fmt.Sprintf("bad inst: %v", err))
		}

		var instOut []byte
		pc, instOut = inst.Execute(regs, pc)
		out = append(out, instOut...)
	}

	return out
}

func solveA(input *Input) string {
	regs := map[string]int{
		"A": input.A,
		"B": input.B,
		"C": input.C,
	}

	out := []string{}
	for _, v := range runProgram(input.Mem, regs) {
		out = append(out, strconv.Itoa(int(v)))
	}
	return strings.Join(out, ",")
}

func outputSame(got, want []byte) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func solveB(input *Input) int64 {
	cands := []int{0}

	for i := len(input.Mem) - 1; i >= 0; i-- {
		want := input.Mem[i:]
		matched := []int{}

		for _, cand := range cands {
			cand <<= 3

			for j := 0; j < 8; j++ {
				a := cand | j

				regs := map[string]int{
					"A": a,
					"B": input.B,
					"C": input.C,
				}

				got := runProgram(input.Mem, regs)
				if outputSame(got, want) {
					matched = append(matched, a)
				}
			}
		}

		cands = matched
	}

	sort.Ints(cands)
	return int64(cands[0])
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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
