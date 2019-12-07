package vm

import "testing"

func assertPanic(t *testing.T, msg string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Error(msg)
		}
	}()
	f()
}

func TestInput(t *testing.T) {
	io := NewIO(1, 2, 3, 4)

	got := []int{}
	for i := 0; i < 4; i++ {
		got = append(got, io.Read())
	}

	assertPanic(t, "read failed to panic", func() { io.Read() })
}
