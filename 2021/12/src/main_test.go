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

package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/simmonmt/aoc/2021/common/logger"
)

func pathsToMap(paths []string) map[string]bool {
	out := map[string]bool{}
	for _, p := range paths {
		out[p] = true
	}
	return out
}

func TestAllPaths(t *testing.T) {
	type TestCase struct {
		input    []string
		paths    []string
		numPaths int
	}

	testCases := []TestCase{
		TestCase{
			input: []string{
				"start-A",
				"start-b",
				"A-c",
				"A-b",
				"b-d",
				"A-end",
				"b-end",
			},
			paths: []string{
				"start,A,b,A,c,A,end",
				"start,A,b,A,end",
				"start,A,b,end",
				"start,A,c,A,b,A,end",
				"start,A,c,A,b,end",
				"start,A,c,A,end",
				"start,A,end",
				"start,b,A,c,A,end",
				"start,b,A,end",
				"start,b,end",
			},
		},
		TestCase{
			input: []string{
				"dc-end",
				"HN-start",
				"start-kj",
				"dc-start",
				"dc-HN",
				"LN-dc",
				"HN-end",
				"kj-sa",
				"kj-HN",
				"kj-dc",
			},
			paths: []string{
				"start,HN,dc,HN,end",
				"start,HN,dc,HN,kj,HN,end",
				"start,HN,dc,end",
				"start,HN,dc,kj,HN,end",
				"start,HN,end",
				"start,HN,kj,HN,dc,HN,end",
				"start,HN,kj,HN,dc,end",
				"start,HN,kj,HN,end",
				"start,HN,kj,dc,HN,end",
				"start,HN,kj,dc,end",
				"start,dc,HN,end",
				"start,dc,HN,kj,HN,end",
				"start,dc,end",
				"start,dc,kj,HN,end",
				"start,kj,HN,dc,HN,end",
				"start,kj,HN,dc,end",
				"start,kj,HN,end",
				"start,kj,dc,HN,end",
				"start,kj,dc,end",
			},
		},
		TestCase{
			input: []string{
				"fs-end",
				"he-DX",
				"fs-he",
				"start-DX",
				"pj-DX",
				"end-zg",
				"zg-sl",
				"zg-pj",
				"pj-he",
				"RW-he",
				"fs-DX",
				"pj-RW",
				"zg-RW",
				"start-pj",
				"he-WI",
				"zg-he",
				"pj-fs",
				"start-RW",
			},
			numPaths: 226,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g, err := parseGraph(tc.input)
			if err != nil {
				t.Fatalf("failed to build graph: %v", err)
			}

			paths := allPaths(g)

			if tc.paths != nil {
				pathsMap := pathsToMap(paths)
				tcPathsMap := pathsToMap(tc.paths)

				for p := range pathsMap {
					if _, found := tcPathsMap[p]; !found {
						t.Errorf("path in output, not in tc: %v", p)
					}
				}

				for p := range tcPathsMap {
					if _, found := pathsMap[p]; !found {
						t.Errorf("path in tc, not in output: %v", p)
					}
				}
			}

			wantNumPaths := tc.numPaths
			if wantNumPaths == 0 {
				wantNumPaths = len(tc.paths)
			}

			if len(paths) != wantNumPaths {
				t.Errorf("wanted %d paths, got %d",
					wantNumPaths, len(paths))
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
