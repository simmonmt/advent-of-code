package game

import (
	"flag"
	"os"
	"strconv"
	"testing"

	"board"
	"logger"
	"object"
)

func TestGame(t *testing.T) {
	type TestCase struct {
		board               *board.Board
		expectedMinNumSteps int
		expectedNumSeen     map[string]int
	}

	testCases := []TestCase{
		TestCase{
			board: board.NewWithElevatorStart(map[object.Object]uint8{
				object.Microchip(1): 3,
				object.Generator(1): 2,
			}, 2),
			expectedMinNumSteps: 2,
			expectedNumSeen: map[string]int{
				"3A3a3": 1,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			minSteps, seen := Play(tc.board)

			seenNum := map[string]int{}
			for k, v := range seen {
				seenNum[k] = len(v)
			}

			seenNumMatches := true
			for k, v := range tc.expectedNumSeen {
				if seenNum[k] != v {
					seenNumMatches = false
					break
				}
			}

			if tc.expectedMinNumSteps != len(minSteps) || !seenNumMatches {
				t.Errorf("Play(_) = %v=%v, %v; want %v, %v",
					len(minSteps), minSteps, seenNum, tc.expectedMinNumSteps, tc.expectedNumSeen)
			}
		})
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	logger.Init(true)

	os.Exit(m.Run())
}
