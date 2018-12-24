package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"intmath"
	"logger"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")

	//groupPattern = regexp.MustCompile(`(\d+) units each with (\d+) hit points(?: \(weak to ([^\)]+)\)) with an attack that does (\d+) ([^ ]+) damage at initiative (\d+)`)
	groupPattern = regexp.MustCompile(`^(\d+) units each with (\d+) hit points (?:\(([^\)]+)\) )?with an attack that does (\d+) ([^ ]+) damage at initiative (\d+)`)
	//17 units each with 5390 hit points (weak to radiation, bludgeoning) with an attack that does 4507 fire damage at initiative 2

)

type Side int

const (
	SIDE_IMMUNE Side = iota
	SIDE_INFECTION
)

func (s Side) String() string {
	switch s {
	case SIDE_IMMUNE:
		return "immune"
	case SIDE_INFECTION:
		return "infection"
	default:
		panic("unknown")
	}
}

type AttackType int

const (
	AT_BLUDGEONING AttackType = iota
	AT_COLD
	AT_FIRE
	AT_RADIATION
	AT_SLASHING
)

func AttackTypeFromString(str string) (AttackType, error) {
	switch strings.ToUpper(str) {
	case "BLUDGEONING":
		return AT_BLUDGEONING, nil
	case "COLD":
		return AT_COLD, nil
	case "FIRE":
		return AT_FIRE, nil
	case "RADIATION":
		return AT_RADIATION, nil
	case "SLASHING":
		return AT_SLASHING, nil
	default:
		return AT_FIRE, fmt.Errorf("unknown attack type %v", str)
	}
}

func (t AttackType) String() string {
	switch t {
	case AT_BLUDGEONING:
		return "bludgeoning"
	case AT_COLD:
		return "cold"
	case AT_FIRE:
		return "fire"
	case AT_RADIATION:
		return "radiation"
	case AT_SLASHING:
		return "slashing"
	default:
		panic("unknown")
	}
}

type Group struct {
	ID           int
	Side         Side
	NumUnits     int
	HP           int
	Initiative   int
	Immunes      []AttackType
	Weaknesses   []AttackType
	AttackDamage int
	AttackType   AttackType
}

func (g Group) String() string {
	return fmt.Sprintf("{%d S:%s Size:%v HP:%v EP:%v, Init:%v Imm:%v Weak:%v AD:%v AT:%v}",
		g.ID, g.Side, g.NumUnits, g.HP, g.EffectivePower(), g.Initiative, g.Immunes,
		g.Weaknesses, g.AttackDamage, g.AttackType)
}

func (g *Group) EffectivePower() int {
	return g.NumUnits * g.AttackDamage
}

type ByPowerAndInitiative []*Group

func (a ByPowerAndInitiative) Len() int      { return len(a) }
func (a ByPowerAndInitiative) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPowerAndInitiative) Less(i, j int) bool {
	epI, epJ := a[i].EffectivePower(), a[j].EffectivePower()
	if epI < epJ {
		return true
	} else if epI > epJ {
		return false
	}
	return a[i].Initiative < a[j].Initiative
}

type ByInitiative []*Group

func (a ByInitiative) Len() int      { return len(a) }
func (a ByInitiative) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByInitiative) Less(i, j int) bool {
	return a[i].Initiative < a[j].Initiative
}

type ByID []*Group

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

type FullID struct {
	Side Side
	ID   int
}

func FullIDFromGroup(g *Group) FullID {
	return FullID{g.Side, g.ID}
}

func parsePowers(str string) (immunes, weaknesses []AttackType, err error) {
	if str == "" {
		return []AttackType{}, []AttackType{}, nil
	}

	for _, spec := range strings.Split(str, "; ") {
		parts := strings.SplitN(spec, " ", 3)
		name := parts[0]
		typeList := parts[2]

		types := []AttackType{}
		for _, s := range strings.Split(typeList, ", ") {
			t, err := AttackTypeFromString(s)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse %v type %v: %v", name, s, err)
			}
			types = append(types, t)
		}

		if name == "weak" {
			weaknesses = types
		} else {
			immunes = types
		}
	}

	return
}

func parseGroup(line string) (*Group, error) {
	parts := groupPattern.FindStringSubmatch(line)
	if parts == nil {
		return nil, fmt.Errorf("match fail")
	}

	numUnits := intmath.AtoiOrDie(parts[1])
	hp := intmath.AtoiOrDie(parts[2])
	immunes, weaknesses, err := parsePowers(parts[3])
	if err != nil {
		return nil, err
	}
	attackDamage := intmath.AtoiOrDie(parts[4])
	attackType, err := AttackTypeFromString(parts[5])
	initiative := intmath.AtoiOrDie(parts[6])
	if err != nil {
		return nil, err
	}

	return &Group{
		NumUnits:     numUnits,
		HP:           hp,
		AttackDamage: attackDamage,
		AttackType:   attackType,
		Initiative:   initiative,
		Immunes:      immunes,
		Weaknesses:   weaknesses,
	}, nil
}

func readInput() (immune, infection []*Group, err error) {
	immune = []*Group{}
	infection = []*Group{}

	readingImmune := true
	scanner := bufio.NewScanner(os.Stdin)
	for lineNum := 1; scanner.Scan(); lineNum++ {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Immune") {
			readingImmune = true
			continue
		} else if strings.HasPrefix(line, "Infection") {
			readingImmune = false
			continue
		}

		group, err := parseGroup(line)
		if err != nil {
			return nil, nil, fmt.Errorf("%d: group fail: %v", lineNum, err)
		}

		if readingImmune {
			group.Side = SIDE_IMMUNE
			group.ID = len(immune) + 1
			immune = append(immune, group)
		} else {
			group.Side = SIDE_INFECTION
			group.ID = len(infection) + 1
			infection = append(infection, group)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("read failed: %v", err)
	}

	return immune, infection, nil
}

func contains(a []AttackType, at AttackType) bool {
	for _, elem := range a {
		if elem == at {
			return true
		}
	}
	return false
}

func calculateDamage(attacker, target *Group) int {
	if contains(target.Immunes, attacker.AttackType) {
		return 0
	}
	damage := attacker.EffectivePower()
	if contains(target.Weaknesses, attacker.AttackType) {
		damage *= 2
	}
	return damage
}

func selectTargets(attackers, targets []*Group) map[FullID]FullID {
	attacking := map[FullID]FullID{}
	beingAttacked := map[int]bool{}

	for _, attacker := range attackers {
		maxDamageCaused := 0
		cands := []*Group{}

		for _, target := range targets {
			if _, found := beingAttacked[target.ID]; found {
				continue
			}

			damageCaused := calculateDamage(attacker, target)
			if damageCaused > maxDamageCaused {
				maxDamageCaused = damageCaused
				cands = []*Group{target}
			} else if damageCaused == maxDamageCaused {
				cands = append(cands, target)
			}

			logger.LogF("%v group %v would deal defending group %v %v damage",
				attacker.Side, attacker.ID, target.ID, damageCaused)
		}

		if len(cands) == 0 || maxDamageCaused == 0 {
			continue
		}

		sort.Sort(sort.Reverse(ByPowerAndInitiative(cands)))
		target := cands[0]

		attacking[FullIDFromGroup(attacker)] = FullIDFromGroup(target)
		beingAttacked[target.ID] = true
	}

	return attacking
}

func dumpGroups(groups []*Group) {
	for _, g := range groups {
		fmt.Println(g)
	}
}

func attack(groups []*Group, toAttacks map[FullID]FullID) {
	groupsByID := map[FullID]*Group{}
	for _, g := range groups {
		groupsByID[FullIDFromGroup(g)] = g
	}

	for _, attacker := range groups {
		if attacker.NumUnits == 0 {
			logger.LogF("attacker already dead %s %v", attacker.Side, attacker.ID)
			continue
		}

		targetID, found := toAttacks[FullIDFromGroup(attacker)]
		if !found {
			logger.LogF("attacker %s %v has no target", attacker.Side, attacker.ID)
			continue
		}

		target := groupsByID[targetID]
		damage := calculateDamage(attacker, target)
		unitsKilled := damage / target.HP
		unitsKilled = intmath.IntMin(unitsKilled, target.NumUnits)

		logger.LogF("%v group %v attacks defending group %v, killing %v units",
			attacker.Side, attacker.ID, target.ID, unitsKilled)

		target.NumUnits -= unitsKilled
		if target.NumUnits == 0 {
			logger.LogF("defending group %v now dead", target.ID)
		}
	}
}

func removeDead(in []*Group) []*Group {
	out := []*Group{}
	for _, g := range in {
		if g.NumUnits > 0 {
			out = append(out, g)
		}
	}
	return out
}

func countUnits(groups []*Group) int {
	num := 0
	for _, g := range groups {
		num += g.NumUnits
	}
	return num
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	immune, infection, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	for len(immune) > 0 && len(infection) > 0 {
		sort.Sort(ByID(immune))
		logger.LogLn("Immune System:")
		for _, g := range immune {
			logger.LogF("Group %v contains %v units", g.ID, g.NumUnits)
		}
		sort.Sort(ByID(infection))
		logger.LogLn("Infection:")
		for _, g := range infection {
			logger.LogF("Group %v contains %v units", g.ID, g.NumUnits)
		}

		sort.Sort(sort.Reverse(ByPowerAndInitiative(infection)))
		sort.Sort(sort.Reverse(ByPowerAndInitiative(immune)))

		infectionAttacks := selectTargets(infection, immune)
		immuneAttacks := selectTargets(immune, infection)
		// logger.LogF("infection attacks: %v", infectionAttacks)
		// logger.LogF("immune attacks: %v", immuneAttacks)

		allAttacks := map[FullID]FullID{}
		for k, v := range infectionAttacks {
			allAttacks[k] = v
		}
		for k, v := range immuneAttacks {
			allAttacks[k] = v
		}

		allGroups := make([]*Group, len(infection)+len(immune))
		copy(allGroups, infection)
		copy(allGroups[len(infection):], immune)

		// By initiative highest to lowest
		sort.Sort(sort.Reverse(ByInitiative(allGroups)))

		attack(allGroups, allAttacks)

		immune = removeDead(immune)
		infection = removeDead(infection)
	}

	fmt.Printf("immune %v units\n", countUnits(immune))
	fmt.Printf("infection %v units\n", countUnits(infection))
}
