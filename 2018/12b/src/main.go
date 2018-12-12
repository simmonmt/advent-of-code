// This solution is a bit of a mess. When they ask for 50B iterations they don't
// actually mean do it 50B times. They mean find a pattern and extrapolate. That
// took longer than I'd like to admit to sink in.
//
// There were two ways to make my 12a solution faster:
//   1. Get rid of some of the loops. I was storing the patterns as arrays of
//      bools, looping through them each time. If I stored them as uints, I
//      could make some sort of map, allowing for constant time lookup of an
//      applicable rule.
//   2. Make the padding dynamic.
//
// I needed to do the second, but wildly overestimated the time savings I'd get
// from the first, and underestimated the complication. Converting the pattern
// from an array of bools to a bitmask cured me of that misconception, so I
// moved onto the second (making the padding dynamic). Dynamic padding was
// sufficient to see the pattern, but I didn't take the time to rip out the
// partial bitmaskification. So that's why pattern is stored as a bitmask, then
// converted back to an array and used as an array.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	numGens = flag.Int("num_gens", -1, "num generations")
	pad     = flag.Int("pad", 10, "pad")
)

type Rule struct {
	Pattern uint
	Result  bool
}

func (r *Rule) Applies(arr []bool, cur int) bool {
	patArr := []bool{
		(r.Pattern & 0x10) != 0,
		(r.Pattern & 0x08) != 0,
		(r.Pattern & 0x04) != 0,
		(r.Pattern & 0x02) != 0,
		(r.Pattern & 0x01) != 0,
	}

	logger.LogF("trying %s at %d of %s", arrToStr(patArr), cur, arrToStr(arr))

	leftOff := len(patArr) / 2
	for i := 0; i < len(patArr); i++ {
		off := cur - leftOff + i

		inArr := false
		if off >= 0 && off < len(arr) {
			inArr = arr[off]
		}

		logger.LogF("i = %v off = %v : %v vs %v", i, off, inArr, patArr[i])

		if inArr != patArr[i] {
			logger.LogLn("match fail")
			return false
		}
	}

	return true
}

func readLines() ([]string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func readInput() ([]bool, []Rule, error) {
	lines, err := readLines()
	if err != nil {
		return nil, nil, err
	}

	initialStr := strings.Split(lines[0], " ")[2]
	initial := []bool{}
	for _, c := range initialStr {
		initial = append(initial, c == '#')
	}

	rules := []Rule{}
	for i := 2; i < len(lines); i++ {
		parts := strings.Split(lines[i], " ")
		patStr := parts[0]
		if len(patStr) != 5 {
			panic("not 5")
		}

		var patVal uint
		for _, c := range patStr {
			patVal = patVal << 1
			if c == '#' {
				patVal |= 1
			}
		}

		result := parts[2] == "#"

		rules = append(rules, Rule{Pattern: patVal, Result: result})
	}

	return initial, rules, nil
}

func dumpGen(gen int, a []bool, zeroOff int) {
	// fmt.Printf("       ")
	// for i := range a {
	// 	if i == zeroOff {
	// 		fmt.Printf("Z")
	// 	} else {
	// 		val := (i - zeroOff) / 10 % 10
	// 		if val < 0 {
	// 			val = -val
	// 		}
	// 		fmt.Printf("%d", val)
	// 	}
	// }
	// fmt.Println()

	// fmt.Printf("       ")
	// for i := range a {
	// 	if i == zeroOff {
	// 		fmt.Printf("Z")
	// 	} else {
	// 		val := (i - zeroOff) % 10
	// 		if val < 0 {
	// 			val = -val
	// 		}
	// 		fmt.Printf("%d", val)
	// 	}
	// }
	// fmt.Println()
	fmt.Printf("%05d: %v\n", gen, arrToStr(a))
}

func patToStr(pat uint) string {
	patOrig := pat
	out := ""
	for i := 0; i < 5; i++ {
		if (pat & 1) == 1 {
			out = "#" + out
		} else {
			out = "." + out
		}
		pat >>= 1
	}
	return fmt.Sprintf("in %x out %v", patOrig, out)
}

func arrToStr(a []bool) string {
	out := ""
	for _, b := range a {
		if b {
			out += "#"
		} else {
			out += "."
		}
	}
	return out
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *numGens == -1 {
		log.Fatalf("--num_gens required")
	}

	initial, rules, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	cur := make([]bool, *pad+len(initial)+*pad)
	copy(cur[*pad:], initial)
	zeroOff := *pad

	lastSum := 0
	dumpGen(0, cur, zeroOff)
	for gen := 1; gen <= *numGens; gen++ {
		next := make([]bool, len(cur))
		for pos := range cur {
			found := false
			hasPlant := false
			for _, r := range rules {
				if r.Applies(cur, pos) {
					logger.LogF("matched rule %v at %v\n", patToStr(r.Pattern), pos)

					if found {
						panic("refind")
					}
					hasPlant = r.Result
					found = true
				}
			}

			next[pos] = hasPlant
		}

		dumpGen(gen, next, zeroOff)

		sum := 0
		for i := range next {
			if !next[i] {
				continue
			}

			//fmt.Printf("pot at %v\n", i-zeroOff)
			sum += (i - zeroOff)
		}

		fmt.Printf("sum %v last %v delta %v\n", sum, lastSum, sum-lastSum)
		lastSum = sum

		numLeadingFalse := 0
		for i := 0; i < len(next); i++ {
			if !next[i] {
				numLeadingFalse++
			} else {
				break
			}
		}

		if numLeadingFalse < *pad {
			newOff := *pad - numLeadingFalse
			newNext := make([]bool, len(next)+newOff)
			copy(newNext[newOff:], next)
			zeroOff += newOff
			next = newNext
		} else if numLeadingFalse > *pad {
			shrinkNum := numLeadingFalse - *pad
			next = next[shrinkNum:]
			zeroOff -= shrinkNum
		}

		numTrailingFalse := 0
		for i := len(next) - 1; i >= 0; i-- {
			if !next[i] {
				numTrailingFalse++
			} else {
				break
			}
		}
		if numTrailingFalse < *pad {
			next = append(next, make([]bool, *pad-numTrailingFalse)...)
		}

		cur = next
	}

}
