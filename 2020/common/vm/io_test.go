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
