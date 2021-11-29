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

package node

import (
	"reflect"
	"testing"
)

func TestSerialize(t *testing.T) {
	type TestCase struct {
		ser   string
		deser *Node
	}

	testCases := []TestCase{
		TestCase{"0,0,pass,", &Node{0, 0, "pass", ""}},
		TestCase{"1,2,pass,UD", &Node{1, 2, "pass", "UD"}},
		TestCase{"-1,-2,pass,UD", &Node{-1, -2, "pass", "UD"}},
	}

	for _, tc := range testCases {
		t.Run(tc.ser, func(t *testing.T) {
			if deser, err := Deserialize(tc.ser); err == nil && !reflect.DeepEqual(deser, tc.deser) {
				t.Errorf(`Deserialize("%v") = %+v, %v, want %+v, nil`, tc.ser, deser, err, tc.deser)
			}

			if ser := tc.deser.Serialize(); ser != tc.ser {
				t.Errorf(`%+v Serialize() = "%v", want "%v"`, tc.deser, ser, tc.ser)
			}
		})
	}
}
