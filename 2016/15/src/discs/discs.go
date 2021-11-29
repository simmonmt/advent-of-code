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

type DiscDesc struct {
	NPos  int
	Start int
}

func Advance(posns []int) {
	for i := range posns {
		posns[i]++
	}
}

func Success(descs []DiscDesc, posns []int) bool {
	for i := range posns {
		if (posns[i]+i+1)%descs[i].NPos != 0 {
			return false
		}
	}

	return true
}
