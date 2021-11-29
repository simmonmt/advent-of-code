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
	"reflect"
	"testing"
)

func TestUniquify(t *testing.T) {
	in := []Pos{Pos{0, 0}, Pos{1, 1}, Pos{2, 2}, Pos{1, 1}}
	expected := []Pos{Pos{0, 0}, Pos{1, 1}, Pos{2, 2}}

	if got := uniquify(in); !reflect.DeepEqual(got, expected) {
		t.Errorf("uniquify(%v) = %v, want %v", in, got, expected)
	}
}
