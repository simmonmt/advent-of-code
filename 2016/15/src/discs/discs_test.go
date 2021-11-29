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

package discs

import "testing"

func TestSuccess(t *testing.T) {
	descs := []DiscDesc{
		DiscDesc{5, 4},
		DiscDesc{2, 1},
	}

	expectedResults := []bool{false, false, false, false, false, true}

	for tm, expected := range expectedResults {
		posns := make([]int, len(descs))
		for i := range descs {
			posns[i] = descs[i].Start + tm
		}

		if res := Success(descs, posns); res != expected {
			t.Errorf("at t=%v, Success(%v, %v) = %v, want %v", tm, descs, posns, expected)
		}
	}
}
