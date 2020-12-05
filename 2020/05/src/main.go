package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

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

func doSearch(max int, dirs []bool) int {
	lo := 0
	hi := max

	logger.LogF("start with lo %v hi %v", lo, hi)
	for _, dir := range dirs {
		if dir {
			// go high
			lo = lo + (hi-lo)/2 + 1
		} else {
			// go low
			hi = lo + (hi-lo)/2
		}

		logger.LogF("after %v, now lo %v hi %v", dir, lo, hi)
	}

	if lo != hi {
		panic("mismatch")
	}
	logger.LogF("done; returning %v", lo)
	return lo
}

func decodeSeat(pass string) (row, col int) {
	rowDirs := []bool{}
	for _, dir := range ([]byte(pass))[0:7] {
		rowDirs = append(rowDirs, dir == 'B')
	}
	row = doSearch(127, rowDirs)

	colDirs := []bool{}
	for _, dir := range ([]byte(pass))[7:10] {
		colDirs = append(colDirs, dir == 'R')
	}
	col = doSearch(7, colDirs)

	return
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

	foundIDs := map[int]bool{}

	maxSeatID := -1
	for _, line := range lines {
		row, col := decodeSeat(line)
		seatID := row*8 + col
		foundIDs[seatID] = true
		if seatID > maxSeatID {
			maxSeatID = seatID
		}

		logger.LogF("%v: row %d col %d, seat ID %v", line, row, col, seatID)
	}

	fmt.Printf("max seat ID: %v\n", maxSeatID)

	for id := range foundIDs {
		_, foundPlus1 := foundIDs[id+1]
		_, foundPlus2 := foundIDs[id+2]

		if !foundPlus1 && foundPlus2 {
			fmt.Printf("hole at %v\n", id+1)
		}
	}
}
