package puzzle

import (
	"strconv"
	"testing"
)

func TestSolve(t *testing.T) {
	type TestCase struct {
		board *Board
		cost  int
	}

	testCases := []TestCase{
		TestCase{NewBoard(map1), 27},
		TestCase{NewBoard(map2), 397},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			allPaths := FindAllPathsFromAllPortals(tc.board)

			start := tc.board.Gate("AA").PortalOut()
			end := tc.board.Gate("ZZ").GateOut()

			//logger.Init(true)

			cost, found := Solve(tc.board, allPaths, start, end)
			if !found || cost != tc.cost {
				t.Errorf("cost, found = %v, %v, want %v, true",
					cost, found, tc.cost)
			}

			//logger.Init(false)
		})
	}

}
