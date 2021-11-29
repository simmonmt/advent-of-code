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

package character

import "logger"
import "utils"

type Character struct {
	Name            string
	HP, Armor, Mana int
	Damage          int
}

func (c *Character) Print() {
	logger.LogF("- %v has %d hit point%s, %d armor, %d mana\n",
		c.Name, c.HP, utils.Pluralize(c.HP != 1), c.Armor, c.Mana)
}

type Player struct {
	Character
}

func NewPlayer(hp, mana int) *Player {
	return &Player{Character{"Player", hp, 0, mana, 0}}
}

type Boss struct {
	Character
}

func NewBoss(hp, damage int) *Boss {
	return &Boss{Character{"Boss", hp, 0, 0, damage}}
}
