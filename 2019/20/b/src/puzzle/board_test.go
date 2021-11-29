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

package puzzle

import (
	"reflect"
	"testing"

	"github.com/simmonmt/aoc/2019/common/pos"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard(map1)

	want := []Gate{
		Gate{
			name: "AA",
			pOut: pos.P2{9, 2},
			gOut: pos.P2{9, 1},
			pIn:  pos.P2{-1, -1},
			gIn:  pos.P2{-1, -1},
		},
		Gate{
			name: "BC",
			pOut: pos.P2{2, 8},
			gOut: pos.P2{1, 8},
			pIn:  pos.P2{9, 6},
			gIn:  pos.P2{9, 7},
		},
		Gate{
			name: "DE",
			pOut: pos.P2{2, 13},
			gOut: pos.P2{1, 13},
			pIn:  pos.P2{6, 10},
			gIn:  pos.P2{7, 10},
		},
		Gate{
			name: "FG",
			pOut: pos.P2{2, 15},
			gOut: pos.P2{1, 15},
			pIn:  pos.P2{11, 12},
			gIn:  pos.P2{11, 11},
		},
		Gate{
			name: "ZZ",
			pOut: pos.P2{13, 16},
			gOut: pos.P2{13, 17},
			pIn:  pos.P2{-1, -1},
			gIn:  pos.P2{-1, -1},
		},
	}

	if got := b.Gates(); !reflect.DeepEqual(got, want) {
		t.Errorf("Gates() = %v, want %v", got, want)
	}
}
