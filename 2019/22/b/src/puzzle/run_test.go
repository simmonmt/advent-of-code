package puzzle

import (
	"reflect"
	"strconv"
	"testing"
)

func TestLarge(t *testing.T) {
	sz := 119315717514047
	inc := 70
	in1 := 119315717513557
	want := 119315717514040

	got := ReverseCommandsForIndex([]*Command{&Command{VERB_DEAL_WITH_INCREMENT, inc}},
		sz, in1)
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestReverseCommandsForIndexSmall(t *testing.T) {
	type TestCase struct {
		in   []int
		cmds []*Command
		out  []int
	}

	testCases := []TestCase{
		TestCase{
			in: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			cmds: []*Command{
				&Command{VERB_DEAL_INTO_NEW_STACK, 0},
			},
			out: []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
		},
		TestCase{
			in: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			cmds: []*Command{
				&Command{VERB_DEAL_WITH_INCREMENT, 3},
			},
			out: []int{0, 7, 4, 1, 8, 5, 2, 9, 6, 3},
		},
		TestCase{
			in: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			cmds: []*Command{
				&Command{VERB_CUT_LEFT, 3},
			},
			out: []int{3, 4, 5, 6, 7, 8, 9, 0, 1, 2},
		},
		TestCase{
			in: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			cmds: []*Command{
				&Command{VERB_CUT_RIGHT, 4},
			},
			out: []int{6, 7, 8, 9, 0, 1, 2, 3, 4, 5},
		},
		TestCase{
			in: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			cmds: []*Command{
				&Command{VERB_DEAL_WITH_INCREMENT, 7},
				&Command{VERB_DEAL_INTO_NEW_STACK, 0},
				&Command{VERB_DEAL_INTO_NEW_STACK, 0},
			},
			out: []int{0, 3, 6, 9, 2, 5, 8, 1, 4, 7},
		},
		TestCase{
			in: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			cmds: []*Command{
				&Command{VERB_CUT_LEFT, 6},
				&Command{VERB_DEAL_WITH_INCREMENT, 7},
				&Command{VERB_DEAL_INTO_NEW_STACK, 0},
			},
			out: []int{3, 0, 7, 4, 1, 8, 5, 2, 9, 6},
		},
		TestCase{
			in: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			cmds: []*Command{
				&Command{VERB_DEAL_INTO_NEW_STACK, 0},
				&Command{VERB_CUT_RIGHT, 2},
				&Command{VERB_DEAL_WITH_INCREMENT, 7},
				&Command{VERB_CUT_LEFT, 8},
				&Command{VERB_CUT_RIGHT, 4},
				&Command{VERB_DEAL_WITH_INCREMENT, 7},
				&Command{VERB_CUT_LEFT, 3},
				&Command{VERB_DEAL_WITH_INCREMENT, 9},
				&Command{VERB_DEAL_WITH_INCREMENT, 3},
				&Command{VERB_CUT_RIGHT, 1},
			},
			out: []int{9, 2, 5, 8, 1, 4, 7, 0, 3, 6},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out := make([]int, len(tc.in))
			for i := range tc.in {
				rev := ReverseCommandsForIndex(tc.cmds, len(tc.in), i)
				fwd := ForwardCommandsForIndex(tc.cmds, len(tc.in), rev)
				if fwd != i {
					t.Errorf("fwd %d != i %d", fwd, i)
				}

				out[i] = tc.in[rev]
			}
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("Reversed = %v, want %v", out, tc.out)
			}
		})
	}
}

func TestCommands(t *testing.T) {
	type TestCase struct {
		in   []int
		cmds []*Command
		out  []int
	}

	testCases := []TestCase{
		TestCase{
			in: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			cmds: []*Command{
				&Command{VERB_DEAL_WITH_INCREMENT, 7},
				&Command{VERB_DEAL_INTO_NEW_STACK, 0},
				&Command{VERB_DEAL_INTO_NEW_STACK, 0},
			},
			out: []int{0, 3, 6, 9, 2, 5, 8, 1, 4, 7},
		},
		TestCase{
			in: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			cmds: []*Command{
				&Command{VERB_CUT_LEFT, 6},
				&Command{VERB_DEAL_WITH_INCREMENT, 7},
				&Command{VERB_DEAL_INTO_NEW_STACK, 0},
			},
			out: []int{3, 0, 7, 4, 1, 8, 5, 2, 9, 6},
		},
		TestCase{
			in: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			cmds: []*Command{
				&Command{VERB_DEAL_INTO_NEW_STACK, 0},
				&Command{VERB_CUT_RIGHT, 2},
				&Command{VERB_DEAL_WITH_INCREMENT, 7},
				&Command{VERB_CUT_LEFT, 8},
				&Command{VERB_CUT_RIGHT, 4},
				&Command{VERB_DEAL_WITH_INCREMENT, 7},
				&Command{VERB_CUT_LEFT, 3},
				&Command{VERB_DEAL_WITH_INCREMENT, 9},
				&Command{VERB_DEAL_WITH_INCREMENT, 3},
				&Command{VERB_CUT_RIGHT, 1},
			},
			out: []int{9, 2, 5, 8, 1, 4, 7, 0, 3, 6},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := RunCommands(tc.in, tc.cmds); !reflect.DeepEqual(got, tc.out) {
				t.Errorf("RunCommands = %v, want %v", got, tc.out)
			}
		})
	}
}

func TestDealIntoNewStack(t *testing.T) {
	in := []int{10, 11, 12, 13, 14, 15}
	want := []int{15, 14, 13, 12, 11, 10}

	if got := DealIntoNewStack(in, 0); !reflect.DeepEqual(got, want) {
		t.Errorf("DealIntoNewStack(%v) = %v, want %v", in, got, want)
	}
}

func TestDealWithIncrement(t *testing.T) {
	in := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []int{0, 7, 4, 1, 8, 5, 2, 9, 6, 3}

	if got := DealWithIncrement(in, 3); !reflect.DeepEqual(got, want) {
		t.Errorf("DealWithIncrement(%v, 3) = %v, want %v", in, got, want)
	}
}

func TestCutLeft(t *testing.T) {
	in := []int{0, 1, 2, 3, 4, 5}
	want := []int{2, 3, 4, 5, 0, 1}

	if got := CutLeft(in, 2); !reflect.DeepEqual(got, want) {
		t.Errorf("CutLeft(%v) = %v, want %v", in, got, want)
	}
}

func TestCutRight(t *testing.T) {
	in := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []int{6, 7, 8, 9, 0, 1, 2, 3, 4, 5}

	if got := CutRight(in, 4); !reflect.DeepEqual(got, want) {
		t.Errorf("CutRight(%v) = %v, want %v", in, got, want)
	}
}
