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
	Pattern []bool
	Center  int
	Result  bool
}

func (r *Rule) Applies(arr []bool, cur int) bool {
	logger.LogF("trying %s at %d of %s\n", arrToStr(r.Pattern), cur, arrToStr(arr))
	leftOff := len(r.Pattern) / 2
	for i := 0; i < len(r.Pattern); i++ {
		off := cur - leftOff + i

		inArr := false
		if off >= 0 && off < len(arr) {
			inArr = arr[off]
		}

		logger.LogF("i = %v off = %v : %v vs %v\n", i, off, inArr, r.Pattern[i])

		if inArr != r.Pattern[i] {
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
		pattern := []bool{}
		for _, c := range patStr {
			pattern = append(pattern, c == '#')
		}

		result := parts[2] == "#"

		if len(pattern)%2 != 1 {
			panic("even pattern")
		}
		center := len(pattern)/2 + 1

		rules = append(rules, Rule{Pattern: pattern, Center: center, Result: result})
	}

	return initial, rules, nil
}

func dumpGen(gen int, a []bool) {
	fmt.Printf("%03d: %v\n", gen, arrToStr(a))
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

	dumpGen(0, cur)
	for gen := 1; gen <= *numGens; gen++ {
		next := make([]bool, len(cur))
		for pos := range cur {
			found := false
			hasPlant := false
			for _, r := range rules {
				if r.Applies(cur, pos) {
					logger.LogF("matched rule %v at %v\n", arrToStr(r.Pattern), pos)

					if found {
						panic("refind")
					}
					hasPlant = r.Result
					found = true
				}
			}

			next[pos] = hasPlant
		}

		dumpGen(gen, next)
		cur = next
	}

	sum := 0
	for i := range cur {
		if !cur[i] {
			continue
		}

		sum += i - zeroOff
	}

	fmt.Println(sum)
}
