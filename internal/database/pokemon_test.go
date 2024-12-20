package database_test

import (
	"context"
	"database/sql"
	"testing"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"

	"github.com/stretchr/testify/require"
)

func TestNewPokemonService(t *testing.T) {
	srv := database.NewPokemonService(database.MustNewWithDatabase(t))
	require.NotNil(t, srv)

	t.Run("Create", func(t *testing.T) {
		pokemon := createTestPokemon(t, srv)
		defer cleanupPokemon(t, srv, pokemon.ID)

		require.Greater(t, pokemon.ID, 100)
	})

	t.Run("Delete", func(t *testing.T) {
		pokemon := createTestPokemon(t, srv)

		require.NoError(t, srv.Delete(context.Background(), pokemon.ID))

		_, err := srv.GetByID(context.Background(), pokemon.ID)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("GetAll", func(t *testing.T) {
		pokemons, err := srv.GetAll(context.Background())
		require.NoError(t, err)
		// There are 100 pokemons in the testdata/01-inserts.sql file
		require.Len(t, pokemons, 100)
	})

	t.Run("GetByID", func(t *testing.T) {
		pokemon, err := srv.GetByID(context.Background(), 1)
		require.NoError(t, err)
		require.Equal(t, 1, pokemon.ID)
		require.Equal(t, "Pikachu", pokemon.Name)
	})

	t.Run("Update", func(t *testing.T) {
		pokemon := createTestPokemon(t, srv)
		defer cleanupPokemon(t, srv, pokemon.ID)

		pokemon.Name = "Test Pikachu"

		require.NoError(t, srv.Update(context.Background(), pokemon))

		pokemon, err := srv.GetByID(context.Background(), pokemon.ID)
		require.NoError(t, err)
		require.Equal(t, "Test Pikachu", pokemon.Name)
	})
}

// createTestPokemon is a helper function to create a pokemon for testing
func createTestPokemon(t *testing.T, srv database.PokemonCRUDService) models.Pokemon {
	t.Helper()

	pokemon := models.Pokemon{
		Name:    "Pikachu",
		Type:    "Electric",
		HP:      100,
		Attack:  55,
		Defense: 40,
	}

	err := srv.Create(context.Background(), &pokemon)
	require.NoError(t, err)

	return pokemon
}

// cleanupPokemon is a helper function to delete a pokemon from the database
func cleanupPokemon(t *testing.T, srv database.PokemonCRUDService, id int) {
	t.Helper()

	err := srv.Delete(context.Background(), id)
	require.NoError(t, err)
}
