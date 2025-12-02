package testutils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSampleIter(t *testing.T) {
	raw := "one.txt\n1\na\ntwo.txt\n2\nb\nc\n"

	type Sample struct {
		Path string
		Body []string
	}

	got := []Sample{}

	for path, body := range NewSampleIter(raw).All() {
		got = append(got, Sample{path, body})
	}

	want := []Sample{
		Sample{"one.txt", []string{"a"}},
		Sample{"two.txt", []string{"b", "c"}},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("parseInput mismatch; -want,+got:\n%s\n", diff)
	}
}
