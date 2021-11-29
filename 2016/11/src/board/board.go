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

package board

import (
	"fmt"
	"sort"

	"object"
)

type Board struct {
	ElevatorFloor uint8
	// rethink layout. it's hard to apply moves with this layout.
	Objs      []object.Object
	ObjFloors []uint8
}

func new(num int) *Board {
	return &Board{
		Objs:      make([]object.Object, num),
		ObjFloors: make([]uint8, num),
	}
}

func NewWithElevatorStart(objsToFloors map[object.Object]uint8, elevatorStart uint8) *Board {
	b := new(len(objsToFloors))
	b.ElevatorFloor = elevatorStart

	objs := []object.Object{}
	for obj := range objsToFloors {
		objs = append(objs, obj)
	}
	sort.Sort(object.Objects(objs))

	for i, obj := range objs {
		b.Objs[i] = obj
		b.ObjFloors[i] = objsToFloors[obj]
	}

	return b
}

func New(objsToFloors map[object.Object]uint8) *Board {
	return NewWithElevatorStart(objsToFloors, 1)
}

func (b *Board) duplicate() *Board {
	nb := new(len(b.Objs))
	nb.ElevatorFloor = b.ElevatorFloor
	copy(nb.Objs, b.Objs)
	copy(nb.ObjFloors, b.ObjFloors)
	return nb
}

func (b *Board) Apply(m *Move) *Board {
	nb := b.duplicate()
	nb.ElevatorFloor = m.dest

	for _, moveObj := range m.objs {
		for i, obj := range b.Objs {
			if obj == moveObj {
				nb.ObjFloors[i] = m.dest
				break
			}
		}
	}

	return nb
}

func (b *Board) Serialize() string {
	out := make([]byte, 1+len(b.Objs)*2)
	out[0] = '0' + byte(b.ElevatorFloor)

	outIdx := 1
	for i, obj := range b.Objs {
		floor := b.ObjFloors[i]

		out[outIdx] = obj.Serialize()
		out[outIdx+1] = '0' + byte(floor)
		outIdx += 2
	}

	return string(out)
}

func Deserialize(str string) (*Board, error) {
	if len(str)%2 != 1 {
		return nil, fmt.Errorf("invalid length %v in %v", len(str), str)
	}

	sArr := []byte(str)

	b := &Board{}
	b.ElevatorFloor = sArr[0] - '0'
	b.Objs = []object.Object{}
	b.ObjFloors = []uint8{}

	for i := 1; i < len(sArr); i += 2 {
		b.Objs = append(b.Objs, object.Deserialize(sArr[i]))
		b.ObjFloors = append(b.ObjFloors, uint8(sArr[i+1]-'0'))
	}

	return b, nil
}

func (b *Board) makeFloorContents() [5][]object.Object {
	floorContents := [5][]object.Object{}
	for i := range floorContents {
		floorContents[i] = []object.Object{}
	}
	for i, obj := range b.Objs {
		objFloor := b.ObjFloors[i]
		floorContents[objFloor] = append(floorContents[objFloor], obj)
	}
	return floorContents
}

func floorIsOk(objs []object.Object) bool {
	gens := map[int]object.Object{}
	chips := map[int]object.Object{}

	for _, obj := range objs {
		if obj.IsMicrochip() {
			chips[obj.Num()] = obj
		} else {
			gens[obj.Num()] = obj
		}
	}

	for num, _ := range chips {
		if _, found := gens[num]; !found {
			// unprotected chip; invalidate if there's
			// another generator
			if len(gens) > 0 {
				//fmt.Printf("invalid move %v to %v: %v would fry\n", cands, destFloor, chip)
				return false
			}
		}
	}

	return true
}

func (b *Board) Print() {
	floorContents := b.makeFloorContents()

	var floor uint8
	for floor = 4; floor > 0; floor-- {
		fmt.Printf("F%d ", floor)
		if floor == b.ElevatorFloor {
			fmt.Printf("E  ")
		} else {
			fmt.Printf(".  ")
		}

		for i, obj := range b.Objs {
			if floor == b.ObjFloors[i] {
				fmt.Printf("%s", obj)
				if i+1 < len(b.Objs) && floor == b.ObjFloors[i+1] && obj.Num() == b.Objs[i+1].Num() {
					fmt.Printf("-")
				} else {
					fmt.Printf(" ")
				}

			} else {
				fmt.Printf(".  ")
			}
		}

		if floorIsOk(floorContents[floor]) {
			fmt.Printf("v")
		} else {
			fmt.Printf("<<<<<<<<<--------- INVALID")
		}
		fmt.Println()
	}
}

func validMove(objsOnSrcFloor, objsOnDestFloor []object.Object, cands ...object.Object) bool {
	newDestFloor := make([]object.Object, len(objsOnDestFloor)+len(cands))
	copy(newDestFloor, objsOnDestFloor)
	copy(newDestFloor[len(objsOnDestFloor):], cands)
	if !floorIsOk(newDestFloor) {
		return false
	}

	newSrcFloor := make([]object.Object, len(objsOnSrcFloor)-len(cands))
	i := 0
	for _, obj := range objsOnSrcFloor {
		foundCand := false
		for _, cand := range cands {
			if obj == cand {
				foundCand = true
			}
		}
		if !foundCand {
			newSrcFloor[i] = obj
			i++
		}
	}
	if !floorIsOk(newSrcFloor) {
		return false
	}

	return true
}

func (b *Board) AllMoves() []*Move {
	floorContents := b.makeFloorContents()
	moves := []*Move{}

	moveableIdx := []int{}
	// All 1-piece moves
	for i, obj := range b.Objs {
		objFloor := b.ObjFloors[i]
		if b.ElevatorFloor != objFloor {
			continue
		}
		moveableIdx = append(moveableIdx, i)

		if b.ElevatorFloor != 4 {
			destFloor := b.ElevatorFloor + 1
			if validMove(floorContents[b.ElevatorFloor], floorContents[destFloor], obj) {
				moves = append(moves, NewMove(destFloor, obj))
			}
		}
		if b.ElevatorFloor != 1 {
			destFloor := b.ElevatorFloor - 1
			if validMove(floorContents[b.ElevatorFloor], floorContents[destFloor], obj) {
				moves = append(moves, NewMove(destFloor, obj))
			}
		}
	}

	// for _, i := range moveableIdx {
	// 	fmt.Printf("moveable: %v on %v\n", b.Objs[i], b.ObjFloors[i])
	// }

	// All 2-piece moves
	for i := 0; i < len(moveableIdx)-1; i++ {
		for j := i + 1; j < len(moveableIdx); j++ {
			obj1 := b.Objs[moveableIdx[i]]
			obj2 := b.Objs[moveableIdx[j]]

			if b.ElevatorFloor != 4 {
				destFloor := b.ElevatorFloor + 1
				if validMove(floorContents[b.ElevatorFloor], floorContents[destFloor], obj1, obj2) {
					moves = append(moves, NewMove(destFloor, obj1, obj2))
				}
			}
			if b.ElevatorFloor != 1 {
				destFloor := b.ElevatorFloor - 1
				if validMove(floorContents[b.ElevatorFloor], floorContents[destFloor], obj1, obj2) {
					moves = append(moves, NewMove(destFloor, obj1, obj2))
				}
			}
		}

	}

	return moves
}

func (b *Board) SuccessBoard() *Board {
	nb := b.duplicate()
	nb.ElevatorFloor = 4
	for i := range nb.ObjFloors {
		nb.ObjFloors[i] = 4
	}
	return nb
}
