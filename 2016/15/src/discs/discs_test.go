package discs

import "testing"

func TestSuccess(t *testing.T) {
	descs := []DiscDesc{
		DiscDesc{5, 4},
		DiscDesc{2, 1},
	}

	expectedResults := []bool{false, false, false, false, false, true}

	for tm, expected := range expectedResults {
		posns := make([]int, len(descs))
		for i := range descs {
			posns[i] = descs[i].Start + tm
		}

		if res := Success(descs, posns); res != expected {
			t.Errorf("at t=%v, Success(%v, %v) = %v, want %v", tm, descs, posns, expected)
		}
	}
}
