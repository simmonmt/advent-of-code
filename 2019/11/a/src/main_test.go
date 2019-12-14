package main

import (
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2019/common/pos"
)

func TestExecuteTurn(t *testing.T) {
	type TestCase struct {
		dir       Dir
		turnRight bool
		newDir    Dir
		newP      pos.P2
	}

	testCases := []TestCase{
		TestCase{DIR_UP, true, DIR_RIGHT, pos.P2{1, 0}},
		TestCase{DIR_UP, false, DIR_LEFT, pos.P2{-1, 0}},
		TestCase{DIR_DOWN, true, DIR_LEFT, pos.P2{-1, 0}},
		TestCase{DIR_DOWN, false, DIR_RIGHT, pos.P2{1, 0}},
		TestCase{DIR_LEFT, true, DIR_UP, pos.P2{0, -1}},
		TestCase{DIR_LEFT, false, DIR_DOWN, pos.P2{0, 1}},
		TestCase{DIR_RIGHT, true, DIR_DOWN, pos.P2{0, 1}},
		TestCase{DIR_RIGHT, false, DIR_UP, pos.P2{0, -1}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			start := pos.P2{0, 0}
			if newDir, newP := executeTurn(tc.dir, start, tc.turnRight); newDir != tc.newDir || !tc.newP.Equals(newP) {
				t.Errorf("executeTurn(%s, %v, %v) = %s, %v, want %s, %v",
					tc.dir, start, tc.turnRight,
					newDir, newP,
					tc.newDir, tc.newP)
			}
		})
	}
}
