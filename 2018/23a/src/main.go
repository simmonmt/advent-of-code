package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"logger"
	"xyzpos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")

	inputPattern = regexp.MustCompile(`pos=<([^>]+)>, r=(\d+)`)
)

type Bot struct {
	Pos    xyzpos.Pos
	Radius int
}

func readInput() ([]*Bot, error) {
	bots := []*Bot{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		parts := inputPattern.FindStringSubmatch(line)
		if parts == nil {
			return nil, fmt.Errorf("failed to parse %v", line)
		}

		pos, err := xyzpos.Parse(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse pos: %v", err)
		}
		radius, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("failed to parse radiusz: %v", err)
		}

		bots = append(bots, &Bot{pos, radius})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return bots, nil
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	bots, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	var strongest *Bot
	for _, bot := range bots {
		if strongest == nil || bot.Radius > strongest.Radius {
			strongest = bot
		}
	}

	fmt.Printf("strongest is %+v\n", *strongest)

	num := 0
	for _, bot := range bots {
		if bot.Pos.Dist(strongest.Pos) <= strongest.Radius {
			logger.LogF("in range %+v\n", *bot)
			num++
		}
	}
	fmt.Println(num)
}
