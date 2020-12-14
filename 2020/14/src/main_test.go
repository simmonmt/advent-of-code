package main

import (
	"strconv"
	"testing"
)

func TestSolveA(t *testing.T) {
	type TestCase struct {
		lines []string
		want  uint64
	}

	testCases := []TestCase{
		TestCase{
			lines: []string{
				"mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X",
				"mem[8] = 11",
				"mem[7] = 101",
				"mem[8] = 0",
			},
			want: 165,
		},
		TestCase{
			lines: []string{
				"mask = 10XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				"mem[8] = 60129542144",
			},
			want: 42949672960,
		},
		TestCase{
			lines: []string{
				"mask = 10XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				"mem[8] = 60129542144",
				"mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X",
				"mem[8] = 101",
			},
			want: 101,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := solveA(tc.lines); got != tc.want {
				t.Errorf("solveA(...) = %v, want %v", got, tc.want)
			}
		})
	}
}
