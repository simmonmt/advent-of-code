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

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Ingredient struct {
	Name   string
	Coeffs []int
}

func parseIngredient(line string) (*Ingredient, error) {
	parts := regexp.MustCompile(`[:,] `).Split(line, -1)
	if parts == nil {
		return nil, fmt.Errorf("unable to split")
	}

	ing := &Ingredient{Name: parts[0]}

	for _, char := range parts[1:] {
		charParts := strings.SplitN(char, " ", 2)

		charName := charParts[0]
		charVal, err := strconv.Atoi(charParts[1])
		if err != nil {
			return nil, fmt.Errorf("bad value for %v: %v", charName, err)
		}

		switch charName {
		case "capacity":
			fallthrough
		case "durability":
			fallthrough
		case "flavor":
			fallthrough
		case "texture":
			ing.Coeffs = append(ing.Coeffs, charVal)
			break

		case "calories":
			// Ignore
			break
		default:
			return nil, fmt.Errorf("unexpected characteristic %v", charName)
		}
	}

	return ing, nil
}

func readInput(r io.Reader) ([]*Ingredient, error) {
	ings := []*Ingredient{}

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		ing, err := parseIngredient(strings.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("%d: failed to parse ingredient: %v", lineNum, err)
		}

		ings = append(ings, ing)
	}

	return ings, nil
}

func eval(qtys []int, ings []*Ingredient) int {
	res := 1

	for coeffNum := 0; coeffNum < len(ings[0].Coeffs); coeffNum++ {
		sum := 0
		for ingNum, ing := range ings {
			sum += qtys[ingNum] * ing.Coeffs[coeffNum]
		}

		if sum < 0 {
			sum = 0
		}

		// fmt.Printf("coeffNum %d sum %d\n", coeffNum, sum)
		res *= sum
	}

	return res
}

func toCheck(qtys []int) [][]int {
	isok := func(qtys []int) bool {
		for _, q := range qtys {
			if q < 0 {
				return false
			}
		}
		return true
	}

	out := [][]int{}
	for toIncIdx := range qtys {
		for toDecIdx := range qtys {
			if toDecIdx == toIncIdx {
				continue
			}

			cand := []int{}
			for i, qty := range qtys {
				newQty := qty
				if i == toIncIdx {
					newQty = qty + 1
				} else if i == toDecIdx {
					newQty = qty - 1
				}
				cand = append(cand, newQty)
			}

			if isok(cand) {
				out = append(out, cand)
			}
		}
	}
	return out
}

func main() {
	ings, err := readInput(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	// for _, ing := range ings {
	// 	fmt.Println(*ing)
	// }

	if len(os.Args) != len(ings)+1 {
		log.Fatalf("expected %d quantities, got %d", len(ings), len(os.Args)-1)
	}

	qtys := []int{}
	for _, qtyStr := range os.Args[1:] {
		qty, err := strconv.Atoi(qtyStr)
		if err != nil {
			log.Fatalf("failed to parse quantity '%v'", qty)
		}
		qtys = append(qtys, qty)
	}

	// for i := 0; i <= 100; i++ {
	// 	tq := []int{i, 100 - i}
	// 	fmt.Printf("%v: %v\n", tq, eval(tq, ings))
	// }
	// log.Fatalf("exit")

	best := eval(qtys, ings)
	fmt.Printf("initial: %v\n", best)

	for i := 0; i < 20; i++ {
		checks := toCheck(qtys)
		fmt.Println(checks)

		newBestIdx := -1
		newBestVal := 0
		for checkIdx, check := range checks {
			checkVal := eval(check, ings)
			fmt.Printf("%v: %v\n", check, checkVal)

			if newBestIdx == -1 || checkVal > newBestVal {
				newBestIdx = checkIdx
				newBestVal = checkVal
			}
		}

		if newBestVal <= best {
			fmt.Printf("no new best; best is %v for %v\n", best, qtys)
			break
		}

		fmt.Printf("new best %v at %v: %v\n", newBestVal, newBestIdx, checks[newBestIdx])

		qtys = checks[newBestIdx]
		best = newBestVal
	}
}
