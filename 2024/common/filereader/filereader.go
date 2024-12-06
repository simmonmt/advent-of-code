// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filereader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func Lines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return linesFromReader(f)
}

func linesFromReader(r io.Reader) ([]string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

func Numbers(path string) ([]int, error) {
	lines, err := Lines(path)
	if err != nil {
		return nil, err
	}

	nums := []int{}
	for i, line := range lines {
		num, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("%d: failed to parse %v: %v",
				i, line, err)
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func BlankSeparatedGroups(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return blankSeparatedGroupsFromReader(f)
}

func blankSeparatedGroupsFromReader(r io.Reader) ([][]string, error) {
	lines, err := linesFromReader(r)
	if err != nil {
		return nil, err
	}

	return BlankSeparatedGroupsFromLines(lines), nil
}

func BlankSeparatedGroupsFromLines(lines []string) [][]string {
	// So we don't have to special-case the loop end
	lines = lines[:]
	lines = append(lines, "")

	groups := [][]string{}
	curGroup := []string{}
	for _, line := range lines {
		if line == "" {
			if len(curGroup) > 0 {
				groups = append(groups, curGroup)
			}
			curGroup = []string{}
			continue
		}

		curGroup = append(curGroup, line)
	}

	return groups
}

func ParseNumbersFromLine(line string) ([]int, error) {
	nums := []int{}
	for _, str := range strings.Split(line, ",") {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", str, err)
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func OneRowOfNumbers(path string) ([]int, error) {
	lines, err := Lines(path)
	if err != nil {
		return nil, err
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("expected one line, got %v", len(lines))
	}

	return ParseNumbersFromLine(lines[0])
}
