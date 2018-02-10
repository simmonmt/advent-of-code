package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	criteriaPath = flag.String("criteria", "", "file containing criteria")
)

func readCriteria(path string) (map[string]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %v: %v", path, err)
	}
	defer file.Close()

	criteria := map[string]int{}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		parts := strings.SplitN(line, ": ", 2)

		name := parts[0]
		val, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse val for %v: %v", name, err)
		}

		criteria[name] = val
	}

	return criteria, nil
}

type Sue struct {
	Name     string
	Criteria map[string]int
}

func parseSue(line string) (*Sue, error) {
	parts := regexp.MustCompile(`[:,] `).Split(line, -1)

	sue := &Sue{Name: parts[0], Criteria: map[string]int{}}
	for i := 1; i < len(parts); i += 2 {
		name := parts[i]
		if i+1 == len(parts) {
			return nil, fmt.Errorf("no val for criteria %v", name)
		}

		val, err := strconv.Atoi(parts[i+1])
		if err != nil {
			return nil, fmt.Errorf("bad number for %v: %v", name, err)
		}

		sue.Criteria[name] = val
	}

	return sue, nil
}

func isMatch(criteria map[string]int, sue *Sue) bool {
	for k, v := range sue.Criteria {
		if criteria[k] != v {
			return false
		}
	}
	return true
}

func main() {
	flag.Parse()

	if *criteriaPath == "" {
		log.Fatalf("--criteria is required")
	}

	criteria, err := readCriteria(*criteriaPath)
	if err != nil {
		log.Fatalf("failed to read criteria: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		sue, err := parseSue(strings.TrimSpace(line))
		if err != nil {
			log.Fatalf("%d: failed to parse sue: %v", lineNum, err)
		}

		if isMatch(criteria, sue) {
			fmt.Println(sue.Name)
		}
	}
}
