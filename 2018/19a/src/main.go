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
	verbose        = flag.Bool("verbose", false, "verbose")
	regFilePattern = regexp.MustCompile(`^[^ ]+ +\[(\d+), (\d+), (\d+), (\d+)\]$`)

	ipReg = 0
)

type InstrLine struct {
	Op      string
	A, B, C int
}

func parseInstr(str string) *InstrLine {
	parts := strings.Split(str, " ")
	if len(parts) != 4 {
		panic("bad parse: " + str)
	}

	il := InstrLine{}
	il.Op = parts[0]
	il.A = intmath.AtoiOrDie(parts[1])
	il.B = intmath.AtoiOrDie(parts[2])
	il.C = intmath.AtoiOrDie(parts[3])

	return &il
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

func readInstrs() (int, []*InstrLine, error) {
	lines, err := readLines()
	if err != nil {
		return -1, nil, err
	}

	ipReg := -1
	instrs := []*InstrLine{}
	for _, line := range lines {
		if strings.HasPrefix(line, "#ip") {
			ipReg = intmath.AtoiOrDie(strings.Split(line, " ")[1])
			continue
		}

		instrs = append(instrs, parseInstr(line))
	}

	if ipReg == -1 {
		panic("no #ip")
	}

	return ipReg, instrs, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	ipReg, instrs, err := readInstrs()
	if err != nil {
		log.Fatal(err)
	}

	descByName := map[string]instr.Desc{}
	for _, desc := range instr.All {
		descByName[desc.Name] = desc
	}

	regFile := reg.File{}

	ip := 0
	for ip < len(instrs) {
		inst := instrs[ip]
		desc, found := descByName[inst.Op]
		if !found {
			panic("bad inst " + inst.Op)
		}

		regFile[ipReg] = ip

		logger.LogF("executing ip %v %v, desc %v", ip, inst, desc)
		desc.F(&regFile, inst.A, inst.B, inst.C)
		logger.LogF("regfile now %v", regFile)

		ip = regFile[ipReg]
		ip++
	}

	fmt.Printf("done")
	fmt.Println(regFile)

}
