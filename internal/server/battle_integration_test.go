package server

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"pokemon-battle/internal/database"
	"pokemon-battle/internal/models"

	"github.com/stretchr/testify/require"
)

func TestBattle_IT(t *testing.T) {
	s := New()
	s.diceSides = 6

	databaseSrv := MustNewWithDatabase(t)
	pokemonSrv := database.NewPokemonService(databaseSrv)
	battleSrv := database.NewBattleService(databaseSrv)

	s.RegisterFiberRoutes(pokemonSrv, battleSrv)

	t.Run("create", func(t *testing.T) {
		t.Run("post-ok", func(t *testing.T) {
			battleReq := battleRequest{
				Pokemon1ID: 1,
				Pokemon2ID: 2,
			}
			body, err := json.Marshal(battleReq)
			require.NoError(t, err)

			req := createAuthenticatedRequest(t, "POST", "/battles", body)

			resp, err := s.App.Test(req, -1) // disable timeout
			require.NoError(t, err)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("expected status Created; got %v", resp.Status)
			}

			var battleResponse models.Battle
			err = json.NewDecoder(resp.Body).Decode(&battleResponse)
			require.NoError(t, err)
			require.Greater(t, battleResponse.ID, 0)
			require.Equal(t, battleResponse.Pokemon1ID, 1)
			require.Equal(t, battleResponse.Pokemon2ID, 2)
		})
	})

	t.Run("get-by-id", func(t *testing.T) {
		b := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
			Turns:      10,
		}

		createBattleUsingDB(t, battleSrv, &b)

		req := createAuthenticatedRequest(t, "GET", "/battles/"+strconv.Itoa(b.ID), nil)

		resp, err := s.App.Test(req, -1) // disable timeout
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		var battleResponse models.Battle
		err = json.NewDecoder(resp.Body).Decode(&battleResponse)
		require.NoError(t, err)
		require.Equal(t, battleResponse.ID, b.ID)
	})

	t.Run("get-all", func(t *testing.T) {
		initialBattles, err := battleSrv.GetAll(context.Background())
		require.NoError(t, err)

		// insert 10 battles
		for i := 0; i < 10; i++ {
			b := models.Battle{
				Pokemon1ID: 1,
				Pokemon2ID: 2,
				WinnerID:   1,
				Turns:      10,
			}

			createBattleUsingDB(t, battleSrv, &b)
		}

		req := createAuthenticatedRequest(t, "GET", "/battles", nil)

		resp, err := s.App.Test(req, -1) // disable timeout
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)

		var battles []models.Battle
		err = json.NewDecoder(resp.Body).Decode(&battles)
		require.NoError(t, err)
		require.Equal(t, len(battles), len(initialBattles)+10)
	})

	t.Run("update", func(t *testing.T) {
		b := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
			Turns:      10,
		}

		createBattleUsingDB(t, battleSrv, &b)

		// update the battle
		b.WinnerID = 2
		b.Turns = 20

		body, err := json.Marshal(b)
		require.NoError(t, err)

		req := createAuthenticatedRequest(t, "PUT", "/battles/"+strconv.Itoa(b.ID), body)

		resp, err := s.App.Test(req, -1) // disable timeout
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)

		// get the battle from the db
		battle, err := battleSrv.GetByID(context.Background(), b.ID)
		require.NoError(t, err)

		require.Equal(t, battle.WinnerID, b.WinnerID)
		require.Equal(t, battle.Turns, b.Turns)
	})

	t.Run("delete", func(t *testing.T) {
		b := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
			Turns:      10,
		}

		createBattleUsingDB(t, battleSrv, &b)

		req := createAuthenticatedRequest(t, "DELETE", "/battles/"+strconv.Itoa(b.ID), nil)

		resp, err := s.App.Test(req, -1) // disable timeout
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusNoContent)

		// get the battle from the db
		battle, err := battleSrv.GetByID(context.Background(), b.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
		require.Equal(t, battle.ID, 0)
	})
}

// createAuthenticatedRequest is a helper function to create an authenticated request,
// passing the correct headers and body: Content-Type: application/json and Authorization: Basic YXNoOmtldGNodW0=,
// which is the base64 encoded string for the username and password used in tests: ash:ketchum
func createAuthenticatedRequest(t *testing.T, method, path string, body []byte) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, path, bytes.NewBuffer(body))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic YXNoOmtldGNodW0=") // base64 for ash:ketchum

	return req
}

// createBattleUsingDB is a helper function to create a battle using the database service
func createBattleUsingDB(t *testing.T, srv database.BattleCRUDService, b *models.Battle) {
	t.Helper()

	err := srv.Create(context.Background(), b)
	require.NoError(t, err)
}
