package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	row = flag.Int("row", -1, "row number")
	col = flag.Int("col", -1, "col number")
)

// calculate the number of the code in the first column for each row

func main() {
	flag.Parse()

	if *row == -1 || *col == -1 {
		log.Fatalf("--row and --col are required")
	}

	curCodeNum := 1
	curRow := 1

	for curRow < *row {
		curRow++
		curCodeNum += (curRow - 1)
	}
	fmt.Printf("row %v codeNum %v\n", curRow, curCodeNum)

	curCol := 1
	colInc := curRow + 1

	for curCol < *col {
		curCodeNum += colInc
		colInc++
		curCol++
	}

	fmt.Printf("row %v col %v codeNum %v\n", curRow, curCol, curCodeNum)
}
