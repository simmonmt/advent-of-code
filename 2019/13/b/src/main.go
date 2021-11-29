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

	"github.com/simmonmt/aoc/2019/13/b/src/puzzle"
	"github.com/simmonmt/aoc/2019/common/logger"
	"github.com/simmonmt/aoc/2019/common/pos"
	"github.com/simmonmt/aoc/2019/common/vm"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	ramPath = flag.String("ram", "", "path to file containing ram values")
)

type Command struct {
	v [3]int
}

type ArcadeIO struct {
	GameState *GameState
	accum     [3]int
	num       int
}

func NewArcadeIO(gameState *GameState) *ArcadeIO {
	return &ArcadeIO{
		GameState: gameState,
	}
}

func (io *ArcadeIO) Read() int64 {
	return int64(io.GameState.Joystick)
}

func (io *ArcadeIO) Write(v int64) {
	io.accum[io.num] = int(v)
	io.num++
	if io.num == 3 {
		io.GameState.Update(&Command{io.accum})
		io.num = 0
	}
}

type GameState struct {
	Board              *puzzle.Board
	BallPos, PaddlePos pos.P2
	Score              int
	Joystick           int
}

func NewGameState(b *puzzle.Board) *GameState {
	return &GameState{
		Board:     b,
		BallPos:   pos.P2{-1, -1},
		PaddlePos: pos.P2{-1, -1},
		Score:     -1,
		Joystick:  0,
	}
}

func (s *GameState) Update(cmd *Command) {
	if cmd.v[0] == -1 && cmd.v[1] == 0 {
		s.Score = cmd.v[2]
	} else {
		x, y := cmd.v[0], cmd.v[1]
		t := puzzle.Tile(cmd.v[2])

		if t == puzzle.TILE_BALL {
			s.BallPos = pos.P2{x, y}
		} else if t == puzzle.TILE_HPADDLE {
			s.PaddlePos = pos.P2{x, y}
		}

		s.Board.Set(x, y, t)
	}

	if s.PaddlePos.X != -1 && s.BallPos.X != -1 {
		if s.PaddlePos.X < s.BallPos.X {
			s.Joystick = 1
		} else if s.PaddlePos.X > s.BallPos.X {
			s.Joystick = -1
		} else {
			s.Joystick = 0
		}
	}
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *ramPath == "" {
		log.Fatalf("--ram is required")
	}

	ram, err := vm.NewRamFromFile(*ramPath)
	if err != nil {
		log.Fatal(err)
	}
	ram.Write(0, 2) // tell it to play rather than just dump

	gameState := NewGameState(puzzle.NewBoard(50, 50))

	io := NewArcadeIO(gameState)

	if err := vm.Run(ram, io); err != nil {
		panic(fmt.Sprintf("vm failed: %v", err))
	}

	fmt.Printf("game over score %d\n", gameState.Score)
}
