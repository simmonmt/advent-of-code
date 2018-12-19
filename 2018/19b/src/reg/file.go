package reg

import (
	"regexp"

	"intmath"
)

var (
	regFilePattern = regexp.MustCompile(`^[^ ]+ +\[(\d+), (\d+), (\d+), (\d+)\]$`)
)

type File [6]int

func ParseFile(str string) *File {
	parts := regFilePattern.FindStringSubmatch(str)
	if parts == nil {
		panic("bad parse: " + str)
	}

	regFile := &File{}
	for i, s := range parts[1:] {
		regFile[i] = intmath.AtoiOrDie(s)
	}

	return regFile
}
