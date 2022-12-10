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

func TestMin(t *testing.T) {
	if got, want := Min(3, 1), 1; got != want {
		t.Errorf("Min(3,1) = %v, want %v", got, want)
	}
	if got, want := Min(1, 3), 1; got != want {
		t.Errorf("Min(1,3) = %v, want %v", got, want)
	}
	if got, want := Min(3, 3), 3; got != want {
		t.Errorf("Min(3,3) = %v, want %v", got, want)
	}

	if got, want := Min(3.5, 1.5), 1.5; got != want {
		t.Errorf("Min(3.5,1.5) = %v, want %v", got, want)
	}
	if got, want := Min(1.5, 3.5), 1.5; got != want {
		t.Errorf("Min(1.5,3.5) = %v, want %v", got, want)
	}
	if got, want := Min(3.5, 3.5), 3.5; got != want {
		t.Errorf("Min(3.5,3.5) = %v, want %v", got, want)
	}
}

func TestMax(t *testing.T) {
	if got, want := Max(3, 1), 3; got != want {
		t.Errorf("Max(3,1) = %v, want %v", got, want)
	}
	if got, want := Max(1, 3), 3; got != want {
		t.Errorf("Max(1,3) = %v, want %v", got, want)
	}
	if got, want := Max(3, 3), 3; got != want {
		t.Errorf("Max(3,3) = %v, want %v", got, want)
	}

	if got, want := Max(3.5, 1.5), 3.5; got != want {
		t.Errorf("Max(3.5,1.5) = %v, want %v", got, want)
	}
	if got, want := Max(1.5, 3.5), 3.5; got != want {
		t.Errorf("Max(1.5,3.5) = %v, want %v", got, want)
	}
	if got, want := Max(3.5, 3.5), 3.5; got != want {
		t.Errorf("Max(3.5,3.5) = %v, want %v", got, want)
	}
}
