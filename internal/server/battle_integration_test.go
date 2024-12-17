package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"

	"github.com/gofiber/fiber/v2"
)

func TestBattle_IT(t *testing.T) {
	app := fiber.New()
	s := &FiberServer{App: app, diceSides: 6}

	databaseSrv := MustNewWithDatabase()
	pokemonSrv := database.NewPokemonService(databaseSrv)
	battleSrv := database.NewBattleService(databaseSrv)

	s.RegisterFiberRoutes(pokemonSrv, battleSrv)

	t.Run("create", func(t *testing.T) {
		t.Run("post-ok", func(t *testing.T) {
			battleReq := battleRequest{
				Pokemon1ID: 1,
				Pokemon2ID: 2,
			}
			body, _ := json.Marshal(battleReq)

			req, err := http.NewRequest("POST", "/battles", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum
			if err != nil {
				t.Fatalf("error creating request. Err: %v", err)
			}

			resp, err := app.Test(req, -1) // disable timeout
			if err != nil {
				t.Fatalf("error making request to server. Err: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("expected status Created; got %v", resp.Status)
			}

			var battleResponse models.Battle
			err = json.NewDecoder(resp.Body).Decode(&battleResponse)
			if err != nil {
				t.Fatalf("error decoding response. Err: %v", err)
			}

			if battleResponse.ID == 0 {
				t.Errorf("expected battle ID to be greater than 0; got %v", battleResponse.ID)
			}
			if battleResponse.Pokemon1ID != 1 || battleResponse.Pokemon2ID != 2 {
				t.Errorf("expected Pokemon1ID to be 1 and Pokemon2ID to be 2; got %v and %v", battleResponse.Pokemon1ID, battleResponse.Pokemon2ID)
			}
		})
	})

	t.Run("get-by-id", func(t *testing.T) {
		b := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
			Turns:      10,
		}

		// use the db layer to insert a battle
		err := battleSrv.Create(context.Background(), &b)
		if err != nil {
			t.Fatalf("error inserting battle. Err: %v", err)
		}

		req, err := http.NewRequest("GET", "/battles/"+strconv.Itoa(b.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req, -1) // disable timeout
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		var battleResponse models.Battle
		err = json.NewDecoder(resp.Body).Decode(&battleResponse)
		if err != nil {
			t.Fatalf("error decoding response. Err: %v", err)
		}

		if battleResponse.ID != b.ID {
			t.Errorf("expected battle ID to be %v; got %v", b.ID, battleResponse.ID)
		}
	})

	t.Run("get-all", func(t *testing.T) {
		initialBattles, err := battleSrv.GetAll(context.Background())
		if err != nil {
			t.Fatalf("error getting battles. Err: %v", err)
		}

		// insert 10 battles
		for i := 0; i < 10; i++ {
			b := models.Battle{
				Pokemon1ID: 1,
				Pokemon2ID: 2,
				WinnerID:   1,
				Turns:      10,
			}

			// use the db layer to insert a battle
			err := battleSrv.Create(context.Background(), &b)
			if err != nil {
				t.Fatalf("error inserting battle. Err: %v", err)
			}
		}

		req, err := http.NewRequest("GET", "/battles", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req, -1) // disable timeout
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		var battles []models.Battle
		err = json.NewDecoder(resp.Body).Decode(&battles)
		if err != nil {
			t.Fatalf("error decoding response. Err: %v", err)
		}

		if len(battles) != len(initialBattles)+10 {
			t.Errorf("expected %v battles; got %v", len(initialBattles)+10, len(battles))
		}
	})

	t.Run("update", func(t *testing.T) {
		b := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
			Turns:      10,
		}

		// use the db layer to insert a battle
		err := battleSrv.Create(context.Background(), &b)
		if err != nil {
			t.Fatalf("error inserting battle. Err: %v", err)
		}

		// update the battle
		b.WinnerID = 2
		b.Turns = 20

		body, err := json.Marshal(b)
		if err != nil {
			t.Fatalf("error marshalling battle. Err: %v", err)
		}

		req, err := http.NewRequest("PUT", "/battles/"+strconv.Itoa(b.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req, -1) // disable timeout
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		// get the battle from the db
		battle, err := battleSrv.GetByID(context.Background(), b.ID)
		if err != nil {
			t.Fatalf("error getting battle. Err: %v", err)
		}

		if battle.WinnerID != b.WinnerID || battle.Turns != b.Turns {
			t.Errorf("expected winner ID to be %v and turns to be %v; got %v and %v", b.WinnerID, b.Turns, battle.WinnerID, battle.Turns)
		}
	})

	t.Run("delete", func(t *testing.T) {
		b := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
			Turns:      10,
		}

		// use the db layer to insert a battle
		err := battleSrv.Create(context.Background(), &b)
		if err != nil {
			t.Fatalf("error inserting battle. Err: %v", err)
		}

		req, err := http.NewRequest("DELETE", "/battles/"+strconv.Itoa(b.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req, -1) // disable timeout
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status NoContent; got %v", resp.Status)
		}

		// get the battle from the db
		battle, err := battleSrv.GetByID(context.Background(), b.ID)
		if err == nil {
			t.Fatalf("expected error getting battle. Err: %v", err)
		}

		if battle.ID != 0 {
			t.Errorf("expected battle ID to be 0; got %v", battle.ID)
		}
	})
}
