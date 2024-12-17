package business_test

import (
	"pokemon-battle/internal/business"
	"pokemon-battle/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
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

		battle, err := business.Fight(10, p1, p2)
		require.NoError(t, err)
		require.Equal(t, battle.WinnerID, p1.ID)
		require.Equal(t, battle.Turns, 1)
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

		battle, err := business.Fight(10, p1, p2)
		require.NoError(t, err)
		require.Equal(t, battle.WinnerID, p2.ID)
		require.Equal(t, battle.Turns, 1)
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

		battle, err := business.Fight(10, p1, p2)
		require.NoError(t, err)
		require.Greater(t, battle.Turns, 1)
	})
}
