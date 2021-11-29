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
	"os"
	"reflect"
	"testing"

	"character"
	"logger"
)

func makePlayers() (player1, player2 character.Player, boss1, boss2 character.Boss) {
	player1 = *character.NewPlayer(1000, 1000)
	player2 = player1

	boss1 = *character.NewBoss(1000, 10)
	boss2 = boss1
	return
}

func TestMagicMissile(t *testing.T) {
	var s Spell = NewMagicMissile()

	if s.IsActive() {
		t.Errorf("initially active")
	}

	expectedPlayer, player, expectedBoss, boss := makePlayers()
	s.Activate(&player, &boss)

	expectedBoss.HP -= 4

	if !reflect.DeepEqual(expectedPlayer, player) {
		t.Errorf("activation changed player; want %+v, got %+v", expectedPlayer, player)
	}

	if !reflect.DeepEqual(expectedBoss, boss) {
		t.Errorf("unexpected boss change; want %+v, got %+v", expectedBoss, boss)
	}

	if s.IsActive() {
		t.Errorf("late active")
	}
}

func TestDrain(t *testing.T) {
	var s Spell = NewDrain()

	if s.IsActive() {
		t.Errorf("initially active")
	}

	expectedPlayer, player, expectedBoss, boss := makePlayers()
	s.Activate(&player, &boss)

	expectedPlayer.HP += 2
	expectedBoss.HP -= 2

	if !reflect.DeepEqual(expectedPlayer, player) {
		t.Errorf("activation changed player; want %+v, got %+v", expectedPlayer, player)
	}
	if !reflect.DeepEqual(expectedBoss, boss) {
		t.Errorf("unexpected boss change; want %+v, got %+v", expectedBoss, boss)
	}

	if s.IsActive() {
		t.Errorf("late active")
	}
}

func TestShield(t *testing.T) {
	var s Spell = NewShield()

	if s.IsActive() {
		t.Errorf("initially active")
	}

	refPlayer, player, refBoss, boss := makePlayers()
	s.Activate(&player, &boss)

	boostedPlayer := refPlayer
	boostedPlayer.Armor += 7

	if !reflect.DeepEqual(boostedPlayer, player) {
		t.Errorf("activation changed player; want %+v, got %+v", boostedPlayer, player)
	}
	if !reflect.DeepEqual(refBoss, boss) {
		t.Errorf("unexpected boss change; want %+v, got %+v", refBoss, boss)
	}

	for i := 1; i <= 6; i++ {
		if !s.IsActive() {
			t.Errorf("spell unexpectedly inactive at round %v")
			break
		}

		s.TurnStart(&player, &boss)
	}

	if s.IsActive() {
		t.Errorf("spell unexpectedly active at end")
	}

	if !reflect.DeepEqual(refPlayer, player) {
		t.Errorf("unrestored player; want %+v, got %+v", boostedPlayer, player)
	}
	if !reflect.DeepEqual(refBoss, boss) {
		t.Errorf("unexpected boss change; want %+v, got %+v", refBoss, boss)
	}
}

func TestPoison(t *testing.T) {
	var s Spell = NewPoison()

	if s.IsActive() {
		t.Errorf("initially active")
	}

	refPlayer, player, expectedBoss, boss := makePlayers()
	s.Activate(&player, &boss)

	if !reflect.DeepEqual(refPlayer, player) {
		t.Errorf("activation changed player; want %+v, got %+v", refPlayer, player)
	}
	if !reflect.DeepEqual(expectedBoss, boss) {
		t.Errorf("unexpected boss change; want %+v, got %+v", expectedBoss, boss)
	}

	for i := 1; i <= 6; i++ {
		if !s.IsActive() {
			t.Errorf("spell unexpectedly inactive at round %v")
			break
		}

		s.TurnStart(&player, &boss)

		expectedBoss.HP -= 3
		if !reflect.DeepEqual(expectedBoss, boss) {
			t.Errorf("unexpected boss change round %v; want %+v, got %+v",
				i, expectedBoss, boss)
		}
	}

	logger.LogF("boss now %+v\n", boss)

	if s.IsActive() {
		t.Errorf("spell unexpectedly active at end")
	}

	if !reflect.DeepEqual(refPlayer, player) {
		t.Errorf("player changed at end; want %+v, got %+v", refPlayer, player)
	}
}

func TestRecharge(t *testing.T) {
	var s Spell = NewRecharge()

	if s.IsActive() {
		t.Errorf("initially active")
	}

	expectedPlayer, player, refBoss, boss := makePlayers()
	s.Activate(&player, &boss)

	if !reflect.DeepEqual(expectedPlayer, player) {
		t.Errorf("activation changed player; want %+v, got %+v", expectedPlayer, player)
	}
	if !reflect.DeepEqual(refBoss, boss) {
		t.Errorf("unexpected boss change; want %+v, got %+v", refBoss, boss)
	}

	for i := 1; i <= 5; i++ {
		if !s.IsActive() {
			t.Errorf("spell unexpectedly inactive at round %v")
			break
		}

		s.TurnStart(&player, &boss)

		expectedPlayer.Mana += 101
		if !reflect.DeepEqual(expectedPlayer, player) {
			t.Errorf("unexpected player change round %v; want %+v, got %+v",
				i, expectedPlayer, player)
		}
		if !reflect.DeepEqual(refBoss, boss) {
			t.Errorf("unexpected boss change round %v; want %+v, got %+v",
				i, refBoss, boss)
		}
	}

	logger.LogF("player now %+v\n", player)

	if s.IsActive() {
		t.Errorf("spell unexpectedly active at end")
	}

	if !reflect.DeepEqual(expectedPlayer, player) {
		t.Errorf("player changed at end; want %+v, got %+v", expectedPlayer, player)
	}
}

func TestNames(t *testing.T) {
	s := []Spell{NewMagicMissile(), NewPoison(), NewRecharge(), NewDrain(), NewShield()}
	expectedNames := []string{"MagicMissile", "Poison", "Recharge", "Drain", "Shield"}

	if names := Names(s); !reflect.DeepEqual(expectedNames, names) {
		t.Errorf("Names(...) = %v, want %v", names, expectedNames)
	}
}

func TestMain(m *testing.M) {
	logger.Init(true)
	os.Exit(m.Run())
}
