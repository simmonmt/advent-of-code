package data

import (
	"testing"

	"util"
)

func TestGrow(t *testing.T) {
	type TestCase struct {
		in       string
		expected string
	}

	testCases := []TestCase{
		TestCase{"1", "100"},
		TestCase{"0", "001"},
		TestCase{"11111", "11111000000"},
		TestCase{"111100001010", "1111000010100101011110000"},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			in := util.StrToBoolArray(tc.in)
			res := util.BoolArrayToStr(Grow(in))
			if res != tc.expected {
				t.Errorf("Grow(e%v) = %v, want %v", tc.in, res, tc.expected)
			}
		})
	}
}

func TestChecksum(t *testing.T) {
	type TestCase struct {
		in       string
		expected string
	}

	testCases := []TestCase{
		TestCase{"110010110100", "100"},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			in := util.StrToBoolArray(tc.in)
			res := util.BoolArrayToStr(Checksum(in))
			if res != tc.expected {
				t.Errorf("Checksum(%v) = %v, want %v", tc.in, res, tc.expected)
			}
		})
	}
}
