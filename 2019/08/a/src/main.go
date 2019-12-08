package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simmonmt/aoc/2019/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
	width   = flag.Int("width", -1, "width")
	height  = flag.Int("height", -1, "height")
)

func readInput(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", fmt.Errorf("read failed: %v", scanner.Err())
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}
	if *width < 0 {
		log.Fatalf("--width is required")
	}
	if *height < 0 {
		log.Fatalf("--height is required")
	}

	line, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	layerSize := *width * *height
	numLayers := (len(line) + layerSize - 1) / layerSize
	layers := make([][]rune, numLayers)
	for i := 0; i < numLayers; i++ {
		layers[i] = make([]rune, layerSize)
	}

	for i, r := range line {
		layers[i/layerSize][i%layerSize] = r
	}

	//fmt.Println(layers)

	fewestIdx := -1
	fewestNum := 0
	for i, l := range layers {
		num := 0
		for _, r := range l {
			if r == '0' {
				num++
			}
		}

		if fewestIdx == -1 || num < fewestNum {
			fewestIdx = i
			fewestNum = num
		}
	}

	fmt.Printf("layer with fewest zeros %d num %d\n", fewestIdx+1, fewestNum)

	numOne, numTwo := 0, 0
	for _, r := range layers[fewestIdx] {
		if r == '1' {
			numOne++
		} else if r == '2' {
			numTwo++
		}
	}

	fmt.Printf("numOne %d, numTwo %d, result %d\n",
		numOne, numTwo, numOne*numTwo)
}
