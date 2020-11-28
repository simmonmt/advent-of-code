package collections

import (
	"testing"

	"github.com/simmonmt/aoc/2020/common/testutils"
)

func TestStack(t *testing.T) {
	s := NewStack()

	if !s.Empty() {
		t.Errorf("s.Empty() = false, want true")
	}

	s.Push(4)
	if got := s.Peek(); got.(int) != 4 {
		t.Errorf("s.Peek() = %v, want 4")
	}
	if s.Empty() {
		t.Errorf("s.Empty() = true, want false")
	}

	s.Push(5)
	if got := s.Peek(); got.(int) != 5 {
		t.Errorf("s.Peek() = %v, want 5")
	}

	if got := s.Pop(); got.(int) != 5 {
		t.Errorf("s.Pop() = %v, want 5")
	}

	if got := s.Peek(); got.(int) != 4 {
		t.Errorf("s.Peek() = %v, want 4")
	}

	if got := s.Pop(); got.(int) != 4 {
		t.Errorf("s.Pop() = %v, want 4")
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
