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

	"github.com/simmonmt/aoc/2020/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	hclPattern = regexp.MustCompile(`^#[0-9a-f]{6}$`)
	validEcls  = map[string]bool{
		"amb": true, "blu": true, "brn": true, "gry": true,
		"grn": true, "hzl": true, "oth": true,
	}
)

func readInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func readPassports(path string) ([]map[string]string, error) {
	lines, err := readInput(path)
	if err != nil {
		return nil, err
	}

	passports := []map[string]string{}
	curPassport := map[string]string{}

	for _, line := range lines {
		if line == "" {
			passports = append(passports, curPassport)
			curPassport = map[string]string{}
			continue
		}

		for _, field := range strings.Split(line, " ") {
			parts := strings.SplitN(field, ":", 2)
			curPassport[parts[0]] = parts[1]
		}
	}
	if len(curPassport) > 0 {
		passports = append(passports, curPassport)
	}

	return passports, nil
}

func validNumber(str string, digits int, min, max uint64) bool {
	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return false
	}

	if num < min || num > max {
		return false
	}

	if digits > 0 && len(str) != digits {
		return false
	}

	return true
}

func validByr(str string) bool {
	return validNumber(str, 4, 1920, 2002)
}

func validIyr(str string) bool {
	return validNumber(str, 4, 2010, 2020)
}

func validEyr(str string) bool {
	return validNumber(str, 4, 2020, 2030)
}

func validHgt(str string) bool {
	if strings.HasSuffix(str, "cm") {
		return validNumber(strings.TrimSuffix(str, "cm"), -1, 150, 193)
	}
	if strings.HasSuffix(str, "in") {
		return validNumber(strings.TrimSuffix(str, "in"), -1, 59, 76)
	}
	return false
}

func validHcl(str string) bool {
	return hclPattern.MatchString(str)
}

func validEcl(str string) bool {
	_, found := validEcls[str]
	return found
}

func validPid(str string) bool {
	return len(str) == 9
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	passports, err := readPassports(*input)
	if err != nil {
		log.Fatal(err)
	}

	wantFields := []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}

	validators := map[string]func(string) bool{
		"byr": validByr,
		"iyr": validIyr,
		"eyr": validEyr,
		"hgt": validHgt,
		"hcl": validHcl,
		"ecl": validEcl,
		"pid": validPid,
	}

	numValid := 0
	for _, pp := range passports {
		valid := true
		for _, fieldName := range wantFields {
			if field, found := pp[fieldName]; !found || !validators[fieldName](field) {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}

		numValid++
		fmt.Println(pp)
	}

	fmt.Printf("num valid: %d\n", numValid)
}
