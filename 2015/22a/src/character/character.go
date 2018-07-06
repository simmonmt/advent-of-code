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
