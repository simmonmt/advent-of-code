package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	pattern = regexp.MustCompile(`^([a-z]+) \(([0-9]+)\)(?: -> (.*))?$`)
)

type Tree struct {
	elems map[string]string
}

func NewTree() *Tree {
	return &Tree{elems: map[string]string{}}
}

func (t *Tree) Insert(bot, top string) {
	t.elems[top] = bot
}

func (t *Tree) Bottom() string {
	var bot string
	for _, bot = range t.elems {
		break
	}

	for {
		next, found := t.elems[bot]
		if !found {
			return bot
		}

		bot = next
	}
}

func main() {
	tree := NewTree()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		matches := pattern.FindStringSubmatch(line)
		if len(matches) == 0 {
			log.Fatalf("failed to parse %v", line)
		}
		bot := matches[1]
		// weight, err := strconv.Atoi(matches[2])
		// if err != nil {
		// 	log.Fatalf("failed to parse weight %v in %v", matches[2], line)
		// }

		topsStr := matches[3]
		if topsStr == "" {
			continue
		}
		tops := strings.Split(topsStr, ", ")

		//fmt.Printf("%v %v\n", bot, tops)

		for _, top := range tops {
			tree.Insert(bot, top)
		}
	}

	fmt.Println(tree.Bottom())
}
