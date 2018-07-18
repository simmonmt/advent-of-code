package code

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	repeatPattern = regexp.MustCompile(`([0-9]+)x([0-9]+)\)`)
)

func parseNum(str string) (int, error) {
	n, err := strconv.ParseUint(str, 10, 32)
	return int(n), err
}

func nextRepeat(chars []byte) (found bool, patternStartIdx, patternLen, numChars, numReps int) {
	found = false

	for i := 0; i < len(chars); i++ {
		c := chars[i]

		if c != '(' {
			continue
		}

		if i+1 == len(chars) {
			return
		}

		matches := repeatPattern.FindSubmatch(chars[i+1:])
		if matches == nil {
			fmt.Println("no match")
			return
		}

		var err error
		numCharsStr := string(matches[1])
		numChars, err = parseNum(numCharsStr)
		if err != nil {
			fmt.Println("bad numchars")
			return
		}

		numRepsStr := string(matches[2])
		numReps, err = parseNum(numRepsStr)
		if err != nil {
			fmt.Println("bad numreps")
			return
		}

		patternLen = 1 + len(matches[0]) // 1+ for leading (
		nextIdx := i + patternLen
		nextEnd := nextIdx + int(numChars)
		if nextEnd > len(chars) {
			fmt.Printf("too long\n")
			return
		}

		return true, i, patternLen, int(numChars), int(numReps)
	}

	return
}

func Decode(in string) (string, error) {
	chars := []byte(in)
	out := []byte{}
	for i := 0; i < len(chars); {
		found, patternStartRelIdx, patternLen, numChars, numReps := nextRepeat(chars[i:])

		if !found {
			out = append(out, chars[i:]...)
			break
		}

		if patternStartRelIdx > 0 {
			out = append(out, chars[i:i+patternStartRelIdx]...)
		}
		i += patternStartRelIdx

		if i+patternLen+numChars > len(chars) {
			return "", fmt.Errorf("repeat %vx%v reaches past end",
				numChars, numReps)
		}

		for j := 0; j < numReps; j++ {
			out = append(out, chars[i+patternLen:i+patternLen+numChars]...)
		}

		i += patternLen + numChars
	}
	return string(out), nil
}

func DecodeLen(in string) (int, error) {
	chars := []byte(in)

	numDecoded := 0
	for i := 0; i < len(chars); {
		found, patternStartRelIdx, patternLen, numChars, numReps := nextRepeat(chars[i:])

		// fmt.Printf("in %v found %v startrel %v len %v numchars %v numreps %v\n",
		// 	string(chars), found, patternStartRelIdx, patternLen, numChars, numReps)

		if !found {
			numDecoded += len(chars[i:])
			break
		}

		if patternStartRelIdx > 0 {
			numDecoded += patternStartRelIdx
			i += patternStartRelIdx
		}

		if i+patternLen+numChars > len(chars) {
			return 0, fmt.Errorf("repeat %vx%v reaches past end",
				numChars, numReps)
		}

		repChars := chars[i+patternLen : i+patternLen+numChars]
		decodedRepCharLen, err := DecodeLen(string(repChars))
		if err != nil {
			return 0, err
		}

		numDecoded += decodedRepCharLen * numReps
		i += patternLen + len(repChars)
	}

	// fmt.Printf("returning num decoded %v\n", numDecoded)
	return numDecoded, nil
}
