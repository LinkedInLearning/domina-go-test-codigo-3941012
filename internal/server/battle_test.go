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

	"github.com/gofiber/fiber/v2"
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
		app := fiber.New()
		s := &FiberServer{App: app}
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that doesn't return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: false}}
		battleRoutes.Post("/", battleServer.CreateBattle)

		battle := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
		}
		body, _ := json.Marshal(battle)

		req, err := http.NewRequest("POST", "/battles", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status Created; got %v", resp.Status)
		}
	})

	t.Run("error", func(t *testing.T) {
		app := fiber.New()
		s := &FiberServer{App: app}
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that returns an error
		battleServer := battleServer{srv: &mockBattleService{hasError: true}}
		battleRoutes.Post("/", battleServer.CreateBattle)

		battle := models.Battle{
			Pokemon1ID: 1,
			Pokemon2ID: 2,
		}
		body, err := json.Marshal(battle)
		if err != nil {
			t.Fatalf("error marshalling battle. Err: %v", err)
		}

		req, err := http.NewRequest("POST", "/battles", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500; got %v", resp.Status)
		}
	})
}

func TestGetAllBattles(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := fiber.New()
		s := &FiberServer{App: app}
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that doesn't return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: false}}
		battleRoutes.Get("/", battleServer.GetAllBattles)

		req, err := http.NewRequest("GET", "/battles", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("error reading response body. Err: %v", err)
		}

		var battles []models.Battle
		err = json.Unmarshal(body, &battles)
		if err != nil {
			t.Fatalf("error unmarshalling response body. Err: %v", err)
		}
		if len(battles) != 2 {
			t.Errorf("expected 2 battles; got %v", len(battles))
		}
	})

	t.Run("error", func(t *testing.T) {
		app := fiber.New()
		s := &FiberServer{App: app}
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that returns an error
		battleServer := battleServer{srv: &mockBattleService{hasError: true}}
		battleRoutes.Get("/", battleServer.GetAllBattles)

		req, err := http.NewRequest("GET", "/battles", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500; got %v", resp.Status)
		}
	})
}

func TestGetBattleByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := fiber.New()
		s := &FiberServer{App: app}
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that doesn't return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: false}}
		battleRoutes.Get("/:id", battleServer.GetBattleByID)

		req, err := http.NewRequest("GET", "/battles/1", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}
	})

	t.Run("error", func(t *testing.T) {
		app := fiber.New()
		s := &FiberServer{App: app}
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that returns an error
		battleServer := battleServer{srv: &mockBattleService{hasError: true}}
		battleRoutes.Get("/:id", battleServer.GetBattleByID)

		req, err := http.NewRequest("GET", "/battles/1", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500; got %v", resp.Status)
		}
	})
}

func TestUpdateBattle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := fiber.New()
		s := &FiberServer{App: app}
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
		if err != nil {
			t.Fatalf("error marshalling battle. Err: %v", err)
		}

		req, err := http.NewRequest("PUT", "/battles/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}
	})

	t.Run("error", func(t *testing.T) {
		app := fiber.New()
		s := &FiberServer{App: app}
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
		if err != nil {
			t.Fatalf("error marshalling battle. Err: %v", err)
		}

		req, err := http.NewRequest("PUT", "/battles/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500; got %v", resp.Status)
		}
	})
}

func TestDeleteBattle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := fiber.New()
		s := &FiberServer{App: app}
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that doesn't return an error
		battleServer := battleServer{srv: &mockBattleService{hasError: false}}
		battleRoutes.Delete("/:id", battleServer.DeleteBattle)

		req, err := http.NewRequest("DELETE", "/battles/1", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status NoContent; got %v", resp.Status)
		}
	})

	t.Run("error", func(t *testing.T) {
		app := fiber.New()
		s := &FiberServer{App: app}
		battleRoutes := s.App.Group("/battles")

		// init the battle routes from a mock battle service that returns an error
		battleServer := battleServer{srv: &mockBattleService{hasError: true}}
		battleRoutes.Delete("/:id", battleServer.DeleteBattle)

		req, err := http.NewRequest("DELETE", "/battles/1", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500; got %v", resp.Status)
		}
	})
}
