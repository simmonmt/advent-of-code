package main

import (
	"flag"
	"fmt"
	"log"
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

func extGCD(a, b int) (s, t int) {
	x, y, u, v := 0, 1, 1, 0
	for a != 0 {
		q, r := b/a, b%a
		m, n := x-u*q, y-v*q
		b, a, x, y, u, v = a, r, u, v, m, n
	}

	return x, y
}

func computeANSlow(a, mod big.Int, n int64) big.Int {
	val := big.NewInt(0)
	val.Exp(&a, big.NewInt(n), &mod)
	return *val
}

// computes a^n*x mod m
func computeANX(a, x *big.Int, n int64, m *big.Int) *big.Int {
	// from
	// https://en.wikipedia.org/w/index.php?title=Modular_exponentiation&action=edit&section=3

	base := &big.Int{}
	base.Set(a)
	base.Mod(base, m)

	result := &big.Int{}
	result.Set(x)
	for n > 0 {
		if (n & 1) != 0 {
			result.Mul(result, base)
			result.Mod(result, m)
		}
		n >>= 1
		base.Mul(base, base)
		base.Mod(base, m)
	}
	return result
}

func computeForward(x *big.Int, n int64, a, c, m *big.Int) *big.Int {
	// fast forwarding and reversing described here
	// https://www.nayuki.io/page/fast-skipping-in-a-linear-congruential-generator

	// We're computing
	//
	//                    (a^n-1) mod (a-1)m
	//  ( (a^n*x mod m) + ------------------ b ) mod m
	//                            a-1

	//	fmt.Printf("x %v n %v a %v c %v m %v\n", x, n, a, c, m)

	aMinus1 := &big.Int{}
	aMinus1.Sub(a, big.NewInt(1))

	num := &big.Int{}
	num.Exp(a, big.NewInt(n), nil)
	num.Sub(num, big.NewInt(1))

	frac := &big.Int{}
	frac.Div(num, aMinus1)
	frac.Mul(frac, c)

	//val := computeANX(a, x, n, m)
	val := &big.Int{}
	val.Exp(a, big.NewInt(n), nil)
	val.Mul(val, x)
	val.Add(val, frac)
	val.Mod(val, m)

	// //fmt.Printf("a-1 %v\n", aMinus1)

	// // num[erator]M = (a-1)m
	// numM := &big.Int{}
	// numM.Set(aMinus1)
	// numM.Mul(numM, m)

	// // (a^n-1) mod (a-1)m
	// //
	// // Which we compute as ((a^n mod (a-1)m) - 1) mod (a-1)m
	// num := computeANX(a, big.NewInt(1), n, numM)
	// num.Sub(num, big.NewInt(1))
	// num.Mod(num, numM)

	// num = computeANX(a, big.NewInt(1), n, m)
	// num.Sub(num, big.NewInt(1))
	// num.Mul(num, c)

	// frac := &big.Int{}
	// frac.Div(num, aMinus1)

	// // a^n*x mod m + the fraction
	// val := computeANX(a, x, n, m)
	// val.Add(val, frac)
	// val.Mod(val, m)

	//fmt.Printf("x %v n %v a %v c %v m %v\n", x, n, a, c, m)
	return val
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

	for i := int64(2); i < 10; i++ {
		slow := computeANSlow(*a, *mod, i)
		fast := computeANX(a, big.NewInt(1), i, mod)
		if slow.Cmp(fast) != 0 {
			log.Fatal("slow vs fast verify failed at %d: %v vs %v",
				i, slow, fast)
		}
	}
	fmt.Println("computeANFast verify succeeded")

	forwardMatch := true
	val := &big.Int{}
	val.Set(seed)
	for i := 0; i < 10; i++ {
		//fmt.Println("---")

		ffwd := computeForward(seed, int64(i), a, c, mod)

		fmt.Printf("%d: %-20v", i, val)
		if ffwd.Cmp(val) == 0 {
			fmt.Println("match")
		} else {
			fmt.Printf("mismatch %v\n", ffwd)
			forwardMatch = false
		}

		val.Mul(val, a)
		val.Add(val, c)
		val.Mod(val, mod)
	}
	if !forwardMatch {
		log.Fatalf("forward test failed")
	}

	s, t := extGCD(int(a.Int64()), int(mod.Int64()))

	fmt.Printf("s %v t %v\n", s, t)

	inv := &big.Int{}
	inv.Mul(a, big.NewInt(int64(t)))
	inv.Mod(a, mod)

	fmt.Printf("as = %v\n", inv)

	// Tried fast forwarding with
	// n=101741582076661, which gave 1608694956433: wrong
	//
	// because n=y tells you where the seed value will end up after that
	// many repetitions.
	//
	// we want to find the y that gives us the result 2020
}
