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

	"grid"
	"node"
)

var (
	sizePattern = regexp.MustCompile(`^/dev/grid/node-x([0-9]+)-y([0-9]+) +([0-9]+)T +([0-9]+)T +[0-9]+T +[0-9]+%$`)
)

func readInput(r io.Reader) (width uint8, height uint8, nodes []node.Node, err error) {
	nodes = []node.Node{}

	var maxX, maxY uint64

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if lineNum < 3 {
			continue
		}

		line = strings.TrimSpace(line)
		matches := sizePattern.FindStringSubmatch(line)
		if matches == nil {
			return 0, 0, nil, fmt.Errorf("%d: failed to parse", lineNum)
		}

		x, err := strconv.ParseUint(matches[1], 10, 8)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("%d: failed to parse x: %v", lineNum, err)
		}
		y, err := strconv.ParseUint(matches[2], 10, 8)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("%d: failed to parse y: %v", lineNum, err)
		}
		size, err := strconv.ParseUint(matches[3], 10, 16)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("%d: failed to parse size: %v", lineNum, err)
		}
		used, err := strconv.ParseUint(matches[4], 10, 16)
		if err != nil {
			return 0, 0, nil, fmt.Errorf("%d: failed to parse used: %v", lineNum, err)
		}

		if x > maxX {
			x = maxX
		}
		if y > maxY {
			y = maxY
		}

		nodes = append(nodes, *node.New(uint16(size), uint16(used)))
	}

	return uint8(maxX + 1), uint8(maxY + 1), nodes, nil
}

func main() {
	width, height, nodes, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	numViable := 0
	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes); j++ {
			if i == j {
				continue
			}

			a, b := nodes[i], nodes[j]
			if a.Used == 0 {
				continue
			}
			if a.Used > (b.Size - b.Used) {
				continue
			}

			numViable++
		}
	}

	fmt.Printf("#viable = %v\n", numViable)

	g, err := grid.New(int(width), int(height), uint8(width-1), 0, nodes)
	if err != nil {
		log.Fatal(err)
	}

	g.Print()
}
