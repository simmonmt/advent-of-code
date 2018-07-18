package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"command"
	"screen"
)

var (
	screenWidth  = flag.Int("width", 50, "screen width")
	screenHeight = flag.Int("height", 6, "screen height")
)

func readInput(r io.Reader) ([]command.Command, error) {
	cmds := []command.Command{}

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		var cmd command.Command
		switch {
		case strings.HasPrefix(line, "rect "):
			cmd, err = command.ParseRect(line)
			break
		case strings.HasPrefix(line, "rotate row "):
			cmd, err = command.ParseRotateRow(line)
			break
		case strings.HasPrefix(line, "rotate column "):
			cmd, err = command.ParseRotateColumn(line)
			break
		default:
			return nil, fmt.Errorf("unknown command %v", line)
		}

		if err != nil {
			return nil, fmt.Errorf("%d: %v", lineNum, err)
		}

		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func main() {
	flag.Parse()

	cmds, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err.Error())
	}

	s := screen.NewScreen(*screenWidth, *screenHeight)
	for _, cmd := range cmds {
		cmd.Execute(s)
		s.Print()
	}

	fmt.Printf("num on = %d\n", s.CountOn())
}
