package pad

import "testing"

func TestMakeHash(t *testing.T) {
	expected := "23734cd52ad4a4fb877d8a1e26e5df5f"
	if res := MakeHash("abc", 1); res != expected {
		t.Errorf(`MakeHash("abc", 1) = "%v", want "%v"`, res, expected)
	}
}
