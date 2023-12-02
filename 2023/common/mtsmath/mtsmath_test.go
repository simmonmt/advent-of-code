package mtsmath

import "testing"

func TestAbs(t *testing.T) {
	if got, want := Abs(1), 1; got != want {
		t.Errorf("Abs(1) = %v, want %v", got, want)
	}
	if got, want := Abs(-1), 1; got != want {
		t.Errorf("Abs(1) = %v, want %v", got, want)
	}
	if got, want := Abs(0), 0; got != want {
		t.Errorf("Abs(1) = %v, want %v", got, want)
	}

	if got, want := Abs(1.5), 1.5; got != want {
		t.Errorf("Abs(1.5) = %v, want %v", got, want)
	}
	if got, want := Abs(-1.5), 1.5; got != want {
		t.Errorf("Abs(1.5) = %v, want %v", got, want)
	}
}
