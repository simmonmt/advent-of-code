package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"tiles"
)

var (
	numRows = flag.Int("num_rows", -1, "number of rows to make")
)

func readInput(r io.Reader) (tiles.Row, error) {
	reader := bufio.NewReader(r)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	row, err := tiles.MakeRow(line)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func main() {
	flag.Parse()

	if *numRows == -1 {
		log.Fatal("--num_rows is required")
	}

	first, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	r := first
	numSafe := 0
	for i := 0; i < *numRows; i++ {
		numSafe += r.NumSafe()
		r = r.Next()
	}

	fmt.Println(numSafe)
}
