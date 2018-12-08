package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"intmath"
	"logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
)

func readInput() ([]string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lines, nil
}

type Node struct {
	Off         int
	Len         int
	Children    []*Node
	MetadataOff int
	NumMetadata int
	Metadata    []int
}

func parseTree(off int, nums []int) *Node {
	logger.LogF("parseTree off %d", off)

	node := &Node{
		Off:         off,
		Children:    make([]*Node, nums[off]),
		NumMetadata: nums[off+1],
	}

	logger.LogF("parseTree partial %+v", node)

	off += 2
	for i := range node.Children {
		logger.LogF("parseTree parsing %d's child %d", node.Off, i)

		child := parseTree(off, nums)
		node.Children[i] = child
		off += child.Len
	}

	logger.LogF("parseTree done with %d's children, metadata at %v", node.Off, off)

	node.MetadataOff = off
	node.Metadata = nums[off : off+node.NumMetadata]
	off += node.NumMetadata
	node.Len = off - node.Off

	logger.LogF("parseTree done with %d: %+v", node.Off, node)

	return node
}

func sumMetadata(node *Node) int {
	sum := 0
	for _, md := range node.Metadata {
		sum += md
	}

	for _, child := range node.Children {
		sum += sumMetadata(child)
	}

	return sum
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	lines, err := readInput()
	if err != nil {
		log.Fatal(err)
	}
	line := lines[0]

	nums := []int{}
	for _, str := range strings.Split(line, " ") {
		nums = append(nums, intmath.AtoiOrDie(str))
	}

	root := parseTree(0, nums)

	sum := sumMetadata(root)

	fmt.Println(sum)
}
