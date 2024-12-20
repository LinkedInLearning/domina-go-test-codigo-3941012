package database_test

import (
	"context"
	"database/sql"
	"testing"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"
)

func TestNewPokemonService(t *testing.T) {
	srv := database.NewPokemonService(database.MustNewWithDatabase())

	if srv == nil {
		t.Fatal("NewPokemonService() returned nil")
	}

	t.Run("Create", func(t *testing.T) {
		pokemon := createTestPokemon(t, srv)
		defer cleanupPokemon(t, srv, pokemon.ID)

		if pokemon.ID <= 100 {
			t.Fatalf("expected ID to be greater than 100, got %d", pokemon.ID)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		pokemon := createTestPokemon(t, srv)

		err := srv.Delete(context.Background(), pokemon.ID)
		if err != nil {
			t.Fatalf("expected Delete() to return nil, got %v", err)
		}

		_, err = srv.GetByID(context.Background(), pokemon.ID)
		if err != sql.ErrNoRows {
			t.Fatalf("expected GetByID() to return sql.ErrNoRows, got %v", err)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		pokemons, err := srv.GetAll(context.Background())
		if err != nil {
			t.Fatalf("expected GetAll() to return nil, got %v", err)
		}

		// There are 100 pokemons in the testdata/01-inserts.sql file
		if len(pokemons) != 100 {
			t.Fatalf("expected GetAll() to return 100 pokemons, got %d", len(pokemons))
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		pokemon, err := srv.GetByID(context.Background(), 1)
		if err != nil {
			t.Fatalf("expected GetByID() to return nil, got %v", err)
		}

		if pokemon.ID != 1 {
			t.Fatalf("expected ID to be 1, got %d", pokemon.ID)
		}

		if pokemon.Name != "Pikachu" {
			t.Fatalf("expected name to be 'Pikachu', got %s", pokemon.Name)
		}
	})

	t.Run("Update", func(t *testing.T) {
		pokemon := createTestPokemon(t, srv)
		defer cleanupPokemon(t, srv, pokemon.ID)

		pokemon.Name = "Test Pikachu"

		err := srv.Update(context.Background(), pokemon)
		if err != nil {
			t.Fatalf("expected Update() to return nil, got %v", err)
		}

		pokemon, err = srv.GetByID(context.Background(), pokemon.ID)
		if err != nil {
			t.Fatalf("expected GetByID() to return nil, got %v", err)
		}

		if pokemon.Name != "Test Pikachu" {
			t.Fatalf("expected name to be 'Test Pikachu', got %s", pokemon.Name)
		}
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
	if err != nil {
		t.Fatalf("expected Create() to return nil, got %v", err)
	}
	return pokemon
}

// cleanupPokemon is a helper function to delete a pokemon from the database
func cleanupPokemon(t *testing.T, srv database.PokemonCRUDService, id int) {
	t.Helper()

	err := srv.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("expected Delete() to return nil, got %v", err)
	}
}
