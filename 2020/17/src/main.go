package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simmonmt/aoc/2020/17/src/board"
	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
	"github.com/simmonmt/aoc/2020/common/pos"
)

var (
	verbose  = flag.Bool("verbose", false, "verbose")
	input    = flag.String("input", "", "input file")
	fourD    = flag.Bool("fourd", false, "use four dimensions")
	numSteps = flag.Int("num_steps", 6, "number of steps")
)

func Dump(b *board.Board) {
	zmin, zmax := b.ZBounds()
	wmin, wmax := b.WBounds()

	for w := wmin; w <= wmax; w++ {
		for z := zmin; z <= zmax; z++ {
			fmt.Printf("z=%d, w=%d\n", z, w)
			b.DumpZW(z, w, os.Stdout)
			fmt.Println()
		}
	}
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

	countSet := func(b *board.Board) int {
		num := 0
		b.Visit(func(p pos.P4, v bool) {
			if v {
				num++
			}
		})
		return num
	}

	b := board.New(lines, *fourD)
	for i := 1; i <= *numSteps; i++ {
		b = b.Evolve()
		if logger.Enabled() {
			Dump(b)
		}

		logger.LogF("step=%v num=%v\n", i, countSet(b))
	}

	fmt.Printf("result: %v\n", countSet(b))
}
