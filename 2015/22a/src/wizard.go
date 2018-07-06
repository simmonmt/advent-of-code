package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"character"
	"game"
	"logger"
	"spells"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	logging            = flag.Bool("verbose", false, "enable logging")
	playerStartingHP   = flag.Int("player_hp", 50, "player starting HP")
	playerStartingMana = flag.Int("player_mana", 500, "player starting mana")
	bossStartingHP     = flag.Int("boss_hp", 51, "boss starting HP")
	bossDamage         = flag.Int("boss_damage", 9, "boss damage")
	seedSpellNamesFlag = flag.String("seed_spells", "", "If specified, use these spells only")
)

func main() {
	flag.Parse()
	logger.Init(*logging)

	player := character.NewPlayer(*playerStartingHP, *playerStartingMana)
	boss := character.NewBoss(*bossStartingHP, *bossDamage)

	allSpells := []spells.Spell{
		spells.NewMagicMissile(),
		spells.NewDrain(),
		spells.NewShield(),
		spells.NewPoison(),
		spells.NewRecharge(),
	}

	allSpellNames := map[string]spells.Spell{}
	for _, spell := range allSpells {
		allSpellNames[spell.Name()] = spell
	}

	var spellProvider spells.Provider
	oneGame := false
	if *seedSpellNamesFlag == "" {
		spellProvider = spells.NewRandProvider(allSpells)
	} else {
		seedSpellNames := strings.Split(*seedSpellNamesFlag, ",")
		seedSpells := []spells.Spell{}
		for _, name := range seedSpellNames {
			if spell, ok := allSpellNames[name]; !ok {
				log.Fatalf("unknown seed spell name %v", name)
			} else {
				seedSpells = append(seedSpells, spell)
			}
		}

		spellProvider = spells.NewSeededProvider(seedSpells)
		oneGame = true
	}

	humanPrinter := message.NewPrinter(language.English)

	minManaUsed := -1
	for i := 0; i == 0 || !oneGame; i++ {
		if i%1000000 == 0 {
			humanPrinter.Printf("game %v\n", i)
		}
		bossDied, manaUsed, spellsCast := game.Run(*player, *boss, allSpells, spellProvider)
		if bossDied && (minManaUsed == -1 || manaUsed < minManaUsed) {
			fmt.Printf("game %v: bossDied := %v, manaUsed := %v, spells := %v\n",
				i, bossDied, manaUsed, strings.Join(spells.Names(spellsCast), ","))
			minManaUsed = manaUsed
		}

		for _, spell := range allSpells {
			spell.Reset()
		}
	}
}
