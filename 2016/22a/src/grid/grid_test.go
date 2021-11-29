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

package grid

import (
	"bytes"
	"reflect"
	"strconv"
	"testing"

	"node"
)

func TestSerializeNode(t *testing.T) {
	type TestCase struct {
		n   *node.Node
		ser []byte
	}

	testCases := []TestCase{
		TestCase{node.New(88, 69), []byte{88, 69}},
		TestCase{node.New(507, 69), []byte{0x81, 0xfb, 69}},
		TestCase{node.New(69, 507), []byte{69, 0x81, 0xfb}},
		TestCase{node.New(507, 490), []byte{0x81, 0xfb, 0x81, 0xea}},
	}

	for i, tc := range testCases {
		t.Run("serialize_"+strconv.Itoa(i), func(t *testing.T) {
			out := make([]byte, 256)
			serLen := serializeNode(tc.n, out)

			if serLen == 0 {
				t.Errorf("serializeNode(%v) == 0, want %v", tc.n, tc.ser)
			} else if serLen != len(tc.ser) || !bytes.Equal(out[0:serLen], tc.ser) {
				t.Errorf("serializeNode(%v) = %v out %v, want %v out %v",
					tc.n, serLen, out[0:serLen], len(tc.ser), tc.ser)
			}
		})
	}

	for i, tc := range testCases {
		t.Run("deserialize_"+strconv.Itoa(i), func(t *testing.T) {
			n, consumed := deserializeNode(tc.ser)

			if n == nil || consumed != len(tc.ser) || !reflect.DeepEqual(n, tc.n) {
				t.Errorf("deserializeNode(%v) = %+v, %v, want %+v, %v",
					tc.ser, n, consumed, tc.n, len(tc.ser))
			}
		})
	}
}

func makeGrid(t *testing.T, w, h int, goalX, goalY uint8, nodes []node.Node) *Grid {
	g, err := New(w, h, goalX, goalY, nodes)
	if err != nil {
		panic(err.Error())
	}
	return g
}

func TestSerializeGrid(t *testing.T) {
	type TestCase struct {
		g   *Grid
		ser []byte
	}

	testCases := []TestCase{
		TestCase{
			g: makeGrid(t, 2, 2, 2, 3, []node.Node{
				*node.New(88, 69), *node.New(507, 490),
				*node.New(32, 12), *node.New(85, 65),
			}),
			ser: []byte{2, 3, 88, 69, 0x81, 0xfb, 0x81, 0xea, 32, 12, 85, 65},
		},
	}

	for i, tc := range testCases {
		t.Run("serialize_"+strconv.Itoa(i), func(t *testing.T) {
			if ser := tc.g.Serialize(); !bytes.Equal([]byte(ser), tc.ser) {
				t.Errorf("Serialize = %v, want %v", []byte(ser), tc.ser)
			}
		})
	}

	for i, tc := range testCases {
		t.Run("deserialize_"+strconv.Itoa(i), func(t *testing.T) {
			if g, err := Deserialize(tc.g.w, tc.g.h, tc.ser); err != nil || !reflect.DeepEqual(g, tc.g) {
				t.Errorf("Deserialize(%v,%v,%v) = %+v,%v, want %+v,nil",
					tc.g.w, tc.g.h, tc.ser, g, err, tc.g)
			}
		})
	}
}
