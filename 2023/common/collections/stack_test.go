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

package collections

import (
	"testing"

	"github.com/simmonmt/aoc/2023/common/testutils"
)

func TestStack(t *testing.T) {
	s := NewStack()

	if !s.Empty() {
		t.Errorf("s.Empty() = false, want true")
	}

	s.Push(4)
	if got := s.Peek(); got.(int) != 4 {
		t.Errorf("s.Peek() = %v, want 4", got)
	}
	if s.Empty() {
		t.Errorf("s.Empty() = true, want false")
	}

	s.Push(5)
	if got := s.Peek(); got.(int) != 5 {
		t.Errorf("s.Peek() = %v, want 5", got)
	}

	if got := s.Pop(); got.(int) != 5 {
		t.Errorf("s.Pop() = %v, want 5", got)
	}

	if got := s.Peek(); got.(int) != 4 {
		t.Errorf("s.Peek() = %v, want 4", got)
	}

	if got := s.Pop(); got.(int) != 4 {
		t.Errorf("s.Pop() = %v, want 4", got)
	}

	testutils.AssertPanic(t, "pop", func() { s.Pop() })
	testutils.AssertPanic(t, "peek", func() { s.Peek() })

	// Test resizing
	for i := 0; i < 100; i++ {
		s.Push(i)
	}
	for i := 99; i >= 0; i-- {
		if got := s.Pop().(int); got != i {
			t.Errorf("s.Pop() = %v, want %v", got, i)
		}
	}
}
