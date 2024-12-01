package testutils

import (
	"bufio"
	"iter"
	"path"
	"strconv"
	"strings"

	"github.com/simmonmt/aoc/2024/common/logger"
)

type SampleTestCase struct {
	File         string
	Body         []string
	WantInput    any
	WantA, WantB int64
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

			numLines := -1
			var err error
			if nlStr, ok := nextLine(); !ok {
				panic("missing num lines")
			} else if numLines, err = strconv.Atoi(nlStr); err != nil {
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
