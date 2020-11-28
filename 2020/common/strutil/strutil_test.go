package strutil

import (
	"reflect"
	"testing"
)

func TestStringToInt64s(t *testing.T) {
	want := []int64{1, -2, 3}
	if got, err := StringToInt64s("1,-2,3"); err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("1,-2,3 = %v, %v, want %v, nil", got, err, want)
	}

	if _, err := StringToInt64s("1,bob,3"); err == nil {
		t.Errorf("1,bob,3 = _, %v, want _, err", err)
	}
}
