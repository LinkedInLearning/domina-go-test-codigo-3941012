package business_test

import (
	"pokemon-battle/internal/business"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFight(t *testing.T) {
	t.Run("weak-second/first-wins", func(t *testing.T) {
		battle := business.Fight(10, strongPokemon, weakPokemon)
		require.Equal(t, strongPokemon.ID, battle.WinnerID)
		require.Equal(t, 1, battle.Turns)
		require.Equal(t, weakPokemon.HP, 1)
		require.Equal(t, strongPokemon.HP, 100)
	})

	t.Run("strong-second/second-wins", func(t *testing.T) {
		battle := business.Fight(10, weakPokemon, strongPokemon)
		require.Equal(t, strongPokemon.ID, battle.WinnerID)
		require.Equal(t, 1, battle.Turns)
		require.Equal(t, weakPokemon.HP, 1)
		require.Equal(t, strongPokemon.HP, 100)
	})

	t.Run("equals", func(t *testing.T) {
		battle := business.Fight(10, strongPokemon, strongPokemon)
		require.Equal(t, strongPokemon.ID, battle.WinnerID)
		require.Greater(t, battle.Turns, 1)
		require.Equal(t, strongPokemon.HP, 100)
	})
}
