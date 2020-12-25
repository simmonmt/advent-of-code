package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/simmonmt/aoc/2020/common/filereader"
	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func findLoopSize(subj int, pubKey int) int {
	num := 1
	for sz := 1; ; sz++ {
		num = num * subj
		num = num % 20201227
		if num == pubKey {
			return sz
		}
	}
}

func transform(subj int, loopSize int) int {
	num := 1
	for sz := 1; sz <= loopSize; sz++ {
		num = num * subj
		num = num % 20201227
	}
	return num
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	nums, err := filereader.Numbers(*input)
	if err != nil {
		log.Fatal(err)
	}

	cardPubKey := nums[0]
	doorPubKey := nums[1]

	cardLoopSize := findLoopSize(7, cardPubKey)
	logger.LogF("card loop size %v", cardLoopSize)

	doorLoopSize := findLoopSize(7, doorPubKey)
	logger.LogF("door loop size %v", doorLoopSize)

	cardEncKey := transform(doorPubKey, cardLoopSize)
	logger.LogF("card enc key %v", cardEncKey)

	doorEncKey := transform(cardPubKey, doorLoopSize)
	logger.LogF("door enc key %v", doorEncKey)

	if cardEncKey != doorEncKey {
		panic("mismatch")
	}

	fmt.Printf("A: %v\n", doorEncKey)
}
