package instr

import (
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

func TestHlf(t *testing.T) {
	i := newHlf(reg.B)
	if s := i.String(); s != "hlf b" {
		t.Errorf("i = \"%s\", want \"hlf b\"", s)
	}

	testSingleRegInstr(t, i, 8, 4)
}

func TestTpl(t *testing.T) {
	i := newTpl(reg.B)
	if s := i.String(); s != "tpl b" {
		t.Errorf("i = \"%s\", want \"tpl b\"", s)
	}

	testSingleRegInstr(t, i, 3, 9)
}

func TestInc(t *testing.T) {
	i := newInc(reg.B)
	if s := i.String(); s != "inc b" {
		t.Errorf("i = \"%v\", want \"inc b\"", s)
	}

	testSingleRegInstr(t, i, 3, 4)
}

func testJmp(t *testing.T) {
	i := newJmp(12)
	if s := i.String(); s != "jmp 12" {
		t.Errorf("i = \"%v\", want \"jmp 12\"", s)
	}

	f := reg.NewFile()
	f.Set(reg.A, 5)
	f.Set(reg.B, 5)

	if rPC := i.Exec(f); rPC != 12 {
		t.Errorf("\"%v\" rPC = %v, want 12", i, rPC)
	}
	verifyFile(t, f, 5, 5)
}

func testJie(t *testing.T) {
	if s := newJie(reg.B, 12).String(); s != "jie 12" {
		t.Errorf("i = \"%v\", want \"jie b, 12\"", s)
	}

	f := reg.NewFile()
	f.Set(reg.A, 5)
	f.Set(reg.B, 6)

	if rPC := newJie(reg.A, 12).Exec(f); rPC != 1 {
		t.Errorf("jie A, 12 = %v, want 1", rPC)
	}
	if rPC := newJie(reg.B, 12).Exec(f); rPC != 12 {
		t.Errorf("jie B, 12 = %v, want 12", rPC)
	}

	verifyFile(t, f, 5, 6)
}

func testJio(t *testing.T) {
	if s := newJio(reg.B, 12).String(); s != "jio 12" {
		t.Errorf("i = \"%v\", want \"jio b, 12\"", s)
	}

	f := reg.NewFile()
	f.Set(reg.A, 5)
	f.Set(reg.B, 1)

	if rPC := newJio(reg.A, 12).Exec(f); rPC != 1 {
		t.Errorf("jio A, 12 = %v, want 1", rPC)
	}
	if rPC := newJio(reg.B, 12).Exec(f); rPC != 12 {
		t.Errorf("jio B, 12 = %v, want 12", rPC)
	}

	verifyFile(t, f, 5, 1)
}

func TestParse(t *testing.T) {
	type testCase struct {
		op, a, b string
		output   string
	}

	testCases := []testCase{
		testCase{"hlf", "a", "", "hlf a"},
		testCase{"hlf", "b", "", "hlf b"},
		testCase{"hlf", "invalid", "", ""},
		testCase{"hlf", "a", "12", ""},

		testCase{"tpl", "a", "", "tpl a"},
		testCase{"inc", "a", "", "inc a"},

		testCase{"jmp", "+1", "", "jmp 1"},
		testCase{"jmp", "1", "", "jmp 1"},
		testCase{"jmp", "-1", "", "jmp -1"},
		testCase{"jmp", "1", "12", ""},

		testCase{"jie", "a", "+1", "jie a, 1"},
		testCase{"jie", "a", "1", "jie a, 1"},
		testCase{"jie", "a", "-1", "jie a, -1"},
		testCase{"jie", "a", "x", ""},

		testCase{"jio", "a", "+1", "jio a, 1"},
		testCase{"jio", "a", "1", "jio a, 1"},
		testCase{"jio", "a", "-1", "jio a, -1"},
		testCase{"jio", "a", "x", ""},
	}

	for _, tc := range testCases {
		i, err := Parse(tc.op, tc.a, tc.b)
		if err != nil && tc.output != "" {
			t.Errorf("%v %v %v: expected nil error, got %v", tc.op, tc.a, tc.b, err)
		} else if err == nil && tc.output == "" {
			t.Errorf("%v %v %v: expected non-nil error, got nil", tc.op, tc.a, tc.b)
		} else if err == nil && i.String() != tc.output {
			t.Errorf("%v %v %v: expected '%v' got '%s'", tc.op, tc.a, tc.b, tc.output, i)
		}
	}
}
