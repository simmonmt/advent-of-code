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
	"io"
	"log"
	"strconv"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

func readInput(path string) (string, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return "", err
	}

	if len(lines) != 1 {
		return "", fmt.Errorf("too many lines: %v, want 1", len(lines))
	}

	return lines[0], err
}

type ArrView struct {
	base []byte
	off  int
}

func NewArrView(a []byte) *ArrView {
	return &ArrView{
		base: a,
		off:  0,
	}
}

func (v *ArrView) Consume(l int) ([]byte, error) {
	if l > len(v.base)-v.off {
		return nil, fmt.Errorf("short read")
	}

	out := v.base[v.off : v.off+l]
	v.off += l
	return out, nil
}

func (v *ArrView) Left() int {
	return len(v.base) - v.off
}

func (v *ArrView) Off() int {
	return v.off
}

func (v *ArrView) Rest() []byte {
	return v.base[v.off:]
}

type Packet struct {
	Version int
	ID      int
	Literal int
	Subs    []*Packet
}

func (p *Packet) DumpTo(w io.Writer) {
	p.dump(w, "")
}

func (p *Packet) dump(w io.Writer, indent string) {
	fmt.Fprintf(w, "%vVer %v ID %v Lit %v\n", indent, p.Version, p.ID, p.Literal)
	if len(p.Subs) > 0 {
		fmt.Fprintf(w, "%v{\n", indent)

		for i, sub := range p.Subs {
			if i > 0 {
				fmt.Fprintf(w, "%v}, {\n", indent)
			}

			sub.dump(w, indent+"  ")
		}

		fmt.Fprintf(w, "%v}\n", indent)
	}
}

func hexToBin(hex string) (string, error) {
	out := ""
	for _, r := range hex {
		d, err := strconv.ParseInt(string(r), 16, 8)
		if err != nil {
			return "", err
		}
		s := strconv.FormatInt(d, 2)
		for len(s) < 4 {
			s = "0" + s
		}

		out += s
	}
	return out, nil
}

func consumeStr(v *ArrView, l int) (string, error) {
	out, err := v.Consume(l)
	if err != nil {
		return "", err
	}

	return string(out), err
}

func consumeInt(v *ArrView, l int) (int, error) {
	s, err := consumeStr(v, l)
	if err != nil {
		return 0, err
	}

	n, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return 0, err
	}

	return int(n), nil
}

func decode(in string) (*Packet, error) {
	in, err := hexToBin(in)
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	v := NewArrView([]byte(in))
	return decodePacket(v)
}

func decodePacket(rest *ArrView) (*Packet, error) {
	logger.LogF("decodePacket %v: %v", rest.Off(), string(rest.Rest()))

	wrapErr := func(where string, sub error) error {
		return fmt.Errorf("%d: %v: %v", rest.Off(), where, sub)
	}

	version, err := consumeInt(rest, 3)
	if err != nil {
		return nil, wrapErr("version", err)
	}

	id, err := consumeInt(rest, 3)
	if err != nil {
		return nil, wrapErr("id", err)
	}

	p := &Packet{
		Version: version,
		ID:      id,
	}

	logger.LogF("packet header %+v rest %v: %v",
		p, rest.Off(), string(rest.Rest()))

	if id == 4 {
		lit, err := decodeLiteral(rest)
		if err != nil {
			return nil, wrapErr("literal", err)
		}

		logger.LogF("decoded literal %v", lit)

		p.Literal = lit
		return p, nil
	}

	lenTypeID, err := consumeInt(rest, 1)
	if err != nil {
		return nil, wrapErr("len type id", err)
	}

	if lenTypeID == 0 {
		var bitLen int
		bitLen, err = consumeInt(rest, 15)
		if err != nil {
			return nil, wrapErr("mode 0 bitlen", err)
		}

		subs := []*Packet{}
		for bitLen > 0 {
			left := rest.Left()
			sub, err := decodePacket(rest)
			if err != nil {
				return nil, wrapErr("mode 0 sub packet", err)
			}

			consumed := left - rest.Left()
			subs = append(subs, sub)

			bitLen -= consumed
			if bitLen < 0 {
				panic("underflow")
			}
		}

		p.Subs = subs
	} else {
		var numPackets int
		numPackets, err = consumeInt(rest, 11)
		if err != nil {
			return nil, wrapErr("mode 1 num packets", err)
		}

		subs := []*Packet{}
		for i := 0; i < numPackets; i++ {
			var sub *Packet
			sub, err = decodePacket(rest)
			if err != nil {
				return nil, wrapErr("mode 1 sub packet", err)
			}

			subs = append(subs, sub)
		}

		p.Subs = subs
	}

	return p, nil
}

func decodeLiteral(rest *ArrView) (int, error) {
	out := 0
	for {
		var n int
		var err error
		n, err = consumeInt(rest, 5)
		if err != nil {
			return 0, err
		}

		out = (out << 4) | (n & 0xf)
		if n&0x10 == 0 {
			break
		}
	}
	return out, nil
}

func versionSum(packet *Packet) int {
	sum := packet.Version
	for _, sub := range packet.Subs {
		sum += versionSum(sub)
	}
	return sum
}

func solveA(p *Packet) {
	fmt.Println("A", versionSum(p))
}

func evaluate(packet *Packet) int {
	switch packet.ID {
	case 0: // sum
		sum := 0
		for _, sub := range packet.Subs {
			sum += evaluate(sub)
		}
		return sum

	case 1: // product
		prod := 1
		for _, sub := range packet.Subs {
			prod *= evaluate(sub)
		}
		return prod

	case 2: // minimum
		min := -1
		for _, sub := range packet.Subs {
			if v := evaluate(sub); min == -1 || v < min {
				min = v
			}
		}
		return min

	case 3: // maximum
		max := -1
		for _, sub := range packet.Subs {
			if v := evaluate(sub); v > max {
				max = v
			}
		}
		return max

	case 4: // literal
		return packet.Literal

	case 5: // greater than
		if len(packet.Subs) != 2 {
			panic("bad gt")
		}

		if evaluate(packet.Subs[0]) > evaluate(packet.Subs[1]) {
			return 1
		}
		return 0

	case 6: // less than
		if len(packet.Subs) != 2 {
			panic("bad lt")
		}

		if evaluate(packet.Subs[0]) < evaluate(packet.Subs[1]) {
			return 1
		}
		return 0

	case 7: // equal
		if len(packet.Subs) != 2 {
			panic("bad eq")
		}

		if evaluate(packet.Subs[0]) == evaluate(packet.Subs[1]) {
			return 1
		}
		return 0

	default:
		panic(fmt.Sprintf("bad id %v", packet.ID))
	}
}

func solveB(p *Packet) {
	fmt.Println("B", evaluate(p))
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	line, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	p, err := decode(line)
	if err != nil {
		log.Fatalf("bad decode: %v", err)
	}

	solveA(p)
	solveB(p)
}
