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
)

// mockPokemonService is used for testing the pokemon routes
// including the ability to return an error so we can test error handling
type mockPokemonService struct {
	hasError bool
}

func (m *mockPokemonService) Create(ctx context.Context, pokemon *models.Pokemon) error {
	if m.hasError {
		return errors.New("mock error")
	}
	return nil
}

func (m *mockPokemonService) Delete(ctx context.Context, id int) error {
	if m.hasError {
		return errors.New("mock error")
	}
	return nil
}

func (m *mockPokemonService) GetAll(ctx context.Context) ([]models.Pokemon, error) {
	if m.hasError {
		return nil, errors.New("mock error")
	}
	return []models.Pokemon{
		{ID: 1, Name: "Pikachu", Type: "Electric", HP: 10, Attack: 10, Defense: 10},
		{ID: 2, Name: "Charmander", Type: "Fire", HP: 10, Attack: 10, Defense: 10},
	}, nil
}

func (m *mockPokemonService) GetByID(ctx context.Context, id int) (models.Pokemon, error) {
	if m.hasError {
		return models.Pokemon{}, errors.New("mock error")
	}
	return models.Pokemon{}, nil
}

func (m *mockPokemonService) Update(ctx context.Context, pokemon models.Pokemon) error {
	if m.hasError {
		return errors.New("mock error")
	}
	return nil
}

func TestGetAllPokemons(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()

		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that doesn't return an error
		pokemonServer := pokemonServer{srv: &mockPokemonService{hasError: false}}
		pokemonRoutes.Get("/", pokemonServer.GetAllPokemons)

		// Create a test HTTP request
		req, err := http.NewRequest("GET", "/pokemons", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}
		// Perform the request
		resp, err := s.App.Test(req)
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

		var pokemons []models.Pokemon
		err = json.Unmarshal(body, &pokemons)
		if err != nil {
			t.Fatalf("error unmarshalling response body. Err: %v", err)
		}
		if len(pokemons) != 2 {
			t.Errorf("expected 2 pokemons; got %v", len(pokemons))
		}

		if pokemons[0].ID != 1 || pokemons[0].Name != "Pikachu" || pokemons[0].Type != "Electric" || pokemons[0].HP != 10 || pokemons[0].Attack != 10 || pokemons[0].Defense != 10 {
			t.Errorf("expected Pikachu; got %v", pokemons[0])
		}
		if pokemons[1].ID != 2 || pokemons[1].Name != "Charmander" || pokemons[1].Type != "Fire" || pokemons[1].HP != 10 || pokemons[1].Attack != 10 || pokemons[1].Defense != 10 {
			t.Errorf("expected Charmander; got %v", pokemons[1])
		}
	})

	t.Run("error", func(t *testing.T) {
		s := New()

		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that returns an error
		pokemonServer := pokemonServer{srv: &mockPokemonService{hasError: true}}
		pokemonRoutes.Get("/", pokemonServer.GetAllPokemons)

		// Create a test HTTP request
		req, err := http.NewRequest("GET", "/pokemons", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}
		// Perform the request
		resp, err := s.App.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500; got %v", resp.Status)
		}
	})
}

func TestCreatePokemon(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that doesn't return an error
		pokemonServer := pokemonServer{srv: &mockPokemonService{hasError: false}}
		pokemonRoutes.Post("/", pokemonServer.CreatePokemon)

		pokemon := models.Pokemon{
			Name:    "Bulbasaur",
			Type:    "Grass",
			HP:      45,
			Attack:  49,
			Defense: 49,
		}
		body, err := json.Marshal(pokemon)
		if err != nil {
			t.Fatalf("error marshalling pokemon. Err: %v", err)
		}

		req, err := http.NewRequest("POST", "/pokemons", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := s.App.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status Created; got %v", resp.Status)
		}
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that returns an error
		pokemonServer := pokemonServer{srv: &mockPokemonService{hasError: true}}
		pokemonRoutes.Post("/", pokemonServer.CreatePokemon)

		pokemon := models.Pokemon{
			Name: "Bulbasaur",
			Type: "Grass",
		}
		body, err := json.Marshal(pokemon)
		if err != nil {
			t.Fatalf("error marshalling pokemon. Err: %v", err)
		}

		req, err := http.NewRequest("POST", "/pokemons", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := s.App.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500; got %v", resp.Status)
		}
	})
}

func TestGetPokemonByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that doesn't return an error
		pokemonServer := pokemonServer{srv: &mockPokemonService{hasError: false}}
		pokemonRoutes.Get("/:id", pokemonServer.GetPokemonByID)

		req, err := http.NewRequest("GET", "/pokemons/1", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := s.App.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that returns an error
		pokemonServer := pokemonServer{srv: &mockPokemonService{hasError: true}}
		pokemonRoutes.Get("/:id", pokemonServer.GetPokemonByID)

		req, err := http.NewRequest("GET", "/pokemons/1", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := s.App.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500; got %v", resp.Status)
		}
	})
}

func TestUpdatePokemon(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that doesn't return an error
		pokemonServer := pokemonServer{srv: &mockPokemonService{hasError: false}}
		pokemonRoutes.Put("/:id", pokemonServer.UpdatePokemon)

		pokemon := models.Pokemon{
			ID:      1,
			Name:    "Bulbasaur",
			Type:    "Grass",
			HP:      45,
			Attack:  49,
			Defense: 49,
		}
		body, _ := json.Marshal(pokemon)

		req, err := http.NewRequest("PUT", "/pokemons/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := s.App.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that returns an error
		pokemonServer := pokemonServer{srv: &mockPokemonService{hasError: true}}
		pokemonRoutes.Put("/:id", pokemonServer.UpdatePokemon)

		pokemon := models.Pokemon{
			ID:   1,
			Name: "Bulbasaur",
			Type: "Grass",
		}
		body, _ := json.Marshal(pokemon)

		req, err := http.NewRequest("PUT", "/pokemons/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := s.App.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500; got %v", resp.Status)
		}
	})
}

func TestDeletePokemon(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that doesn't return an error
		pokemonServer := pokemonServer{srv: &mockPokemonService{hasError: false}}
		pokemonRoutes.Delete("/:id", pokemonServer.DeletePokemon)

		req, err := http.NewRequest("DELETE", "/pokemons/1", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := s.App.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected status NoContent; got %v", resp.Status)
		}
	})

	t.Run("error", func(t *testing.T) {
		s := New()
		pokemonRoutes := s.App.Group("/pokemons")

		// init the pokemon routes from a mock pokemon service that returns an error
		pokemonServer := pokemonServer{srv: &mockPokemonService{hasError: true}}
		pokemonRoutes.Delete("/:id", pokemonServer.DeletePokemon)

		req, err := http.NewRequest("DELETE", "/pokemons/1", nil)
		if err != nil {
			t.Fatalf("error creating request. Err: %v", err)
		}

		resp, err := s.App.Test(req)
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500; got %v", resp.Status)
		}
	})
}
