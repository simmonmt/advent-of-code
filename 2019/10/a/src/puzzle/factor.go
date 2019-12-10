package puzzle

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/intmath"
)

var (
	kPrimes = []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97,
	}
)

func Factor(n, d int) (int, int) {
	if intmath.Abs(n) > kPrimes[len(kPrimes)-1] || intmath.Abs(d) > kPrimes[len(kPrimes)-1] {
		panic(fmt.Sprintf("too large n=%d d=%d", n, d))
	}

	for {
		changed := false
		for _, p := range kPrimes {
			if p > intmath.Abs(n) && p > intmath.Abs(d) {
				break
			}

			if n%p == 0 && d%p == 0 {
				n /= p
				d /= p
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	return n, d
}
