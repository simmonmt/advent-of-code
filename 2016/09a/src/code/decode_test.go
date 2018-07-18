package code

import "testing"

func TestDecode(t *testing.T) {
	type TestCase struct {
		in, expected string
	}

	testCases := []TestCase{
		TestCase{"ADVENT", "ADVENT"},
		TestCase{"A(1x5)BC", "ABBBBBC"},
		TestCase{"(3x3)XYZ", "XYZXYZXYZ"},
		TestCase{"A(2x2)BCD(2x2)EFG", "ABCBCDEFEFG"},
		TestCase{"(6x1)(1x3)A", "(1x3)A"},
		TestCase{"X(8x2)(3x3)ABCY", "X(3x3)ABC(3x3)ABCY"},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			out, err := Decode(tc.in)
			if err != nil || tc.expected != out {
				t.Errorf(`Decode("%v") = "%v", %v, want "%v", nil`, tc.in, out, err, tc.expected)
			}
		})
	}
}
