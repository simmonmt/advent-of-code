package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/simmonmt/aoc/2025/common/filereader"
	"github.com/simmonmt/aoc/2025/common/logger"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Command struct {
	Dir  byte
	Dist int
}

type Input struct {
	Commands []Command
}

func parseInput(lines []string) (*Input, error) {
	cmds := []Command{}

	for i, line := range lines {
		if line[0] != 'L' && line[0] != 'R' {
			return nil, fmt.Errorf("%d: bad dir", i)
		}
		dist, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, fmt.Errorf("%d: bad dist: %v", i, err)
		}

		cmds = append(cmds, Command{line[0], dist})
	}

	return &Input{cmds}, nil
}

func solveA(input *Input) int {
	pos := 50
	num := 0

	for _, cmd := range input.Commands {
		if cmd.Dir == 'R' {
			pos = (pos + cmd.Dist) % 100
		} else {
			pos -= cmd.Dist % 100
			if pos < 0 {
				pos += 100
			}
		}

		if pos == 0 {
			num++
		}
	}

	return num
}

func solveB(input *Input) int {
	pos := 50
	num := 0

	for _, cmd := range input.Commands {
		dist := cmd.Dist
		num += dist / 100
		dist -= (dist / 100) * 100

		if dist == 0 {
			continue
		}

		if cmd.Dir == 'R' {
			pos += dist
			if pos >= 100 {
				num++
				pos -= 100
			}
		} else {
			was := pos
			pos -= dist
			if pos < 0 {
				pos += 100
				if was != 0 {
					num++
				}
			} else if pos == 0 {
				num++
			}
		}

		logger.Infof("%s %d %d %d", string(cmd.Dir), cmd.Dist, pos, num)
	}

	return num
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
