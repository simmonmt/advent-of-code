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
	"character"
	"logger"
	"spells"
)

func bossAttack(player *character.Player, boss *character.Boss) {
	damage := boss.Damage - player.Armor
	if damage < 1 {
		damage = 1
	}
	logger.LogF("Boss attacks for %v (boss damage %v, player armor %v)\n", damage, boss.Damage, player.Armor)
	player.HP -= damage
}

func nextSpell(spellProvider spells.Provider, availMana int) spells.Spell {
	spell := spellProvider.NextSpell(availMana)
	if spell == nil {
		return nil
	}
	if spell.Cost() > availMana {
		panic("expensive spell")
	}
	if spell.IsActive() {
		panic("already active")
	}
	return spell
}

func Run(player character.Player, boss character.Boss, allSpells []spells.Spell, spellProvider spells.Provider) (bossDead bool, manaUsed int, spellsCast []spells.Spell) {
	manaUsed = 0
	spellsCast = []spells.Spell{}

	for i := 0; ; i++ {
		if i > 0 {
			logger.LogLn()
		}

		playerTurn := i%2 == 0

		if playerTurn {
			logger.LogLn("-- Player turn")
		} else {
			logger.LogLn("-- Boss turn")
		}

		player.Print()
		boss.Print()

		for _, spell := range allSpells {
			if spell.IsActive() {
				spell.TurnStart(&player, &boss)
			}
		}

		if playerTurn {
			spellToCast := nextSpell(spellProvider, player.Mana)
			if spellToCast == nil {
				logger.LogF("Player mana %v too low; killing player.", player.Mana)
				player.HP = 0
			} else {
				player.Mana -= spellToCast.Cost()
				manaUsed += spellToCast.Cost()
				spellsCast = append(spellsCast, spellToCast)
				spellToCast.Activate(&player, &boss)
			}
		} else { // boss turn
			if boss.HP > 0 {
				bossAttack(&player, &boss)
			}
		}

		if player.HP <= 0 {
			logger.LogLn("Player is dead.")
			bossDead = false
			return
		} else if boss.HP <= 0 {
			logger.LogLn("Boss is dead.")
			bossDead = true
			return
		}
	}

	panic("unreached")
}
