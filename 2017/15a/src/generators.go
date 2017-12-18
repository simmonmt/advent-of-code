package main

import "fmt"

type Generator struct {
	prev, factor int
}

func NewGenerator(prev, factor int) *Generator {
	return &Generator{prev, factor}
}

func (g *Generator) Next() int {
	next := (g.prev * g.factor) % 2147483647
	g.prev = next
	return next
}

func main() {
	genAFactor := 16807
	genBFactor := 48271

	genA := NewGenerator(699, genAFactor)
	genB := NewGenerator(124, genBFactor)

	numMatches := 0
	lastMatch := -1
	for i := 0; i < 40000000; i++ {
		genAVal := genA.Next() & 0xffff
		genBVal := genB.Next() & 0xffff

		if genAVal != genBVal {
			continue
		}

		numMatches++

		fmt.Printf("%10d ", i)
		if lastMatch != -1 {
			fmt.Printf("%d", i-lastMatch)
		}
		lastMatch = i
		fmt.Println()
	}

	fmt.Printf("matches: %v\n", numMatches)
}
