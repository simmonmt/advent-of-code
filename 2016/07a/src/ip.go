package main

import (
	"bufio"
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

func hasABBA(s string) bool {
	chars := []byte(s)

	for i := 0; i < len(chars)-3; i++ {
		if chars[i] == chars[i+3] && chars[i+1] == chars[i+2] && chars[i] != chars[i+1] {
			return true
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
		fmt.Printf("%+v\n", addr)
		hyperABBAFound := false
		for _, hyper := range addr.Hyper {
			if hasABBA(hyper) {
				hyperABBAFound = true
				break
			}
		}
		if hyperABBAFound {
			continue
		}

		normalABBAFound := false
		for _, normal := range addr.Normal {
			if hasABBA(normal) {
				normalABBAFound = true
				break
			}
		}
		if !normalABBAFound {
			continue
		}

		numGood++
	}

	fmt.Println(numGood)
}
