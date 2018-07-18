package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	pattern = regexp.MustCompile(`([a-z]*)(?:\[([a-z]*)\])?`)
)

type Addr struct {
	Normal []string
	Hyper  []string
}

func readInput(r io.Reader) ([]Addr, error) {
	addrs := []Addr{}

	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		allMatches := pattern.FindAllStringSubmatch(line, -1)
		if allMatches == nil {
			return nil, fmt.Errorf("can't parse %v", line)
		}

		addr := Addr{
			Normal: []string{},
			Hyper:  []string{},
		}
		for _, matches := range allMatches {
			if matches[1] != "" {
				addr.Normal = append(addr.Normal, matches[1])
			}
			if matches[2] != "" {
				addr.Hyper = append(addr.Hyper, matches[2])
			}
		}
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

func findABAs(s string) [][3]byte {
	abas := [][3]byte{}
	chars := []byte(s)

	for i := 0; i < len(chars)-2; i++ {
		if chars[i] == chars[i+2] && chars[i] != chars[i+1] {
			aba := [3]byte{}
			copy(aba[:], chars[i:i+3])
			abas = append(abas, aba)
		}
	}

	return abas
}

func hasBAB(abas [][3]byte, s string) bool {
	chars := []byte(s)

	for _, aba := range abas {
		bab := []byte{aba[1], aba[0], aba[1]}

		for i := 0; i < len(chars)-2; i++ {
			if bytes.Equal(chars[i:i+3], bab) {
				return true
			}
		}
	}

	return false
}

func main() {
	addrs, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err.Error())
	}

	numGood := 0
	for _, addr := range addrs {
		//fmt.Printf("%+v\n", addr)

		allABAs := [][3]byte{}

		for _, s := range addr.Normal {
			allABAs = append(allABAs, findABAs(s)...)
		}
		if len(allABAs) == 0 {
			continue
		}

		//fmt.Printf("found ABAs: %v\n", allABAs)

		foundBAB := false
		for _, s := range addr.Hyper {
			if hasBAB(allABAs, s) {
				//fmt.Println("BAB found")
				foundBAB = true
				break
			}
		}
		if !foundBAB {
			continue
		}

		numGood++
	}

	fmt.Println(numGood)
}
