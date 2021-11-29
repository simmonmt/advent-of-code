// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import "fmt"

type Generator struct {
	prev, factor, div int
}

func NewGenerator(prev, factor, div int) *Generator {
	return &Generator{prev, factor, div}
}

func (g *Generator) gen() int {
	next := (g.prev * g.factor) % 2147483647
	g.prev = next
	return next
}

func (g *Generator) Next() int {
	for {
		next := g.gen()
		if next%g.div == 0 {
			return next
		}
	}
}

func main() {
	// genAStarting := 65
	// genBStarting := 8921

	genAStarting := 699
	genBStarting := 124

	genAFactor := 16807
	genBFactor := 48271

	genA := NewGenerator(genAStarting, genAFactor, 4)
	genB := NewGenerator(genBStarting, genBFactor, 8)

	// for i := 0; i < 5; i++ {
	// 	fmt.Printf("%10d %10d\n", genA.Next(), genB.Next())
	// }

	numMatches := 0
	lastMatch := -1
	for i := 0; i < 5000000; i++ {
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
