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

	posns := []Position{Position{0, 0}, Position{0, 0}}
	active := 0

	presents := map[Position]int{}
	presents[posns[0]]++
	presents[posns[1]]++

	for _, command := range commands {
		oldPos := posns[active]
		var newPos Position

		//fmt.Sprintf("command '%c'\n", command)
		switch command {
		case '<':
			newPos = Position{oldPos.X - 1, oldPos.Y}
			break
		case '>':
			newPos = Position{oldPos.X + 1, oldPos.Y}
			break
		case '^':
			newPos = Position{oldPos.X, oldPos.Y - 1}
			break
		case 'v':
			newPos = Position{oldPos.X, oldPos.Y + 1}
			break
		default:
			panic(fmt.Sprintf("unknown commmand '%c'", command))
		}

		posns[active] = newPos
		presents[newPos]++
		active = (active + 1) % 2
	}

	fmt.Println(presents)
	fmt.Println(len(presents))
}
