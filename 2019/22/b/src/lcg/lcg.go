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

// calculate ainv such that a*ainv mod m == 1
//
// extended euclidean calculates x and y
//    ax+by=gcd(a,b)
//
// for a given a and b. in this case, a=a and b=m from the top-level
// program. gcd is expected to be 1 since both values are prime. (note the
// variable names in reciprocalMod are .. confusing. x is gcd -- i think). the
// returned value, which comes from a in the routine, is x above in the
// euclidean formula. in ax+by=gcd(a,b) we want to find x such that ax=gcd=1
// (which also means b must be 0).

// in this case a and b are prime, so gcd(a,b) = 1

func reciprocalMod(xIn, m *big.Int) *big.Int {
	y := &big.Int{}
	y.Set(xIn)

	x := &big.Int{}
	x.Set(m)

	a := big.NewInt(0)
	b := big.NewInt(1)

	zero, one := big.NewInt(0), big.NewInt(1)

	for y.Cmp(zero) != 0 {
		tmp := &big.Int{}
		tmp.Div(x, y)
		tmp.Mul(tmp, b)
		tmp.Sub(a, tmp)

		a.Set(b)
		b.Set(tmp)

		tmp.Mod(x, y)
		x.Set(y)
		y.Set(tmp)
	}

	if x.Cmp(one) == 0 {
		r := &big.Int{}
		r.Mod(a, m)
		return r
	} else {
		return nil
	}
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

	one := big.NewInt(1)

	aMinus1 := &big.Int{}
	aMinus1.Sub(a, one)

	numMod := &big.Int{}
	numMod.Set(aMinus1)
	numMod.Mul(numMod, m)

	num := computeANX(a, one, n, numMod)
	num.Sub(num, one)
	num.Mod(num, numMod)

	frac := &big.Int{}
	frac.Div(num, aMinus1)
	frac.Mul(frac, c)

	val := computeANX(a, x, n, m)
	val.Add(val, frac)
	val.Mod(val, m)
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

	bigN := int64(101741582076661)
	fmt.Printf("ffwd %v = %v\n", bigN, computeForward(seed, bigN, a, c, mod))

	ainv := reciprocalMod(a, mod)

	{
		tmp := &big.Int{}
		tmp.Mul(a, ainv)
		tmp.Mod(tmp, mod)

		fmt.Printf("a %v * ainv %v = %v\n", a, ainv, tmp)
		if tmp.Cmp(big.NewInt(1)) != 0 {
			log.Fatalf("inverse verification failed")
		}
	}

	// Tried fast forwarding with
	// n=101741582076661, which gave 1608694956433: wrong
	//
	// because n=y tells you where the seed value will end up after that
	// many repetitions.
	//
	// we want to find the y that gives us the result 2020
}
