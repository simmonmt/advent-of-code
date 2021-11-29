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

package instr

import (
	"strconv"
	"testing"
)

func TestSwapPos(t *testing.T) {
	inst := newSwapPos(1, 3)
	if res := inst.String(); res != "swap [1] with [3]" {
		t.Errorf(`swap pos String() = "%v", want "%v"`, res, "swap [1] with [3]")
	}

	str := []byte("abcdef")
	if ok := inst.Exec(str); !ok || string(str) != "adcbef" {
		t.Errorf(`swap 1,3 Exec = "%v", want "%v"`, string(str), "adcbef")
	}
}

func TestSwapChar(t *testing.T) {
	inst := newSwapChar('c', 'b')
	if res := inst.String(); res != "swap c with b" {
		t.Errorf(`swap char String() = "%v", want "%v"`, res, "swap c with b")
	}

	str := []byte("abcdef")
	if ok := inst.Exec(str); !ok || string(str) != "acbdef" {
		t.Errorf(`swap c,b Exec = "%v", want "%v"`, string(str), "acbdef")
	}
}

type TestCase struct {
	inst        Instr
	str         string
	expected    bool
	expectedStr string
}

func testTestCases(t *testing.T, testCases []TestCase) {
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			arr := make([]byte, len(tc.str))
			copy(arr, []byte(tc.str))
			if ok := tc.inst.Exec(arr); ok != tc.expected || string(arr) != tc.expectedStr {
				t.Errorf(`%v Exec("%v") = %v (str now "%v"), want %v, str="%v"`,
					tc.inst, tc.str, ok, string(arr), tc.expected, tc.expectedStr)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	if res := newRotate(true, 3); res.String() != "rotate left 3" {
		t.Errorf(`got "%v", want "rotate left 3"`)
	}
	if res := newRotate(false, 3); res.String() != "rotate right 3" {
		t.Errorf(`got "%v", want "rotate right 3"`)
	}

	testTestCases(t, []TestCase{
		TestCase{newRotate(true, 1), "abcdef", true, "bcdefa"},
		TestCase{newRotate(true, 3), "abcdef", true, "defabc"},
		TestCase{newRotate(true, 21), "abcdef", true, "defabc"},

		TestCase{newRotate(false, 1), "abcdef", true, "fabcde"},
		TestCase{newRotate(false, 3), "abcdef", true, "defabc"},
		TestCase{newRotate(false, 21), "abcdef", true, "defabc"},
	})
}

func TestRotateMagic(t *testing.T) {
	if res := newRotateMagic('c'); res.String() != "rotate magic c" {
		t.Errorf(`got "%v", want "rotate magic c"`, res)
	}

	testTestCases(t, []TestCase{
		TestCase{newRotateMagic('c'), "abcdef", true, "defabc"},
		TestCase{newRotateMagic('d'), "abcdef", true, "cdefab"},
		TestCase{newRotateMagic('b'), "abdec", true, "ecabd"}, // from sample
		TestCase{newRotateMagic('d'), "ecabd", true, "decab"}, // from sample
	})
}

func TestReverse(t *testing.T) {
	if res := newReverse(1, 4); res.String() != "reverse 1 4" {
		t.Errorf(`got "%v", want "reverse 1 4"`, res)
	}

	testTestCases(t, []TestCase{
		TestCase{newReverse(2, 5), "abcdefg", true, "abfedcg"},
		TestCase{newReverse(0, 4), "edcba", true, "abcde"}, // from sample
	})
}

func TestMove(t *testing.T) {
	if res := newMove(1, 4); res.String() != "move 1 4" {
		t.Errorf(`got "%v", want "move 1 4"`, res)
	}

	testTestCases(t, []TestCase{
		TestCase{newMove(1, 4), "bcdea", true, "bdeac"}, // from sample
		TestCase{newMove(3, 0), "bdeac", true, "abdec"}, // from sample
	})
}

func TestParse(t *testing.T) {
	type TestCase struct {
		in     string
		output string
	}

	testCases := []TestCase{
		TestCase{"swap position 4 with position 0", "swap [4] with [0]"},
		TestCase{"swap letter d with letter b", "swap d with b"},
		TestCase{"rotate left 1 step", "rotate left 1"},
		TestCase{"rotate right 2 steps", "rotate right 2"},
		TestCase{"rotate based on position of letter b", "rotate magic b"},
		TestCase{"reverse positions 0 through 4", "reverse 0 4"},
		TestCase{"move position 1 to position 4", "move 1 4"},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			inst, err := Parse(tc.in)
			if err != nil && tc.output != "" {
				t.Errorf("%v: expected nil error, got %v", tc.in, err)
			} else if err == nil && tc.output == "" {
				t.Errorf("%v: expected non-nil error, got nil", tc.in)
			} else if err == nil && inst.String() != tc.output {
				t.Errorf("%v: expected '%v' got '%s'", tc.in, tc.output, inst)
			}
		})
	}
}
