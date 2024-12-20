package database_test

import (
	"context"
	"database/sql"
	"testing"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"

	"github.com/stretchr/testify/require"
)

func TestNewBattleService(t *testing.T) {
	dbService := database.MustNewWithDatabase(t)

	srv := database.NewBattleService(dbService)
	require.NotNil(t, srv)

	t.Run("Create", func(t *testing.T) {
		battle := createTestBattle(t, srv)
		defer cleanupBattle(t, srv, battle.ID)

		require.Equal(t, 1, battle.ID)
	})

	t.Run("Delete", func(t *testing.T) {
		battle := createTestBattle(t, srv)

		err := srv.Delete(context.Background(), battle.ID)
		require.NoError(t, err)

		_, err = srv.GetByID(context.Background(), battle.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("GetAll/zero", func(t *testing.T) {
		battles, err := srv.GetAll(context.Background())
		require.NoError(t, err)
		require.Equal(t, 0, len(battles))
	})

	t.Run("GetAll/many", func(t *testing.T) {
		pokemonSrv := database.NewPokemonService(dbService)

		pokemon1 := models.Pokemon{
			Name:    "Pikachu",
			Type:    "Electric",
			HP:      100,
			Attack:  55,
			Defense: 40,
		}

		err := pokemonSrv.Create(context.Background(), &pokemon1)
		require.NoError(t, err)
		defer func() {
			err = pokemonSrv.Delete(context.Background(), pokemon1.ID)
			require.NoError(t, err)
		}()

		pokemon2 := models.Pokemon{
			Name:    "Charizard",
			Type:    "Fire",
			HP:      100,
			Attack:  55,
			Defense: 40,
		}

		err = pokemonSrv.Create(context.Background(), &pokemon2)
		require.NoError(t, err)
		defer func() {
			err = pokemonSrv.Delete(context.Background(), pokemon2.ID)
			require.NoError(t, err)
		}()

		count := 100
		for i := 1; i <= count; i++ {
			battle := &models.Battle{
				Pokemon1ID: pokemon1.ID,
				Pokemon2ID: pokemon2.ID,
				WinnerID:   pokemon2.ID, // Charizard always wins
				Turns:      10,
			}
			err := srv.Create(context.Background(), battle)
			require.NoError(t, err)
			defer func() {
				err = srv.Delete(context.Background(), battle.ID)
				require.NoError(t, err)
			}()
		}

		battles, err := srv.GetAll(context.Background())
		require.NoError(t, err)
		require.Equal(t, count, len(battles))
	})

	t.Run("GetByID", func(t *testing.T) {
		b := createTestBattle(t, srv)
		defer cleanupBattle(t, srv, b.ID)

		battle, err := srv.GetByID(context.Background(), b.ID)
		require.NoError(t, err)
		require.Equal(t, b.ID, battle.ID)
		require.Equal(t, b.Turns, battle.Turns)
	})

	t.Run("Update", func(t *testing.T) {
		battle := createTestBattle(t, srv)
		defer cleanupBattle(t, srv, battle.ID)
		battle.WinnerID = 2
		battle.Turns = 5

		err := srv.Update(context.Background(), battle)
		require.NoError(t, err)

		battle, err = srv.GetByID(context.Background(), battle.ID)
		require.NoError(t, err)
		require.Equal(t, 2, battle.WinnerID)
		require.Equal(t, 5, battle.Turns)
	})
}

// createTestBattle is a helper function to create a battle for testing
func createTestBattle(t *testing.T, srv database.BattleCRUDService) models.Battle {
	t.Helper()

	battle := models.Battle{
		Pokemon1ID: 1,
		Pokemon2ID: 2,
		WinnerID:   1,
		Turns:      10,
	}

	err := srv.Create(context.Background(), &battle)
	require.NoError(t, err)

	return battle
}

// cleanupBattle is a helper function to delete a battle from the database
func cleanupBattle(t *testing.T, srv database.BattleCRUDService, id int) {
	t.Helper()

	err := srv.Delete(context.Background(), id)
	require.NoError(t, err)
}
