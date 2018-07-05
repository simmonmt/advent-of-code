package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	logging            = flag.Bool("verbose", false, "enable logging")
	seedSpellNamesFlag = flag.String("seed_spells", "", "If specified, use these spells only")
)

func logLn(a ...interface{}) {
	if *logging {
		fmt.Println(a...)
	}
}

func logF(msg string, a ...interface{}) {
	if *logging {
		fmt.Printf(msg, a...)
	}
}

func pluralize(s bool) string {
	if s {
		return "s"
	} else {
		return ""
	}
}

type Character struct {
	Name            string
	HP, Armor, Mana int
	Damage          int
}

func (c *Character) Print() {
	logF("- %v has %d hit point%s, %d armor, %d mana\n",
		c.Name, c.HP, pluralize(c.HP != 1), c.Armor, c.Mana)
}

type Spell interface {
	Name() string
	Cost() int

	Activate(player *Character)
	IsActive() bool
	Reset()

	PlayerTurnStart(player, boss *Character)
	BossTurnStart(player, boss *Character)
	TurnEnd(player *Character)
}

type spellBase struct {
	active bool
}

func (s *spellBase) Activate(player *Character) {
	s.active = true
}

func (s *spellBase) IsActive() bool { return s.active }
func (s *spellBase) Reset()         { s.active = false }

func (s *spellBase) playerTurnStart(player *Character) {
	if !s.active {
		panic("PlayerTurn while inactive")
	}
	s.active = false
}

func (s *spellBase) BossTurnStart(player, boss *Character) {
	// Spells should never be active during the boss turn; they
	// should deactivate after player's turn.
	panic("spellBase BossTurn")
}

func (s *spellBase) TurnEnd(player *Character) {}

type SpellMagicMissile struct {
	spellBase
}

func (s *SpellMagicMissile) Name() string { return "MagicMissile" }
func (s *SpellMagicMissile) Cost() int    { return 53 }

func (s *SpellMagicMissile) PlayerTurnStart(player, boss *Character) {
	s.spellBase.playerTurnStart(player)

	damage := 4
	logF("Player casts MagicMissile, dealing %v damage.\n", damage)
	boss.HP -= damage
}

type SpellDrain struct {
	spellBase
}

func (s *SpellDrain) Name() string { return "Drain" }
func (s *SpellDrain) Cost() int    { return 73 }

func (s *SpellDrain) PlayerTurnStart(player, boss *Character) {
	s.spellBase.playerTurnStart(player)

	damage := 2
	heal := 2
	logF("Player casts Drain, dealing %v damage and healing %v hit points.\n", damage, heal)
	boss.HP -= damage
	player.HP += heal
}

type effectBase struct {
	numTurnsLeft int
	starting     bool
}

func (e *effectBase) activate(numTurnsLeft int) {
	e.numTurnsLeft = numTurnsLeft
	e.starting = true
}

func (e *effectBase) IsActive() bool { return e.numTurnsLeft > 0 }

func (e *effectBase) Reset() {
	e.numTurnsLeft = 0
	e.starting = false
}

func (e *effectBase) playerTurnStart() bool {
	if e.starting {
		e.starting = false
		return false
	}

	if e.numTurnsLeft <= 0 {
		panic("playerTurn with no turns left")
	}

	return true
}

func (e *effectBase) BossTurnStart(player, boss *Character) {}

type EffectShield struct {
	effectBase
}

func (e *EffectShield) Name() string { return "Shield" }
func (e *EffectShield) Cost() int    { return 113 }

func (e *EffectShield) armorChange() int { return 7 }

func (e *EffectShield) Activate(player *Character) {
	e.activate(7)

	logF("Player casts Shield, increasing armor by %v.\n", e.armorChange())
	player.Armor += e.armorChange()
}

func (e *EffectShield) PlayerTurnStart(player, boss *Character) {
	if !e.playerTurnStart() {
		return
	}
}

func (e *EffectShield) TurnEnd(player *Character) {
	e.numTurnsLeft--
	logF("Shield's timer is now %v.\n", e.numTurnsLeft)
	if e.numTurnsLeft == 0 {
		logF("Shield wears off, decreasing armor by %v.\n", e.armorChange())
		player.Armor -= e.armorChange()
	}
}

type EffectPoison struct {
	effectBase
}

func (e *EffectPoison) Name() string { return "Poison" }
func (e *EffectPoison) Cost() int    { return 173 }

func (e *EffectPoison) Activate(player *Character) {
	e.activate(7)
	logF("Player casts Poison.\n")
}

func (e *EffectPoison) act(boss *Character) {
	damage := 3
	logF("Poison deals %v damage.\n", damage)
	boss.HP -= damage
}

func (e *EffectPoison) PlayerTurnStart(player, boss *Character) {
	if !e.playerTurnStart() {
		return
	}

	e.act(boss)
}

func (e *EffectPoison) BossTurnStart(player, boss *Character) {
	e.act(boss)
}

func (e *EffectPoison) TurnEnd(player *Character) {
	e.numTurnsLeft--
	logF("Poison's timer is now %v.\n", e.numTurnsLeft)
	if e.numTurnsLeft == 0 {
		logF("Poison wears off.")
	}
}

type EffectRecharge struct {
	effectBase
}

func (e *EffectRecharge) Name() string { return "Recharge" }
func (e *EffectRecharge) Cost() int    { return 229 }

func (e *EffectRecharge) Activate(player *Character) {
	e.activate(6)
	logF("Player casts Recharge.\n")
}

func (e *EffectRecharge) act(player *Character) {
	manaIncr := 101
	logF("Recharge provides %v mana.\n", manaIncr)
	player.Mana += manaIncr
}

func (e *EffectRecharge) PlayerTurnStart(player, boss *Character) {
	if !e.playerTurnStart() {
		return
	}

	e.act(player)
}

func (e *EffectRecharge) BossTurnStart(player, boss *Character) {
	e.act(player)
}

func (e *EffectRecharge) TurnEnd(player *Character) {
	e.numTurnsLeft--
	logF("Recharge's timer is now %v.\n", e.numTurnsLeft)
	if e.numTurnsLeft == 0 {
		logF("Recharge wears off.\n")
	}
}

type SpellProvider interface {
	NextSpell(availMana int) Spell
}

type RandSpellProvider struct {
	allSpells  []Spell
	randSource *rand.Rand
}

func (p *RandSpellProvider) NextSpell(availMana int) Spell {
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

func NewRandSpellProvider(allSpells []Spell) SpellProvider {
	return &RandSpellProvider{
		allSpells:  allSpells,
		randSource: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type SeededSpellProvider struct {
	spells    []Spell
	nextSpell int
}

func (p *SeededSpellProvider) NextSpell(availMana int) Spell {
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

func NewSeededSpellProvider(spells []Spell) SpellProvider {
	return &SeededSpellProvider{
		spells:    spells,
		nextSpell: 0,
	}
}

func bossAttack(player, boss *Character) {
	damage := boss.Damage - player.Armor
	if damage < 1 {
		damage = 1
	}
	logF("Boss attacks for %v (boss damage %v, player armor %v)\n", damage, boss.Damage, player.Armor)
	player.HP -= damage
}

func nextSpell(spellProvider SpellProvider, availMana int) Spell {
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

func runOneGame(player, boss Character, allSpells []Spell, spellProvider SpellProvider) (bossDead bool, manaUsed int, spellsCast []Spell) {
	manaUsed = 0
	spellsCast = []Spell{}

	for i := 0; ; i++ {
		if i > 0 {
			logLn()
		}

		playerTurn := i%2 == 0

		if playerTurn {
			logLn("-- Player turn")
		} else {
			logLn("-- Boss turn")
		}

		player.Print()
		boss.Print()

		if playerTurn {
			spellToCast := nextSpell(spellProvider, player.Mana)
			if spellToCast == nil {
				logF("Player mana %v too low; killing player.", player.Mana)
				player.HP = 0
			} else {
				player.Mana -= spellToCast.Cost()
				manaUsed += spellToCast.Cost()
				spellsCast = append(spellsCast, spellToCast)
				spellToCast.Activate(&player)

				for _, spell := range allSpells {
					if spell.IsActive() {
						spell.PlayerTurnStart(&player, &boss)
					}
				}
			}
		} else {
			for _, spell := range allSpells {
				if spell.IsActive() {
					spell.BossTurnStart(&player, &boss)
				}
			}
			if boss.HP > 0 {
				bossAttack(&player, &boss)
			}
		}
		for _, spell := range allSpells {
			if spell.IsActive() {
				spell.TurnEnd(&player)
			}
		}

		if player.HP <= 0 {
			logLn("Player is dead.")
			bossDead = false
			return
		} else if boss.HP <= 0 {
			logLn("Boss is dead.")
			bossDead = true
			return
		}
	}

	panic("unreached")
}

func spellNames(spells []Spell) []string {
	names := make([]string, len(spells))
	for i := range spells {
		names[i] = spells[i].Name()
	}
	return names
}

func testSpellProvider(allSpells []Spell, spellProvider SpellProvider) {
	fmt.Println("-- 10 spells, no limit")
	for i := 0; i < 10; i++ {
		spell := spellProvider.NextSpell(9999999)
		fmt.Println(spell.Name())
	}

	fmt.Println("-- 10 spells, limit 150")
	for i := 0; i < 10; i++ {
		spell := spellProvider.NextSpell(150)
		fmt.Println(spell.Name())
	}

	zero := allSpells[0]
	zero.Activate(&Character{})
	one := allSpells[1]
	one.Activate(&Character{})
	fmt.Printf("-- 10 spells, no limit activated %v %v\n", zero.Name(), one.Name())
	for i := 0; i < 10; i++ {
		spell := spellProvider.NextSpell(9999999)
		fmt.Println(spell.Name())
	}

	fmt.Printf("-- 10 spells, 150 activated %v %v\n", zero.Name(), one.Name())
	for i := 0; i < 10; i++ {
		spell := spellProvider.NextSpell(150)
		fmt.Println(spell.Name())
	}

	for {
	}
}

func main() {
	flag.Parse()

	player := Character{"Player", 50, 0, 500, 0}
	boss := Character{"Boss", 51, 0, 0, 9}

	allSpells := []Spell{
		&SpellMagicMissile{},
		&SpellDrain{},
		&EffectShield{},
		&EffectPoison{},
		&EffectRecharge{},
	}

	allSpellNames := map[string]Spell{}
	for _, spell := range allSpells {
		allSpellNames[spell.Name()] = spell
	}

	var spellProvider SpellProvider
	oneGame := false
	if *seedSpellNamesFlag == "" {
		spellProvider = NewRandSpellProvider(allSpells)
	} else {
		seedSpellNames := strings.Split(*seedSpellNamesFlag, ",")
		seedSpells := []Spell{}
		for _, name := range seedSpellNames {
			if spell, ok := allSpellNames[name]; !ok {
				log.Fatalf("unknown seed spell name %v", name)
			} else {
				seedSpells = append(seedSpells, spell)
			}
		}

		spellProvider = NewSeededSpellProvider(seedSpells)
		oneGame = true
	}

	//testSpellProvider(allSpells, spellProvider)

	humanPrinter := message.NewPrinter(language.English)

	minManaUsed := -1
	for i := 0; i == 0 || !oneGame; i++ {
		if i%1000000 == 0 {
			humanPrinter.Printf("game %v\n", i)
		}
		bossDied, manaUsed, spellsCast := runOneGame(player, boss, allSpells, spellProvider)
		if bossDied && (minManaUsed == -1 || manaUsed < minManaUsed) {
			fmt.Printf("game %v: bossDied := %v, manaUsed := %v, spells := %v\n",
				i, bossDied, manaUsed, strings.Join(spellNames(spellsCast), ","))
			minManaUsed = manaUsed
		}

		for _, spell := range allSpells {
			spell.Reset()
		}
	}

	//testSpell(&player, &boss, &SpellMagicMissile{})
	//testSpell(&player, &boss, &SpellDrain{})
	//testSpell(&player, &boss, &EffectShield{})
	//testSpell(&player, &boss, &EffectPoison{})
	//testSpell(&player, &boss, &EffectRecharge{})
}
