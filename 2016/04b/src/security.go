package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	pattern = regexp.MustCompile(`^([a-z-]+)-([0-9]+)\[([a-z]+)\]$`) //\[[a-z]\])$`)
)

type Room struct {
	Name     string
	Sector   int
	Checksum string
}

func readInput(r io.Reader) ([]Room, error) {
	rooms := []Room{}

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		matches := pattern.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("%d: unexpected format", lineNum)
		}

		name := matches[1]
		sectorStr := matches[2]
		cksum := matches[3]

		sector, err := strconv.ParseUint(sectorStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%d: invalid sector %v: %v", lineNum, sectorStr, err)
		}

		rooms = append(rooms, Room{Name: name, Sector: int(sector), Checksum: cksum})
	}

	return rooms, nil
}

type CharCount struct {
	C string
	N int
}

type CharCounts []*CharCount

func (cc CharCounts) Less(i, j int) bool {
	if cc[i].N < cc[j].N {
		return false
	} else if cc[i].N > cc[j].N {
		return true
	} else {
		return cc[i].C < cc[j].C
	}
}

func (cc CharCounts) Len() int { return len(cc) }

func (cc CharCounts) Swap(i, j int) {
	cc[i], cc[j] = cc[j], cc[i]
}

func countChars(s string) []*CharCount {
	charToCount := map[string]int{}
	for _, r := range s {
		charToCount[string(r)]++
	}
	delete(charToCount, "-")

	counts := []*CharCount{}
	for c, n := range charToCount {
		counts = append(counts, &CharCount{C: c, N: n})
	}

	sort.Sort(CharCounts(counts))
	return counts
}

func verifyChecksum(charCounts []*CharCount, cksum string) bool {
	i := 0
	for _, r := range cksum {
		if i == len(charCounts) {
			//fmt.Printf("checksum too long\n")
			return false // checksum too long
		}

		s := string(r)
		if charCounts[i].C != s {
			//fmt.Printf("mismatch; wanted %v, got %v\n", charCounts[i].C, s)
			return false // mismatch
		}

		i++
	}

	return true
}

func decrypt(in string, sector int) string {
	out := []rune{}
	for _, r := range in {
		if r == '-' {
			out = append(out, ' ')
		} else {
			r2 := rune((((int(r) - int('a')) + sector) % 26) + 'a')
			out = append(out, r2)
		}
	}
	return string(out)
}

func main() {
	rooms, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, room := range rooms {
		charCounts := countChars(room.Name)

		good := verifyChecksum(charCounts, room.Checksum)
		if good {
			fmt.Printf("Room: %v Sector: %v\n", decrypt(room.Name, room.Sector), room.Sector)
		}
	}
}
