package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Position struct {
	X, Y int
}

func readCommands(in io.Reader) string {
	reader := bufio.NewReader(in)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func main() {
	commands := readCommands(os.Stdin)

	pos := Position{0, 0}
	presents := map[Position]int{}
	presents[pos]++

	for _, command := range commands {
		//fmt.Sprintf("command '%c'\n", command)
		switch command {
		case '<':
			pos = Position{pos.X - 1, pos.Y}
			break
		case '>':
			pos = Position{pos.X + 1, pos.Y}
			break
		case '^':
			pos = Position{pos.X, pos.Y - 1}
			break
		case 'v':
			pos = Position{pos.X, pos.Y + 1}
			break
		default:
			panic(fmt.Sprintf("unknown commmand '%c'", command))
		}

		presents[pos]++
	}

	fmt.Println(presents)
	fmt.Println(len(presents))
}
