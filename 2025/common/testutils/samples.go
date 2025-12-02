package testutils

import (
	"bufio"
	"fmt"
	"iter"
	"path"
	"strings"

	"github.com/simmonmt/aoc/2025/common/logger"
)

type SampleTestCase struct {
	File         string
	Body         []string
	WantInput    any
	WantA, WantB any
}

type SampleIter struct {
	s *bufio.Scanner
}

func NewSampleIter(raw string) *SampleIter {
	return &SampleIter{
		s: bufio.NewScanner(strings.NewReader(raw)),
	}
}

func (si *SampleIter) All() iter.Seq2[string, []string] {
	nextLine := func() (string, bool) {
		if !si.s.Scan() {
			if err := si.s.Err(); err != nil {
				panic("bad read")
			}
			return "", false
		}

		return si.s.Text(), true
	}

	return func(yield func(string, []string) bool) {
		for {
			fullPath, ok := nextLine()
			if !ok {
				return
			}

			nlStr, ok := nextLine()
			if !ok {
				panic("missing num lines")
			}

			var numLines int
			if _, err := fmt.Sscanf(nlStr, "%d", &numLines); err != nil {
				panic("bad num lines")
			}

			out := make([]string, numLines)
			for i := range numLines {
				var ok bool
				if out[i], ok = nextLine(); !ok {
					panic("bad body")
				}
			}

			if !yield(fullPath, out) {
				return
			}
		}
	}
}

func PopulateTestCases(raw string, testCases []SampleTestCase) {
	i := 0
	for fullPath, body := range NewSampleIter(raw).All() {
		if i >= len(testCases) {
			logger.Fatalf("too many samples; have %d cases",
				len(testCases))
		}

		tc := &testCases[i]
		tc.File = path.Base(fullPath)
		tc.Body = body

		i++
	}

	if i != len(testCases) {
		logger.Fatalf("not enough samples; have %d cases",
			len(testCases))
	}
}
