package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2019/07/b/src/amp"
	"github.com/simmonmt/aoc/2019/07/b/src/vm"
	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	ramPath = flag.String("ram", "", "path to file containing ram values")
	input   = flag.String("input", "", "input values")
)

func readRam(path string) (vm.Ram, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var line string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	ram := vm.NewRam()
	for i, str := range strings.Split(line, ",") {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", str, err)
		}
		ram.Write(i, val)
	}

	return ram, nil
}

func parseInput(inputStr string) ([]int, error) {
	out := []int{}
	for _, s := range strings.Split(inputStr, ",") {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse input value %v: %v", s, err)
		}

		out = append(out, v)
	}
	return out, nil
}

func tryPhaseCombination(phases []int, ram vm.Ram) (int, error) {
	logger.LogF("trying phase combination %v\n", phases)

	amps := make([]*amp.Amp, 5)
	for i, phase := range phases {
		amps[i] = amp.Start(phase, ram.Clone())
	}

	amps[0].In <- &vm.ChanIOMessage{Val: 0}

	last := 0
	for {
		for i := range amps {
			msg, ok := <-amps[i].Out
			if !ok {
				logger.LogF("amp %d stopping show for %v result %v\n",
					i, phases, last)
				return last, nil
			}

			if msg.Err != nil {
				return 0, fmt.Errorf("amp %d error: %v", i, msg.Err)
			}
			signal := msg.Val

			if i == 4 {
				last = signal
			}

			dest := (i + 1) % 5
			logger.LogF("amp %d out %v to amp %d\n", i, msg, dest)
			amps[dest].In <- msg
		}
	}
}

func tryAllPhases(ram vm.Ram) ([]int, int, error) {
	phases := []int{5, 5, 5, 5, 5}
	var maxPhases [5]int
	maxResult := 0
	for {
		var phaseCounts [10]bool
		hasRepeats := false
		for _, v := range phases {
			if phaseCounts[v] {
				hasRepeats = true
				break
			} else {
				phaseCounts[v] = true
			}
		}
		if !hasRepeats {
			result, err := tryPhaseCombination(phases[:], ram)
			if err != nil {
				return nil, 0, err
			}

			if result > maxResult {
				maxResult = result
				copy(maxPhases[:], phases)
			}
		}

		var i int
		for i = 0; i < len(phases); i++ {
			phases[i]++
			if phases[i] == 10 {
				phases[i] = 5
			} else {
				break
			}
		}
		if i == len(phases) {
			break
		}
	}

	return maxPhases[:], maxResult, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *ramPath == "" {
		log.Fatalf("--ram is required")
	}

	ram, err := readRam(*ramPath)
	if err != nil {
		log.Fatal(err)
	}

	if *input != "" {
		inputValues, err := parseInput(*input)
		if err != nil {
			log.Fatal(err)
		}

		result, err := tryPhaseCombination(inputValues, ram)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("result %v\n", result)
	} else {
		seq, result, err := tryAllPhases(ram)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("seq %v result %v\n", seq, result)
	}
}
