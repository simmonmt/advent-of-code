package solver

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"grid"
	"node"
)

func makeGrid(t *testing.T, w, h, goalX, goalY uint8, nodes []node.Node) *grid.Grid {
	g, err := grid.New(w, h, goalX, goalY, nodes)
	if err != nil {
		panic(err.Error())
	}
	return g
}

// func TestFindNewGoals(t *testing.T) {
// 	type TestCase struct {
// 		goalX, goalY uint8
// 		successes    [][2]uint8
// 	}

// 	testCases := []TestCase{
// 		TestCase{0, 0, [][2]uint8{[2]uint8{1, 0}, [2]uint8{0, 1}}},
// 		TestCase{1, 1, [][2]uint8{[2]uint8{0, 1}, [2]uint8{1, 0}, [2]uint8{2, 1}, [2]uint8{1, 2}}},
// 		TestCase{2, 2, [][2]uint8{[2]uint8{1, 2}, [2]uint8{2, 1}}},
// 	}

// 	for i, tc := range testCases {
// 		t.Run(strconv.Itoa(i), func(t *testing.T) {
// 			g := makeGrid(t, 3, 3, tc.goalX, tc.goalY, []node.Node{
// 				*node.New(1, 1), *node.New(2, 1), *node.New(3, 1),
// 				*node.New(1, 2), *node.New(2, 2), *node.New(3, 2),
// 				*node.New(1, 3), *node.New(2, 3), *node.New(3, 3),
// 			})

// 			newGs := findNewGoals(g)
// 			newGoals := [][2]uint8{}
// 			for _, newG := range newGs {
// 				x, y := newG.Goal()
// 				newGoals = append(newGoals, [2]uint8{x, y})
// 			}

// 			if !reflect.DeepEqual(newGoals, tc.successes) {
// 				t.Errorf("got new goals %v; want %v", newGoals, tc.successes)
// 			}
// 		})
// 	}
// }

func TestFindTransfers(t *testing.T) {
	g := makeGrid(t, 4, 4, 3, 3, []node.Node{
		*node.New(10, 2), *node.New(10, 2), *node.New(99, 99), *node.New(99, 99),
		*node.New(10, 2), *node.New(10, 2), *node.New(99, 99), *node.New(99, 99),
		*node.New(99, 99), *node.New(99, 99), *node.New(10, 3), *node.New(10, 3),
		*node.New(99, 99), *node.New(99, 99), *node.New(10, 4), *node.New(10, 4),
	})

	type TestCase struct {
		x, y              uint8
		expectedTransfers []transferDesc
	}

	testCases := []TestCase{
		TestCase{0, 0, []transferDesc{transferDesc{1, 0}, transferDesc{0, 1}}},
		TestCase{3, 3, []transferDesc{transferDesc{2, 3}, transferDesc{3, 2}}},
		TestCase{1, 1, []transferDesc{transferDesc{0, 1}, transferDesc{1, 0}}},
		TestCase{2, 2, []transferDesc{transferDesc{3, 2}, transferDesc{2, 3}}},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if transfers := findTransfers(g, tc.x, tc.y); !reflect.DeepEqual(transfers, tc.expectedTransfers) {
				t.Errorf("findTransfers(g, %v, %v) = %+v, want %+v", tc.x, tc.y, transfers, tc.expectedTransfers)
			}
		})
	}
}

func TestAllNeighbors(t *testing.T) {
	var w uint8 = 4
	var h uint8 = 4

	g := makeGrid(t, w, h, 3, 3, []node.Node{
		*node.New(2, 2), *node.New(10, 5), *node.New(99, 99), *node.New(99, 99),
		*node.New(99, 99), *node.New(99, 99), *node.New(99, 99), *node.New(99, 99),
		*node.New(99, 99), *node.New(99, 99), *node.New(99, 99), *node.New(99, 99),
		*node.New(99, 99), *node.New(99, 99), *node.New(10, 5), *node.New(2, 2),
	})

	n1 := makeGrid(t, w, h, 3, 3, []node.Node{
		*node.New(2, 0), *node.New(10, 7), *node.New(99, 99), *node.New(99, 99),
		*node.New(99, 99), *node.New(99, 99), *node.New(99, 99), *node.New(99, 99),
		*node.New(99, 99), *node.New(99, 99), *node.New(99, 99), *node.New(99, 99),
		*node.New(99, 99), *node.New(99, 99), *node.New(10, 5), *node.New(2, 2),
	})

	n2 := makeGrid(t, w, h, 2, 3, []node.Node{
		*node.New(2, 2), *node.New(10, 5), *node.New(99, 99), *node.New(99, 99),
		*node.New(99, 99), *node.New(99, 99), *node.New(99, 99), *node.New(99, 99),
		*node.New(99, 99), *node.New(99, 99), *node.New(99, 99), *node.New(99, 99),
		*node.New(99, 99), *node.New(99, 99), *node.New(10, 7), *node.New(2, 0),
	})

	expectedNeighbors := []*grid.Grid{n1, n2}

	helper := newHelper(w, h)
	serNeighbors := helper.AllNeighbors(string(g.Serialize()))
	neighbors := []*grid.Grid{}
	for i, sn := range serNeighbors {
		n, err := grid.Deserialize(w, h, []byte(sn))
		if err != nil {
			t.Errorf("failed to deserialize neighbor %v of %v: %v", i, len(serNeighbors), err)
			continue
		}
		neighbors = append(neighbors, n)
	}

	if !reflect.DeepEqual(neighbors, expectedNeighbors) {
		for i, n := range neighbors {
			fmt.Printf("Got %v:\n", i)
			n.Print()
		}

		t.Errorf("expected neighbors %+v, got %+v", neighbors, expectedNeighbors)
	}
}

func TestEstimateDistance(t *testing.T) {
	var w uint8 = 4
	var h uint8 = 4

	g1 := makeGrid(t, w, h, 3, 3, []node.Node{
		*node.New(1, 1), *node.New(1, 1), *node.New(1, 1), *node.New(1, 1),
		*node.New(1, 1), *node.New(1, 1), *node.New(1, 1), *node.New(1, 1),
		*node.New(1, 1), *node.New(1, 1), *node.New(1, 1), *node.New(1, 1),
		*node.New(1, 1), *node.New(1, 1), *node.New(1, 1), *node.New(1, 1),
	})

	g2 := makeGrid(t, w, h, 1, 0, []node.Node{
		*node.New(1, 1), *node.New(1, 1), *node.New(1, 1), *node.New(1, 1),
		*node.New(1, 1), *node.New(1, 1), *node.New(1, 1), *node.New(1, 1),
		*node.New(1, 1), *node.New(1, 1), *node.New(1, 1), *node.New(1, 1),
		*node.New(1, 1), *node.New(1, 1), *node.New(1, 1), *node.New(1, 1),
	})

	helper := newHelper(w, h)
	if dist := helper.EstimateDistance(string(g1.Serialize()), string(g2.Serialize())); dist != 5 {
		t.Errorf("EstimateDistance(g1, g2) = %v, want %v", dist, 5)
	}
	if dist := helper.EstimateDistance(string(g2.Serialize()), string(g1.Serialize())); dist != 5 {
		t.Errorf("EstimateDistance(g2, g1) = %v, want %v", dist, 5)
	}
}
