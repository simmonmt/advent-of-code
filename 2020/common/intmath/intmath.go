package intmath

import (
	"fmt"
	"strconv"
)

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func Uint64Max(a, b uint64) uint64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func AtoiOrDie(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %v: %v", s, err))
	}
	return val
}

var (
	kPrimes = []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97,
	}
)

func GCD(vs ...int) int {
	d := 1

	for _, p := range kPrimes {
		for {
			all := true
			for _, v := range vs {
				if v < p || v%p != 0 {
					all = false
				}
			}

			if !all {
				break
			}

			d *= p
			for i := range vs {
				vs[i] /= p
			}
		}
	}

	for _, v := range vs {
		if v > kPrimes[len(kPrimes)-1] {
			panic(fmt.Sprintf("%d too big", v))
		}
	}

	return d
}
