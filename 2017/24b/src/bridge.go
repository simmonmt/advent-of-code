package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Port struct {
	Name string
	L, R int
}

type PortDirectory struct {
	portMap map[int][]Port
}

func newPortDirectory(ports []Port) *PortDirectory {
	portMap := map[int][]Port{}

	for _, port := range ports {
		if _, found := portMap[port.L]; !found {
			portMap[port.L] = []Port{}
		}
		portMap[port.L] = append(portMap[port.L], port)

		if _, found := portMap[port.R]; !found {
			portMap[port.R] = []Port{}
		}
		portMap[port.R] = append(portMap[port.R], port)
	}

	return &PortDirectory{portMap}
}

func (pd *PortDirectory) Find(num int) []Port {
	if ports, found := pd.portMap[num]; found {
		return ports
	}
	return []Port{}
}

type PortStack struct {
	names map[string]int
}

func newPortStack(ports []Port) *PortStack {
	names := map[string]int{}
	for _, port := range ports {
		names[port.Name]++
	}
	return &PortStack{names}
}

func (ps *PortStack) Alloc(name string) bool {
	if ps.names[name] > 0 {
		ps.names[name]--
		return true
	}
	return false
}

func (ps *PortStack) Free(name string) {
	ps.names[name]++
}

func readPorts(in io.Reader) ([]Port, error) {
	ports := []Port{}

	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		name := line
		parts := strings.Split(name, "/")

		left, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse left in %v: %v",
				line, err)
		}

		right, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse right in %v: %v",
				line, err)
		}

		ports = append(ports, Port{name, left, right})
	}

	return ports, nil
}

func portsCost(ports []Port) int {
	sum := 0
	for _, port := range ports {
		sum += port.L + port.R
	}
	return sum
}

func findMaxCost(start int, dir *PortDirectory, stack *PortStack) (int, []Port) {
	maxCostPorts := []Port{}
	maxCost := 0

	for _, port := range dir.Find(start) {
		if !stack.Alloc(port.Name) {
			continue
		}

		next := port.L
		if port.L == start {
			next = port.R
		}

		cands := []Port{port}

		subCost, subCands := findMaxCost(next, dir, stack)
		if subCost != 0 {
			cands = append(cands, subCands...)
		}

		cost := portsCost(cands)
		if len(cands) > len(maxCostPorts) ||
			(len(cands) == len(maxCostPorts) && cost > maxCost) {
			maxCost = cost
			maxCostPorts = cands
		}

		stack.Free(port.Name)
	}

	return maxCost, maxCostPorts
}

func main() {
	ports, err := readPorts(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read ports: %v", err)
	}

	dir := newPortDirectory(ports)
	stack := newPortStack(ports)

	maxCost, parts := findMaxCost(0, dir, stack)

	fmt.Println(maxCost, parts)
}
