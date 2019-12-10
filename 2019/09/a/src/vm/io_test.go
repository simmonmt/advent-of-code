package vm

import (
	"testing"

	"github.com/simmonmt/aoc/2019/09/a/src/testutils"
)

func TestInput(t *testing.T) {
	io := NewIO(1, 2, 3, 4)

	got := []int{}
	for i := 0; i < 4; i++ {
		got = append(got, io.Read())
	}

	testutils.AssertPanic(t, "read failed to panic", func() { io.Read() })
}
