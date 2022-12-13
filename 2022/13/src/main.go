// Copyright 2022 Google LLC
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
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/simmonmt/aoc/2022/common/filereader"
	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type Elem struct {
	Subs []*Elem
	Val  int
}

func (e *Elem) IsList() bool {
	return e.Subs != nil
}

func (e *Elem) String() string {
	if e.IsList() {
		out := make([]string, len(e.Subs))
		for i, subElem := range e.Subs {
			out[i] = subElem.String()
		}
		return fmt.Sprintf("[%s]", strings.Join(out, ","))
	} else {
		return fmt.Sprintf("%d", e.Val)
	}
}

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(path)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func parseList(line []byte) (elem *Elem, consumed int, err error) {
	elem = &Elem{Subs: []*Elem{}}

	if len(line) < 2 {
		return nil, -1, fmt.Errorf("short list")
	}
	if line[0] != '[' {
		return nil, -1, fmt.Errorf("missing [ at start of list")
	}
	if line[1] == ']' {
		return elem, 2, nil
	}

	for i := 1; i < len(line); {
		sub, n, err := parseListElem(line[i:])
		if err != nil {
			return nil, -1, err
		}
		elem.Subs = append(elem.Subs, sub)
		if i+n >= len(line) {
			break
		}
		i += n

		if line[i] == ',' {
			i++
			continue
		} else if line[i] == ']' {
			return elem, i + 1, nil
		} else {
			return nil, -1, fmt.Errorf("unexpected %s", string(line[i]))
		}
	}

	return nil, -1, fmt.Errorf("unterminated")
}

func parseListElem(line []byte) (elem *Elem, consumed int, err error) {
	if line[0] == '[' {
		l, n, err := parseList(line)
		if err != nil {
			return nil, -1, err
		}
		return l, n, nil
	}

	num := 0
	for i := 0; i < len(line); i++ {
		if line[i] >= '0' && line[i] <= '9' {
			num = num*10 + int(line[i]-'0')
		} else if line[i] == ',' || line[i] == ']' {
			if i == 0 {
				return nil, -1, fmt.Errorf("unexpected %s", string(line[i]))
			}
			return &Elem{Val: num}, i, nil
		} else {
			return nil, -1, fmt.Errorf("bad char %s", string(line[i]))
		}
	}

	return nil, -1, fmt.Errorf("unterminated")
}

func parsePacket(line string) (*Elem, error) {
	elem, n, err := parseList([]byte(line))
	if err != nil {
		return nil, err
	}
	if n != len(line) {
		return nil, fmt.Errorf("unexpected extra; left: '%s'", line[n:])
	}

	return elem, nil
}

type CompareResult int

const (
	LESS_THAN    CompareResult = 0
	EQUALS       CompareResult = 1
	GREATER_THAN CompareResult = 2
)

func packetCompare(a *Elem, b *Elem) CompareResult {
	if !a.IsList() && !b.IsList() {
		if a.Val < b.Val {
			return LESS_THAN
		} else if a.Val == b.Val {
			return EQUALS
		} else {
			return GREATER_THAN
		}
	}

	if a.IsList() != b.IsList() {
		if !a.IsList() {
			a = &Elem{Subs: []*Elem{a}}
		} else {
			b = &Elem{Subs: []*Elem{b}}
		}
	}

	logger.LogF("comparing %v and %v", a, b)

	i := 0
	for ; i < len(a.Subs); i++ {
		if i >= len(b.Subs) {
			// a > b because b ran out first
			return GREATER_THAN
		}
		if result := packetCompare(a.Subs[i], b.Subs[i]); result != EQUALS {
			return result
		}
	}

	if len(a.Subs) == len(b.Subs) {
		return EQUALS
	}

	return LESS_THAN // a ran out first so it's less
}

func parseInput(lines []string) ([][2]*Elem, error) {
	pairs := [][2]*Elem{}
	for i := 0; i < len(lines); i += 3 {
		one, err := parsePacket(lines[i])
		if err != nil {
			return nil, fmt.Errorf("bad packet line %d: %v",
				i, err)
		}

		two, err := parsePacket(lines[i+1])
		if err != nil {
			return nil, fmt.Errorf("bad packet line %d: %v",
				i+1, err)
		}

		pairs = append(pairs, [2]*Elem{one, two})
	}

	return pairs, nil
}

func solveA(pairs [][2]*Elem) int {
	out := 0

	for i, pair := range pairs {
		logger.LogF("pairs %v %v %v", pair[0], pair[1], packetCompare(pair[0], pair[1]))

		if packetCompare(pair[0], pair[1]) == LESS_THAN {
			logger.LogF("pair %d in the right order", i+1)
			out += (i + 1)
		}
	}

	return out
}

func solveB(pairs [][2]*Elem) int {
	divider2, _ := parsePacket("[[2]]")
	divider6, _ := parsePacket("[[6]]")

	packets := []*Elem{divider2, divider6}
	for _, pair := range pairs {
		packets = append(packets, pair[0], pair[1])
	}

	sort.Slice(packets, func(i, j int) bool {
		return packetCompare(packets[i], packets[j]) == LESS_THAN
	})

	out := 1
	found2, found6 := false, false
	for i, packet := range packets {
		if !found2 && packetCompare(packet, divider2) == EQUALS {
			out *= (i + 1)
			found2 = true
		}
		if !found6 && packetCompare(packet, divider6) == EQUALS {
			out *= (i + 1)
			found6 = true
		}
	}

	return out
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	lines, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	pairs, err := parseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("A", solveA(pairs))
	fmt.Println("B", solveB(pairs))
}
