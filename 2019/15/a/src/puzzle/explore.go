package puzzle

import (
	"fmt"

	"github.com/simmonmt/aoc/2019/common/pos"
)

type ExploreState struct {
	goBack Dir
	pos    pos.P2
}

type ExploreStack struct {
	s []*ExploreState
}

func NewExploreStack() *ExploreStack {
	return &ExploreStack{
		s: []*ExploreState{},
	}
}

func (s *ExploreStack) Push(state *ExploreState) {
	s.s = append(s.s, state)
}

func (s *ExploreStack) Peek() *ExploreState {
	if len(s.s) == 0 {
		return nil
	}
	return s.s[len(s.s)-1]
}

func (s *ExploreStack) Pop() {
	if len(s.s) == 0 {
		panic("stack empty")
	}
	s.s = s.s[0 : len(s.s)-1]
}

func (s *ExploreStack) Depth() int {
	return len(s.s)
}

type Candidate struct {
	dir Dir
	pos pos.P2
}

var (
	allDirs = []Dir{DIR_NORTH, DIR_SOUTH, DIR_EAST, DIR_WEST}
)

func findACandidate(b *Board, p pos.P2) *Candidate {
	for _, d := range allDirs {
		newPos := d.From(p)
		if b.Get(newPos) == TILE_UNKNOWN {
			return &Candidate{d, newPos}
		}
	}

	return nil
}

func findNeighbors(b *Board, p pos.P2, want Tile) []pos.P2 {
	found := []pos.P2{}
	for _, d := range allDirs {
		newPos := d.From(p)
		if b.Get(newPos) == want {
			found = append(found, newPos)
		}
	}

	return found
}

func Explore(b *Board, start pos.P2, moveTo func(curPos pos.P2, dir Dir) (pos.P2, Tile)) pos.P2 {
	goalPos := pos.P2{-1, -1}

	stateStack := NewExploreStack()
	stateStack.Push(&ExploreState{DIR_UNKNOWN, start})

	for step := 1; ; step++ {
		//fmt.Printf("step %d (depth %d) start:\n", step, stateStack.Depth())
		curState := stateStack.Peek()
		//PrintBoard(b, curState.pos)

		cand := findACandidate(b, curState.pos)
		if cand == nil {
			if curState.goBack == DIR_UNKNOWN {
				return goalPos
			}

			newPos, _ := moveTo(curState.pos, curState.goBack)
			stateStack.Pop()
			if !stateStack.Peek().pos.Equals(newPos) {
				panic(fmt.Sprintf("backed up to different pos"))
			}
			continue
		}

		p, newTile := moveTo(curState.pos, cand.dir)
		if newTile == TILE_GOAL {
			fmt.Printf("found goal with depth %d\n", stateStack.Depth())
		}
		b.Set(cand.pos, newTile)
		if newTile == TILE_GOAL {
			goalPos = cand.pos
		}
		if newTile != TILE_WALL {
			stateStack.Push(&ExploreState{
				goBack: cand.dir.Reverse(),
				pos:    p,
			})
		}
	}

	panic("unreachable")
}

func Fill(b *Board, start pos.P2) int {
	visited := map[pos.P2]int{start: 1}

	queue := []pos.P2{start}

	var round int
	for round = 0; ; round++ {
		//fmt.Printf("round %d: queue %v\n", round, queue)

		qlen := len(queue)
		for i := 0; i < qlen; i++ {
			if i >= len(queue) {
				break
			}

			cur := queue[i]
			visited[cur] = round

			//fmt.Printf("cur %v neighbors %v\n", cur, findNeighbors(b, cur, TILE_OPEN))
			for _, n := range findNeighbors(b, cur, TILE_OPEN) {
				if _, found := visited[n]; found {
					continue
				}

				queue = append(queue, n)
			}
		}

		queue = queue[qlen:]
		if len(queue) == 0 {
			break
		}
	}

	return round
}
