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
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	pattern = regexp.MustCompile(`^([a-z]+) \(([0-9]+)\)(?: -> (.*))?$`)
)

type Tree struct {
	veryBottom string
	elems      map[string][]string
	weights    map[string]int
}

func NewTree() *Tree {
	return &Tree{
		elems:   map[string][]string{},
		weights: map[string]int{},
	}
}

func (t *Tree) Insert(bot, top string, weight int) {
	if t.veryBottom == "" || t.veryBottom == top {
		t.veryBottom = bot
	}

	if top != "" {
		if _, found := t.elems[bot]; found {
			t.elems[bot] = append(t.elems[bot], top)
		} else {
			t.elems[bot] = []string{top}
		}
	}

	t.weights[bot] = weight
}

func (t *Tree) Dump() {
	for _, bot := range t.AllBottoms() {
		fmt.Printf("%v (%d)", bot, t.weights[bot])
		if tops, found := t.elems[bot]; found {
			fmt.Printf("-> %v", tops)
		}
		fmt.Printf("\n")
	}
}

func (t *Tree) AllBottoms() []string {
	bots := []string{}
	for bot, _ := range t.weights {
		bots = append(bots, bot)
	}
	sort.Strings(bots)
	return bots
}

func (t *Tree) VeryBottom() string {
	return t.veryBottom
}

func (t *Tree) Tops(bot string) []string {
	return t.elems[bot]
}

func (t *Tree) IsDescendent(bot, cand string) bool {
	for _, top := range t.elems[bot] {
		if cand == top || t.IsDescendent(top, cand) {
			return true
		}
	}
	return false
}

func (t *Tree) Weight(bot string) int {
	return t.weights[bot]
}

// Determines whether a tree rooted at 'bot' is balanced. Returns true
// if the tree is balanced, or false if it isn't, along with the total
// weight of the tree.
func (t *Tree) IsBalanced(bot string) (bool, int) {
	tops, found := t.elems[bot]
	if !found {
		return true, t.weights[bot]
	}

	_, goal := t.IsBalanced(tops[0])
	//fmt.Printf("goal for %v is %v\n", bot, goal)

	topTotalWeight := goal
	balanced := true
	for _, top := range tops[1:] {
		_, topWeight := t.IsBalanced(top)
		topTotalWeight += topWeight
		//fmt.Printf("topTotalWeight %v\n", topTotalWeight)
		if topWeight != goal {
			balanced = false
		}
	}
	return balanced, topTotalWeight + t.weights[bot]
}

// Given a set of candidate bottoms, determine which is bottom-most in
// the tree (which one isn't a descendent of any of the others).
func findBottomMost(tree *Tree, cands []string) string {
	for len(cands) > 1 {
		adjusted := false
		for _, cand := range cands[1:] {
			if tree.IsDescendent(cands[0], cand) {
				cands = cands[1:]
				adjusted = true
				break
			}
		}

		if !adjusted {
			break
		}
	}

	return cands[0]
}

// Given an unbalanced subtree, find the single unbalanced node that
// needs to be adjusted to bring the subtree into balance.
//
// Assumption: There exists only one node that needs to be adjusted.
func findUnbalanced(tree *Tree, unbalanced string) int {
	if isBalanced, _ := tree.IsBalanced(unbalanced); isBalanced {
		log.Fatalf("root %v is balanced", unbalanced)
	}

	for {
		fmt.Printf("Unbalanced %v\n", unbalanced)
		tops := tree.Tops(unbalanced)
		topWeights := map[int][]string{}

		newUnbalanced := ""
		for _, top := range tops {
			isBalanced, topWeight := tree.IsBalanced(top)

			if _, found := topWeights[topWeight]; !found {
				topWeights[topWeight] = []string{}
			}
			topWeights[topWeight] = append(topWeights[topWeight], top)

			if !isBalanced {
				newUnbalanced = top
				break
			}
		}

		if newUnbalanced != "" {
			// One of the subtrees is unbalanced. Thanks to the
			// assumption above, that there's only one node needing
			// adjustment, we can shift our attention to that
			// subtree.
			unbalanced = newUnbalanced
			continue
		}

		// The subtrees are balanced, which means that one of the
		// children of 'unbalanced' has an incorrect weight. If we group
		// them by subtree weight, one will be an outlier (either
		// because there's >2 subtrees, so majority wins, or because
		// there's 2 subtrees, and we arbitrarily pick the first as the
		// outlier).
		if len(topWeights) != 2 {
			log.Fatalf("found %d topWeights for %v\n", len(topWeights), unbalanced)
		}

		wrongTop := ""
		rightTop := ""
		wrongWeight := 0
		rightWeight := 0
		for topWeight, tops := range topWeights {
			if len(tops) == 1 {
				wrongTop = tops[0]
				wrongWeight = topWeight
			} else {
				rightTop = tops[0]
				rightWeight = topWeight
			}
		}
		difference := rightWeight - wrongWeight

		fmt.Printf("bot %v wrongTop %v wrongWeight %v rightTop %v rightWeight %v diff %v\n",
			unbalanced, wrongTop, wrongWeight, rightTop, rightWeight, difference)

		// Return the weight the wrong node should've had
		return tree.Weight(wrongTop) + difference
	}
}

func main() {
	tree := NewTree()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		matches := pattern.FindStringSubmatch(line)
		if len(matches) == 0 {
			log.Fatalf("failed to parse %v", line)
		}
		bot := matches[1]
		weight, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatalf("failed to parse weight %v in %v", matches[2], line)
		}

		topsStr := matches[3]
		if topsStr == "" {
			tree.Insert(bot, "", weight)
			continue
		}
		tops := strings.Split(topsStr, ", ")

		//fmt.Printf("%v %v\n", bot, tops)

		for _, top := range tops {
			tree.Insert(bot, top, weight)
		}
	}

	//tree.Dump()

	unbalanced := []string{}
	for _, bot := range tree.AllBottoms() {
		if isBalanced, _ := tree.IsBalanced(bot); !isBalanced {
			unbalanced = append(unbalanced, bot)
		}
	}
	fmt.Printf("unbalanced: %v\n", unbalanced)

	bottomUnbalanced := findBottomMost(tree, unbalanced)
	fmt.Printf("bottom-most: %v\n", bottomUnbalanced)

	fmt.Printf("unbalanced: %v\n", findUnbalanced(tree, bottomUnbalanced))
}
