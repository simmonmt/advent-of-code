package main

import (
	"reflect"
	"testing"
)

func TestUniquify(t *testing.T) {
	in := []Pos{Pos{0, 0}, Pos{1, 1}, Pos{2, 2}, Pos{1, 1}}
	expected := []Pos{Pos{0, 0}, Pos{1, 1}, Pos{2, 2}}

	if got := uniquify(in); !reflect.DeepEqual(got, expected) {
		t.Errorf("uniquify(%v) = %v, want %v", in, got, expected)
	}
}
