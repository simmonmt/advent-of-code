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
