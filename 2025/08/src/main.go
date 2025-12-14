package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"sort"

	"github.com/simmonmt/aoc/2025/common/filereader"
	"github.com/simmonmt/aoc/2025/common/logger"
	"github.com/simmonmt/aoc/2025/common/pos"
)

var (
	verbose    = flag.Bool("verbose", false, "verbose")
	input      = flag.String("input", "", "input file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

type Input struct {
	Ps []pos.P3
}

func parseInput(lines []string) (*Input, error) {
	ps := []pos.P3{}
	for i, line := range lines {
		p, err := pos.P3FromString(line)
		if err != nil {
			return nil, fmt.Errorf("%d: bad pos: %v", i+1, err)
		}
		ps = append(ps, p)
	}
	return &Input{Ps: ps}, nil
}

type Pair struct {
	A, B pos.P3
}

func Dist(a, b pos.P3) float64 {
	return math.Sqrt(
		math.Pow(float64(a.X)-float64(b.X), 2) +
			math.Pow(float64(a.Y)-float64(b.Y), 2) +
			math.Pow(float64(a.Z)-float64(b.Z), 2))
}

func solveA(input *Input, num int) int {
	pairs := []Pair{}
	for i, a := range input.Ps {
		for j := i + 1; j < len(input.Ps); j++ {
			b := input.Ps[j]
			pairs = append(pairs, Pair{a, b})
		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return Dist(pairs[i].A, pairs[i].B) < Dist(pairs[j].A, pairs[j].B)
	})

	nextCircuit := 1
	circuits := map[pos.P3]int{}

	for i := 0; i < num; i++ {
		a, b := pairs[i].A, pairs[i].B
		//fmt.Println(a, b)

		ac, fa := circuits[a]
		bc, fb := circuits[b]

		if fa && fb {
			// merge bc into ac
			for p, n := range circuits {
				if n == bc {
					circuits[p] = ac
				}
			}

		} else if fa {
			circuits[b] = ac
		} else if fb {
			circuits[a] = bc
		} else {
			circuits[a], circuits[b] = nextCircuit, nextCircuit
			nextCircuit++
		}

	}

	ps := map[int][]pos.P3{}
	for p, n := range circuits {
		ps[n] = append(ps[n], p)
	}

	// for i := range nextCircuit {
	// 	if l, found := ps[i]; found {
	// 		fmt.Printf("size %d: %v\n", len(l), l)
	// 	}
	// }

	szs := []int{}
	for _, l := range ps {
		szs = append(szs, len(l))
	}
	sort.Ints(szs)

	tot := 1
	for i := len(szs) - 1; i >= 0 && i > len(szs)-1-3; i-- {
		tot *= szs[i]
	}
	return tot

}

func solveB(input *Input) int {
	pairs := []Pair{}
	for i, a := range input.Ps {
		for j := i + 1; j < len(input.Ps); j++ {
			b := input.Ps[j]
			pairs = append(pairs, Pair{a, b})
		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return Dist(pairs[i].A, pairs[i].B) < Dist(pairs[j].A, pairs[j].B)
	})

	nextCircuit := 1
	numCircuits := 0
	circuits := map[pos.P3]int{}

	var last Pair
	for _, pair := range pairs {
		if len(circuits) == len(input.Ps) && numCircuits == 1 {
			break
		}

		a, b := pair.A, pair.B
		//fmt.Println(a, b)

		ac, fa := circuits[a]
		bc, fb := circuits[b]

		if fa && fb {
			if ac != bc {
				// merge bc into ac
				//fmt.Println("merging", ac, bc)
				for p, n := range circuits {
					if n == bc {
						circuits[p] = ac
					}
				}
				numCircuits--
			}

		} else if fa {
			circuits[b] = ac
		} else if fb {
			circuits[a] = bc
		} else {
			circuits[a], circuits[b] = nextCircuit, nextCircuit
			nextCircuit++
			numCircuits++
		}

		//fmt.Println(pair, numCircuits)
		last = pair
	}

	//fmt.Println(last)
	return last.A.X * last.B.X

	// ps := map[int][]pos.P3{}
	// for p, n := range circuits {
	// 	ps[n] = append(ps[n], p)
	// }

	// for i := range nextCircuit {
	// 	if l, found := ps[i]; found {
	// 		fmt.Printf("size %d: %v\n", len(l), l)
	// 	}
	// }

	// szs := []int{}
	// for _, l := range ps {
	// 	szs = append(szs, len(l))
	// }
	// fmt.Println(szs)
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

	fmt.Println("A", solveA(input, 1000))
	fmt.Println("B", solveB(input))
}
