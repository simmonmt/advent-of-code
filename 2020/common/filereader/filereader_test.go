package filereader

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestLinesFromReader(t *testing.T) {
	in := "a\nb\n\nd\n"
	r := strings.NewReader(in)
	want := []string{"a", "b", "", "d"}
	if got, err := linesFromReader(r); err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf(`linesFromReader("%v") = %v, %v, want %v, nil`,
			strconv.Quote(in), got, err, want)
	}
}

func TestBlankSeparatedGroups(t *testing.T) {
	in := "a\nb\n\nd\n"
	r := strings.NewReader(in)
	want := [][]string{[]string{"a", "b"}, []string{"d"}}
	if got, err := blankSeparatedGroupsFromReader(r); err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf(`blankSeparatedGroupsFromReader("%v") = %v, %v, want %v, nil`,
			strconv.Quote(in), got, err, want)
	}
}
