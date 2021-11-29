// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package board

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/simmonmt/aoc/2020/common/strutil"
)

func TestDumpZW(t *testing.T) {
	spec := []string{".#.", "..#", "###"}

	b := New(spec, false)

	buf := &bytes.Buffer{}

	want := strings.Join(spec, "\n") + "\n"
	b.DumpZW(0, 0, buf)
	got := buf.String()
	if want != got {
		t.Errorf("DumpZ(0, 0, _) = %v, want %v",
			strconv.Quote(got), strconv.Quote(want))
	}
}

func TestEvolve3D(t *testing.T) {
	spec := []string{".#.", "..#", "###"}
	b := New(spec, false)

	nb := b.Evolve()

	if min, max := nb.ZBounds(); min != -1 || max != 1 {
		t.Errorf("nb.ZBounds() = %v, %v, want -1, 1", min, max)
	}

	if min, max := nb.WBounds(); min != 0 || max != 0 {
		t.Errorf("nb.WBounds() = %v, %v, want 0, 0", min, max)
	}

	wants := map[int][]string{
		-1: []string{"#..", "..#", ".#."},
		0:  []string{"#.#", ".##", ".#."},
		1:  []string{"#..", "..#", ".#."},
	}

	for z, want := range wants {
		buf := &bytes.Buffer{}
		nb.DumpZW(z, 0, buf)
		got := buf.String()

		wantStr := strings.Join(want, "\n") + "\n"

		if diff, msg := strutil.StringDiff(wantStr, got); diff {
			t.Errorf("nb.DumpZW(%v, 0, _) mismatch: %v\nwant\n%v\ngot\n%v",
				z, msg, wantStr, got)
		}
	}
}

func TestEvolve4D(t *testing.T) {
	spec := []string{".#.", "..#", "###"}
	b := New(spec, true)

	nb := b.Evolve()

	if min, max := nb.ZBounds(); min != -1 || max != 1 {
		t.Errorf("nb.ZBounds() = %v, %v, want -1, 1", min, max)
	}

	if min, max := nb.WBounds(); min != -1 || max != 1 {
		t.Errorf("nb.WBounds() = %v, %v, want -1, 1", min, max)
	}

	type ZW struct {
		Z, W int
	}

	wants := map[ZW][]string{
		ZW{-1, -1}: []string{"#..", "..#", ".#."},
		ZW{0, -1}:  []string{"#..", "..#", ".#."},
		ZW{1, -1}:  []string{"#..", "..#", ".#."},
		ZW{-1, 0}:  []string{"#..", "..#", ".#."},
		ZW{0, 0}:   []string{"#.#", ".##", ".#."},
		ZW{1, 0}:   []string{"#..", "..#", ".#."},
		ZW{-1, 1}:  []string{"#..", "..#", ".#."},
		ZW{0, 1}:   []string{"#..", "..#", ".#."},
		ZW{1, 1}:   []string{"#..", "..#", ".#."},
	}

	for zw, want := range wants {
		buf := &bytes.Buffer{}
		nb.DumpZW(zw.Z, zw.W, buf)
		got := buf.String()

		wantStr := strings.Join(want, "\n") + "\n"

		if diff, msg := strutil.StringDiff(wantStr, got); diff {
			t.Errorf("nb.DumpZW(%v, _) mismatch: %v\nwant\n%v\ngot\n%v",
				zw, msg, wantStr, got)
		}
	}
}
