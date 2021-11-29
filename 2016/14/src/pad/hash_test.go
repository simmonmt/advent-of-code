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

package pad

import "testing"

func TestNormalHasher(t *testing.T) {
	hasher := NormalHasher{}
	expected := "23734cd52ad4a4fb877d8a1e26e5df5f"
	if res := hasher.MakeHash("abc", 1); res != expected {
		t.Errorf(`NormalHasher.MakeHash("abc", 1) = "%v", want "%v"`, res, expected)
	}
}

func TestStretchedHasher(t *testing.T) {
	hasher := StretchedHasher{}
	expected := "a107ff634856bb300138cac6568c0f24"
	if res := hasher.MakeHash("abc", 0); res != expected {
		t.Errorf(`StretchedHasher.MakeHash("abc", 1) = "%v", want "%v"`, res, expected)
	}
}
