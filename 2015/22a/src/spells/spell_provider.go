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

package spells

import (
	"math/rand"
	"time"
)

type Provider interface {
	NextSpell(availMana int) Spell
}

type RandProvider struct {
	allSpells  []Spell
	randSource *rand.Rand
}

func (p *RandProvider) NextSpell(availMana int) Spell {
	spell := p.allSpells[p.randSource.Int31n(int32(len(p.allSpells)))]
	if spell.Cost() <= availMana && !spell.IsActive() {
		return spell
	}

	usableSpells := make([]Spell, len(p.allSpells))
	numAdded := 0
	for _, spell := range p.allSpells {
		if spell.Cost() <= availMana && !spell.IsActive() {
			//fmt.Printf("# added %v %v (%v)\n", spell.Name(), spell.Cost(), availMana)
			usableSpells[numAdded] = spell
			numAdded++
		} else {
			//fmt.Printf("# skipped %v %v (%v)\n", spell.Name(), spell.Cost(), availMana)
		}
	}

	if numAdded == 0 {
		return nil
	}

	return usableSpells[p.randSource.Int31n(int32(numAdded))]
}

func NewRandProvider(allSpells []Spell) Provider {
	return &RandProvider{
		allSpells:  allSpells,
		randSource: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type SeededProvider struct {
	spells    []Spell
	nextSpell int
}

func (p *SeededProvider) NextSpell(availMana int) Spell {
	if p.nextSpell >= len(p.spells) {
		panic("seeded provider ran out of spells")
	}

	spell := p.spells[p.nextSpell]
	p.nextSpell++

	if spell.Cost() > availMana {
		panic("seeded provider returning expensive spell")
	}

	return spell
}

func NewSeededProvider(spells []Spell) Provider {
	return &SeededProvider{
		spells:    spells,
		nextSpell: 0,
	}
}
