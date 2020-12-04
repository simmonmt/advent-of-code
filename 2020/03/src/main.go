package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simmonmt/aoc/2020/common/logger"
	"github.com/simmonmt/aoc/2020/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Table struct {
	arr  []bool
	w, h int
}

func (t *Table) Height() int {
	return t.h
}

func (t *Table) Get(p pos.P2) bool {
	x := p.X % t.w
	off := p.Y*t.w + x
	return t.arr[off]
}

func (t *Table) Dump() {
	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			if t.Get(pos.P2{x, y}) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func newTable(strs []string) *Table {
	h, w := len(strs), len(strs[0])
	arr := make([]bool, h*w)

	for r, str := range strs {
		for c, v := range str {
			arr[r*w+c] = v == '#'
		}
	}

	return &Table{
		arr: arr,
		w:   w,
		h:   h,
	}
}

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

func trySlope(tbl *Table, slope pos.P2) int {
	numTrees := 0
	p := pos.P2{0, 0}
	for p.Y < tbl.Height() {
		hasTree := tbl.Get(p)
		logger.LogF("%s: %v", p, hasTree)
		if hasTree {
			numTrees++
		}
		p.Add(slope)
	}
	return numTrees
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

	tbl := newTable(lines)
	if logger.Enabled() {
		tbl.Dump()
	}

	slopes := []pos.P2{
		pos.P2{1, 1},
		pos.P2{3, 1},
		pos.P2{5, 1},
		pos.P2{7, 1},
		pos.P2{1, 2},
	}

	result := 1
	for _, slope := range slopes {
		numTrees := trySlope(tbl, slope)
		result *= numTrees
		fmt.Printf("slope %s: num trees: %v\n", slope, numTrees)
	}
	fmt.Printf("result: %v\n", result)
}
