package spells

import (
	"character"
	"logger"
)

type Spell interface {
	Name() string
	Cost() int

	Activate(player *character.Player, boss *character.Boss)
	IsActive() bool
	Reset()

	TurnStart(player *character.Player, boss *character.Boss)
}

type spellBase struct {
	name string
	cost int
}

func (s *spellBase) Name() string { return s.name }
func (s *spellBase) Cost() int    { return s.cost }

func (s *spellBase) IsActive() bool { return false }
func (s *spellBase) Reset()         {}

func (s *spellBase) TurnStart(player *character.Player, boss *character.Boss) {}

type magicMissile struct {
	spellBase
	bossDamage int
}

func NewMagicMissile() Spell {
	return &magicMissile{
		spellBase:  spellBase{name: "MagicMissile", cost: 53},
		bossDamage: 4,
	}
}

func (s *magicMissile) Activate(player *character.Player, boss *character.Boss) {
	logger.LogF("Player casts MagicMissile, dealing %v damage.\n", s.bossDamage)
	boss.HP -= s.bossDamage
}

type drain struct {
	spellBase
	bossDamage int
	playerHeal int
}

func NewDrain() Spell {
	return &drain{
		spellBase:  spellBase{name: "Drain", cost: 73},
		bossDamage: 2,
		playerHeal: 2,
	}
}

func (s *drain) Activate(player *character.Player, boss *character.Boss) {
	logger.LogF("Player casts Drain, dealing %v damage and healing %v hit points.\n",
		s.bossDamage, s.playerHeal)
	boss.HP -= s.bossDamage
	player.HP += s.playerHeal
}

type effectBase struct {
	spellBase
	activationDuration int

	numTurnsLeft int
}

func (e *effectBase) Activate(player *character.Player, boss *character.Boss) {
	logger.LogF("Player casts %v.\n", e.Name())
	e.numTurnsLeft = e.activationDuration
}

func (e *effectBase) IsActive() bool { return e.numTurnsLeft > 0 }

func (e *effectBase) Reset() {
	e.numTurnsLeft = 0
}

func (e *effectBase) decrementTurnsLeft() {
	if e.numTurnsLeft <= 0 {
		panic("playerTurn with no turns left")
	}

	e.numTurnsLeft--
	logger.LogF("%v's timer is now %v.\n", e.Name(), e.numTurnsLeft)
	if e.numTurnsLeft == 0 {
		logger.LogF("%v wears off.\n", e.Name())
	}
}

type shield struct {
	effectBase
	playerArmorBoost int
}

func NewShield() Spell {
	return &shield{
		effectBase: effectBase{
			spellBase:          spellBase{"Shield", 113},
			activationDuration: 6,
		},
		playerArmorBoost: 7,
	}
}

func (e *shield) Activate(player *character.Player, boss *character.Boss) {
	e.effectBase.Activate(player, boss)
	logger.LogF("Shield activation increases armor by %v.\n", e.playerArmorBoost)
	player.Armor += e.playerArmorBoost
}

func (e *shield) TurnStart(player *character.Player, boss *character.Boss) {
	e.decrementTurnsLeft()
	if e.numTurnsLeft == 0 {
		logger.LogF("Deactivated shield decreases armor by %v.\n", e.playerArmorBoost)
		player.Armor -= e.playerArmorBoost
	}
}

type poison struct {
	effectBase
	bossDamage int
}

func NewPoison() Spell {
	return &poison{
		effectBase: effectBase{
			spellBase:          spellBase{"Poison", 173},
			activationDuration: 6,
		},
		bossDamage: 3,
	}
}

func (e *poison) TurnStart(player *character.Player, boss *character.Boss) {
	logger.LogF("Poison deals %v damage.\n", e.bossDamage)
	boss.HP -= e.bossDamage
	e.decrementTurnsLeft()
}

type recharge struct {
	effectBase
	manaBoost int
}

func NewRecharge() Spell {
	return &recharge{
		effectBase: effectBase{
			spellBase:          spellBase{"Recharge", 229},
			activationDuration: 5,
		},
		manaBoost: 101,
	}
}

func (e *recharge) TurnStart(player *character.Player, boss *character.Boss) {
	logger.LogF("Recharge provides %v mana.\n", e.manaBoost)
	player.Mana += e.manaBoost
	e.decrementTurnsLeft()
}

func Names(spells []Spell) []string {
	names := make([]string, len(spells))
	for i := range spells {
		names[i] = spells[i].Name()
	}
	return names
}
