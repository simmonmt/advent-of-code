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
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/simmonmt/aoc/2021/common/filereader"
	"github.com/simmonmt/aoc/2021/common/logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")

	inputPattern = regexp.MustCompile(`^Player (\d+) starting position: (\d+)$`)
)

type Player struct {
	Num   int
	Pos   int
	Score int
}

func (p *Player) String() string {
	return fmt.Sprintf("<%d p:%d s:%d>", p.Num, p.Pos, p.Score)
}

func (p *Player) Advance(num int) {
	p.Pos = (((p.Pos - 1) + num) % 10) + 1
}

type DetDie struct {
	next     int
	numRolls int
}

func NewDetDie() *DetDie {
	return &DetDie{
		next:     1,
		numRolls: 0,
	}
}

func (d *DetDie) Roll() int {
	ret := d.next

	d.next++
	if d.next == 101 {
		d.next = 1
	}
	d.numRolls++
	return ret
}

func (d *DetDie) NumRolls() int {
	return d.numRolls
}

func readInput(path string) ([]*Player, error) {
	lines, err := filereader.Lines(*input)
	if err != nil {
		return nil, err
	}

	out := []*Player{}
	for i, line := range lines {
		parts := inputPattern.FindStringSubmatch(line)
		if parts == nil {
			return nil, fmt.Errorf("%d: parse failure", i+1)
		}

		playerNum, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("%d: bad player number: %v", i+1, err)
		}

		pos, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("%d: bad pos: %v", i+1, err)
		}

		out = append(out, &Player{
			Num:   playerNum,
			Pos:   pos,
			Score: 0,
		})
	}

	return out, err
}

func dupPlayers(players []*Player) []*Player {
	tmp := make([]*Player, len(players))
	for i, p := range players {
		p2 := *p
		tmp[i] = &p2
	}
	return tmp
}

func playPlayer(die *DetDie, player *Player) bool {
	roll := 0
	for i := 0; i < 3; i++ {
		roll += die.Roll()
	}

	player.Advance(roll)
	player.Score += player.Pos

	//logger.LogF("Player %v rolls %v and moves to space %v for a total score of %v",
	//	player.Num, roll, player.Pos, player.Score)

	return player.Score >= 1000
}

func playRound(die *DetDie, players []*Player) bool {
	for _, p := range players {
		if playPlayer(die, p) {
			return true
		}
	}

	return false
}

func solveA(players []*Player) {
	players = dupPlayers(players)
	die := NewDetDie()

	for {
		if playRound(die, players) {
			break
		}
	}

	loserScore := 0
	for _, p := range players {
		if p.Score < 1000 {
			loserScore = p.Score
			break
		}
	}

	fmt.Println("A", loserScore*die.NumRolls())
}

type HashKey struct {
	start1, start2 int
	score1, score2 int
}

func recurse(level int, players [2]Player, cache map[HashKey][2]int64) [2]int64 {
	//logger.LogF("%d recurse %v", level, players)

	wins := [2]int64{0, 0}

	for n1 := 0; n1 < 27; n1++ {
		roll1 := (n1%3 + 1) + ((n1/3)%3 + 1) + ((n1/9)%3 + 1)
		if roll1 < 3 || roll1 > 9 {
			panic("roll1 oob")
		}

		np := players
		np[0].Advance(roll1)
		np[0].Score += np[0].Pos
		if np[0].Score >= 21 {
			wins[0]++
			continue
		}

		for n2 := 0; n2 < 27; n2++ {
			roll2 := (n2%3 + 1) + ((n2/3)%3 + 1) + ((n2/9)%3 + 1)
			if roll2 < 3 || roll2 > 9 {
				panic("roll2 oob")
			}

			np[1] = players[1]
			np[1].Advance(roll2)
			np[1].Score += np[1].Pos
			if np[1].Score >= 21 {
				wins[1]++
				continue
			}

			key := HashKey{
				start1: np[0].Pos,
				start2: np[1].Pos,
				score1: np[0].Score,
				score2: np[1].Score,
			}

			res, found := cache[key]
			if !found {
				res = recurse(level+1, np, cache)
				cache[key] = res
			}

			wins[0] += res[0]
			wins[1] += res[1]
		}
	}

	return wins
}

func solveB(players []*Player) {
	parr := [2]Player{*players[0], *players[1]}
	cache := map[HashKey][2]int64{}
	res := recurse(1, parr, cache)
	logger.LogF("result %v", res)

	max := int64(0)
	for _, r := range res {
		if r > max {
			max = r
		}
	}

	fmt.Println("B", max)
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		log.Fatalf("--input is required")
	}

	players, err := readInput(*input)
	if err != nil {
		log.Fatal(err)
	}

	solveA(players)
	solveB(players)
}
