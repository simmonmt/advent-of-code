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

package parse

import (
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2020/common/logger"
)

func TestParse(t *testing.T) {
	type TestCase struct {
		in      []string
		want    string
		matches []string
		fails   []string
	}

	testCases := []TestCase{
		TestCase{
			in: []string{
				`0: 1 2`,
				`1: "a"`,
				`2: 1 3 | 3 1`,
				`3: "b"`,
			},
			want:    "a(?:ab|ba)",
			matches: []string{"aab", "aba"},
			fails:   []string{"aaa"},
		},
		TestCase{
			in: []string{
				`0: 4 1 5`,
				`1: 2 3 | 3 2`,
				`2: 4 4 | 5 5`,
				`3: 4 5 | 5 4`,
				`4: "a"`,
				`5: "b"`,
			},
			want: "a(?:(?:aa|bb)(?:ab|ba)|(?:ab|ba)(?:aa|bb))b",
			matches: []string{
				"aaaabb",
				"aaabab",
				"aabaab",
				"aabbbb",
				"abaaab",
				"ababbb",
				"abbabb",
				"abbbab",
			},
			fails: []string{
				"bababa", "aaabbb", "aaaabbb",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := Parse(tc.in, 0)
			if err != nil || got != tc.want {
				t.Errorf("Parse(_, 0) = %v, %v; want %v, nil",
					got, err, tc.want)
			}

			pat, err := regexp.Compile(got)
			if err != nil {
				t.Errorf("got failed to compile: %v", err)
			}

			for _, wantMatch := range tc.matches {
				if sz := len(pat.FindString(wantMatch)); sz != len(wantMatch) {
					t.Errorf("failed to match %v", wantMatch)
				}
			}

			for _, wantFail := range tc.fails {
				if sz := len(pat.FindString(wantFail)); sz == len(wantFail) {
					t.Errorf("bad match %v", wantFail)
				}
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
