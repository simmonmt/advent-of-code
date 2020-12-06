package filereader

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	// So we don't have to special-case the loop end
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

	return groups, nil
}
