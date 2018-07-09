package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("missing code number")
	}

	codeNum, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("failed to parse code number %v", os.Args[1])
	}

	var curCodeNum uint64 = 1
	var curCode uint64 = 20151125

	for curCodeNum <= codeNum {
		if curCodeNum > 0 && curCodeNum%1000000 == 0 {
			fmt.Printf("-- %v\n", curCodeNum)
		}

		fmt.Printf("num %v code %v\n", curCodeNum, curCode)
		curCode = (curCode * uint64(252533)) % uint64(33554393)
		curCodeNum++
	}
}
