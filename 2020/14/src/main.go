// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// 15470975766942 too high
// 77524251746 too low

package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	cmdPattern = regexp.MustCompile(`mem\[([0-9]+)\] = ([0-9]+)$`)
)

func parseMask(str string) (or, and uint64, notBits []int) {
	or, and = 0, 0xfffffffff
	notBits = []int{}

	if len(str) != 36 {
		panic("bad mask len")
	}
	for i, r := range str {
		if r == '1' {
			or |= 1 << (35 - i)
		} else if r == '0' {
			//fmt.Printf("r=0 at %d;\n  %v =>\n  %v\n",
			//	i, decToBin(and), decToBin(and^(1<<(35-i))))
			and ^= 1 << (35 - i)
		} else if r == 'X' {
			notBits = append(notBits, 35-i)
		}
	}

	return
}

func decToBin(val uint64) string {
	out := ""
	for i := 0; i < 36; i++ {
		if (val & 1) == 1 {
			out = "1" + out
		} else {
			out = "0" + out
		}
		val >>= 1
	}
	return out
}

func dumpStep(value, result uint64, mask string) {
	fmt.Printf("value:  %v\n", decToBin(value))
	fmt.Printf("mask:   %v\n", mask)
	fmt.Printf("result: %v\n", decToBin(result))
}

func dumpMask(mask string, or, and uint64) {
	fmt.Printf("mask: %v\n", mask)
	fmt.Printf("or:   %v\n", decToBin(or))
	fmt.Printf("and:  %v\n", decToBin(and))
}

func solveA(lines []string) uint64 {
	mem := map[uint64]uint64{}
	orMask, andMask := uint64(0), uint64(0xfffffffff)
	maskStr := ""
	for lineNo, line := range lines {
		if strings.HasPrefix(line, "mask = ") {
			maskStr = lines[lineNo][len("mask = "):]
			orMask, andMask, _ = parseMask(maskStr)
			logger.LogF("%d: ormask %x andmask %x",
				lineNo, orMask, andMask)
			//dumpMask(maskStr, orMask, andMask)
			continue
		}

		parts := cmdPattern.FindStringSubmatch(line)
		if parts == nil {
			log.Fatalf(`%d: bad match on "%v"`, lineNo, line)
		}

		loc := uint64(intmath.AtoiOrDie(parts[1])) & 0xfffffffff
		val := uint64(intmath.AtoiOrDie(parts[2])) & 0xfffffffff

		newVal := ((val | orMask) & andMask) & 0xfffffffff

		logger.LogF("line %d: mem %x val %d/%x => %d/%x",
			lineNo, loc, val, val, newVal, newVal)
		//dumpStep(val, newVal, maskStr)

		mem[loc] = newVal
	}

	sum := uint64(0)
	for _, val := range mem {
		sum += val
	}
	return sum
}

func allAddrs(cmdAddr, orMask, andMask uint64, notBits []int) []uint64 {
	out := []uint64{}
	base := (cmdAddr | orMask) // & andMask

	logger.LogF("cmdAddr %v\nor      %v\nand     %v\nnot %v",
		decToBin(cmdAddr), decToBin(orMask), decToBin(andMask), notBits)
	for i := 0; i < (1 << len(notBits)); i++ {
		addr := base
		logger.LogF("i=%v", i)
		logger.LogF("start    %v", decToBin(addr))
		for j, bit := range notBits {
			mask := uint64(1) << bit
			logger.LogF("bit %v => %v", bit, decToBin(mask))

			addr &= ^mask
			if (1<<j)&i != 0 {
				logger.LogF("bit %v on", bit)
				addr |= mask
			}
		}
		logger.LogF("result   %v", decToBin(addr))
		out = append(out, addr)
	}

	logger.LogF("=> %v", out)

	return out
}

func solveB(lines []string) uint64 {
	mem := map[uint64]uint64{}
	orMask, andMask := uint64(0), uint64(0xfffffffff)
	notBits := []int{}
	maskStr := ""
	for lineNo, line := range lines {
		if strings.HasPrefix(line, "mask = ") {
			maskStr = lines[lineNo][len("mask = "):]
			orMask, andMask, notBits = parseMask(maskStr)
			logger.LogF("%d: ormask %x andmask %x",
				lineNo, orMask, andMask)
			continue
		}

		parts := cmdPattern.FindStringSubmatch(line)
		if parts == nil {
			log.Fatalf(`%d: bad match on "%v"`, lineNo, line)
		}

		cmdAddr := uint64(intmath.AtoiOrDie(parts[1])) & 0xfffffffff
		cmdVal := uint64(intmath.AtoiOrDie(parts[2])) & 0xfffffffff

		for _, addr := range allAddrs(cmdAddr, orMask, andMask, notBits) {
			mem[addr] = cmdVal
		}
	}

	sum := uint64(0)
	for _, val := range mem {
		sum += val
	}
	return sum
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

	fmt.Printf("A: %v\n", solveA(lines))
	fmt.Printf("B: %v\n", solveB(lines))
}
