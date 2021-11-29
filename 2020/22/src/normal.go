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
	"github.com/simmonmt/aoc/2020/common/logger"
)

func PlayNormalRound(roundNum int, decks *[2]*Deck) {
	if logger.Enabled() {
		logger.LogF("-- Round %d -- ", roundNum)
		for _, deck := range decks {
			logger.LogF("Player %v's deck: %v", deck.Name(), deck.Cards())
		}
	}

	cards := [2]int{
		decks[0].Pop(),
		decks[1].Pop(),
	}

	if logger.Enabled() {
		for i, deck := range decks {
			logger.LogF("Player %v plays: %d", deck.Name(), cards[i])
		}
	}

	winner := 0
	if cards[1] > cards[0] {
		winner = 1
	}

	logger.LogF("Player %v wins the round", decks[winner].Name())

	decks[winner].Push(cards[winner])
	decks[winner].Push(cards[1-winner])
}

func PlayNormalGame(decks *[2]*Deck) {
	for i := 1; ; i++ {
		PlayNormalRound(i, decks)

		for _, deck := range decks {
			if deck.Empty() {
				return
			}
		}
	}
}
