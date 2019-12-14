package pos

import "testing"

func TestP2FromStringOK(t *testing.T) {
	want := P2{1, -2}
	if got, err := P2FromString("1,-2"); err != nil || !got.Equals(want) {
		t.Errorf(`P2FromString("1,-2") = %v, %v, want %v, nil`,
			got, err, want)
	}
}

func TestP2FromStringBad(t *testing.T) {
	if _, err := P2FromString("1,bob"); err == nil {
		t.Errorf(`P2FromString("1,bob") = _, %v, want _, non-nil`, err)
	}
}

func TestP3FromStringOK(t *testing.T) {
	want := P3{1, -2, 3}
	if got, err := P3FromString("1,-2,3"); err != nil || !got.Equals(want) {
		t.Errorf(`P3FromString("1,-2,3") = %v, %v, want %v, nil`,
			got, err, want)
	}
}

func TestP3FromStringBad(t *testing.T) {
	if _, err := P3FromString("1,bob,3"); err == nil {
		t.Errorf(`P3FromString("1,bob,3") = _, %v, want _, non-nil`, err)
	}
}
