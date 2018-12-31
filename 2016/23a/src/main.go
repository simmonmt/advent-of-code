package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"instr"
	"logger"
	"reg"
)

var (
	instrPattern = regexp.MustCompile(`^(...) ([^ ]+)(?: ([^ ]+))?$`)

	regInit = flag.String("reg_init", "", "initial values for registers as a=1,b=2")
	verbose = flag.Bool("verbose", false, "verbose")

	rewriteRules = map[string]string{
		"inc": "dec",
		"dec": "inc",
		"tgl": "inc",
		"jnz": "cpy",
		"cpy": "jnz",
	}
)

func readInput(r io.Reader) ([]string, error) {
	lines := []string{}

	reader := bufio.NewReader(os.Stdin)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		lines = append(lines, strings.TrimSpace(line))
	}

	return lines, nil
}

func parseInstr(str string) (instr.Instr, error) {
	matches := instrPattern.FindStringSubmatch(str)
	if matches == nil {
		return nil, fmt.Errorf("failed to parse instr: %v", str)
	}

	op := matches[1]
	a := matches[2]
	b := matches[3]

	i, err := instr.Parse(op, a, b)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func rewriteInstr(in string) string {
	parts := strings.SplitN(in, " ", 2)
	if rep, found := rewriteRules[parts[0]]; found {
		parts[0] = rep
	}
	return strings.Join(parts, " ")
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	lines, err := readInput(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	regFile := reg.NewFile()

	if *regInit != "" {
		for _, init := range strings.Split(*regInit, ",") {
			parts := strings.SplitN(init, "=", 2)
			reg, err := reg.FromString(parts[0])
			if err != nil {
				log.Fatalf("bad register %v: %v", parts[0], err)
			}

			val, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				log.Fatalf("bad val %v: %v", parts[1], err)
			}

			regFile.Set(reg, int(val))
		}

		regFile.Print()
	}

	pc := 0
	for {
		if pc >= len(lines) {
			break
		}

		line := lines[pc]
		logger.LogF("%v: %v", pc, line)

		i, err := parseInstr(line)
		switch {
		case err != nil:
			{
				//log.Fatalf("failed to parse instr '%v' at pc %d: %v",
				// line, pc, err)
				logger.LogF("failed to parse instr '%v' at pc %d: %v",
					line, pc, err)
				pc++
			}

		case i.IsTgl():
			{
				off := i.Exec(regFile)
				modAddr := pc + off
				if modAddr < len(lines) && modAddr >= 0 {
					mod := rewriteInstr(lines[modAddr])
					logger.LogF("modify %d; was %s now %s\n",
						modAddr, lines[modAddr], mod)
					lines[modAddr] = mod

					if *verbose {
						for a, l := range lines {
							logger.LogF("  %d: %s", a, l)
						}
					}
				} else {
					logger.LogF("toggle wants %d, which is out of bounds",
						modAddr)
				}
				pc++
			}
		default:
			pc += i.Exec(regFile)
		}

		if *verbose {
			regFile.Print()
		}
	}

	fmt.Println("done")
	regFile.Print()
}
