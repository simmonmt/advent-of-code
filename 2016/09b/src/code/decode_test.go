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

package code

import "testing"

func TestDecode(t *testing.T) {
	type TestCase struct {
		in, expected string
	}

	testCases := []TestCase{
		TestCase{"ADVENT", "ADVENT"},
		TestCase{"A(1x5)BC", "ABBBBBC"},
		TestCase{"(3x3)XYZ", "XYZXYZXYZ"},
		TestCase{"A(2x2)BCD(2x2)EFG", "ABCBCDEFEFG"},
		TestCase{"(6x1)(1x3)A", "(1x3)A"},
		TestCase{"X(8x2)(3x3)ABCY", "X(3x3)ABC(3x3)ABCY"},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			out, err := Decode(tc.in)
			if err != nil || tc.expected != out {
				t.Errorf(`Decode("%v") = "%v", %v, want "%v", nil`, tc.in, out, err, tc.expected)
			}
		})
	}
}

func TestDecodeLen(t *testing.T) {
	type TestCase struct {
		in          string
		expectedLen int
	}

	testCases := []TestCase{
		TestCase{"ADVENT", 6},
		TestCase{"AA(1x30)A", 32},
		TestCase{"(3x3)XYZ", 9},
		TestCase{"X(8x2)(3x3)ABCY", 20},
		TestCase{"(27x12)(20x12)(13x14)(7x10)(1x12)A", 241920},
		TestCase{"(25x3)(3x3)ABC(2x3)XY(5x2)PQRSTX(18x9)(3x2)TWO(5x7)SEVEN", 445},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			decodedLen, err := DecodeLen(tc.in)
			if err != nil || tc.expectedLen != decodedLen {
				t.Errorf(`DecodedLen("%v") = %v, %v, want %v, nil`, tc.in, decodedLen, err, tc.expectedLen)
			}
		})
	}
}
