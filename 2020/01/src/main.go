package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simmonmt/aoc/2020/common/intmath"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func findTwo(small []int, large []int) {
	for _, sn := range small {
		for _, ln := range large {
			if sn+ln == 2020 {
				fmt.Printf("%d * %d = %d\n", sn, ln, sn*ln)
				return
			}
		}
	}
}

func findThree(all []int) {
	for i1, n1 := range all {
		for i2, n2 := range all {
			if i1 == i2 {
				break
			}

			sum := n1 + n2
			if sum > 2020 {
				continue
			}

			for i3, n3 := range all {
				if i3 == i1 || i3 == i2 {
					break
				}

				if sum+n3 == 2020 {
					fmt.Printf("%d * %d * %d = %d\n", n1, n2, n3, n1*n2*n3)
					return
				}
			}
		}
	}
	fmt.Println("no three found")
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	all, large, small := []int{}, []int{}, []int{}
	smallest := -1
	for _, line := range lines {
		num := intmath.AtoiOrDie(line)
		all = append(all, num)
		if num < 1000 {
			small = append(small, num)
		} else {
			large = append(large, num)
		}
		if smallest == -1 || num < smallest {
			smallest = num
		}
	}

	findTwo(small, large)

	findThree(all)
}
