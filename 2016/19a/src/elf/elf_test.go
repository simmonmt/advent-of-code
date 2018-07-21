package elf

import (
	"flag"
	"os"
	"strconv"
	"testing"

	"logger"
)

func TestInitElves(t *testing.T) {
	elves := InitElves(5)
	if res := elves.Len(); res != 5 {
		t.Errorf("InitElves(5).Len() = %v, want %v", res, 5)
	}
}

func TestPlay(t *testing.T) {
	type TestCase struct {
		num     int
		allName uint
	}

	testCases := []TestCase{
		TestCase{5, 3},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			elves := InitElves(tc.num)

			if res := Play(elves); res != tc.allName {
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
