package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reindeerPattern = regexp.MustCompile(`^([^ ]+) can fly ([0-9]+) km/s for ([0-9]+) seconds, but then must rest for ([0-9]+) seconds\.$`)
)

type Reindeer struct {
	Name                     string
	FlyRate, FlyLen, RestLen int

	Flying bool
	Left   int

	Dist, Points int
}

func advance(deer *Reindeer) {
	if deer.Left == 0 {
		if deer.Flying {
			deer.Left = deer.RestLen
		} else {
			deer.Left = deer.FlyLen
		}
		deer.Flying = !deer.Flying
	}

	deer.Left--
	if deer.Flying {
		deer.Dist += deer.FlyRate
	}
}

func summarize(deer *Reindeer) string {
	return fmt.Sprintf("%v: dist %v pts %v", deer.Name, deer.Dist, deer.Points)
}

func parseReindeer(line string) (*Reindeer, error) {
	matches := reindeerPattern.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("failed to parse reindeer")
	}

	name := matches[1]
	flyRate, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse reindeer fly rate %v: %v", matches[2], err)
	}
	flyLen, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("failed to parse reindeer fly len %v: %v", matches[3], err)
	}
	restLen, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse reindeer rest len %v: %v", matches[4], err)
	}

	return &Reindeer{
		Name:    name,
		FlyRate: flyRate,
		FlyLen:  flyLen,
		RestLen: restLen,
	}, nil
}

func readReindeers(r io.Reader) ([]*Reindeer, error) {
	deers := []*Reindeer{}

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		deer, err := parseReindeer(strings.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("%d: %v", lineNum, err)
		}

		deers = append(deers, deer)
	}

	return deers, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v num_secs", os.Args[0])
	}
	numSecs, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("failed to parse num_secs %v: %v", os.Args[1], err)
	}

	deers, err := readReindeers(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= numSecs; i++ {
		maxDist := -1
		for _, deer := range deers {
			advance(deer)

			if deer.Dist > maxDist {
				maxDist = deer.Dist
			}
		}

		for _, deer := range deers {
			if deer.Dist == maxDist {
				deer.Points++
			}
		}
	}

	for _, deer := range deers {
		fmt.Println(summarize(deer))
	}
}
