package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"sleigh"
)

func calcQE(values []int) int {
	out := 1
	for _, val := range values {
		out *= val
	}
	return out
}

func main() {
	values := []int{}
	for i := 1; i < len(os.Args); i++ {
		if val, err := strconv.ParseInt(os.Args[i], 10, 32); err != nil {
			log.Fatalf("failed to parse %v: %v", os.Args[i], err)
		} else {
			values = append(values, int(val))
		}
	}

	for groupOneSize := 1; groupOneSize < len(values)-2; groupOneSize++ {
		fmt.Printf("== size %v\n", groupOneSize)
		groups := sleigh.FindWithGroupOneSize(values, groupOneSize)
		if len(groups) == 0 {
			fmt.Println("no groups found")
			continue
		}

		minQE := -1
		for _, group := range groups {
			if qe := calcQE(group); minQE == -1 || qe < minQE {
				minQE = qe
			}
		}

		fmt.Printf("found %v => %v\n", groups, minQE)
		break
	}
}
