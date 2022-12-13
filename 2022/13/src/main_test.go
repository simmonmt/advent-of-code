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
	_ "embed"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2022/common/logger"
)

var (
	//go:embed sample.txt
	rawSample   string
	sampleLines []string
)

func TestParsePacket(t *testing.T) {
	type TestCase struct {
		in   string
		want *Elem
	}

	testCases := []TestCase{
		TestCase{"[3]", &Elem{Subs: []*Elem{&Elem{Val: 3}}}},
		TestCase{"[3,1]", &Elem{Subs: []*Elem{&Elem{Val: 3}, &Elem{Val: 1}}}},
		TestCase{"[[1],[2,3]]",
			&Elem{Subs: []*Elem{
				&Elem{Subs: []*Elem{&Elem{Val: 1}}}, // [1]
				&Elem{Subs: []*Elem{
					&Elem{Val: 2}, &Elem{Val: 3}}}, // [2,3]
			}}},
		TestCase{"[]", &Elem{Subs: []*Elem{}}},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			packet, err := parsePacket(tc.in)
			if err != nil || !reflect.DeepEqual(packet, tc.want) {
				t.Errorf("parsePacket('%s') = %v, %v, want %v, nil",
					tc.in, packet, err, tc.want)
			}

			if got := packet.String(); got != tc.in {
				t.Errorf("String() = %v, want %v", got, tc.in)
			}
		})
	}
}

func TestPacketCompare(t *testing.T) {
	type TestCase struct {
		a, b string
		want CompareResult
	}

	testCases := []TestCase{
		TestCase{"[1]", "[1]", EQUALS},
		TestCase{"[1,1,3,1,1]", "[1,1,5,1,1]", LESS_THAN},
		TestCase{"[[1],[2,3,4]]", "[[1],4]", LESS_THAN},
		TestCase{"[9]", "[[8,7,6]]", GREATER_THAN},
		TestCase{"[[4,4],4,4]", "[[4,4],4,4,4]", LESS_THAN},
		TestCase{"[7,7,7,7]", "[7,7,7]", GREATER_THAN},
		TestCase{"[]", "[3]", LESS_THAN},
		TestCase{"[[[]]]", "[[]]", GREATER_THAN},
		TestCase{"[1,[2,[3,[4,[5,6,7]]]],8,9]", "[1,[2,[3,[4,[5,6,0]]]],8,9]",
			GREATER_THAN},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a, err := parsePacket(tc.a)
			if err != nil {
				t.Fatalf("can't parse %v: %v", tc.a, err)
			}
			b, err := parsePacket(tc.b)
			if err != nil {
				t.Fatalf("can't parse %v: %v", tc.b, err)
			}

			logger.LogF("test packetCompare %v %v = %v", a, b, packetCompare(a, b))

			if got := packetCompare(a, b); got != tc.want {
				t.Errorf("packetLessThan(%s, %s) = %v, want %v",
					a, b, got, tc.want)
			}
		})
	}
}

func TestSolveA(t *testing.T) {
	pairs, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveA(pairs), 13; got != want {
		t.Errorf("solveA(sample) = %v, want %v", got, want)
	}
}

func TestSolveB(t *testing.T) {
	pairs, err := parseInput(sampleLines)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := solveB(pairs), 140; got != want {
		t.Errorf("solveB(sample) = %v, want %v", got, want)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)

	sampleLines = strings.Split(rawSample, "\n")
	if len(sampleLines) > 0 && sampleLines[len(sampleLines)-1] == "" {
		sampleLines = sampleLines[0 : len(sampleLines)-1]
	}

	os.Exit(m.Run())
}
