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
	"github.com/simmonmt/aoc/2021/common/pos"
)

func TestBox(t *testing.T) {
	type TestCase struct {
		p1, p2 pos.P2
		ins    []pos.P2
		outs   []pos.P2
	}

	testCases := []TestCase{
		TestCase{
			p1: pos.P2{X: 3, Y: 6},
			p2: pos.P2{X: 8, Y: 10},
			ins: []pos.P2{
				pos.P2{X: 4, Y: 7},
			},
			outs: []pos.P2{
				pos.P2{X: 0, Y: 0},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b := NewBox(tc.p1, tc.p2)

			ins := []pos.P2{tc.p1, tc.p2}
			ins = append(ins, tc.ins...)

			for _, in := range ins {
				if got := b.Contains(in); !got {
					t.Errorf("Containts(%v) = %v, want true", in, got)
				}
			}

			for _, out := range tc.outs {
				if got := b.Contains(out); got {
					t.Errorf("Containts(%v) = %v, want false", out, got)
				}
			}
		})
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
