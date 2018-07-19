package game

import (
	"board"
	"fmt"
	"logger"
)

var (
	globalPlayNum = 0
)

type SeenVal struct {
	done  bool
	moves []*boardMove
}

func doPlay(b *board.Board, seen map[string][]*board.Move, playNum int) []*board.Move {
	var minMoves []*board.Move

	if logger.Enabled() {
		logger.LogLn(b.Serialize())
		b.Print()
	}

	possibleMoves := b.AllMoves()
	unseenMoves := []*board.Move{}
	unseenBoards := []*board.Board{}

	logger.LogF("%v: possible moves\n", playNum)
	for _, move := range possibleMoves {
		nb := b.Apply(move)

		logger.LogF("%v:  %+v (to %v)", playNum, move, nb.Serialize())

		serialized := nb.Serialize()
		if seenMoves, found := seen[serialized]; found {
			if len(seenMoves) == 0 {
				logger.LogLn("  (already seen; in progress)")
			} else {
				logger.LogF("  (already seen; success in %d)\n", seenMoves)
				if minMoves == nil || len(seenMoves)+1 < len(minMoves) {
					minMoves = []*board.Move{move}
					minMoves = append(minMoves, seenMoves...)
				}
			}
			continue
		}
		logger.LogLn()

		unseenMoves = append(unseenMoves, move)
		unseenBoards = append(unseenBoards, nb)
	}

	for i, move := range unseenMoves {
		nb := unseenBoards[i]
		logger.LogF("%v: playing %+v (%v)\n", playNum, move, nb.Serialize())

		if nb.Success() {
			logger.LogF("%v: success\n", playNum)
			minMoves = []*board.Move{move}
			continue
		}

		serialized := nb.Serialize()
		seen[serialized] = []*board.Move{}

		globalPlayNum++
		logger.LogF("%v: recursing with playNum %v\n", playNum, globalPlayNum)

		successMoves := doPlay(nb, seen, globalPlayNum)
		if successMoves != nil {
			if minMoves == nil || len(successMoves)+1 < len(minMoves) {
				minMoves = []*board.Move{move}
				minMoves = append(minMoves, successMoves...)
			}

			seen[serialized] = successMoves
		}
	}

	logger.LogF("%v: returning; %v\n", playNum, minMoves)
	return minMoves
}

func Play(b *board.Board) ([]*board.Move, map[string][]*board.Move) {
	// Tracks parts of the move space that we've already
	// visited. If the value is -1, we don't know whether this
	// board leads to success. Otherwise, it's the number of moves
	// to get to success from that board.
	seen := map[string][]*board.Move{}
	seen[b.Serialize()] = []*board.Move{}

	moves := doPlay(b, seen, globalPlayNum)
	if moves != nil {
		seen[b.Serialize()] = moves
	}

	return moves, seen
}

func Audit(b *board.Board, moves []*board.Move) {
	b.Print()
	fmt.Println()
	for i, move := range moves {
		fmt.Printf("%d: %+v\n", i+1, move)
		b = b.Apply(move)
		b.Print()
		fmt.Println()
	}
}
