package mtsmath

import (
	"fmt"
	"testing"
)

func TestAbs(t *testing.T) {
	if got, want := Abs(1), 1; got != want {
		t.Errorf("Abs(1) = %v, want %v", got, want)
	}
	if got, want := Abs(-1), 1; got != want {
		t.Errorf("Abs(1) = %v, want %v", got, want)
	}
	if got, want := Abs(0), 0; got != want {
		t.Errorf("Abs(1) = %v, want %v", got, want)
	}

	if got, want := Abs(1.5), 1.5; got != want {
		t.Errorf("Abs(1.5) = %v, want %v", got, want)
	}
	if got, want := Abs(-1.5), 1.5; got != want {
		t.Errorf("Abs(1.5) = %v, want %v", got, want)
	}
}

func TestGCD(t *testing.T) {
	type TestCase struct {
		vs []int64
		d  int64
	}

	testCases := []TestCase{
		TestCase{[]int64{4, 6}, 2},
		TestCase{[]int64{8, 12}, 4},
		TestCase{[]int64{12, 18}, 6},
		TestCase{[]int64{9, 9}, 9},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.vs), func(t *testing.T) {
			if got := GCD(tc.vs...); got != tc.d {
				t.Errorf("GCD(%v) = %d, want %d", tc.vs, got, tc.d)
			}
		})
	}
}

func TestLCM(t *testing.T) {
	type TestCase struct {
		vs []int64
		d  int64
	}

	testCases := []TestCase{
		TestCase{[]int64{4, 6}, 12},
		TestCase{[]int64{8, 12}, 24},
		TestCase{[]int64{5, 18}, 90},
		TestCase{[]int64{9, 9}, 9},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.vs), func(t *testing.T) {
			if got := LCM(tc.vs...); got != tc.d {
				t.Errorf("LCM(%v) = %d, want %d", tc.vs, got, tc.d)
			}
		})
	}
}
