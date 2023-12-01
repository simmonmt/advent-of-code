// Copyright 2023 Google LLC
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

package strutil

import (
	"reflect"
	"testing"
)

func TestStringToInt64s(t *testing.T) {
	want := []int64{1, -2, 3}
	if got, err := StringToInt64s("1,-2,3"); err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("1,-2,3 = %v, %v, want %v, nil", got, err, want)
	}

	if _, err := StringToInt64s("1,bob,3"); err == nil {
		t.Errorf("1,bob,3 = _, %v, want _, err", err)
	}
}
