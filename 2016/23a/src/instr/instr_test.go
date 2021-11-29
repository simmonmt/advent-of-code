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

	"reg"
)

func verifyFile(t *testing.T, f *reg.File, expA, expB uint32) {
	if f.Get(reg.A) != expA || f.Get(reg.B) != expB {
		t.Errorf("reg file unexpected, got A=%v b=%v, want A=%v b=%v",
			f.Get(reg.A), f.Get(reg.B), expA, expB)
	}
}

func testSingleRegInstr(t *testing.T, i Instr, initB, expB uint32) {
	f := reg.NewFile()
	f.Set(reg.A, 99)
	f.Set(reg.B, initB)

	if rPC := i.Exec(f); rPC != 1 {
		t.Errorf("\"%v\" rPC = %v, want 1", rPC)
	}
	verifyFile(t, f, 99, expB)
}

func TestInc(t *testing.T) {
	i := newInc(reg.B)
	if s := i.String(); s != "inc b" {
		t.Errorf("i = \"%s\", want \"inc b\"", s)
	}

	testSingleRegInstr(t, i, 8, 9)
}

func TestDec(t *testing.T) {
	i := newDec(reg.B)
	if s := i.String(); s != "dec b" {
		t.Errorf("i = \"%s\", want \"inc b\"", s)
	}

	testSingleRegInstr(t, i, 8, 7)
}

func TestJnzReg(t *testing.T) {
	if s := newJnzReg(reg.A, 12).String(); s != "jnz a, 12" {
		t.Errorf("i = \"%v\", want \"jnz a, 12\"", s)
	}

	f := reg.NewFile()
	f.Set(reg.A, 5)
	f.Set(reg.B, 0)

	if rPC := newJnzReg(reg.A, 12).Exec(f); rPC != 12 {
		t.Errorf("jnz a, 12 rPC = %v, want 12", rPC)
	}
	if rPC := newJnzReg(reg.B, 12).Exec(f); rPC != 1 {
		t.Errorf("jnz b, 12 rPC = %v, want 1", rPC)
	}
	verifyFile(t, f, 5, 0)
}

func TestJnzImm(t *testing.T) {
	if s := newJnzImm(1, 12).String(); s != "jnz 1, 12" {
		t.Errorf("i = \"%v\", want \"jnz 1, 12\"", s)
	}

	f := reg.NewFile()
	f.Set(reg.A, 5)
	f.Set(reg.B, 0)

	if rPC := newJnzImm(5, 12).Exec(f); rPC != 12 {
		t.Errorf("jnz 5, 12 rPC = %v, want 12", rPC)
	}
	if rPC := newJnzImm(0, 12).Exec(f); rPC != 1 {
		t.Errorf("jnz 0, 12 rPC = %v, want 1", rPC)
	}
	verifyFile(t, f, 5, 0)
}

func TestParse(t *testing.T) {
	type testCase struct {
		op, a, b string
		output   string
	}

	testCases := []testCase{
		testCase{"cpy", "a", "b", "cpy a, b"},
		testCase{"cpy", "a", "", ""},
		testCase{"cpy", "a", "invalid", ""},
		testCase{"cpy", "invalid", "b", ""},

		testCase{"cpy", "12", "a", "cpy 12, a"},
		testCase{"cpy", "-12", "a", ""},
		testCase{"cpy", "a", "12", ""},

		testCase{"inc", "a", "", "inc a"},
		testCase{"dec", "a", "", "dec a"},

		testCase{"jnz", "a", "+1", "jnz a, 1"},
		testCase{"jnz", "a", "1", "jnz a, 1"},
		testCase{"jnz", "a", "-1", "jnz a, -1"},
		testCase{"jnz", "a", "b", ""},

		testCase{"jnz", "1", "1", "jnz 1, 1"},
		testCase{"jnz", "1", "a", ""},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			inst, err := Parse(tc.op, tc.a, tc.b)
			if err != nil && tc.output != "" {
				t.Errorf("%v %v %v: expected nil error, got %v", tc.op, tc.a, tc.b, err)
			} else if err == nil && tc.output == "" {
				t.Errorf("%v %v %v: expected non-nil error, got nil", tc.op, tc.a, tc.b)
			} else if err == nil && inst.String() != tc.output {
				t.Errorf("%v %v %v: expected '%v' got '%s'", tc.op, tc.a, tc.b, tc.output, inst)
			}
		})
	}
}
