package main

import (
	"fmt"
	"knot"
	"log"
	"os"
	"strconv"
)

type Coord struct {
	x, y int
}

func findStart(grid [][]bool) *Coord {
	for y, row := range grid {
		for x, val := range row {
			if val {
				return &Coord{x, y}
			}
		}
	}

	return nil
}

func gridVal(grid [][]bool, pos Coord) bool {
	if pos.x < 0 || pos.y < 0 {
		return false
	}
	if pos.y >= len(grid) || pos.x >= len(grid[0]) {
		return false
	}

	return grid[pos.y][pos.x]
}

func findRegion(grid [][]bool) bool {
	start := findStart(grid)
	if start == nil {
		return false
	}

	toExamine := []Coord{*start}

	for len(toExamine) > 0 {
		pos := toExamine[0]
		toExamine = toExamine[1:len(toExamine)]
		grid[pos.y][pos.x] = false

		cand := Coord{pos.x - 1, pos.y}
		if gridVal(grid, cand) {
			toExamine = append(toExamine, cand)
		}
		cand = Coord{pos.x + 1, pos.y}
		if gridVal(grid, cand) {
			toExamine = append(toExamine, cand)
		}
		cand = Coord{pos.x, pos.y - 1}
		if gridVal(grid, cand) {
			toExamine = append(toExamine, cand)
		}
		cand = Coord{pos.x, pos.y + 1}
		if gridVal(grid, cand) {
			toExamine = append(toExamine, cand)
		}
	}

	return true

}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v key", os.Args[0])
	}
	key := os.Args[1]

	grid := make([][]bool, 128)
	for i := range grid {
		row := make([]bool, 128)

		rowHashKey := fmt.Sprintf("%s-%d", key, i)
		rowHash := knot.Hash(rowHashKey)

		// rowHash = "a0c2017" + rowHash[7:]

		for i, r := range rowHash {
			num, err := strconv.ParseUint(string(r), 16, 64)
			if err != nil {
				log.Fatalf("failed to parse %v\n", r)
			}

			// fmt.Printf("row[i*4]=(num&8) row[%d]=%d&8=%d=%v\n",
			// 	i*4, num, (num & 0x8), (num&0x8) != 0)

			row[i*4] = (num & 0x8) != 0
			row[i*4+1] = (num & 0x4) != 0
			row[i*4+2] = (num & 0x2) != 0
			row[i*4+3] = (num & 0x1) != 0
		}

		grid[i] = row
	}

	for i := 0; i < 8; i++ {
		row := grid[i]
		for j := 0; j < 8; j++ {
			if row[j] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println("")
	}

	numRegions := 0
	for findRegion(grid) {
		numRegions++
	}

	fmt.Println(numRegions)
}
