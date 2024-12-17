package business_test

import (
	"pokemon-battle/internal/business"
	"pokemon-battle/internal/models"
	"testing"
)

func TestFight(t *testing.T) {
	t.Run("weak-second/first-wins", func(t *testing.T) {
		p1 := models.Pokemon{
			ID:      1,
			Name:    "Pikachu",
			Type:    "Electric",
			HP:      100,
			Attack:  55,
			Defense: 40,
		}

		p2 := models.Pokemon{
			ID:      2,
			Name:    "Charmander",
			Type:    "Fire",
			HP:      1,
			Attack:  1,
			Defense: 1,
		}

		battle := business.Fight(10, p1, p2)
		if battle.WinnerID != p1.ID {
			t.Fatalf("expected winner ID to be %d, got %d", p1.ID, battle.WinnerID)
		}
		if battle.Turns != 1 {
			t.Fatalf("expected turns to be 1, got %d", battle.Turns)
		}
	})

	t.Run("strong-second/second-wins", func(t *testing.T) {
		p1 := models.Pokemon{
			ID:      2,
			Name:    "Charmander",
			Type:    "Fire",
			HP:      1,
			Attack:  1,
			Defense: 1,
		}

		p2 := models.Pokemon{
			ID:      1,
			Name:    "Pikachu",
			Type:    "Electric",
			HP:      100,
			Attack:  55,
			Defense: 40,
		}

		battle := business.Fight(10, p1, p2)
		if battle.WinnerID != p2.ID {
			t.Fatalf("expected winner ID to be %d, got %d", p2.ID, battle.WinnerID)
		}
		if battle.Turns != 1 {
			t.Fatalf("expected turns to be 1, got %d", battle.Turns)
		}
	})

	t.Run("equals", func(t *testing.T) {
		p1 := models.Pokemon{
			ID:      2,
			Name:    "Charmander",
			Type:    "Fire",
			HP:      100,
			Attack:  55,
			Defense: 40,
		}

		p2 := models.Pokemon{
			ID:      1,
			Name:    "Pikachu",
			Type:    "Electric",
			HP:      100,
			Attack:  55,
			Defense: 40,
		}

		battle := business.Fight(10, p1, p2)
		if battle.Turns <= 1 {
			t.Fatalf("expected turns to be greater than 1, got %d", battle.Turns)
		}
		if p1.HP != 100 {
			t.Fatalf("expected p1 HP to be 100, got %d", p1.HP)
		}
		if p2.HP != 100 {
			t.Fatalf("expected p2 HP to be 100, got %d", p2.HP)
		}
	})
}
