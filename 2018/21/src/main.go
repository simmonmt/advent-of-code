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
	regZero = flag.Int64("reg_zero", 0, "register 0 initial value")

	commentPattern = regexp.MustCompile(`#.*$`)
	ipReg          = 0
)

type InstrLine struct {
	Op      string
	A, B, C int64
}

func parseInstr(str string) *InstrLine {
	parts := strings.Split(str, " ")
	if len(parts) != 4 {
		panic("bad parse: " + str)
	}

	il := InstrLine{}
	il.Op = parts[0]
	il.A = intmath.Atoi64OrDie(parts[1])
	il.B = intmath.Atoi64OrDie(parts[2])
	il.C = intmath.Atoi64OrDie(parts[3])

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

		match := commentPattern.FindStringIndex(line)
		if match != nil {
			line = strings.TrimRight(line[0:match[0]], " ")
		}
		if line == "" {
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
	regFile[0] = *regZero

	var ip int64
	for ip < int64(len(instrs)) {
		inst := instrs[ip]
		desc, found := descByName[inst.Op]
		if !found {
			panic("bad inst " + inst.Op)
		}

		regFile[ipReg] = ip

		if ip == 24 {
			fmt.Println(regFile)
		}

		oldFile := regFile
		desc.F(&regFile, inst.A, inst.B, inst.C)
		logger.LogF("executed ip %v\t%-20v\t%-40s => %-40s",
			ip, fmt.Sprint(*inst), fmt.Sprint(oldFile),
			fmt.Sprint(regFile))
		//fmt.Println(regFile)

		ip = regFile[ipReg]
		ip++
	}

	fmt.Println("done")
	fmt.Println(regFile)

}
