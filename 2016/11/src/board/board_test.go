package board

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"testing"

	"logger"
	"object"
)

func TestValidMove(t *testing.T) {
	type TestCase struct {
		onSrcFloor, onDestFloor []object.Object
		cands                   []object.Object
		expectedResult          bool
	}

	testCases := []TestCase{
		TestCase{
			onSrcFloor:     []object.Object{object.Microchip(1)},
			onDestFloor:    []object.Object{object.Generator(1)},
			cands:          []object.Object{object.Microchip(1)},
			expectedResult: true,
		},
		TestCase{
			onSrcFloor:     []object.Object{object.Microchip(1)},
			onDestFloor:    []object.Object{},
			cands:          []object.Object{object.Microchip(1)},
			expectedResult: true,
		},
		TestCase{
			onSrcFloor:     []object.Object{object.Generator(1)},
			onDestFloor:    []object.Object{object.Microchip(1)},
			cands:          []object.Object{object.Generator(1)},
			expectedResult: true,
		},
		TestCase{
			onSrcFloor:     []object.Object{object.Microchip(1)},
			onDestFloor:    []object.Object{object.Generator(2)},
			cands:          []object.Object{object.Microchip(1)},
			expectedResult: false,
		},
		TestCase{
			onSrcFloor:     []object.Object{object.Microchip(1), object.Microchip(2)},
			onDestFloor:    []object.Object{object.Generator(1)},
			cands:          []object.Object{object.Microchip(1), object.Microchip(2)},
			expectedResult: false,
		},
		TestCase{
			onSrcFloor:     []object.Object{object.Microchip(1), object.Generator(2)},
			onDestFloor:    []object.Object{},
			cands:          []object.Object{object.Microchip(1), object.Generator(2)},
			expectedResult: false,
		},
		TestCase{
			onSrcFloor:     []object.Object{object.Microchip(1), object.Generator(1)},
			onDestFloor:    []object.Object{},
			cands:          []object.Object{object.Microchip(1), object.Generator(1)},
			expectedResult: true,
		},
		TestCase{
			onSrcFloor:     []object.Object{object.Microchip(1), object.Generator(1)},
			onDestFloor:    []object.Object{},
			cands:          []object.Object{object.Microchip(1)},
			expectedResult: true,
		},
		TestCase{
			onSrcFloor:     []object.Object{object.Microchip(1), object.Generator(1), object.Generator(2)},
			onDestFloor:    []object.Object{},
			cands:          []object.Object{object.Generator(1)},
			expectedResult: false,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := validMove(tc.onSrcFloor, tc.onDestFloor, tc.cands...)
			if r != tc.expectedResult {
				t.Errorf("validMove(%v, %v, %v) = %v, want %v",
					tc.onSrcFloor, tc.onDestFloor, tc.cands, r, tc.expectedResult)
			}
		})
	}
}

func TestBoard(t *testing.T) {
	type TestCase struct {
		board         *Board
		expectedMoves []*Move
	}

	testCases := []TestCase{
		TestCase{
			board: New(map[object.Object]uint8{
				object.Microchip(1): 1,
				object.Generator(1): 1,
				object.Microchip(2): 1,
				object.Generator(2): 1,
			}),
			expectedMoves: []*Move{
				NewMove(2, object.Microchip(1)),
				NewMove(2, object.Microchip(2)),
				NewMove(2, object.Generator(1), object.Microchip(1)),
				NewMove(2, object.Generator(1), object.Generator(2)),
				NewMove(2, object.Microchip(1), object.Microchip(2)),
				NewMove(2, object.Generator(2), object.Microchip(2)),
				// Does not include G2=>2 or G1=>2, both of
				// which would leave their respective microchips
				// without protection.
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			moves := tc.board.AllMoves()
			if !reflect.DeepEqual(moves, tc.expectedMoves) {
				t.Errorf("move mismatch: got %+v, want %+v", moves, tc.expectedMoves)
			}
		})
	}
}

func TestSerialize(t *testing.T) {
	in := New(map[object.Object]uint8{
		object.Microchip(1): 2,
		object.Generator(1): 3,
		object.Microchip(2): 4,
		object.Generator(2): 4,
	})

	ser := in.Serialize()
	if ser != "1A3a2B4b4" {
		in.Print()
		t.Errorf("serialize: want \"1A3a2B4b4\", got \"%v\"", ser)
	}

	deser, err := Deserialize(ser)
	if err != nil || !reflect.DeepEqual(deser, in) {
		t.Errorf("Deserialize(%v) = %+v, %v, want %+v, nil", ser, deser, err, in)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	logger.Init(true)

	os.Exit(m.Run())
}
