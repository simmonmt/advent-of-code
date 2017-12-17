package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s num", os.Args[0])
	}

	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to parse '%s' as num", os.Args[1])
	}

	if num == 1 {
		fmt.Println(0)
		return
	}

	spiralNum := 1
	spiralLow := 2
	base := 3
	for base*base < num {
		spiralLow = base*base + 1
		base += 2
		spiralNum++
	}

	spiralHigh := base * base

	fmt.Printf("spiral low %d high %d\n", spiralLow, spiralHigh)

	x, y := 0, 0
	sideLen := base - 1
	for side := 0; side < 4; side++ {
		sideStart := spiralLow + (base-1)*side
		sideEnd := sideStart + sideLen - 1

		// Offset from sideStart that is at the same
		// coordinate as the center.
		sideZero := sideLen/2 - 1

		fmt.Printf("side %d start %d end %d zero %d\n",
			side, sideStart, sideEnd, sideZero)

		if num >= sideStart && num <= sideEnd {
			fmt.Printf("on side\n")
		} else {
			continue
		}

		switch side {
		case 0: // right side
			x = spiralNum
			y = (num - sideStart) - sideZero
			break
		case 1: // top side
			x = -((num - sideStart) - sideZero)
			y = spiralNum
			break
		case 2: // left side
			x = -spiralNum
			y = -((num - sideStart) - sideZero)
			break
		case 3: // bottom side
			x = (num - sideStart) - sideZero
			y = -spiralNum
			break
		}
	}

	fmt.Printf("x %d y %d\n", x, y)

	fmt.Printf("distance %d\n", abs(x)+abs(y))
}
