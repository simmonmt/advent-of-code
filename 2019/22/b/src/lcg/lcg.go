package main

import (
	"flag"
	"fmt"
	"math/big"
)

var (
	aFlag = flag.Int64("a", -100444802940351, "a value")

	// Negative of the cut value from the shortened command list
	cFlag   = flag.Int64("c", 115490606888493, "c value")
	modFlag = flag.Int64("mod", 119315717514047, "mod value")
	seed    = flag.Int64("seed", 2020, "seed")
)

// with params
//   a = -100444802940351
//   c = 115490606888493
//   mod = 119315717514047
//   seed = 2020
//
// we should have these values:
// 0: 2020
// 99990: 119171824157945
// 99991: 2490423324226
// 99992: 64058352024721
// 99993: 63295764867845
// 99994: 73553578208733
// 99995: 102325460379295
// 99996: 63212872553673
// 99997: 59490684611703
// 99998: 49742805710047
// 99999: 87378128651512

func computeAN(a *big.Int, n int64) *big.Int {
	// This is horribly slow for large n, need to use
	// https://en.wikipedia.org/wiki/Modular_exponentiation instead
	an := big.NewInt(0)
	an.Exp(a, big.NewInt(n), nil)
	return an
}

func main() {
	mod := big.NewInt(*modFlag)

	a := big.NewInt(*aFlag)
	for a.Cmp(big.NewInt(0)) < 0 {
		a.Add(a, mod)
	}

	c := big.NewInt(*cFlag)
	for c.Cmp(big.NewInt(0)) < 0 {
		c.Add(c, mod)
	}

	seed := big.NewInt(*seed)

	val := big.NewInt(0)
	*val = *seed
	for i := 0; i < 10; i++ {
		fmt.Printf("%d: %v\n", i, val)

		val.Mul(val, a)
		val.Add(val, c)
		val.Mod(val, mod)
	}

	x := big.NewInt(0)
	*x = *seed

	n := int64(8)

	// fast forwarding and reversing described here
	// https://www.nayuki.io/page/fast-skipping-in-a-linear-congruential-generator

	an := computeAN(a, n)

	num, numMod, den, frac := big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)

	numMod.Sub(a, big.NewInt(1))
	numMod.Mul(numMod, mod)

	num.Sub(an, big.NewInt(1))
	num.Mod(num, numMod)

	den.Sub(a, big.NewInt(1))
	frac.Div(num, den)
	frac.Mul(frac, c)

	result := big.NewInt(0)
	result.Mul(an, seed)
	result.Mod(result, mod)
	result.Add(result, frac)
	result.Mod(result, mod)

	// Tried fast forwarding with
	// n=101741582076661, which gave 1608694956433: wrong
	//
	// because n=y tells you where the seed value will end up after that
	// many repetitions.
	//
	// we want to find the y that gives us the result 2020

	fmt.Printf("%v\n", result)
}
