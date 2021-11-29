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

package main

import "fmt"

type Character struct {
	Health, Damage, Armor int
}

func NewCharacter(health, armor, damage int) *Character {
	return &Character{
		Health: health,
		Armor:  armor,
		Damage: damage,
	}
}

type Value struct {
	Name                string
	Cost, Damage, Armor int
}

var (
	weapons = []Value{
		Value{"Dagger", 8, 4, 0},
		Value{"Shortsword", 10, 5, 0},
		Value{"Warhammer", 25, 6, 0},
		Value{"Longsword", 40, 7, 0},
		Value{"Greataxe", 74, 8, 0},
	}

	armors = []Value{
		Value{"None", 0, 0, 0},
		Value{"Leather", 13, 0, 1},
		Value{"Chainmail", 31, 0, 2},
		Value{"Splintmail", 53, 0, 3},
		Value{"Bandedmail", 75, 0, 4},
		Value{"Platemail", 102, 0, 5},
	}

	rings = []Value{
		Value{"None", 0, 0, 0},
		Value{"Damage+1", 25, 1, 0},
		Value{"Damage+2", 50, 2, 0},
		Value{"Damage+3", 100, 3, 0},
		Value{"Defense+1", 20, 0, 1},
		Value{"Defense+2", 40, 0, 2},
		Value{"Defense+3", 80, 0, 3},
	}
)

type Sequence struct {
	cur    int
	name   string
	values []Value
}

func NewSequence(name string, values []Value) *Sequence {
	return &Sequence{
		cur:    0,
		name:   name,
		values: values,
	}
}

func (s *Sequence) Next() (done bool) {
	if s.cur == len(s.values)-1 {
		done = true
		return
	}

	s.cur++
	done = false
	return
}

func (s *Sequence) Reset() {
	s.cur = 0
}

func (s *Sequence) Name() string {
	return s.name
}

func (s *Sequence) Value() *Value {
	return &s.values[s.cur]
}

type Seqs []*Sequence

func (s *Seqs) Next() (done bool) {
	for _, seq := range *s {
		if done := seq.Next(); !done {
			return false
		}
		seq.Reset()
	}
	return true
}

func (s *Seqs) Print() {
	for _, seq := range *s {
		fmt.Printf("%s: %s ", seq.Name(), seq.Value().Name)
	}
	fmt.Println()
}

func Fight(boss, player Character) bool {
	for {
		playerDealsDamage := player.Damage - boss.Armor
		bossDealsDamage := boss.Damage - player.Armor

		boss.Health -= playerDealsDamage
		// fmt.Printf("The player deals %v damage, boss down to %v hp\n", playerDealsDamage, boss.Health)
		if boss.Health <= 0 {
			return true
		}

		player.Health -= bossDealsDamage
		// fmt.Printf("The boss deals %v damage, player down to %v hp\n", bossDealsDamage, player.Health)
		if player.Health <= 0 {
			return false
		}
	}
}

func main() {
	weaponSeq := NewSequence("Weapon", weapons)
	armorSeq := NewSequence("Armor", armors)
	ring1Seq := NewSequence("Ring1", rings)
	ring2Seq := NewSequence("Ring2", rings)

	seqs := Seqs([]*Sequence{weaponSeq, armorSeq, ring1Seq, ring2Seq})

	boss := NewCharacter(103, 2, 9)

	lowestCost := 0
	for done := false; !done; done = seqs.Next() {
		if ring1Seq.Value().Name == ring2Seq.Value().Name {
			continue
		}

		seqs.Print()

		config := make([]*Value, len(seqs))
		for i, seq := range seqs {
			config[i] = seq.Value()
		}

		cost, armor, damage := 0, 0, 0
		for _, val := range config {
			cost += val.Cost
			armor += val.Armor
			damage += val.Damage
		}

		health := 100
		result := Fight(*boss, *NewCharacter(health, armor, damage))
		if result && (lowestCost == 0 || cost < lowestCost) {
			fmt.Printf("newest low cost %v\n", cost)
			lowestCost = cost
		}
	}

	fmt.Printf("finished iteration; lowest cost %v\n", lowestCost)
}
