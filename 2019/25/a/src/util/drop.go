package main

import (
	"flag"
	"fmt"
)

var (
	seq = flag.Int("seq", -1, "seq")
)

func main() {
	flag.Parse()

	objs := []string{
		"pointer",    // 1
		"hypercube",  // 2
		"cake",       // 4
		"tambourine", // 8
		"mouse",      // 16
		"coin",       // 32
		"mug",        // 64
		"monolith",   // 128
	}

	i := *seq
	for j := 0; j < 8; j++ {
		if (i & (1 << j)) != 0 {
			fmt.Printf("drop %s\n", objs[j])
		}
	}
}
