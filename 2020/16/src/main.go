// NOTE: This doesn't even pretend to be an efficient solution. This
// is a "I have some time before my meetings start" solution. The size
// of the problem is such that optimization might even be a waste of
// time.

package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	fieldSpecPattern = regexp.MustCompile(
		`^([^:]+): ([0-9]+)-([0-9]+) or ([0-9]+)-([0-9]+)$`)
)

type FieldSpec struct {
	name       string
	a, b, c, d int
}

func newFieldSpec(name string, a, b, c, d int) *FieldSpec {
	return &FieldSpec{
		name: name,
		a:    a,
		b:    b,
		c:    c,
		d:    d,
	}
}

func (fs *FieldSpec) Name() string {
	return fs.name
}

func (fs *FieldSpec) Valid(num int) bool {
	return (num >= fs.a && num <= fs.b) || (num >= fs.c && num <= fs.d)
}

func (fs *FieldSpec) String() string {
	return fmt.Sprintf("%v:[%d-%d],[%d-%d]",
		fs.name, fs.a, fs.b, fs.c, fs.d)
}

func parseTicket(str string) ([]int, error) {
	parts := strings.Split(str, ",")
	out := make([]int, len(parts))

	for i := 0; i < len(parts); i++ {
		val, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, err
		}
		out[i] = val
	}

	return out, nil
}

func except(l []string, remove string) []string {
	out := []string{}
	for _, e := range l {
		if e == remove {
			continue
		}
		out = append(out, e)
	}
	return out
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		log.Fatal(err)
	}

	specs := map[string]*FieldSpec{}
	lineNo := 0
	for ; lines[lineNo] != ""; lineNo++ {
		line := lines[lineNo]

		parts := fieldSpecPattern.FindStringSubmatch(line)
		if parts == nil {
			log.Fatalf("%d: failed to parse", lineNo)
		}

		name := parts[1]
		a := intmath.AtoiOrDie(parts[2])
		b := intmath.AtoiOrDie(parts[3])
		c := intmath.AtoiOrDie(parts[4])
		d := intmath.AtoiOrDie(parts[5])
		specs[name] = newFieldSpec(name, a, b, c, d)
	}
	lineNo++ // skip blank

	if lines[lineNo] != "your ticket:" {
		log.Fatal("no your ticket")
	}
	lineNo++

	myTicket, err := parseTicket(lines[lineNo])
	if err != nil {
		log.Fatal("bad my ticket")
	}
	lineNo += 2

	if lines[lineNo] != "nearby tickets:" {
		log.Fatal("no nearby tickets")
	}
	lineNo++

	validTickets := [][]int{}
	invalidRate := 0
	for ; lineNo < len(lines); lineNo++ {
		line := lines[lineNo]

		nums, err := parseTicket(line)
		if err != nil {
			log.Fatalf("%d: %v", lineNo, err)
		}

		logger.LogF("checking nearby %v => %v", line, nums)
		validTicket := true
		for _, num := range nums {
			validField := false
			for _, fs := range specs {
				if fs.Valid(num) {
					logger.LogF("%v valid for %v", num, fs)
					validField = true
					break
				}
			}

			if !validField {
				invalidRate += num
				validTicket = false
			}
		}

		if validTicket {
			validTickets = append(validTickets, nums)
		}
	}

	fmt.Printf("A: invalid rate %v\n", invalidRate)

	logger.LogF("%d valid tickets", len(validTickets))

	// field number to possible names
	cands := map[int][]string{}
	names := []string{}
	for name := range specs {
		names = append(names, name)
	}

	for i := range validTickets[0] {
		n := make([]string, len(names))
		copy(n, names)
		cands[i] = n
	}
	logger.LogF("cands: %v", cands)

	for _, nums := range validTickets {
		logger.LogF("examining %v", nums)
		for i, num := range nums {
			for _, fs := range specs {
				if !fs.Valid(num) {
					logger.LogF(
						"field %d: %v not valid for %v",
						i, num, fs.Name())
					cands[i] = except(cands[i], fs.Name())
				} else {
					logger.LogF(
						"field %d: %v valid for %v",
						i, num, fs.Name())
				}
			}
		}
		logger.LogF("cands: %v", cands)
	}

	assignments := map[string]int{}

	for len(cands) > 0 {
		one := -1
		for fieldNum, names := range cands {
			if len(names) == 1 {
				one = fieldNum
			}
		}

		if one == -1 {
			panic("no one found")
		}

		name := cands[one][0]
		assignments[name] = one

		delete(cands, one)
		fieldNums := []int{}
		for i := range cands {
			fieldNums = append(fieldNums, i)
		}
		for _, fieldNum := range fieldNums {
			cands[fieldNum] = except(cands[fieldNum], name)
		}
	}

	logger.LogF("assignments %v", assignments)

	result := 1
	for name, num := range assignments {
		if strings.HasPrefix(name, "departure") {
			result *= myTicket[num]
		}
	}

	fmt.Printf("B: %v\n", result)
}
