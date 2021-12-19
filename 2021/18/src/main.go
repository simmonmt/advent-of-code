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
	"flag"
	"fmt"
	"log"
	"strconv"
	"unicode"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
	"github.com/simmonmt/aoc/2021/common/strutil"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type SNumber struct {
	Lit         int
	Left, Right *SNumber
	Parent      *SNumber
}

func (n *SNumber) String() string {
	if n.Left != nil {
		return fmt.Sprintf("[%v,%v]", n.Left.String(),
			n.Right.String())
	}
	return strconv.Itoa(n.Lit)
}

func (n *SNumber) ReplaceSelf(rep *SNumber) {
	logger.LogF("replacing %v with %v", n, rep)

	if n.Parent == nil {
		panic("can't replace root")
	}
	if rep.Parent != nil {
		panic("rep has parent")
	}

	rep.Parent = n.Parent
	if n.Parent.Left == n {
		n.Parent.Left = rep
	} else if n.Parent.Right == n {
		n.Parent.Right = rep
	} else {
		panic("bad child pointer")
	}
}

func readInt(buf *strutil.StrBuf) (int, error) {
	numStr := ""
	for {
		s, found := buf.Consume(1)
		if !found {
			return -1, fmt.Errorf("short read")
		}
		numStr += s

		r, found := buf.Peek()
		if !found {
			break
		}

		if !unicode.IsDigit(r) {
			break
		}
	}

	//logger.LogF("readInt %v", numStr)

	num, err := strconv.ParseUint(numStr, 10, 32)
	if err != nil {
		return -1, err
	}

	return int(num), nil
}

func doParseSNumber(buf *strutil.StrBuf) (*SNumber, error) {
	//logger.LogF("doParseSNumber %v", string(buf.Rest()))

	start := buf.Off()
	r, found := buf.Peek()
	if !found {
		return nil, fmt.Errorf("%d: short read", start)
	}

	if unicode.IsDigit(r) {
		lit, err := readInt(buf)
		if err != nil {
			return nil, fmt.Errorf("%d: %v", start, err)
		}
		return &SNumber{Lit: lit}, nil
	}

	if !buf.ConsumeIf('[') {
		return nil, fmt.Errorf("%d: unexpected %v", start, string(r))
	}

	left, err := doParseSNumber(buf)
	if err != nil {
		return nil, err
	}

	if !buf.ConsumeIf(',') {
		return nil, fmt.Errorf("%d: missing comma", buf.Off())
	}

	right, err := doParseSNumber(buf)
	if err != nil {
		return nil, err
	}

	if !buf.ConsumeIf(']') {
		return nil, fmt.Errorf("%d: missing ]", buf.Off())
	}

	sn := &SNumber{Left: left, Right: right}
	sn.Left.Parent = sn
	sn.Right.Parent = sn

	return sn, nil
}

func parseSNumber(in string) (*SNumber, error) {
	buf := strutil.NewStrBuf(in)
	return doParseSNumber(buf)
}

func readInput(path string) ([]string, error) {
	lines, err := filereader.Lines(*input)
	return lines, err
}

func explodeFinder(sn *SNumber, depth int) *SNumber {
	if depth == 5 && sn.Left != nil {
		return sn
	}

	if sn.Left != nil {
		if n := explodeFinder(sn.Left, depth+1); n != nil {
			return n
		}
		if n := explodeFinder(sn.Right, depth+1); n != nil {
			return n
		}
	}

	return nil
}

func neighborFinder(sn *SNumber, leftNeighbor bool) *SNumber {
	for {
		p := sn.Parent
		if p == nil {
			return nil
		}

		if leftNeighbor && p.Left != sn {
			sn = p.Left
			break
		} else if !leftNeighbor && p.Right != sn {
			sn = p.Right
			break
		}

		sn = p
	}

	if leftNeighbor {
		for sn.Right != nil {
			sn = sn.Right
		}
		return sn
	} else {
		for sn.Left != nil {
			sn = sn.Left
		}
		return sn
	}
}

func explodeSNumber(sn *SNumber) bool {
	logger.LogF("explodeSNumber %v", sn)

	toExplode := explodeFinder(sn, 1)
	if toExplode == nil {
		return false
	}

	logger.LogF("finder found %v", toExplode)

	leftNeighbor := neighborFinder(toExplode, true)
	rightNeighbor := neighborFinder(toExplode, false)

	logger.LogF("left %v right %v", leftNeighbor, rightNeighbor)

	if leftNeighbor != nil {
		leftNeighbor.Lit += toExplode.Left.Lit
	}
	if rightNeighbor != nil {
		rightNeighbor.Lit += toExplode.Right.Lit
	}

	toExplode.ReplaceSelf(&SNumber{Lit: 0})
	return true
}

func splitFinder(sn *SNumber) *SNumber {
	if sn.Left != nil {
		if n := splitFinder(sn.Left); n != nil {
			return n
		}
		if n := splitFinder(sn.Right); n != nil {
			return n
		}
	} else {
		if sn.Lit >= 10 {
			return sn
		}
	}

	return nil
}

func splitSNumber(sn *SNumber) bool {
	logger.LogF("splitSNumber %v", sn)

	toSplit := splitFinder(sn)
	if toSplit == nil {
		return false
	}

	logger.LogF("finder found %v in %v", toSplit, sn)

	left := toSplit.Lit / 2
	nSn := &SNumber{
		Left:  &SNumber{Lit: left},
		Right: &SNumber{Lit: toSplit.Lit - left},
	}
	nSn.Left.Parent = nSn
	nSn.Right.Parent = nSn

	toSplit.ReplaceSelf(nSn)

	logger.LogF("nsn %v; sn now %v", nSn, sn)

	return true
}

func reduceSNumber(sn *SNumber) {
	for i := 0; i < 1000; i++ {
		if explodeSNumber(sn) {
			continue
		}

		if splitSNumber(sn) {
			continue
		}

		return
	}

	panic("infinite loop?")
}

func addSNumber(n1, n2 *SNumber) *SNumber {
	result := &SNumber{
		Left:  n1,
		Right: n2,
	}
	n1.Parent = result
	n2.Parent = result

	logger.LogF("added %v to %v", n1, n2)

	reduceSNumber(result)
	return result
}

func addSNumbers(ins []*SNumber) *SNumber {
	out := ins[0]
	for i := 1; i < len(ins); i++ {
		out = addSNumber(out, ins[i])
	}
	return out
}

func magnitude(sn *SNumber) int {
	if sn.Left != nil {
		return 3*magnitude(sn.Left) + 2*magnitude(sn.Right)
	}
	return sn.Lit
}

func solveA(lines []string) {
	sns := []*SNumber{}
	for i, line := range lines {
		sn, err := parseSNumber(line)
		if err != nil {
			log.Fatalf("%d: parse failure: %v", i, err)
		}
		sns = append(sns, sn)
	}

	fmt.Println("A", magnitude(addSNumbers(sns)))
}

func solveB(lines []string) {
	maxMag := 0

	for i, l1 := range lines {
		for j, l2 := range lines {
			if i == j {
				continue
			}

			sn1, _ := parseSNumber(l1)
			sn2, _ := parseSNumber(l2)

			mag := magnitude(addSNumber(sn1, sn2))
			if mag > maxMag {
				maxMag = mag
			}
		}
	}

	fmt.Println("B", maxMag)
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

	solveA(lines)
	solveB(lines)
}
