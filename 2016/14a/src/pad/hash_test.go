package pad

import "testing"

func TestNormalHasher(t *testing.T) {
	hasher := NormalHasher{}
	expected := "23734cd52ad4a4fb877d8a1e26e5df5f"
	if res := hasher.MakeHash("abc", 1); res != expected {
		t.Errorf(`NormalHasher.MakeHash("abc", 1) = "%v", want "%v"`, res, expected)
	}
}

func TestStretchedHasher(t *testing.T) {
	hasher := StretchedHasher{}
	expected := "a107ff634856bb300138cac6568c0f24"
	if res := hasher.MakeHash("abc", 0); res != expected {
		t.Errorf(`StretchedHasher.MakeHash("abc", 1) = "%v", want "%v"`, res, expected)
	}
}
