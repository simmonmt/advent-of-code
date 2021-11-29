// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	valuePattern = regexp.MustCompile(`^value ([0-9]+) goes to bot ([0-9]+)$`)
	botPattern   = regexp.MustCompile(
		`^bot ([0-9]+) gives (low|high) to (bot|output) ([0-9]+) and (low|high) to (bot|output) ([0-9]+)$`)
)

type DestType int

const (
	DEST_BOT DestType = iota
	DEST_OUTPUT
)

func (d DestType) String() string {
	switch d {
	case DEST_BOT:
		return "bot"
	case DEST_OUTPUT:
		return "output"
	default:
		panic(fmt.Sprintf("unknown dest %d", d))
	}
}

func ParseDest(str string) DestType {
	switch str {
	case "bot":
		return DEST_BOT
	case "output":
		return DEST_OUTPUT
	default:
		panic(fmt.Sprintf("unknown dest %v", str))
	}
}

type LevelType int

const (
	LEVEL_LOW LevelType = iota
	LEVEL_HIGH
)

func (l LevelType) String() string {
	switch l {
	case LEVEL_LOW:
		return "low"
	case LEVEL_HIGH:
		return "high"
	default:
		panic(fmt.Sprintf("unknown level %d", l))
	}
}

func ParseLevel(str string) LevelType {
	switch str {
	case "low":
		return LEVEL_LOW
	case "high":
		return LEVEL_HIGH
	default:
		panic(fmt.Sprintf("unknown level %v", str))
	}
}

type Bot struct {
	num            int
	in1Val, in2Val int

	outHighDest DestType
	outHighNum  int
	outLowDest  DestType
	outLowNum   int
}

type Bots []*Bot

func (b Bots) Len() int      { return len(b) }
func (b Bots) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

type ByBotNum struct{ Bots }

func (b ByBotNum) Less(i, j int) bool { return b.Bots[i].num < b.Bots[j].num }

func parseNum(str string) (int, error) {
	n, err := strconv.ParseUint(str, 10, 32)
	return int(n), err
}

func parseBot(command string) (*Bot, error) {
	matches := botPattern.FindStringSubmatch(command)
	if matches == nil {
		return nil, fmt.Errorf("failed to parse")
	}

	botNum, err := parseNum(matches[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse bot num %v: %v", matches[1], err)
	}

	type giver struct {
		level LevelType
		dest  DestType
		num   int
	}

	parseGiver := func(level, dest, num string) (*giver, error) {
		n, err := parseNum(num)
		if err != nil {
			return nil, err
		}

		return &giver{
			level: ParseLevel(level),
			dest:  ParseDest(dest),
			num:   n,
		}, nil
	}

	givers := [2]*giver{}
	givers[0], err = parseGiver(matches[2], matches[3], matches[4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse first giver: %v", err)
	}
	givers[1], err = parseGiver(matches[5], matches[6], matches[7])
	if err != nil {
		return nil, fmt.Errorf("failed to parse first giver: %v", err)
	}

	var lowGiver, highGiver *giver
	if givers[0].level == LEVEL_LOW {
		lowGiver, highGiver = givers[0], givers[1]
	} else {
		lowGiver, highGiver = givers[1], givers[0]
	}

	return &Bot{
		num:    botNum,
		in1Val: -1,
		in2Val: -1,

		outHighDest: highGiver.dest,
		outHighNum:  highGiver.num,
		outLowDest:  lowGiver.dest,
		outLowNum:   lowGiver.num,
	}, nil
}

func parseValue(command string) (val, botNum int, err error) {
	matches := valuePattern.FindStringSubmatch(command)
	if matches == nil {
		return 0, 0, fmt.Errorf("failed to parse")
	}

	val, err = parseNum(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse value num %v: %v", matches[1], err)
	}

	botNum, err = parseNum(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse bot num %v: %v", matches[2], err)
	}

	return val, botNum, nil
}

func readInput(r io.Reader) ([]*Bot, map[int][]int, error) {
	bots := []*Bot{}
	initial := map[int][]int{}

	reader := bufio.NewReader(r)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "bot ") {
			bot, err := parseBot(line)
			if err != nil {
				return nil, nil, fmt.Errorf("%d: failed to parse bot: %v", lineNum, err)
			}
			bots = append(bots, bot)
		} else if strings.HasPrefix(line, "value ") {
			val, botNum, err := parseValue(line)
			if err != nil {
				return nil, nil, fmt.Errorf("%d: failed to parse value: %v", lineNum, err)
			}

			if _, found := initial[botNum]; !found {
				initial[botNum] = []int{}
			}
			initial[botNum] = append(initial[botNum], val)
		} else {
			return nil, nil, fmt.Errorf("%d: unknown command", lineNum)
		}
	}

	return bots, initial, nil
}

func addBotInputValue(bot *Bot, val int) {
	if bot.in1Val == -1 {
		bot.in1Val = val
	} else if bot.in2Val == -1 {
		bot.in2Val = val
	} else {
		panic(fmt.Sprintf("bot %v is full", bot.num))
	}
}

func pushBotValues(bot *Bot, bots map[int]*Bot) {
	high, low := bot.in1Val, bot.in2Val
	if high < low {
		high, low = low, high
	}

	if bot.outHighDest == DEST_BOT {
		if destBot, found := bots[bot.outHighNum]; !found {
			panic(fmt.Sprintf("bot %v has bad high dest bot %v", bot.num, bot.outHighNum))
		} else {
			addBotInputValue(destBot, high)
		}
	}
	if bot.outLowDest == DEST_BOT {
		if destBot, found := bots[bot.outLowNum]; !found {
			panic(fmt.Sprintf("bot %v has bad low dest bot %v", bot.num, bot.outLowNum))
		} else {
			addBotInputValue(destBot, low)
		}
	}
}

func populateBots(bots map[int]*Bot) bool {
	doneBots := map[int]bool{}

	lastRoundNumDone := 0
	for roundNum := 0; ; roundNum++ {
		for _, bot := range bots {
			if _, found := doneBots[bot.num]; found {
				continue
			}

			if bot.in1Val == -1 || bot.in2Val == -1 {
				continue
			}

			pushBotValues(bot, bots)
			doneBots[bot.num] = true
		}

		if len(doneBots) == len(bots) {
			fmt.Printf("round %v: done\n", roundNum)
			return true
		} else if lastRoundNumDone != len(doneBots) {
			fmt.Printf("round %v: num done %v of %v\n", roundNum, len(doneBots), len(bots))
		} else {
			fmt.Printf("round %v: no termination\n", roundNum)
			return false
		}

		lastRoundNumDone = len(doneBots)
	}

	panic("unreachable")
}

func botCompares(bot *Bot, a, b int) bool {
	return (a == bot.in1Val && b == bot.in2Val) || (a == bot.in2Val && b == bot.in1Val)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %v want1 want2", os.Args[0])
	}
	want1, err := parseNum(os.Args[1])
	if err != nil {
		log.Fatalf("failed to parse want1: %v", err)
	}
	want2, err := parseNum(os.Args[2])
	if err != nil {
		log.Fatalf("failed to parse want2: %v", err)
	}

	allBots, initial, err := readInput(os.Stdin)
	if err != nil {
		log.Fatal(err.Error())
	}

	sort.Sort(ByBotNum{allBots})

	bots := map[int]*Bot{}
	for _, bot := range allBots {
		fmt.Printf("bot %+v\n", bot)
		bots[bot.num] = bot
	}
	fmt.Printf("initial: %+v\n", initial)

	for botNum, vals := range initial {
		if len(vals) > 2 {
			log.Fatal("initial conditions for bot %v have %v values %v, wanted <=2",
				botNum, len(vals), vals)
		}

		if bot, found := bots[botNum]; !found {
			log.Fatal("initial conditions for unknown bot %v", botNum)
		} else {
			bot.in1Val = vals[0]
			if len(vals) == 2 {
				bot.in2Val = vals[1]
			}
		}
	}

	// for _, bot := range allBots {
	// 	fmt.Printf("bot %+v\n", bot)
	// }

	populated := populateBots(bots)
	// fmt.Println("post-population:")
	// for _, bot := range allBots {
	// 	fmt.Printf("bot %+v\n", bot)
	// }

	if !populated {
		log.Fatalf("did not finish population")
	}

	for _, bot := range allBots {
		if botCompares(bot, want1, want2) {
			fmt.Printf("bot %d matches\n", bot.num)
			break
		}
	}
}
