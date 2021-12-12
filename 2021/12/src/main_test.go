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
		input                    []string
		pathsOne, pathsTwo       []string
		numPathsOne, numPathsTwo int
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
			pathsOne: []string{
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
			pathsTwo: []string{
				"start,A,b,A,b,A,c,A,end",
				"start,A,b,A,b,A,end",
				"start,A,b,A,b,end",
				"start,A,b,A,c,A,b,A,end",
				"start,A,b,A,c,A,b,end",
				"start,A,b,A,c,A,c,A,end",
				"start,A,b,A,c,A,end",
				"start,A,b,A,end",
				"start,A,b,d,b,A,c,A,end",
				"start,A,b,d,b,A,end",
				"start,A,b,d,b,end",
				"start,A,b,end",
				"start,A,c,A,b,A,b,A,end",
				"start,A,c,A,b,A,b,end",
				"start,A,c,A,b,A,c,A,end",
				"start,A,c,A,b,A,end",
				"start,A,c,A,b,d,b,A,end",
				"start,A,c,A,b,d,b,end",
				"start,A,c,A,b,end",
				"start,A,c,A,c,A,b,A,end",
				"start,A,c,A,c,A,b,end",
				"start,A,c,A,c,A,end",
				"start,A,c,A,end",
				"start,A,end",
				"start,b,A,b,A,c,A,end",
				"start,b,A,b,A,end",
				"start,b,A,b,end",
				"start,b,A,c,A,b,A,end",
				"start,b,A,c,A,b,end",
				"start,b,A,c,A,c,A,end",
				"start,b,A,c,A,end",
				"start,b,A,end",
				"start,b,d,b,A,c,A,end",
				"start,b,d,b,A,end",
				"start,b,d,b,end",
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
			pathsOne: []string{
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
			numPathsTwo: 103,
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
			numPathsOne: 226,
			numPathsTwo: 3509,
		},
	}

	for i, tc := range testCases {
		g, err := parseGraph(tc.input)
		if err != nil {
			t.Fatalf("failed to build graph: %v", err)
		}

		t.Run(strconv.Itoa(i)+"_1", func(t *testing.T) {
			runAllPathsTest(t, g, 1, tc.pathsOne, tc.numPathsOne)
		})

		t.Run(strconv.Itoa(i)+"_2", func(t *testing.T) {
			runAllPathsTest(t, g, 2, tc.pathsTwo, tc.numPathsTwo)
		})
	}
}

func runAllPathsTest(t *testing.T, g *Graph, smallMax int, tcPaths []string, tcNumPaths int) {
	paths := allPaths(g, smallMax)

	if tcPaths != nil {
		pathsMap := pathsToMap(paths)
		tcPathsMap := pathsToMap(tcPaths)

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

	if tcNumPaths == 0 {
		tcNumPaths = len(tcPaths)
	}

	if len(paths) != tcNumPaths {
		t.Errorf("wanted %d paths, got %d", tcNumPaths, len(paths))
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
