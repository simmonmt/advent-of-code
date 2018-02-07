package main

import (
	"fmt"
	"log"
	"os"
)

const (
	InvalidNone = iota
	InvalidNoStraight
	InvalidBadChar
	InvalidNoTwoPairs
)

func ReasonToString(reason int) string {
	switch reason {
	case InvalidNone:
		return "None"
	case InvalidNoStraight:
		return "NoStraight"
	case InvalidBadChar:
		return "BadChar"
	case InvalidNoTwoPairs:
		return "NoTwoPairs"
	default:
		panic(fmt.Sprintf("unknown reason %v", reason))
	}
}

type Password [8]rune

func (p *Password) ToString() string {
	return string(p[0:8])
}

func Validate(pw Password) (bool, int) {
	firstPairStart, secondPairStart := -1, -1
	foundStraight := false

	for i, c := range pw {
		if c == 'i' || c == 'o' || c == 'l' {
			return false, InvalidBadChar
		}

		if secondPairStart == -1 {
			if i > 0 && c == pw[i-1] {
				if firstPairStart == -1 {
					firstPairStart = i - 1
				} else {
					if i-1 != firstPairStart+1 {
						secondPairStart = i - 1
					}
				}
			}
		}

		if i > 1 && !foundStraight {
			if pw[i-2]+1 == pw[i-1] && pw[i-1]+1 == c {
				foundStraight = true
			}
		}
	}

	if secondPairStart == -1 {
		return false, InvalidNoTwoPairs
	}
	if !foundStraight {
		return false, InvalidNoStraight
	}

	return true, 0
}

func Add(pw *Password, addend [8]uint8) {
	for i := 7; i >= 0; i-- {
		r := int(pw[i]) - 'a' + int(addend[i])
		pw[i] = rune('a' + r%26)
		if i > 0 && r > 25 {
			pw[i-1] = rune(int(pw[i-1]) + r/26)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v seed", os.Args[0])
	}

	if len(os.Args[1]) != 8 {
		log.Fatalf("seed must be 8 chars")
	}
	var pw Password
	for i, c := range os.Args[1] {
		if c < 'a' || c > 'z' {
			log.Fatalf("seed chars must be a-z")
		}
		pw[i] = c
	}

	for {
		Add(&pw, [8]uint8{0, 0, 0, 0, 0, 0, 0, 1})

		if valid, _ := Validate(pw); valid {
			fmt.Println(pw.ToString())
			break
		}
	}
}
