package elf

import (
	"flag"
	"os"
	"strconv"
	"testing"

	"logger"
)

func TestPlay(t *testing.T) {
	type TestCase struct {
		num     int
		allName int
	}

	testCases := []TestCase{
		TestCase{5, 2},
		TestCase{12, 3},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if res := Play(tc.num); res != tc.allName {
				t.Errorf("Play(%v) = %v, want %v", tc.num, res, tc.allName)
			}
		})
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	logger.Init(true)

	os.Exit(m.Run())
}
