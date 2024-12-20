package database_test

import (
	"context"
	"database/sql"
	"testing"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"
)

func TestNewBattleService(t *testing.T) {
	dbService := database.MustNewWithDatabase()

	srv := database.NewBattleService(dbService)

	if srv == nil {
		t.Fatal("NewBattleService() returned nil")
	}

	t.Run("Create", func(t *testing.T) {
		battle := createTestBattle(t, srv)
		defer cleanupBattle(t, srv, battle.ID)

		if battle.ID != 1 {
			t.Fatalf("expected ID to be 1, got %d", battle.ID)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		battle := createTestBattle(t, srv)

		err := srv.Delete(context.Background(), battle.ID)
		if err != nil {
			t.Fatalf("expected Delete() to return nil, got %v", err)
		}

		_, err = srv.GetByID(context.Background(), battle.ID)
		if err != sql.ErrNoRows {
			t.Fatalf("expected GetByID() to return sql.ErrNoRows, got %v", err)
		}
	})

	t.Run("GetAll/zero", func(t *testing.T) {
		battles, err := srv.GetAll(context.Background())
		if err != nil {
			t.Fatalf("expected GetAll() to return nil, got %v", err)
		}

		if len(battles) != 0 {
			t.Fatalf("expected GetAll() to return 0 battles, got %d", len(battles))
		}
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
		if err != nil {
			t.Fatalf("expected Create() to return nil, got %v", err)
		}
		defer func() {
			err = pokemonSrv.Delete(context.Background(), pokemon1.ID)
			if err != nil {
				t.Fatalf("expected Delete() to return nil, got %v", err)
			}
		}()

		pokemon2 := models.Pokemon{
			Name:    "Charizard",
			Type:    "Fire",
			HP:      100,
			Attack:  55,
			Defense: 40,
		}

		err = pokemonSrv.Create(context.Background(), &pokemon2)
		if err != nil {
			t.Fatalf("expected Create() to return nil, got %v", err)
		}
		defer func() {
			err = pokemonSrv.Delete(context.Background(), pokemon2.ID)
			if err != nil {
				t.Fatalf("expected Delete() to return nil, got %v", err)
			}
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
			if err != nil {
				t.Fatalf("expected Create() to return nil, got %v", err)
			}
			defer func() {
				err = srv.Delete(context.Background(), battle.ID)
				if err != nil {
					t.Fatalf("expected Delete() to return nil, got %v", err)
				}
			}()
		}

		battles, err := srv.GetAll(context.Background())
		if err != nil {
			t.Fatalf("expected GetAll() to return nil, got %v", err)
		}

		if len(battles) != count {
			t.Fatalf("expected GetAll() to return %d battles, got %d", count, len(battles))
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		b := createTestBattle(t, srv)
		defer cleanupBattle(t, srv, b.ID)
		battle, err := srv.GetByID(context.Background(), b.ID)
		if err != nil {
			t.Fatalf("expected GetByID() to return nil, got %v", err)
		}

		if battle.ID != b.ID {
			t.Fatalf("expected ID to be %d, got %d", b.ID, battle.ID)
		}

		if battle.Turns != b.Turns {
			t.Fatalf("expected Turns to be %d, got %d", b.Turns, battle.Turns)
		}
	})

	t.Run("Update", func(t *testing.T) {
		battle := createTestBattle(t, srv)
		defer cleanupBattle(t, srv, battle.ID)
		battle.WinnerID = 2
		battle.Turns = 5

		err := srv.Update(context.Background(), battle)
		if err != nil {
			t.Fatalf("expected Update() to return nil, got %v", err)
		}

		battle, err = srv.GetByID(context.Background(), battle.ID)
		if err != nil {
			t.Fatalf("expected GetByID() to return nil, got %v", err)
		}

		if battle.WinnerID != 2 {
			t.Fatalf("expected winnerID to be 2, got %d", battle.WinnerID)
		}

		if battle.Turns != 5 {
			t.Fatalf("expected Turns to be 5, got %d", battle.Turns)
		}
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
	if err != nil {
		t.Fatalf("expected Create() to return nil, got %v", err)
	}
	return battle
}

// cleanupBattle is a helper function to delete a battle from the database
func cleanupBattle(t *testing.T, srv database.BattleCRUDService, id int) {
	t.Helper()

	err := srv.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("expected Delete() to return nil, got %v", err)
	}
}
