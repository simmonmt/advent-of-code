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

package vm

import (
	"testing"

	"github.com/simmonmt/aoc/2020/common/testutils"
)

func TestInput(t *testing.T) {
	io := NewSaverIO(1, 2, 3, 4)

	got := []int64{}
	for i := int64(0); i < 4; i++ {
		got = append(got, io.Read())
	}

	testutils.AssertPanic(t, "read failed to panic", func() { io.Read() })
}
