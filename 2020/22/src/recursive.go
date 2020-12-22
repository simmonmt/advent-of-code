package main

import "github.com/simmonmt/aoc/2020/common/logger"

var (
	lastRecursiveGameNum = 0
)

func nextGameNum() int {
	lastRecursiveGameNum++
	return lastRecursiveGameNum
}

func PlayRecursiveRound(gameNum, roundNum int, decks *[2]*Deck) {
	cards := [2]int{
		decks[0].Pop(),
		decks[1].Pop(),
	}

	if logger.Enabled() {
		for i, deck := range decks {
			logger.LogF("Player %v plays: %d", deck.Name(), cards[i])
		}
	}

	shouldRecurse := (decks[0].Num() >= cards[0]) &&
		(decks[1].Num() >= cards[1])

	winner := -1
	if shouldRecurse {
		logger.LogF("Playing a sub-game to determine the winner...\n")

		subGameNum := nextGameNum()
		subDecks := [2]*Deck{
			decks[0].CloneFirstN(cards[0]),
			decks[1].CloneFirstN(cards[1]),
		}

		winner = PlayRecursiveGame(subGameNum, &subDecks)

		logger.LogF("... anyway, back to game %v.", gameNum)
	} else {
		winner = 0
		if cards[1] > cards[0] {
			winner = 1
		}
	}
	logger.LogF("Player %v wins round %v of game %v!\n",
		decks[winner].Name(), roundNum, gameNum)

	decks[winner].Push(cards[winner])
	decks[winner].Push(cards[1-winner])
}

func doPlayRecursiveGame(gameNum int, decks *[2]*Deck) int {
	roundCache := map[string]int{}

	for roundNum := 1; ; roundNum++ {
		if logger.Enabled() {
			logger.LogF("-- Round %d (Game %d) -- ", roundNum, gameNum)
			for _, deck := range decks {
				logger.LogF("Player %v's deck: %v",
					deck.Name(), deck)
			}
		}

		roundCacheKey := decks[0].String() + "/" + decks[1].String()
		if prevRound, found := roundCache[roundCacheKey]; found {
			logger.LogF("Repeat of round %v; player %v wins",
				prevRound, decks[0].Name())
			return 0
		}
		roundCache[roundCacheKey] = roundNum

		PlayRecursiveRound(gameNum, roundNum, decks)

		for i, deck := range decks {
			if deck.Empty() {
				winner := 1 - i
				logger.LogF("The winner of game %v is player %v!\n",
					gameNum, decks[winner].Name())
				return winner
			}
		}

		// if roundNum >= 1000 {
		// 	panic("too many rounds")
		// }
	}
}

var (
	recursiveGameCache = map[string]int{}
)

func PlayRecursiveGame(gameNum int, decks *[2]*Deck) int {
	logger.LogF("=== Game %d ===\n", gameNum)

	gameCacheKey := decks[0].String() + "/" + decks[1].String()
	if winner, found := recursiveGameCache[gameCacheKey]; found {
		logger.LogF("Game was cached; winner %v", decks[winner].Name())
		return winner
	}

	winner := doPlayRecursiveGame(gameNum, decks)
	cacheWinner, found := recursiveGameCache[gameCacheKey]
	if found && cacheWinner != winner {
		panic("cache disagree")
	}

	recursiveGameCache[gameCacheKey] = winner
	return winner
}

func PlayRecursiveCards(decks *[2]*Deck) {
	PlayRecursiveGame(nextGameNum(), decks)
}
