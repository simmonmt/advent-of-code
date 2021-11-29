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

package game

import (
	"os"
	"reflect"
	"testing"

	"character"
	"logger"
	"spells"
)

// First example game from AoC
func TestFirstExample(t *testing.T) {
	player := character.NewPlayer(10, 250)
	boss := character.NewBoss(13, 8)

	magicMissile := spells.NewMagicMissile()
	poison := spells.NewPoison()

	allSpells := []spells.Spell{
		magicMissile,
		poison,
	}

	toCast := []spells.Spell{poison, magicMissile}

	spellProvider := spells.NewSeededProvider(toCast)

	bossDead, manaUsed, spellsCast := Run(*player, *boss, allSpells, spellProvider)
	if !bossDead || manaUsed != 226 || !reflect.DeepEqual(spellsCast, toCast) {
		t.Errorf("Run(...) = %v, %v, %v, want %v, %v, %v",
			bossDead, manaUsed, spells.Names(spellsCast),
			true, 226, spells.Names(toCast))
	}
}

// Second example game from AoC
func TestSecondExample(t *testing.T) {
	player := character.NewPlayer(10, 250)
	boss := character.NewBoss(14, 8)

	drain := spells.NewDrain()
	magicMissile := spells.NewMagicMissile()
	poison := spells.NewPoison()
	recharge := spells.NewRecharge()
	shield := spells.NewShield()

	allSpells := []spells.Spell{
		drain,
		magicMissile,
		poison,
		recharge,
		shield,
	}

	toCast := []spells.Spell{recharge, shield, drain, poison, magicMissile}

	spellProvider := spells.NewSeededProvider(toCast)

	bossDead, manaUsed, spellsCast := Run(*player, *boss, allSpells, spellProvider)
	if !bossDead || manaUsed != 641 || !reflect.DeepEqual(spellsCast, toCast) {
		t.Errorf("Run(...) = %v, %v, %v, want %v, %v, %v",
			bossDead, manaUsed, spells.Names(spellsCast),
			true, 641, spells.Names(toCast))
	}
}

// Recast an effect in the same round as it expires
func TestImmediateReuse(t *testing.T) {
	player := character.NewPlayer(9999, 9999)
	boss := character.NewBoss(28, 8)

	magicMissile := spells.NewMagicMissile()
	poison := spells.NewPoison()

	allSpells := []spells.Spell{
		magicMissile,
		poison,
	}

	toCast := []spells.Spell{poison, magicMissile, magicMissile, poison}

	spellProvider := spells.NewSeededProvider(toCast)

	bossDead, manaUsed, spellsCast := Run(*player, *boss, allSpells, spellProvider)
	if !bossDead || manaUsed != 452 || !reflect.DeepEqual(toCast, spellsCast) {
		t.Errorf("Run(...) = %v, %v, %v, want %v, %v, %v",
			bossDead, manaUsed, spells.Names(spellsCast),
			true, 641, spells.Names(toCast))
	}
}

func TestPlayerDies(t *testing.T) {
	player := character.NewPlayer(8, 250)
	boss := character.NewBoss(13, 8)

	magicMissile := spells.NewMagicMissile()
	poison := spells.NewPoison()

	allSpells := []spells.Spell{
		magicMissile,
		poison,
	}

	toCast := []spells.Spell{poison, magicMissile}
	expectedSpellsCast := []spells.Spell{poison}

	spellProvider := spells.NewSeededProvider(toCast)

	bossDead, manaUsed, spellsCast := Run(*player, *boss, allSpells, spellProvider)
	if bossDead || manaUsed != 173 || !reflect.DeepEqual(expectedSpellsCast, spellsCast) {
		t.Errorf("Run(...) = %v, %v, %v, want %v, %v, %v",
			bossDead, manaUsed, spells.Names(spellsCast),
			false, 173, spells.Names(expectedSpellsCast))
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
