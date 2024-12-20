package business_test

import (
	"pokemon-battle/internal/business"
	"testing"
)

func TestFight(t *testing.T) {
	t.Run("weak-second/first-wins", func(t *testing.T) {
		battle := business.Fight(10, strongPokemon, weakPokemon)
		if battle.WinnerID != strongPokemon.ID {
			t.Fatalf("expected winner ID to be %d, got %d", strongPokemon.ID, battle.WinnerID)
		}
		if battle.Turns != 1 {
			t.Fatalf("expected turns to be 1, got %d", battle.Turns)
		}

		// verify that the pokemons return in the same state as they were before the fight
		if weakPokemon.HP != 1 {
			t.Fatalf("expected weakPokemon HP to be 1, got %d", weakPokemon.HP)
		}
		if strongPokemon.HP != 100 {
			t.Fatalf("expected strongPokemon HP to be 100, got %d", strongPokemon.HP)
		}
	})

	t.Run("strong-second/second-wins", func(t *testing.T) {
		battle := business.Fight(10, weakPokemon, strongPokemon)
		if battle.WinnerID != strongPokemon.ID {
			t.Fatalf("expected winner ID to be %d, got %d", strongPokemon.ID, battle.WinnerID)
		}
		if battle.Turns != 1 {
			t.Fatalf("expected turns to be 1, got %d", battle.Turns)
		}

		// verificar que los pokemons retornen en el mismo estado
		// que antes de la batalla
		if weakPokemon.HP != 1 {
			t.Fatalf("expected weakPokemon HP to be 1, got %d", weakPokemon.HP)
		}
		if strongPokemon.HP != 100 {
			t.Fatalf("expected strongPokemon HP to be 100, got %d", strongPokemon.HP)
		}
	})

	t.Run("equals", func(t *testing.T) {
		battle := business.Fight(10, strongPokemon, strongPokemon)
		if battle.Turns <= 1 {
			t.Fatalf("expected turns to be greater than 1, got %d", battle.Turns)
		}

		// verificar que los pokemons retornen en el mismo estado
		// que antes de la batalla
		if strongPokemon.HP != 100 {
			t.Fatalf("expected p1 HP to be 100, got %d", strongPokemon.HP)
		}
	})
}
