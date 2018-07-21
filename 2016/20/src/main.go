package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"extent"
)

var (
	rangeMax = flag.Int("range_max", -1, "range max")
)

func readInput(r io.Reader) (extent.Extents, error) {
	exts := extent.Extents{}

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if ext, err := extent.Parse(line); err != nil {
			return nil, fmt.Errorf("%d: %v", lineNum, err)
		} else {
			exts = append(exts, ext)
		}
	}

	return exts, nil
}

func main() {
	flag.Parse()

	if *rangeMax == -1 {
		log.Fatal("--range_max is required")
	}

	exts, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err.Error())
	}

	sort.Sort(exts)
	exts = exts.Merge()

	for i := 0; i < 5 && i < len(exts); i++ {
		fmt.Printf("  %s\n", exts[i])
	}

	var cur, numAllowed uint64
	for _, e := range exts {
		if cur < e.Start {
			numAllowed += e.Start - cur
		}
		cur = e.End + 1
	}
	fmt.Printf("allowed = %d\n", numAllowed)
}
