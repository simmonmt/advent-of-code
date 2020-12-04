package pos

import (
	"strconv"
	"testing"
)

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

func TestP2(t *testing.T) {
	type TestCase struct {
		p1, p2            P2
		lessThan          bool
		manhattanDistance int
	}

	testCases := []TestCase{
		TestCase{P2{1, 1}, P2{4, 2}, true, 4},
		TestCase{P2{1, 1}, P2{1, 2}, true, 1},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := tc.p1.LessThan(tc.p2); got != tc.lessThan {
				t.Errorf("want %v < %v, got %v", tc.p1, tc.p2, got)
			}
			if got := tc.p2.LessThan(tc.p1); got == tc.lessThan {
				t.Errorf("want %v < %v, got %v", tc.p2, tc.p1, got)
			}
			if got := tc.p1.LessThan(tc.p1); got == true {
				t.Errorf("want %v !< %v, got %v", tc.p1, tc.p1, got)
			}

			if got := tc.p1.ManhattanDistance(tc.p2); got != tc.manhattanDistance {
				t.Errorf("manhattan distance %v, %v = %v, want %v",
					tc.p1, tc.p2, got, tc.manhattanDistance)
			}
			if got := tc.p2.ManhattanDistance(tc.p1); got != tc.manhattanDistance {
				t.Errorf("manhattan distance %v, %v = %v, want %v",
					tc.p2, tc.p1, got, tc.manhattanDistance)
			}

			p := tc.p1
			p.Add(P2{2, 1})
			want := P2{tc.p1.X + 2, tc.p1.Y + 1}
			if !p.Equals(want) {
				t.Errorf("%v.Add({2,1}) = %v, got %v", want, p)
			}
		})
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
