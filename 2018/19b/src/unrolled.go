// This is the Go version of the assembly program in ../input.txt

package main

import (
	"flag"
	"fmt"
)

var (
	runlong = flag.Bool("runlong", false, "")
)

func main() {
	flag.Parse()

	//reg := [6]int{} //0, 0, 1, 0, 987, 151}

	var c int
	if *runlong {
		c = 10551387
		//d = 10550400
	} else {
		c = 987
		//d = 151
	}

	out := 0
	for a := 1; a <= c; a++ {
		for b := 1; b <= c; b++ {
			if a*b == c {
				out += a
			}
			// This conditional added after analysis of
			// transliterated assembly. Without this, the
			// loop is way too expensive to run with
			// --runlong
			if a*b > c {
				break
			}
		}
	}

	fmt.Println(out)
}
