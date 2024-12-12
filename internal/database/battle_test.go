package database_test

import (
	"context"
	"database/sql"
	"testing"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"
)

func TestNewBattleService(t *testing.T) {
	srv := database.NewBattleService()

	if srv == nil {
		t.Fatal("NewBattleService() returned nil")
	}

	t.Run("Create", func(t *testing.T) {
		battle := &models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
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

		// There is no battle in the testdata/01-inserts.sql file, so the ID should be 1
		if battle.ID != 1 {
			t.Fatalf("expected ID to be 1, got %d", battle.ID)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		battle := &models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
		}

		err := srv.Create(context.Background(), battle)
		if err != nil {
			t.Fatalf("expected Create() to return nil, got %v", err)
		}

		err = srv.Delete(context.Background(), battle.ID)
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

		// There are no battles in the testdata/01-inserts.sql file, so the length should be 0
		if len(battles) != 0 {
			t.Fatalf("expected GetAll() to return 0 battles, got %d", len(battles))
		}
	})

	t.Run("GetAll/many", func(t *testing.T) {
		pokemonSrv := database.NewPokemonService()

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

		// There are no battles in the testdata/01-inserts.sql file, so the length should be 0
		if len(battles) != count {
			t.Fatalf("expected GetAll() to return %d battles, got %d", count, len(battles))
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		b := &models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
		}

		err := srv.Create(context.Background(), b)
		if err != nil {
			t.Fatalf("expected Create() to return nil, got %v", err)
		}
		defer func() {
			err = srv.Delete(context.Background(), b.ID)
			if err != nil {
				t.Fatalf("expected Delete() to return nil, got %v", err)
			}
		}()

		battle, err := srv.GetByID(context.Background(), b.ID)
		if err != nil {
			t.Fatalf("expected GetByID() to return nil, got %v", err)
		}

		if battle.ID != b.ID {
			t.Fatalf("expected ID to be %d, got %d", b.ID, battle.ID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		battle := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
		}

		err := srv.Create(context.Background(), &battle)
		if err != nil {
			t.Fatalf("expected Create() to return nil, got %v", err)
		}
		defer func() {
			err = srv.Delete(context.Background(), battle.ID)
			if err != nil {
				t.Fatalf("expected Delete() to return nil, got %v", err)
			}
		}()

		battle.WinnerID = 2

		err = srv.Update(context.Background(), battle)
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
	})
}
