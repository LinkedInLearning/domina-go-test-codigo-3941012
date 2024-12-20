package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"pokemon-battle/internal/models"

	"github.com/stretchr/testify/require"
)

// mockBattleService is used for testing the battle routes
// including the ability to return an error so we can test error handling
type mockBattleService struct {
	hasError bool
}

func (m *mockBattleService) Create(ctx context.Context, battle *models.Battle) error {
	if m.hasError {
		return errors.New("mock error")
	}
	battle.ID = 1
	return nil
}

func (m *mockBattleService) GetAll(ctx context.Context) ([]models.Battle, error) {
	if m.hasError {
		return nil, errors.New("mock error")
	}
	return []models.Battle{
		{ID: 1, Pokemon1ID: 1, Pokemon2ID: 2, WinnerID: 1},
		{ID: 2, Pokemon1ID: 3, Pokemon2ID: 4, WinnerID: 4},
	}, nil
}

func (m *mockBattleService) GetByID(ctx context.Context, id int) (models.Battle, error) {
	if m.hasError {
		return models.Battle{}, errors.New("mock error")
	}
	return models.Battle{ID: id, Pokemon1ID: 1, Pokemon2ID: 2, WinnerID: 1}, nil
}

func (m *mockBattleService) Update(ctx context.Context, battle models.Battle) error {
	if m.hasError {
		return errors.New("mock error")
	}
	return nil
}

func (m *mockBattleService) Delete(ctx context.Context, id int) error {
	if m.hasError {
		return errors.New("mock error")
	}
	return nil
}

func TestCreateBattle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()

		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that doesn't return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: false}, pokemonSrv: &mockPokemonService{hasError: false}, diceSides: 6}
		battleRoutes.Post("/", battleServer.CreateBattle)

		battle := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
		}
		body, _ := json.Marshal(battle)

		req, err := http.NewRequest("POST", "/battles", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusCreated)
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that returns an error
		// the pokemon service is mocked to not return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: true}, pokemonSrv: &mockPokemonService{hasError: false}, diceSides: 6}
		battleRoutes.Post("/", battleServer.CreateBattle)

		battle := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
		}
		body, err := json.Marshal(battle)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/battles", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)

		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("error/pokemon-failed", func(t *testing.T) {
		s := New()
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that returns an error
		// the pokemon service is mocked to not return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: false}, pokemonSrv: &mockPokemonService{hasError: true}, diceSides: 6}
		battleRoutes.Post("/", battleServer.CreateBattle)

		battle := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
		}
		body, err := json.Marshal(battle)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/battles", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)

		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}

func TestGetAllBattles(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that doesn't return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: false}}
		battleRoutes.Get("/", battleServer.GetAllBattles)

		req, err := http.NewRequest("GET", "/battles", nil)
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var battles []models.Battle
		err = json.Unmarshal(body, &battles)
		require.NoError(t, err)
		require.Equal(t, len(battles), 2)
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that returns an error
		battleServer := battleServer{srv: &mockBattleService{hasError: true}}
		battleRoutes.Get("/", battleServer.GetAllBattles)

		req, err := http.NewRequest("GET", "/battles", nil)
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}

func TestGetBattleByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that doesn't return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: false}}
		battleRoutes.Get("/:id", battleServer.GetBattleByID)

		req, err := http.NewRequest("GET", "/battles/1", nil)
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that returns an error
		battleServer := battleServer{srv: &mockBattleService{hasError: true}}
		battleRoutes.Get("/:id", battleServer.GetBattleByID)

		req, err := http.NewRequest("GET", "/battles/1", nil)
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)

		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}

func TestUpdateBattle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that doesn't return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: false}}
		battleRoutes.Put("/:id", battleServer.UpdateBattle)

		battle := models.Battle{
			ID:         1,
			Pokemon1ID: 1,
			Pokemon2ID: 2,
			WinnerID:   1,
		}
		body, err := json.Marshal(battle)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", "/battles/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusOK)
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that returns an error
		battleServer := battleServer{srv: &mockBattleService{hasError: true}}
		battleRoutes.Put("/:id", battleServer.UpdateBattle)

		battle := models.Battle{
			ID:         1,
			Pokemon1ID: 1,
			Pokemon2ID: 2,
		}
		body, err := json.Marshal(battle)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", "/battles/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)

		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}

func TestDeleteBattle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that doesn't return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: false}}
		battleRoutes.Delete("/:id", battleServer.DeleteBattle)

		req, err := http.NewRequest("DELETE", "/battles/1", nil)
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, resp.StatusCode, http.StatusNoContent)
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that returns an error
		battleServer := battleServer{srv: &mockBattleService{hasError: true}}
		battleRoutes.Delete("/:id", battleServer.DeleteBattle)

		req, err := http.NewRequest("DELETE", "/battles/1", nil)
		require.NoError(t, err)

		resp, err := s.App.Test(req)
		require.NoError(t, err)

		require.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}
