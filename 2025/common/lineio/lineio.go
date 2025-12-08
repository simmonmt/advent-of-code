package lineio

import (
	"fmt"
	"strconv"
	"strings"
)

func Numbers(lines []string) ([]int, error) {
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

func BlankSeparatedGroups(lines []string) [][]string {
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

func NumbersFromLine(line string, sep string) ([]int, error) {
	nums := []int{}
	for _, str := range strings.Split(line, sep) {
		if str == "" {
			continue
		}
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %v", str, err)
		}
		nums = append(nums, num)
	}
	return nums, nil
}
