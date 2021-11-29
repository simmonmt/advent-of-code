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

var (
	scorePattern = regexp.MustCompile(`^([^ ]+) would (gain|lose) ([0-9]+) happiness units by sitting next to ([^.]+)\.$`)
)

type Permuter struct {
	elems []string
	last  []int
}

func NewPermuter(elems []string) *Permuter {
	return &Permuter{elems: elems}
}

func (p *Permuter) inc() {
	p.last[len(p.elems)-1]++
	for i := len(p.elems) - 1; i >= 0; i-- {
		if p.last[i] >= len(p.elems) && i > 0 {
			p.last[i] = 0
			p.last[i-1]++
		} else {
			break
		}
	}
}

func (p *Permuter) valid() bool {
	for i := 0; i < len(p.elems)-1; i++ {
		for j := i + 1; j < len(p.elems); j++ {
			if p.last[i] == p.last[j] {
				//fmt.Printf("match %v %v in %v\n", i, j, p.last)
				return false
			}
		}
	}
	return true
}

func (p *Permuter) nextIdxs() []int {
	if p.last == nil {
		p.last = make([]int, len(p.elems))
		for i := 0; i < len(p.elems); i++ {
			p.last[i] = i
		}
		return p.last
	}

	for {
		p.inc()
		if p.last[0] >= len(p.elems) {
			return nil
		}
		if p.valid() {
			return p.last
		}
	}
}

func (p *Permuter) Next() []string {
	idxs := p.nextIdxs()
	if idxs == nil {
		return nil
	}

	out := make([]string, len(idxs))
	for i := range idxs {
		out[i] = p.elems[idxs[i]]
	}
	return out
}

type PersonPair struct {
	A, B string
}

func readScore(line string) (*PersonPair, int, error) {
	matches := scorePattern.FindStringSubmatch(line)
	if matches == nil {
		return nil, 0, fmt.Errorf("failed to parse")
	}

	pairA := matches[1]
	dir := matches[2]
	pairB := matches[4]

	score, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse score %v", matches[2])
	}

	if dir == "lose" {
		score = -score
	}

	return &PersonPair{pairA, pairB}, score, nil
}

func readInput(r io.Reader) (map[PersonPair]int, error) {
	br := bufio.NewReader(r)
	scores := map[PersonPair]int{}
	for lineNum := 1; ; lineNum++ {
		line, err := br.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		if pair, score, err := readScore(line); err != nil {
			return nil, fmt.Errorf("%d: failed to read score: %v", lineNum, err)
		} else {
			scores[*pair] = score
		}
	}

	return scores, nil
}

func scoreForPair(a, b string, scores map[PersonPair]int) int {
	score := scores[PersonPair{a, b}]
	score += scores[PersonPair{b, a}]
	return score
}

func scoreForArrangement(arrangement []string, scores map[PersonPair]int) int {
	score := scoreForPair(arrangement[0], arrangement[len(arrangement)-1], scores)
	for i := 0; i < len(arrangement)-1; i++ {
		score += scoreForPair(arrangement[i], arrangement[i+1], scores)
	}
	return score
}

func main() {
	scores, err := readInput(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	namesMap := map[string]bool{}
	for pair := range scores {
		namesMap[pair.A] = true
		namesMap[pair.B] = true
	}

	names := []string{}
	for name := range namesMap {
		names = append(names, name)
	}
	names = append(names, "Self")

	perm := NewPermuter(names)

	var bestArrangement []string
	var bestScore int

	for {
		arrangement := perm.Next()
		if arrangement == nil {
			break
		}

		score := scoreForArrangement(arrangement, scores)

		if bestArrangement == nil || bestScore < score {
			bestArrangement = arrangement
			bestScore = score
		}
	}
	fmt.Printf("%d: %v\n", bestScore, bestArrangement)
}
