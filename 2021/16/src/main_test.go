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
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2021/common/logger"
)

func TestHexToBin(t *testing.T) {
	type TestCase struct {
		in   string
		want string
	}

	testCases := []TestCase{
		TestCase{"D2FE28", "110100101111111000101000"},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			if got, err := hexToBin(tc.in); err != nil || got != tc.want {
				t.Errorf(`hexToBin("%v") = %v, %v, want %v, nil`,
					tc.in, got, err, tc.want)
			}
		})
	}
}

func TestArrView(t *testing.T) {
	in := []byte{'a', 'b', 'c', 'd'}
	v := NewArrView(in)

	if got := v.Left(); got != 4 {
		t.Fatalf("v.Left() = %v, want 4", got)
	}
	if got := v.Off(); got != 0 {
		t.Fatalf("v.Left() = %v, want 0", got)
	}

	if got, err := v.Consume(2); err != nil || !reflect.DeepEqual(got, []byte{'a', 'b'}) {
		t.Fatalf(`v.Consume(2) = %v, %v, want ['a','b'], nil`, got, err)
	}

	if got := v.Left(); got != 2 {
		t.Fatalf("v.Left() = %v, want 2", got)
	}
	if got := v.Off(); got != 2 {
		t.Fatalf("v.Left() = %v, want 2", got)
	}

	if got, err := v.Consume(1); err != nil || !reflect.DeepEqual(got, []byte{'c'}) {
		t.Fatalf(`v.Consume(2) = %v, %v, want ['c'], nil`, got, err)
	}

	if got := v.Left(); got != 1 {
		t.Fatalf("v.Left() = %v, want 1", got)
	}
	if got := v.Off(); got != 3 {
		t.Fatalf("v.Left() = %v, want 3", got)
	}

	if got, err := v.Consume(9); err == nil {
		t.Fatalf(`v.Consume(9) = %v, %v, want _, non-nil`, got, err)
	}
}

func TestConsumeInt(t *testing.T) {
	v := NewArrView([]byte("1101110"))
	if got, err := consumeInt(v, 3); err != nil || got != 6 {
		t.Fatalf(`consumeInt("11011", 3) = %v, %v; want 6, nil`, got, err)
	}

	if got := v.Left(); got != 4 {
		t.Fatalf("v.Left() = %v, want 4", got)
	}
	if got := v.Off(); got != 3 {
		t.Fatalf("v.Left() = %v, want 3", got)
	}

	if got, err := consumeInt(v, 3); err != nil || got != 7 {
		t.Fatalf(`consumeInt("111", 3) = %v, %v; want 6, nil`, got, err)
	}
}

func TestDecodeLiteral(t *testing.T) {
	type TestCase struct {
		in   string
		want int
	}

	testCases := []TestCase{
		TestCase{in: "00101", want: 5},
		TestCase{in: "1010101100", want: 0x5c},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			in := NewArrView([]byte(tc.in))
			if got, err := decodeLiteral(in); err != nil || got != tc.want {
				t.Errorf(`decodeLiteral("%v") = %v, %v, want %v, nil`,
					tc.in, got, err, tc.want)
			}
			if in.Left() != 0 {
				t.Errorf("in un-empty: %+v", in)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type TestCase struct {
		in         string
		versionSum int
	}

	testCases := []TestCase{
		TestCase{
			in:         "D2FE28",
			versionSum: 6,
		},
		TestCase{
			in:         "38006F45291200",
			versionSum: 9,
		},
		TestCase{
			in:         "EE00D40C823060",
			versionSum: 7 + 2 + 4 + 1,
		},
		TestCase{
			in:         "8A004A801A8002F478",
			versionSum: 4 + 1 + 5 + 6,
		},
		TestCase{
			in:         "620080001611562C8802118E34",
			versionSum: 12,
		},
		TestCase{
			in:         "C0015000016115A2E0802F182340",
			versionSum: 23,
		},
		TestCase{
			in:         "A0016C880162017C3686B18A3D4780",
			versionSum: 31,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			packet, err := decode(tc.in)
			if err != nil {
				t.Fatalf(`decode("%v") = _, %v, want _, nil`,
					tc.in, err)
			}

			packet.DumpTo(os.Stdout)

			if got := versionSum(packet); got != tc.versionSum {
				t.Errorf(`versionSum(_) = %v, want %v, nil`,
					got, tc.versionSum)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	type TestCase struct {
		in   string
		want int
	}

	testCases := []TestCase{
		TestCase{"C200B40A82", 3},
		TestCase{"04005AC33890", 54},
		TestCase{"880086C3E88112", 7},
		TestCase{"CE00C43D881120", 9},
		TestCase{"D8005AC2A8F0", 1},
		TestCase{"F600BC2D8F", 0},
		TestCase{"9C005AC2F8F0", 0},
		TestCase{"9C0141080250320F1802104A08", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			packet, err := decode(tc.in)
			if err != nil {
				t.Fatalf(`decode("%v") = _, %v, want _, nil`,
					tc.in, err)
			}

			fmt.Printf("packet for %v:\n", tc.in)
			packet.DumpTo(os.Stdout)

			if got := evaluate(packet); got != tc.want {
				t.Errorf(`evaluate(_) = %v, want %v`, got, tc.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
