// tried 20: Outcome: 31 * 1364 = 42284

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"

	"lib"
	"logger"
)

var (
	verbose        = flag.Bool("verbose", false, "verbose")
	numTurns       = flag.Int("num_turns", 1, "num turns")
	elfAttackPower = flag.Int("elf_attack_power", 3, "elf attack power")
	validateBoard  = flag.Bool("validate", true, "validate")
)

type Result int

const (
	RESULT_CONTINUE Result = iota
	RESULT_NOTHINGTODO
)

func (r Result) String() string {
	switch r {
	case RESULT_CONTINUE:
		return "continue"
	case RESULT_NOTHINGTODO:
		return "nothing_to_do"
	default:
		panic("unknown result")
	}
}

func readInput() (*lib.Board, error) {
	lines := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for y := 0; scanner.Scan(); y++ {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return lib.NewBoard(lines, *elfAttackPower, *validateBoard), nil
}

func getOthers(board *lib.Board, self int) []lib.Char {
	others := []lib.Char{}
	for _, char := range board.Chars() {
		if char.Num != self {
			others = append(others, char)
		}
	}
	return others
}

func findEnemies(cands []lib.Char, selfIsElf bool) []lib.Char {
	enemies := []lib.Char{}
	for _, other := range cands {
		if other.IsElf != selfIsElf {
			enemies = append(enemies, other)
		}
	}
	return enemies
}

func neighborsToAttack(board *lib.Board, self lib.Char) []lib.Char {
	cands := []lib.Char{}
	for _, neighbor := range board.SurroundingCharacters(self.P) {
		if neighbor.IsElf != self.IsElf {
			cands = append(cands, neighbor)
		}
	}
	return cands
}

func findInRange(board *lib.Board, self lib.Char, enemies []lib.Char) []lib.Pos {
	inRangeMap := map[lib.Pos]bool{}

	for _, enemy := range enemies {
		for _, pos := range board.OpenSurroundingPositions(enemy.P) {
			inRangeMap[pos] = true
		}
	}

	inRange := make([]lib.Pos, len(inRangeMap))
	i := 0
	for p := range inRangeMap {
		inRange[i] = p
		i++
	}
	return inRange

}

func findReachable(board *lib.Board, self lib.Char, targets []lib.Pos) map[lib.Pos][]lib.Pos {
	results := map[lib.Pos][]lib.Pos{}
	for _, target := range targets {
		if result := board.ShortestPath(self.P, target); result != nil {
			results[target] = result
		}
	}
	return results
}

func shortestPaths(paths map[lib.Pos][]lib.Pos) (map[lib.Pos][]lib.Pos, int) {
	shortestLen := math.MaxInt32
	for _, path := range paths {
		if len(path) < shortestLen {
			shortestLen = len(path)
		}
	}

	shortest := map[lib.Pos][]lib.Pos{}
	for target, path := range paths {
		if len(path) == shortestLen {
			shortest[target] = path
		}
	}

	return shortest, shortestLen
}

func neighborShortestPaths(board *lib.Board, from, to lib.Pos, wantLen int) []lib.Pos {
	options := []lib.Pos{}

	for _, neighbor := range board.OpenSurroundingPositions(from) {
		var path []lib.Pos

		if neighbor == to {
			path = []lib.Pos{}
		} else {
			path = board.ShortestPath(neighbor, to)
			if path == nil {
				continue
			}
		}

		logger.LogF("%v to %v path %v (want len %v)", neighbor, to, path, wantLen)

		if len(path)+1 == wantLen {
			options = append(options, neighbor)
		}
	}

	return options
}

func findNextMove(board *lib.Board, self lib.Char, enemies []lib.Char) (lib.Pos, bool) {
	inRange := findInRange(board, self, enemies)
	if *verbose {
		fmt.Println("in range:")
		board.DumpWithDecorations(inRange, '?')
	}

	if len(inRange) == 0 {
		logger.LogLn("nothing in range")
		return lib.Pos{}, false
	}

	reachable := findReachable(board, self, inRange)
	if *verbose {
		posns := []lib.Pos{}
		for posn, _ := range reachable {
			posns = append(posns, posn)
		}

		fmt.Println("Reachable:")
		board.DumpWithDecorations(posns, '@')
		//fmt.Println(reachable)
	}

	if len(reachable) == 0 {
		logger.LogLn("nothing reachable")
		return lib.Pos{}, false
	}

	shortestPaths, shortestLen := shortestPaths(reachable)
	if *verbose {
		posns := []lib.Pos{}
		for posn, _ := range shortestPaths {
			posns = append(posns, posn)
		}

		fmt.Println("Nearest:")
		board.DumpWithDecorations(posns, '!')
	}

	nearest := []lib.Pos{}
	for target := range shortestPaths {
		nearest = append(nearest, target)
	}
	sort.Sort(lib.PosByReadingOrder(nearest))
	chosen := nearest[0]
	if *verbose {
		fmt.Println("Chosen:")
		board.DumpWithDecoration(chosen, '+')
	}

	neighborChoices := neighborShortestPaths(board, self.P, chosen, shortestLen)
	sort.Sort(lib.PosByReadingOrder(neighborChoices))
	if *verbose {
		fmt.Printf("options: %v\n", neighborChoices)
	}

	return neighborChoices[0], true
}

func chooseVictim(cands []lib.Char) lib.Char {
	lowestHP := math.MaxInt32
	lowHPCands := []lib.Char{}
	for _, c := range cands {
		if c.HP == lowestHP {
			lowHPCands = append(lowHPCands, c)
		} else if c.HP < lowestHP {
			lowHPCands = []lib.Char{c}
			lowestHP = c.HP
		}
	}

	sort.Sort(lib.CharByReadingOrder(lowHPCands))

	return lowHPCands[0]
}

func charAttack(board *lib.Board, self lib.Char, neighbors []lib.Char) (Result, *int) {
	victim := chooseVictim(neighbors)
	logger.LogF("chose victim %v", victim)

	victim, isDead := board.Attack(self, victim)
	if isDead {
		logger.LogF("victim now dead: %v", victim)
		board.RemoveChar(victim)
		num := victim.Num
		return RESULT_CONTINUE, &num
	}

	logger.LogF("victim current status %v", victim)
	return RESULT_CONTINUE, nil
}

func charTurn(board *lib.Board, self lib.Char) (Result, *int) {
	logger.LogF("\n-- character turn %v", self)

	others := getOthers(board, self.Num)
	enemies := findEnemies(others, self.IsElf)
	if len(enemies) == 0 {
		return RESULT_NOTHINGTODO, nil
	}

	if neighbors := neighborsToAttack(board, self); len(neighbors) != 0 {
		return charAttack(board, self, neighbors)
	}

	nextPos, hasMove := findNextMove(board, self, enemies)
	if !hasMove {
		logger.LogF("can't attack, can't move")
		return RESULT_CONTINUE, nil
	}
	logger.LogF("moving to %v", nextPos)

	self = board.MoveChar(self, nextPos)
	if *verbose {
		fmt.Println("After move:")
		board.Dump()
	}

	if neighbors := neighborsToAttack(board, self); len(neighbors) != 0 {
		return charAttack(board, self, neighbors)
	}

	return RESULT_CONTINUE, nil
}

func playTurn(board *lib.Board) Result {
	// Character numbers of dead characters
	deadChars := map[int]bool{}

	chars := board.Chars()
	sort.Sort(lib.CharByReadingOrder(chars))

	for _, char := range chars {
		if _, found := deadChars[char.Num]; found {
			logger.LogF("-- dead char %v skipped", char.Num)
			continue
		}

		// The character may have been attacked since the
		// beginning of the round, so fetch its state
		// again. Characters only move under their own power
		// so it's safe to request by location.
		char = board.GetChar(char.P)

		result, victimNum := charTurn(board, char)
		if victimNum != nil {
			logger.LogF("recording deadChars %v\n", *victimNum)
			deadChars[*victimNum] = true
		}

		if result != RESULT_CONTINUE {
			return result
		}
	}

	return RESULT_CONTINUE
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	board, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	var numElvesOrig, numGoblinsOrig int
	for _, char := range board.Chars() {
		if char.IsElf {
			numElvesOrig++
		} else {
			numGoblinsOrig++
		}
	}

	gameOver := false
	lastFullTurn := 0
	for turnNo := 1; !gameOver && (*numTurns == -1 || turnNo <= *numTurns); turnNo++ {
		if *verbose {
			fmt.Printf("start turn %d:\n", turnNo)
			board.Dump()
		}

		result := playTurn(board)
		if result == RESULT_NOTHINGTODO {
			gameOver = true
			break
		}

		lastFullTurn = turnNo
	}

	fmt.Println()
	if gameOver {
		fmt.Printf("Combat ends after %d full rounds\n", lastFullTurn)

		hpLeft := 0
		whoWon := ""
		numElves := 0
		numGoblins := 0
		for _, char := range board.Chars() {
			if char.IsElf {
				whoWon = "Elves"
				numElves++
			} else {
				whoWon = "Goblins"
				numGoblins++
			}
			hpLeft += char.HP
		}

		fmt.Printf("%v win with %d total hit points left\n", whoWon, hpLeft)
		fmt.Printf("Elves alive %d (lost %d) Goblins alive %d (lost %d)\n",
			numElves, numElvesOrig-numElves, numGoblins, numGoblinsOrig-numGoblins)
		fmt.Printf("Outcome: %v * %v = %v\n", lastFullTurn, hpLeft, lastFullTurn*hpLeft)

	} else {
		fmt.Println("game terminated due to turn restriction")
	}
}
