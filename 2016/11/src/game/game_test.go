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

func TestPlay(t *testing.T) {
	type TestCase struct {
		board               *board.Board
		expectedMinNumSteps int
		// expectedSeenNums    map[string]int
	}

	testCases := []TestCase{
		TestCase{
			board: board.NewWithElevatorStart(map[object.Object]uint8{
				object.Microchip(1): 3,
				object.Generator(1): 2,
			}, 3),
			expectedMinNumSteps: 3,
			// expectedSeenNums: map[string]int{
			// 	"2A2a2": 2,
			// 	"3A2a3": 3,
			// 	"3A3a3": 1,
			// },
		},
		TestCase{
			board: board.NewWithElevatorStart(map[object.Object]uint8{
				object.Microchip(1): 3,
				object.Generator(1): 3,
				object.Generator(2): 3,
				object.Microchip(2): 3,
			}, 3),
			expectedMinNumSteps: 5,
			// expectedSeenNums: map[string]int{
			// 	"3A3a3B3b3": 5,
			// },
		},
		TestCase{
			// Sample input from 11a
			board: board.New(map[object.Object]uint8{
				object.Microchip(1): 1, // hydrogen
				object.Generator(1): 2, // hydrogen
				object.Microchip(2): 1, // lithium
				object.Generator(2): 3, // lithium
			}),
			expectedMinNumSteps: 11,
			// expectedSeenNums:    map[string]int{},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			minSteps := Play(tc.board)

			// Audit(tc.board, minSteps)

			// seenNums := map[string]int{}
			// for k, v := range seen {
			// 	fmt.Printf("%v = %v\n", k, v)

			// 	if v.inProgress == false && v.moves != nil {
			// 		seenNums[k] = len(v.moves)
			// 	}
			// }

			// foundAllExpectedSeenNums := true
			// for k, v := range tc.expectedSeenNums {
			// 	if seenNums[k] != v {
			// 		foundAllExpectedSeenNums = false
			// 	}
			// }

			// fmt.Printf("started at %v, min steps = %v %v\n", tc.board.Serialize(), len(minSteps), minSteps)
			// for s, v := range seen {
			// 	fmt.Printf("%v: %+v\n", s, v)
			// }

			if tc.expectedMinNumSteps != len(minSteps)-1 {
				t.Errorf("Play(_) = %v=%v; want %v",
					len(minSteps), minSteps, tc.expectedMinNumSteps)
			}
		})
	}
}

// func TestDoPlay(t *testing.T) {
// 	type TestCase struct {
// 		board               *board.Board
// 		existingSeen        map[string]*SeenVal
// 		expectedMinNumSteps int
// 		expectedSeenNums    map[string]int
// 	}

// 	chip1 := object.Microchip(1)
// 	gen1 := object.Generator(1)

// 	testCases := []TestCase{
// 		TestCase{
// 			board: board.NewWithElevatorStart(map[object.Object]uint8{
// 				chip1: 1,
// 				gen1:  1,
// 			}, 1),
// 			existingSeen: map[string]*SeenVal{
// 				"2A2a2": &SeenVal{
// 					inProgress: false,
// 					moves: []*board.Move{
// 						board.NewMove(3, chip1, gen1),
// 						board.NewMove(4, chip1, gen1),
// 					},
// 				},
// 			},
// 			expectedMinNumSteps: 3,
// 			expectedSeenNums:    map[string]int{
// 			// "2A2a2": 2,
// 			// "3A2a3": 3,
// 			// "3A3a3": 1,
// 			},
// 		},
// 	}

// 	for i, tc := range testCases {
// 		t.Run(strconv.Itoa(i), func(t *testing.T) {
// 			seen := map[string]*SeenVal{}
// 			for k, v := range tc.existingSeen {
// 				seen[k] = v.Duplicate()
// 			}

// 			minSteps := doPlay(tc.board, seen, 1)

// 			Audit(tc.board, minSteps)

// 			seenNums := map[string]int{}
// 			for k, v := range seen {
// 				fmt.Printf("%v = %v\n", k, v)

// 				if v.inProgress == false && v.moves != nil {
// 					seenNums[k] = len(v.moves)
// 				}
// 			}

// 			foundAllExpectedSeenNums := true
// 			for k, v := range tc.expectedSeenNums {
// 				if seenNums[k] != v {
// 					foundAllExpectedSeenNums = false
// 				}
// 			}

// 			fmt.Printf("started at %v, min steps = %v %v\n", tc.board.Serialize(), len(minSteps), minSteps)
// 			for s, v := range seen {
// 				fmt.Printf("%v: %+v\n", s, v)
// 			}

// 			if tc.expectedMinNumSteps != len(minSteps) || !foundAllExpectedSeenNums {
// 				t.Errorf("Play(_) = %v=%v, %v; want %v, %v",
// 					len(minSteps), minSteps, seenNums, tc.expectedMinNumSteps, tc.expectedSeenNums)
// 			}
// 		})
// 	}
// }

func TestMain(m *testing.M) {
	flag.Parse()
	logger.Init(true)

	os.Exit(m.Run())
}
