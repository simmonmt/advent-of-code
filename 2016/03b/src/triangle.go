package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	inputPattern = regexp.MustCompile(`\w+`)
)

func readInput(r io.Reader) ([][3]int, error) {
	out := [][3]int{}

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		matches := inputPattern.FindAllString(line, -1)
		if matches != nil && len(matches) != 3 {
			return nil, fmt.Errorf("%d: expected 3 words, found %v", lineNum, matches)
		}

		dims := [3]int{}
		for i, match := range matches {
			match = strings.TrimSpace(match)

			dim, err := strconv.ParseUint(match, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("failed to parse dim %v: %v", match, err)
			}
			dims[i] = int(dim)
		}

		out = append(out, dims)
	}

	return out, nil
}

func rotateInput(input [][3]int) ([][3]int, error) {
	if len(input)%3 != 0 {
		return nil, fmt.Errorf("input has non-multiple-of-3 rows; found %v", len(input))
	}

	out := [][3]int{}
	for row := 0; row < len(input); row += 3 {
		for col := 0; col < len(input[0]); col++ {
			dim := [3]int{}
			for i := 0; i < 3; i++ {
				dim[i] = input[row+i][col]
			}
			out = append(out, dim)
		}
	}

	return out, nil
}

func main() {
	var triangles [][3]int
	var err error

	triangles, err = readInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	triangles, err = rotateInput(triangles)
	if err != nil {
		log.Fatalf(err.Error())
	}

	numPossible := 0
	for _, tri := range triangles {
		possible := tri[0]+tri[1] > tri[2] && tri[0]+tri[2] > tri[1] && tri[1]+tri[2] > tri[0]
		if possible {
			numPossible++
		}
	}

	fmt.Printf("num possible = %v\n", numPossible)
}
