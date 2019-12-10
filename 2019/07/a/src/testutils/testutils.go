package testutils

import "testing"

func AssertPanic(t *testing.T, msg string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Error(msg)
		}
	}()
	f()
}
