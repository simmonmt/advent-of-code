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
