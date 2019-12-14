package vm

import (
	"reflect"
	"strings"
	"testing"
)

func TestRam(t *testing.T) {
	ram := NewRam()
	ram.Write(4, 3)

	if got := ram.Read(4); got != 3 {
		t.Errorf("Read(4) = %v, want 3", got)
	}

	// uninitialized
	if got := ram.Read(3); got != 0 {
		t.Errorf("Read(3) = %v, want 0", got)
	}
}

func TestRamWithData(t *testing.T) {
	ram := NewRam(0, 0, 10, 11, 12, 13)
	ram.Write(4, 99)

	got := []int64{}
	for i := int64(0); i <= 6; i++ {
		got = append(got, ram.Read(i))
	}

	want := []int64{0, 0, 10, 11, 99, 13, 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Ram got %v, want %v", got, want)
	}
}

func TestClone(t *testing.T) {
	vals := []int64{0, 0, 10, 11, 12, 13}
	ram := NewRam(vals...)
	clone := ram.Clone()

	for addr, val := range vals {
		if got := clone.Read(int64(addr)); got != val {
			t.Errorf("clone %v = %v, want %v", addr, got, val)
		}
	}
}

func TestRamFromReader(t *testing.T) {
	in := "10,11,12,13"
	ram, err := NewRamFromReader(strings.NewReader(in))
	if err != nil {
		t.Errorf("NewRamFromReader() = _, %v, want _, nil", err)
		return
	}

	CheckRam(t, ram, []int64{10, 11, 12, 13})
}
