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
	Name     string
	Coeffs   []int
	Calories int
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
			ing.Calories = charVal
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

func countCals(qtys []int, ings []*Ingredient) int {
	cals := 0
	for ingNum, ing := range ings {
		cals += qtys[ingNum] * ing.Calories
	}
	return cals
}

func eval(qtys []int, ings []*Ingredient) (score, cals int) {
	cals = countCals(qtys, ings)

	score = 1
	for coeffNum := 0; coeffNum < len(ings[0].Coeffs); coeffNum++ {
		sum := 0
		for ingNum, ing := range ings {
			sum += qtys[ingNum] * ing.Coeffs[coeffNum]
		}

		if sum < 0 {
			sum = 0
		}

		//fmt.Printf("coeffNum %d sum %d\n", coeffNum, sum)
		score *= sum
	}

	return
}

type QuantityIterator struct {
	Num, Max int
	Last     []int
}

func NewQuantityIterator(num, max int) *QuantityIterator {
	return &QuantityIterator{Num: num, Max: max}
}

func (q *QuantityIterator) isValid() bool {
	sum := 0
	for _, val := range q.Last {
		sum += val
	}
	return sum == q.Max
}

func (q *QuantityIterator) Next() []int {
	for {
		if q.Last == nil {
			q.Last = make([]int, q.Num)
		} else {
			for i := range q.Last {
				if i == 0 {
					q.Last[i]++
				}
				if q.Last[i] > q.Max {
					if i == len(q.Last)-1 {
						return nil
					}
					q.Last[i] = 0
					q.Last[i+1]++
				}
			}
		}

		if q.isValid() {
			return q.Last
		}
	}
}

func main() {
	ings, err := readInput(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	bestScore := -1

	qtyIter := NewQuantityIterator(len(ings), 100)
	for {
		qtys := qtyIter.Next()
		if qtys == nil {
			break
		}

		score, cals := eval(qtys, ings)
		if cals != 500 {
			continue
		}

		if score > bestScore {
			bestScore = score
		}
	}

	fmt.Println(bestScore)
}
